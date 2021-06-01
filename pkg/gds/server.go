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
	"github.com/trisacrypto/directory/pkg"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v1"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
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

	// Create the Email Manager with SendGrid API client
	if s.email, err = emails.New(conf.Email); err != nil {
		return nil, err
	}

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
	email  *emails.EmailManager
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

	// Compute the common name from the TRISA endpoint if not specified
	if vasp.CommonName == "" && vasp.TrisaEndpoint != "" {
		if vasp.CommonName, _, err = net.SplitHostPort(in.TrisaEndpoint); err != nil {
			log.Warn().Err(err).Msg("could not parse common name from endpoint")
			return nil, status.Error(codes.InvalidArgument, "no common name supplied, could not parse common name from endpoint")
		}
	}

	// Validate partial VASP record to ensure that it can be registered.
	if err = vasp.Validate(true); err != nil {
		log.Warn().Err(err).Msg("invalid or incomplete VASP registration")
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %s", err)
	}

	// TODO: create legal entity hash to detect a repeat registration without ID
	// TODO: add signature to leveldb indices
	// TODO: check already exists and uniqueness constraints
	if vasp.Id, err = s.db.Create(vasp); err != nil {
		// Assuming uniqueness is the primary constraint here
		// TODO: better database error checking or handling
		log.Warn().Err(err).Msg("could not register VASP in database")
		return nil, status.Error(codes.AlreadyExists, "could not complete registration, uniqueness constraints violated")
	}

	// Log successful registration
	name, _ := vasp.Name()
	log.Info().Str("name", name).Str("id", vasp.Id).Msg("registered VASP")

	// Begin verification process by sending emails to all contacts in the VASP record.
	// TODO: add to processing queue to return sooner/parallelize work
	// Create the verification tokens and save the VASP back to the database
	var contacts = []*pb.Contact{
		vasp.Contacts.Technical,
		vasp.Contacts.Administrative,
		vasp.Contacts.Billing,
		vasp.Contacts.Legal,
	}

	for idx, contact := range contacts {
		if contact != nil && contact.Email != "" {
			if err = models.SetContactVerification(contact, CreateToken(48), false); err != nil {
				log.Error().Err(err).Int("index", idx).Str("vasp", vasp.Id).Msg("could not set contact verification token")
				return nil, status.Error(codes.Aborted, "could not send contact verification emails")
			}
		}
	}

	if err = s.db.Update(vasp); err != nil {
		log.Error().Err(err).Str("vasp", vasp.Id).Msg("could not update vasp with contact verification tokens")
		return nil, status.Error(codes.Aborted, "could not send contact verification emails")
	}

	// Send contacts with updated tokens
	if err = s.email.SendVerifyContacts(vasp); err != nil {
		log.Error().Err(err).Str("vasp", vasp.Id).Msg("could not sennd verify contacts emails")
		return nil, status.Error(codes.Aborted, "could not send contact verification emails")
	}

	// Log successful contact verification emails sent
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
	if err = s.secret.With(certRequest.Id).CreateSecret(ctx, secretType); err != nil {
		log.Error().Err(err).Str("vasp", vasp.Id).Msg("could not create new secret for pkcs12 password")
		return nil, status.Error(codes.FailedPrecondition, "internal error with registration, please contact admins")
	}
	if err = s.secret.With(certRequest.Id).AddSecretVersion(ctx, secretType, []byte(password)); err != nil {
		log.Error().Err(err).Str("vasp", vasp.Id).Msg("unable to add secret version for pkcs12 password")
		return nil, status.Error(codes.FailedPrecondition, "internal error with registration, please contact admins")
	}

	if err = s.db.SaveCertRequest(certRequest); err != nil {
		log.Error().Err(err).Str("vasp", vasp.Id).Msg("could not save certificate request")
		return nil, status.Error(codes.FailedPrecondition, "internal error with registration, please contact admins")
	}

	out = &api.RegisterReply{
		Id:                  vasp.Id,
		RegisteredDirectory: vasp.RegisteredDirectory,
		CommonName:          vasp.CommonName,
		Status:              vasp.VerificationStatus,
		Message:             "a verification code has been sent to contact emails, please check spam folder if it has not arrived; pkcs12 password attached, this is the only time it will be available -- do not lose!",
		Pkcs12Password:      password,
	}
	return out, nil
}

// Lookup a VASP entity by name or ID to get full details including the TRISA certification
// if it exists and the entity has been verified.
func (s *Server) Lookup(ctx context.Context, in *api.LookupRequest) (out *api.LookupReply, err error) {
	var vasp *pb.VASP
	switch {
	case in.Id != "":
		// TODO: add registered directory to lookup
		if vasp, err = s.db.Retrieve(in.Id); err != nil {
			log.Warn().Err(err).Str("id", in.Id).Str("registered_directory", in.RegisteredDirectory).Msg("could not find VASP by ID")
			return nil, status.Error(codes.NotFound, "could not find VASP by ID")
		}
	case in.CommonName != "":
		var vasps []*pb.VASP
		if vasps, err = s.db.Search(map[string]interface{}{"name": in.CommonName}); err != nil {
			log.Warn().Err(err).Str("common_name", in.CommonName).Msg("could not search for common name")
			return nil, status.Error(codes.NotFound, "could not find VASP by common name")
		}

		if len(vasps) != 1 {
			log.Warn().Str("common_name", in.CommonName).Int("nresults", len(vasps)).Msg("wrong number of VASPs returned from search")
			return nil, status.Error(codes.NotFound, "could not find VASP by common name")
		}

		vasp = vasps[0]
	default:
		log.Warn().Str("rpc", "lookup").Msg("no arguments supplied")
		return nil, status.Error(codes.InvalidArgument, "please supply ID and registered directory or common name for lookup")
	}

	// TODO: should lookups only return verified peers?
	out = &api.LookupReply{
		Id:                  vasp.Id,
		RegisteredDirectory: vasp.RegisteredDirectory,
		CommonName:          vasp.CommonName,
		Endpoint:            vasp.TrisaEndpoint,
		IdentityCertificate: vasp.IdentityCertificate,
		Country:             vasp.Entity.CountryOfRegistration,
		VerifiedOn:          vasp.VerifiedOn,
	}

	// Ignore errors on name lookup
	out.Name, _ = vasp.Name()

	// TODO: how do we determine which signing certificate to send?
	// Currently sending the last certificate in the array so that to update a
	// signing certificate, a new cert just has to be appended to the slice.
	if len(vasp.SigningCertificates) > 0 {
		out.SigningCertificate = vasp.SigningCertificates[len(vasp.SigningCertificates)-1]
	}

	log.Info().Str("id", vasp.Id).Str("common_name", vasp.CommonName).Msg("VASP lookup succeeded")
	return out, nil
}

// Search for VASP entity records by name or by country in order to perform more detailed
// Lookup requests. The search process is purposefully simplistic at the moment.
func (s *Server) Search(ctx context.Context, in *api.SearchRequest) (out *api.SearchReply, err error) {
	// Create search query to send to database
	query := make(map[string]interface{})
	query["name"] = in.Name
	query["website"] = in.Website
	query["country"] = in.Country

	// Build categories query
	categories := make([]string, 0, len(in.BusinessCategory)+len(in.VaspCategory))
	for _, category := range in.BusinessCategory {
		categories = append(categories, category.String())
	}
	categories = append(categories, in.VaspCategory...)
	query["category"] = categories

	var vasps []*pb.VASP
	if vasps, err = s.db.Search(query); err != nil {
		log.Error().Err(err).Msg("vasp search failed")
		return nil, status.Error(codes.Aborted, err.Error())
	}

	// Build search results to return
	out = &api.SearchReply{
		Results: make([]*api.SearchReply_Result, 0, len(vasps)),
	}
	for _, vasp := range vasps {
		out.Results = append(out.Results, &api.SearchReply_Result{
			Id:                  vasp.Id,
			RegisteredDirectory: vasp.RegisteredDirectory,
			CommonName:          vasp.CommonName,
			Endpoint:            vasp.TrisaEndpoint,
		})
	}

	log.Info().
		Strs("name", in.Name).
		Strs("websites", in.Website).
		Strs("country", in.Country).
		Strs("categories", categories).
		Int("results", len(out.Results)).
		Msg("search succeeded")

	return out, nil
}

// Status returns the status of a VASP including its verification and service status if
// the directory service is performing health check monitoring.
func (s *Server) Verification(ctx context.Context, in *api.VerificationRequest) (out *api.VerificationReply, err error) {
	var vasp *pb.VASP
	switch {
	case in.Id != "":
		// TODO: add registered directory to retrieve
		if vasp, err = s.db.Retrieve(in.Id); err != nil {
			log.Warn().Err(err).Str("id", in.Id).Str("registered_directory", in.RegisteredDirectory).Msg("could not find VASP by ID")
			return nil, status.Error(codes.NotFound, "could not find VASP by ID")
		}
	case in.CommonName != "":
		var vasps []*pb.VASP
		if vasps, err = s.db.Search(map[string]interface{}{"name": in.CommonName}); err != nil {
			log.Warn().Err(err).Str("common_name", in.CommonName).Msg("could not search for common name")
			return nil, status.Error(codes.NotFound, "could not find VASP by common name")
		}

		if len(vasps) != 1 {
			log.Warn().Str("common_name", in.CommonName).Int("nresults", len(vasps)).Msg("wrong number of VASPs returned from search")
			return nil, status.Error(codes.NotFound, "could not find VASP by common name")
		}

		vasp = vasps[0]
	default:
		log.Warn().Str("rpc", "lookup").Msg("no arguments supplied")
		return nil, status.Error(codes.InvalidArgument, "please supply ID and registered directory or common name for verification")
	}

	// TODO: also return RevokedOn, which needs to be stored on the VASP
	out = &api.VerificationReply{
		VerificationStatus: vasp.VerificationStatus,
		ServiceStatus:      vasp.ServiceStatus,
		VerifiedOn:         vasp.VerifiedOn,
		FirstListed:        vasp.FirstListed,
		LastUpdated:        vasp.LastUpdated,
	}
	log.Info().Str("id", vasp.Id).Str("common_name", vasp.CommonName).Msg("verification status check")
	return out, nil
}

// VerifyEmail checks the contact tokens for the specified VASP and registers the
// contact email verification. If successful, this method then sends the verification
// request to the TRISA Admins for review.
func (s *Server) VerifyContact(ctx context.Context, in *api.VerifyContactRequest) (out *api.VerifyContactReply, err error) {
	// Retrieve VASP associated with contact from the database.
	var vasp *pb.VASP
	if vasp, err = s.db.Retrieve(in.Id); err != nil {
		log.Error().Err(err).Str("id", in.Id).Msg("could not retrieve vasp")
		return nil, status.Error(codes.NotFound, "could not find associated VASP record by ID")
	}

	// Search through the contacts to determine the contacts verified by the supplied token.
	prevVerified := 0
	found := false
	contacts := []*pb.Contact{
		vasp.Contacts.Technical,
		vasp.Contacts.Administrative,
		vasp.Contacts.Billing,
		vasp.Contacts.Legal,
	}
	for idx, contact := range contacts {
		// Ignore empty contacts
		if contact == nil {
			continue
		}

		// Get the verification status
		token, verified, err := models.GetContactVerification(contact)
		if err != nil {
			log.Error().Err(err).Msg("could not retrieve verification from contact extra data field")
			return nil, status.Error(codes.Aborted, "could not verify contact")
		}

		// Perform token check and if token matches, mark contact as verified
		if token == in.Token {
			found = true
			log.Info().Str("vasp", vasp.Id).Int("index", idx).Msg("contact email verified")
			if err = models.SetContactVerification(contact, "", true); err != nil {
				log.Error().Err(err).Msg("could not set verification on contact extra data field")
				return nil, status.Error(codes.Aborted, "could not verify contact")
			}

		} else if verified {
			// Determine the total number of contacts previously verified, not including
			// the current contact that was just verified. This will help prevent
			// sending multiple emails to the TRISA Admins for review.
			prevVerified++
		}
	}

	// Check if we haven't managed to verify the contact
	if !found {
		log.Error().Err(err).Str("vasp", vasp.Id).Msg("could not find contact with token")
		return nil, status.Error(codes.NotFound, "could not find contact with the specified token")
	}

	// Ensures that we only send the verification email to the admins once.
	// If we have previously verified contacts, assume that we've already sent the
	// registration review email and do nothing.
	if prevVerified > 0 && vasp.VerificationStatus > pb.VerificationState_SUBMITTED {
		// Save the updated contact
		if err = s.db.Update(vasp); err != nil {
			log.Error().Err(err).Msg("could not update VASP record after contact verification")
			return nil, status.Error(codes.Internal, "could not update contact after verification")
		}

		return &api.VerifyContactReply{
			Status:  vasp.VerificationStatus,
			Message: "email successfully verified",
		}, nil
	}

	// Since we have one successful email verification at this point, begin the
	// registration review process by sending an email to the TRISA admins.
	// Step 1: mark the VASP as email verified and create an admin token.
	vasp.VerificationStatus = pb.VerificationState_EMAIL_VERIFIED

	// Create verification token for admin and update database
	// TODO: replace with actual authentication
	if err = models.SetAdminVerificationToken(vasp, CreateToken(48)); err != nil {
		log.Error().Err(err).Msg("could not create admin verification token")
		return nil, status.Error(codes.FailedPrecondition, "there was a problem submitting your registration review request, please contact the admins")
	}
	if err = s.db.Update(vasp); err != nil {
		log.Error().Err(err).Msg("could not save admin verification token")
		return nil, status.Error(codes.FailedPrecondition, "there was a problem submitting your registration review request, please contact the admins")
	}

	// Step 2: send review request email to the TRISA admins.
	if err = s.email.SendReviewRequest(vasp); err != nil {
		log.Error().Err(err).Msg("could not send verification review email")
		return nil, status.Error(codes.FailedPrecondition, "there was a problem submitting your registration review request, please contact the admins")
	}

	// Step 3: if the review email has been successfully sent, mark as pending review.
	vasp.VerificationStatus = pb.VerificationState_PENDING_REVIEW

	// Save the VASP and newly created certificate request
	if err = s.db.Update(vasp); err != nil {
		log.Error().Err(err).Msg("could not update vasp status to pending review")
		return nil, status.Error(codes.Internal, "there was a problem submitting your registration review request, please contact the admins")
	}

	return &api.VerifyContactReply{
		Status:  vasp.VerificationStatus,
		Message: "email successfully verified and verification review sent to TRISA admins",
	}, nil
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
