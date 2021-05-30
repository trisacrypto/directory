package gds

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/sendgrid-go"
	"github.com/trisacrypto/directory/pkg"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v1"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/sectigo"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// New creates a TRISA Directory Service with the specified configuration and prepares
// it to listen for and serve GRPC requests.
func New(conf config.Config) (s *Server, err error) {
	// Load the default configuration from the environment
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Set the global level
	zerolog.SetGlobalLevel(zerolog.Level(conf.LogLevel))

	// Create the server and open the connection to the database
	s = &Server{conf: conf, echan: make(chan error, 1)}
	if s.db, err = store.Open(conf.DatabaseURL); err != nil {
		return nil, err
	}

	// Create the Sectigo API client
	if s.certs, err = sectigo.New(conf.Sectigo.Username, conf.Sectigo.Password); err != nil {
		return nil, err
	}

	// Ensure the certificate storage can be reached
	if _, err = s.getCertStorage(); err != nil {
		return nil, err
	}

	// Create the SendGrid API client
	s.email = sendgrid.NewSendClient(conf.SendGridAPIKey)

	// Configuration complete!

	if s.secret, err = NewSecretManager(conf.Secrets); err != nil {
		return s, nil
	}
	return s, nil
}

// Server implements the GRPC TRISADirectoryService.
type Server struct {
	api.UnimplementedTRISADirectoryServer
	admin.UnimplementedDirectoryAdministrationServer
	db     store.Store
	srv    *grpc.Server
	conf   config.Config
	certs  *sectigo.Sectigo
	email  *sendgrid.Client
	secret *SecretManager
	echan  chan error
}

// Serve GRPC requests on the specified address.
func (s *Server) Serve() (err error) {
	// Initialize the gRPC server
	s.srv = grpc.NewServer(grpc.UnaryInterceptor(s.serverInterceptor))
	api.RegisterTRISADirectoryServer(s.srv, s)
	admin.RegisterDirectoryAdministrationServer(s.srv, s)

	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.echan <- s.Shutdown()
	}()

	// Start the certificate manager go routine process
	go s.CertManager()

	// Start the backup manager go routine process
	go s.BackupManager()

	if s.conf.Maintenance {
		log.Warn().Msg("starting server in maintenance mode")
	}

	// Listen for TCP requests on the specified address and port
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		return fmt.Errorf("could not listen on %q", s.conf.BindAddr)
	}
	defer sock.Close()

	// Run the server
	go func() {
		log.Info().
			Str("listen", s.conf.BindAddr).
			Str("version", pkg.Version()).
			Msg("server started")

		if err := s.srv.Serve(sock); err != nil {
			s.echan <- err
		}
	}()

	// Listen for any errors that might have occurred and wait for all go routines to finish
	if err = <-s.echan; err != nil {
		return err
	}
	return nil
}

// Shutdown the TRISA Directory Service gracefully
func (s *Server) Shutdown() (err error) {
	log.Info().Msg("gracefully shutting down")
	s.srv.GracefulStop()
	if err = s.db.Close(); err != nil {
		log.Error().Err(err).Msg("could not shutdown database")
		return err
	}
	log.Debug().Msg("successful shutdown")
	return nil
}

// Register a new VASP entity with the directory service. After registration, the new
// entity must go through the verification process to get issued a certificate. The
// status of verification can be obtained by using the lookup RPC call.
// Register generates a PKCS12 password, provided in the RPC response which can be
// used to access the certificate private keys when they're emailed.
func (s *Server) Register(ctx context.Context, in *api.RegisterRequest) (out *api.RegisterReply, err error) {
	out = &api.RegisterReply{}
	vasp := &pb.VASP{
		RegisteredDirectory: s.conf.DirectoryID,
		Entity:              in.Entity,
		Contacts:            in.Contacts,
		TrisaEndpoint:       in.TrisaEndpoint,
		CommonName:          in.CommonName,
		Website:             in.Website,
		BusinessCategory:    in.BusinessCategory,
		VaspCategories:      in.VaspCategories,
		EstablishedOn:       in.EstablishedOn,
		Trixo:               in.Trixo,
		VerificationStatus:  pb.VerificationState_SUBMITTED,
		Version:             &pb.Version{Version: 1},
	}

	// Compute the common name from the trisa endpoint if not specified
	if vasp.CommonName == "" && vasp.TrisaEndpoint != "" {
		if vasp.CommonName, _, err = net.SplitHostPort(in.TrisaEndpoint); err != nil {
			log.Warn().Err(err).Msg("could not parse common name from endpoint")
			out.Error = &api.Error{
				Code:    400,
				Message: err.Error(),
			}
			return out, nil
		}
	}

	// Validate partial VASP record to ensure that it can be registered.
	if err = vasp.Validate(true); err != nil {
		log.Warn().Err(err).Msg("invalid or incomplete VASP registration")
		out.Error = &api.Error{
			Code:    400,
			Message: err.Error(),
		}
		return out, nil
	}

	// TODO: create legal entity hash to detect a repeat registration without ID
	// TODO: add signature to leveldb indices
	if out.Id, err = s.db.Create(vasp); err != nil {
		log.Warn().Err(err).Msg("could not register VASP")
		out.Error = &api.Error{
			Code:    400,
			Message: err.Error(),
		}
		return out, nil
	}

	name, _ := vasp.Name()
	log.Info().Str("name", name).Str("id", vasp.Id).Msg("registered VASP")

	// Begin verification process by sending emails to all contacts in the VASP record.
	// TODO: add to processing queue to return sooner/parallelize work
	if err = s.VerifyContactEmail(vasp); err != nil {
		log.Error().Err(err).Msg("could not verify contacts")
		out.Error = &api.Error{
			Code:    500,
			Message: err.Error(),
		}
		return out, nil
	}
	log.Info().Msg("contact email verifications sent")

	// Create PKCS12 password along with certificate request.
	password := CreateToken(16)
	certRequest := &models.CertificateRequest{
		Id:         uuid.New().String(),
		Vasp:       vasp.Id,
		CommonName: vasp.CommonName,
		Status:     models.CertificateRequestState_INITIALIZED,
	}

	// Make a new secret of type "password"
	secretType := "password"
	if err = s.secret.With(certRequest.Id, secretType).CreateSecret(ctx, secretType); err != nil {
		log.Error().Err(err).Msg("could not create new secret on registration")
		out.Error = &api.Error{Code: 500, Message: "internal error with registration, please contact admins"}
		return out, nil
	}
	if err = s.secret.With(certRequest.Id, secretType).AddSecretVersion(ctx, secretType, []byte(password)); err != nil {
		log.Error().Err(err).Msg("unable to add secret version on registration")
		out.Error = &api.Error{Code: 500, Message: "internal error during registration, please contact admins"}
		return out, nil
	}

	if err = s.db.SaveCertRequest(certRequest); err != nil {
		log.Error().Err(err).Msg("could not save certificate request for VASP")
		out.Error = &api.Error{
			Code:    500,
			Message: err.Error(),
		}
		return out, nil
	}

	out.Id = vasp.Id
	out.RegisteredDirectory = vasp.RegisteredDirectory
	out.CommonName = vasp.CommonName
	out.Status = vasp.VerificationStatus
	out.Message = "verification code sent to contact emails, please check spam folder if not arrived; pkcs12 password attached, this is the only time it will be available - do not lose!"
	out.Pkcs12Password = password
	return out, nil
}

// Lookup a VASP entity by name or ID to get full details including the TRISA certification
// if it exists and the entity has been verified.
func (s *Server) Lookup(ctx context.Context, in *api.LookupRequest) (out *api.LookupReply, err error) {
	var vasp *pb.VASP
	out = &api.LookupReply{}

	if in.Id != "" {
		// TODO: add registered directory to lookup
		if vasp, err = s.db.Retrieve(in.Id); err != nil {
			out.Error = &api.Error{
				Code:    404,
				Message: err.Error(),
			}
		}

	} else if in.CommonName != "" {
		var vasps []*pb.VASP
		if vasps, err = s.db.Search(map[string]interface{}{"name": in.CommonName}); err != nil {
			out.Error = &api.Error{
				Code:    404,
				Message: err.Error(),
			}
		}

		if len(vasps) == 1 {
			vasp = vasps[0]
		} else {
			out.Error = &api.Error{
				Code:    404,
				Message: "not found",
			}
		}
	} else {
		out.Error = &api.Error{
			Code:    400,
			Message: "no lookup query provided",
		}
		return out, nil
	}

	if out.Error == nil {
		// TODO: should lookups only return verified peers?
		out.Id = vasp.Id
		out.RegisteredDirectory = vasp.RegisteredDirectory
		out.CommonName = vasp.CommonName
		out.Endpoint = vasp.TrisaEndpoint
		out.IdentityCertificate = vasp.IdentityCertificate
		out.Name, _ = vasp.Name()
		out.Country = vasp.Entity.CountryOfRegistration
		out.VerifiedOn = vasp.VerifiedOn

		// TODO: how do we determine which signing certificate to send?
		// Currently sending the last certificate in the array so that to update a
		// signing certificate, a new cert just has to be appended to the slice.
		if len(vasp.SigningCertificates) > 0 {
			out.SigningCertificate = vasp.SigningCertificates[len(vasp.SigningCertificates)-1]
		}

		log.Info().Str("id", vasp.Id).Msg("VASP lookup succeeded")
	} else {
		log.Warn().
			Err(out.Error).
			Str("id", in.Id).
			Str("name", in.CommonName).
			Msg("could not lookup VASP")
	}
	return out, nil
}

// Search for VASP entity records by name or by country in order to perform more detailed
// Lookup requests. The search process is purposefully simplistic at the moment.
func (s *Server) Search(ctx context.Context, in *api.SearchRequest) (out *api.SearchReply, err error) {
	out = &api.SearchReply{}
	query := make(map[string]interface{})
	query["name"] = in.Name
	query["website"] = in.Website
	query["country"] = in.Country

	// Build category query
	categories := make([]string, 0, len(in.BusinessCategory)+len(in.VaspCategory))
	for _, category := range in.BusinessCategory {
		categories = append(categories, category.String())
	}
	categories = append(categories, in.VaspCategory...)
	query["category"] = categories

	var vasps []*pb.VASP
	if vasps, err = s.db.Search(query); err != nil {
		out.Error = &api.Error{
			Code:    400,
			Message: err.Error(),
		}
	}

	out.Results = make([]*api.SearchReply_Result, 0, len(vasps))
	for _, vasp := range vasps {
		out.Results = append(out.Results, &api.SearchReply_Result{
			Id:                  vasp.Id,
			RegisteredDirectory: vasp.RegisteredDirectory,
			CommonName:          vasp.CommonName,
			Endpoint:            vasp.TrisaEndpoint,
		})
	}

	entry := log.With().
		Strs("name", in.Name).
		Strs("websites", in.Website).
		Strs("country", in.Country).
		Strs("categories", categories).
		Int("results", len(out.Results)).
		Logger()

	if out.Error != nil {
		entry.Warn().Err(out.Error).Msg("unsuccessful search")
	} else {
		entry.Info().Msg("search succeeded")
	}
	return out, nil
}

// Status returns the status of a VASP including its verification and service status if
// the directory service is performing health check monitoring.
func (s *Server) Verification(ctx context.Context, in *api.VerificationRequest) (out *api.VerificationReply, err error) {
	var vasp *pb.VASP
	out = &api.VerificationReply{}

	if in.Id != "" {
		// TODO: add registered directory to lookup
		if vasp, err = s.db.Retrieve(in.Id); err != nil {
			log.Error().Err(err).Str("id", in.Id).Msg("could not retrieve vasp")
			return nil, status.Error(codes.NotFound, err.Error())
		}

	} else if in.CommonName != "" {
		// TODO: change lookup to unique common name lookup
		var vasps []*pb.VASP
		if vasps, err = s.db.Search(map[string]interface{}{"name": in.CommonName}); err != nil {
			log.Error().Err(err).Str("name", in.CommonName).Msg("could not retrieve vasp")
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if len(vasps) == 1 {
			vasp = vasps[0]
		} else {
			log.Warn().Int("results", len(vasps)).Msg("vasp not found")
			return nil, status.Error(codes.NotFound, "VASP not found")
		}
	} else {
		log.Warn().Msg("no lookup query provided")
		return nil, status.Error(codes.InvalidArgument, "no lookup query provided")
	}

	// TODO: should lookups only return verified peers?
	out.VerificationStatus = vasp.VerificationStatus
	out.ServiceStatus = vasp.ServiceStatus
	out.VerifiedOn = vasp.VerifiedOn
	out.FirstListed = vasp.FirstListed
	out.LastUpdated = vasp.LastUpdated
	log.Info().Str("id", vasp.Id).Msg("VASP status succeeded")
	return out, nil
}

// VerifyEmail checks the contact tokens for the specified VASP and registers the
// contact email verification. If successful, this method then sends the verification
// request to the TRISA Admins for review.
func (s *Server) VerifyContact(ctx context.Context, in *api.VerifyContactRequest) (out *api.VerifyContactReply, err error) {
	out = &api.VerifyContactReply{}

	var vasp *pb.VASP
	if vasp, err = s.db.Retrieve(in.Id); err != nil {
		log.Error().Err(err).Str("id", in.Id).Msg("could not retrieve vasp")
		out.Error = &api.Error{
			Code:    404,
			Message: err.Error(),
		}
		return out, nil
	}

	nVerified := 0
	found := false
	contacts := []*pb.Contact{
		vasp.Contacts.Technical,
		vasp.Contacts.Administrative,
		vasp.Contacts.Billing,
		vasp.Contacts.Legal,
	}
	for _, contact := range contacts {
		if contact == nil {
			continue
		}

		token, verified, err := models.GetContactVerification(contact)
		if err != nil {
			log.Error().Err(err).Msg("could not retrieve verification from contact extra data field")
			return nil, status.Error(codes.Internal, "could not verify contact")
		}

		if token == in.Token {
			found = true
			log.Info().Str("email", contact.Email).Msg("contact email verified")
			if err = models.SetContactVerification(contact, "", true); err != nil {
				log.Error().Err(err).Msg("could not set verification on contact extra data field")
				return nil, status.Error(codes.Internal, "could not verify contact")
			}

		}
		if verified {
			nVerified++
		}
	}

	if !found || nVerified == 0 {
		log.Error().Err(err).Str("token", in.Token).Msg("could not find contact with token")
		out.Error = &api.Error{
			Code:    404,
			Message: "could not find contact with specified token",
		}
		return out, nil
	}

	// Ensures that we only send the verification email to the admins once.
	// NOTE: we will only generate the password on the first email verification.
	if nVerified > 1 {
		// Save the updated contact
		if err = s.db.Update(vasp); err != nil {
			log.Error().Err(err).Msg("could not update vasp after contact verification")
			out.Error = &api.Error{
				Code:    500,
				Message: err.Error(),
			}
			return out, nil
		}

		out.Status = vasp.VerificationStatus
		out.Message = "email successfully verified; verification review already sent to TRISA admins"
		return out, nil
	}

	// Note that this status will get updated in the review request email
	vasp.VerificationStatus = pb.VerificationState_EMAIL_VERIFIED

	// If this is the first verification, generate the PKCS12 password and send verification review email
	// TODO: make this better
	if err = s.ReviewRequestEmail(vasp); err != nil {
		log.Error().Err(err).Msg("could not send verification review email")
		out.Error = &api.Error{
			Code:    500,
			Message: "could not send verification review email",
		}
		return out, nil
	}

	// Now that the email has been sent out the vasp is pending review
	vasp.VerificationStatus = pb.VerificationState_PENDING_REVIEW

	// Save the VASP and newly created certificate request
	if err = s.db.Update(vasp); err != nil {
		log.Error().Err(err).Msg("could not update vasp status to pending review")
		out.Error = &api.Error{
			Code:    500,
			Message: err.Error(),
		}
		return out, nil
	}

	out.Status = vasp.VerificationStatus
	out.Message = "email successfully verified and verification review sent to TRISA admins"
	return out, nil
}

func (s *Server) Status(ctx context.Context, in *api.HealthCheck) (out *api.ServiceState, err error) {
	log.Info().
		Uint32("attempts", in.Attempts).
		Str("last_checked_at", in.LastCheckedAt).
		Msg("status check")

	// Request another health check between 30-60 min from now
	now := time.Now()

	// Default service state is healthy.
	out = &api.ServiceState{
		Status:    api.ServiceState_HEALTHY,
		NotBefore: now.Add(30 * time.Minute).Format(time.RFC3339),
		NotAfter:  now.Add(60 * time.Minute).Format(time.RFC3339),
	}

	// If we're in maintenance mode, update the service state.
	if s.conf.Maintenance {
		out.Status = api.ServiceState_MAINTENANCE
	}

	return out, nil
}
