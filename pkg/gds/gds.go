package gds

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/utils"
	activity "github.com/trisacrypto/directory/pkg/utils/activity"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewGDS creates a new GDS server derived from a parent Service.
func NewGDS(svc *Service) (gds *GDS, err error) {
	gds = &GDS{
		svc:  svc,
		conf: &svc.conf.GDS,
		db:   svc.db,
	}

	// Initialize the gRPC server
	gds.srv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(svc.UnaryInterceptors()...),
		grpc.ChainStreamInterceptor(svc.StreamInterceptors()...),
	)
	api.RegisterTRISADirectoryServer(gds.srv, gds)
	return gds, nil
}

// GDS implements the TRISADirectoryService as defined by the v1beta1 or later TRISA
// protocol buffers. This service is the primary interaction point with TRISA service
// implementations that lookup information from the directory service, and this service
// also allows users to register and verify with the directory.
//
// SEE FIRST: Service as defined in service.go (the main entrypoint of the server)
type GDS struct {
	api.UnimplementedTRISADirectoryServer
	svc  *Service          // The parent Service GDS uses to interact with other components
	srv  *grpc.Server      // The gRPC server that listens on its own independent port
	conf *config.GDSConfig // The GDS service specific configuration (helper alias to s.svc.conf.GDS)
	db   store.Store       // Database connection for loading objects (helper alias to s.svc.db)
}

// Serve gRPC requests on the specified address.
func (s *GDS) Serve() (err error) {
	if !s.conf.Enabled {
		log.Warn().Msg("trisa directory service is not enabled")
		return nil
	}

	// This service must run even if we're in maintenance mode to send service state
	// MAINTENANCE in status replies.
	if s.svc.conf.Maintenance {
		log.Debug().Msg("starting GDS server in maintenance mode")
	}

	// Listen for TCP requests on the specified address and port
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		return fmt.Errorf("could not listen on %q", s.conf.BindAddr)
	}

	// Run the server
	go s.Run(sock)
	log.Info().Str("listen", s.conf.BindAddr).Str("version", pkg.Version()).Msg("trisa directory server started")

	// Now that the go routine is started return nil, meaning the service has started
	// successfully with no problems.
	return nil
}

// Run the gRPC server. This method is extracted from the Serve function so that it can
// be run in its own go routine and to allow tests to Run a bufconn server without
// starting a live server with all of the various go routines and channels running.
func (s *GDS) Run(sock net.Listener) {
	defer sock.Close()
	if err := s.srv.Serve(sock); err != nil {
		s.svc.echan <- err
	}
}

// Shutdown the TRISA Directory Service gracefully
func (s *GDS) Shutdown() (err error) {
	log.Debug().Msg("gracefully shutting down GDS server")
	s.srv.GracefulStop()
	log.Debug().Msg("successful shutdown of GDS server")
	return nil
}

//===========================================================================
// GDS Server Methods
//===========================================================================

// Register a new VASP entity with the directory service. After registration, the new
// entity must go through the verification process to get issued a certificate. The
// status of verification can be obtained by using the lookup RPC call.
// Register generates a PKCS12 password, provided in the RPC response which can be
// used to access the certificate private keys when they're emailed.
func (s *GDS) Register(ctx context.Context, in *api.RegisterRequest) (out *api.RegisterReply, err error) {
	// send register request activity to network activity handler
	activity.Register().Add()

	vasp := &pb.VASP{
		RegisteredDirectory: s.svc.conf.DirectoryID,
		Entity:              in.Entity,
		Contacts:            in.Contacts,
		TrisaEndpoint:       in.TrisaEndpoint,
		CommonName:          in.CommonName,
		Website:             in.Website,
		BusinessCategory:    in.BusinessCategory,
		VaspCategories:      in.VaspCategories,
		EstablishedOn:       in.EstablishedOn,
		Trixo:               in.Trixo,
		VerificationStatus:  pb.VerificationState_NO_VERIFICATION,
		Version:             &pb.Version{Version: 1},
	}

	// Validate TRISA endpoint
	if in.TrisaEndpoint == "" {
		sentry.Warn(ctx).Err(err).Msg("missing endpoint in request")
		return nil, status.Error(codes.InvalidArgument, "no endpoint supplied")
	}

	if err = validateEndpoint(in.TrisaEndpoint); err != nil {
		sentry.Warn(ctx).Err(err).Str("endpoint", in.TrisaEndpoint).Msg("invalid endpoint")
		return nil, status.Error(codes.InvalidArgument, "invalid endpoint supplied")
	}

	// Compute the common name from the TRISA endpoint if not specified
	if vasp.CommonName == "" {
		if vasp.CommonName, _, err = net.SplitHostPort(in.TrisaEndpoint); err != nil {
			sentry.Warn(ctx).Err(err).Msg("could not parse common name from endpoint")
			return nil, status.Error(codes.InvalidArgument, "no common name supplied, could not parse common name from endpoint")
		}
	} else {
		// Validate common name if supplied
		if err = utils.ValidateCommonName(vasp.CommonName); err != nil {
			sentry.Warn(ctx).Err(err).Str("common_name", vasp.CommonName).Msg("invalid common name")
			return nil, status.Error(codes.InvalidArgument, "invalid common name supplied")
		}
	}

	// Validate partial VASP record to ensure that it can be registered.
	if err = models.ValidateVASP(vasp, true); err != nil {
		sentry.Warn(ctx).Err(err).Msg("invalid or incomplete VASP registration")
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %s", err)
	}

	// Set any zero valued contacts to nil to ensure empty records aren't created.
	if vasp.Contacts.Administrative != nil && vasp.Contacts.Administrative.IsZero() {
		vasp.Contacts.Administrative = nil
	}
	if vasp.Contacts.Technical != nil && vasp.Contacts.Technical.IsZero() {
		vasp.Contacts.Technical = nil
	}
	if vasp.Contacts.Billing != nil && vasp.Contacts.Billing.IsZero() {
		vasp.Contacts.Billing = nil
	}
	if vasp.Contacts.Legal != nil && vasp.Contacts.Legal.IsZero() {
		vasp.Contacts.Legal = nil
	}

	// Retrieve email address from one of the supplied contacts.
	var email string
	if email := GetContactEmail(vasp); email == "" {
		sentry.Error(ctx).Err(errors.New("no contact email address found")).Msg("incorrect access on validated VASP")
		return nil, status.Error(codes.InvalidArgument, "no email address in supplied VASP contacts")
	}

	// Set verification status to SUBMITTED.
	if err := models.UpdateVerificationStatus(vasp, pb.VerificationState_SUBMITTED, "register request received", email); err != nil {
		sentry.Warn(ctx).Err(err).Msg("could not update VASP verification status")
		return nil, status.Error(codes.Aborted, "could not add new entry to VASP audit log")
	}

	// TODO: create legal entity hash to detect a repeat registration without ID
	// TODO: add signature to leveldb indices
	// TODO: check already exists and uniqueness constraints
	if vasp.Id, err = s.db.CreateVASP(ctx, vasp); err != nil {
		// Assuming uniqueness is the primary constraint here
		// TODO: better database error checking or handling
		sentry.Warn(ctx).Err(err).Msg("could not register VASP in database")
		return nil, status.Error(codes.AlreadyExists, "could not complete registration, uniqueness constraints violated")
	}

	// Log successful registration
	vaspName, _ := vasp.Name()
	log.Info().Str("name", vaspName).Str("id", vasp.Id).Msg("registered VASP")

	// Begin verification process by sending emails to all contacts in the VASP record.
	var contacts *models.Contacts
	if contacts, err = s.db.VASPContacts(ctx, vasp); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not create or retrieve vasp contacts")
		return nil, status.Error(codes.Internal, "could not retrieve vasp contacts from database")
	}

	// TODO: add to processing queue to return sooner/parallelize work
	// Create the verification tokens and save the VASP back to the database
	sent := 0
	iter := contacts.NewIterator(models.SkipDuplicates())
	for iter.Next() {
		contact := iter.Contact()
		if contact.Email == nil {
			// Create the contact email record in the database.
			contact.Email = &models.Email{
				Name:     contact.Contact.Name,
				Email:    contact.Contact.Email,
				Vasps:    []string{vasp.Id},
				Verified: false,
				Token:    secrets.CreateToken(models.VerificationTokenLength),
			}
			contacts.Emails = append(contacts.Emails, contact.Email)

			if _, err = s.db.CreateEmail(ctx, contact.Email); err != nil {
				sentry.Error(ctx).Err(err).Msg("could not create email record for contact")
				return nil, status.Error(codes.Internal, "could not finalize registration")
			}
		} else {
			// Update the contact email with the VASP
			if contact.Contact.Name != "" {
				contact.Email.Name = contact.Contact.Name
			}
			contact.Email.AddVASP(vasp.Id)

			if err = s.db.UpdateEmail(ctx, contact.Email); err != nil {
				sentry.Error(ctx).Err(err).Str("email", contact.Contact.Email).Msg("could not update contact email")
			}
		}

		if !contact.IsVerified() {
			// If the contact already existed but is still not verified, send a
			// verification email a second time. If the contact was created, then send
			// the verification email for the first time.
			// TODO: add to processing queue to return sooner/parallelize work
			if err = s.svc.email.SendVerifyContact(vasp.Id, contact); err != nil {
				// If there is an error sending contact verification emails, alert admins who
				// can resend emails later, do not abort processing the registration.
				sentry.Error(ctx).Err(err).Str("vasp", vasp.Id).Int("sent", sent).Msg("could not send verify contacts emails")
			} else {
				// Log successful contact verification emails sent
				log.Info().Msg("contact email verification sent")
				sent++

				// Update the model contact record to update the email log
				if err = s.db.UpdateEmail(ctx, contact.Email); err != nil {
					sentry.Error(ctx).Err(err).Str("contact", contact.Email.Email).Msg("could not update email logs on contact")
				}
			}
		}
	}

	// If at least one of the contacts is verified, begin review process
	iter = contacts.NewIterator(models.SkipUnverified(), models.SkipDuplicates())
	for iter.Next() {
		contact := iter.Contact()
		if contact.Email.Verified {
			if err = s.beginReview(ctx, vasp, contact.Email.Email); err != nil {
				sentry.Error(ctx).Err(err).Str("contact", contact.Email.Email).Msg("could not begin review process with verified contact")
				return nil, status.Error(codes.Internal, "could not complete registration process")
			}
			break
		}
	}

	// Create PKCS12 password along with certificate request.
	var certRequest *models.CertificateRequest
	password := secrets.CreateToken(16)
	if certRequest, err = models.NewCertificateRequest(vasp); err != nil {
		sentry.Error(ctx).Err(err).Str("vasp", vasp.Id).Msg("could not create certificate request")
		return nil, status.Error(codes.Internal, "internal error with registration, please contact admins")
	}

	if err = models.UpdateCertificateRequestStatus(certRequest, models.CertificateRequestState_INITIALIZED, "created certificate request", email); err != nil {
		sentry.Error(ctx).Err(err).Str("vasp", vasp.Id).Msg("could not update certificate request status")
		return nil, status.Error(codes.Internal, "internal error with registration, please contact admins")
	}

	// Make a new secret of type "password"
	secretType := "password"
	if err = s.svc.secret.With(certRequest.Id).CreateSecret(ctx, secretType); err != nil {
		sentry.Error(ctx).Err(err).Str("vasp", vasp.Id).Msg("could not create new secret for pkcs12 password")
		return nil, status.Error(codes.Internal, "internal error with registration, please contact admins")
	}
	if err = s.svc.secret.With(certRequest.Id).AddSecretVersion(ctx, secretType, []byte(password)); err != nil {
		sentry.Error(ctx).Err(err).Str("vasp", vasp.Id).Msg("unable to add secret version for pkcs12 password")
		return nil, status.Error(codes.Internal, "internal error with registration, please contact admins")
	}

	// Create certificate request in the database.
	if err = s.db.UpdateCertReq(ctx, certRequest); err != nil {
		sentry.Error(ctx).Err(err).Str("vasp", vasp.Id).Msg("could not save certificate request")
		return nil, status.Error(codes.Internal, "internal error with registration, please contact admins")
	}

	// Add the CertificateRequest to the VASP
	if err = models.AppendCertReqID(vasp, certRequest.Id); err != nil {
		sentry.Error(ctx).Err(err).Str("vasp", vasp.Id).Msg("could not add cert request to VASP")
		return nil, status.Error(codes.Internal, "internal error with registration, please contact admins")
	}

	// Store VASP with updated certificate requests
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(ctx).Err(err).Str("vasp", vasp.Id).Msg("could not update vasp with certificate request ID")
		return nil, status.Error(codes.Internal, "internal error with registration, please contact admins")
	}

	out = &api.RegisterReply{
		Id:                  vasp.Id,
		RegisteredDirectory: vasp.RegisteredDirectory,
		CommonName:          vasp.CommonName,
		Status:              vasp.VerificationStatus,
		Message:             "a verification code has been sent to the contact email, please check spam folder if it has not arrived; pkcs12 password attached, this is the only time it will be available -- do not lose!",
		Pkcs12Password:      password,
	}
	return out, nil
}

// beginReview starts the registration review process by sending an email to the TRISA admins.
// This method does not update the passed vasp, the caller should be aware that they will need to
// update the vasp record. contactEmail will be used to update the vasp's verification status. The
// returned error will be a gRPC status error.
func (s *GDS) beginReview(ctx context.Context, vasp *pb.VASP, contactEmail string) (err error) {
	// Step 1: mark the VASP as email verified and create an admin token.
	if err := models.UpdateVerificationStatus(vasp, pb.VerificationState_EMAIL_VERIFIED, "completed email verification", contactEmail); err != nil {
		sentry.Warn(ctx).Err(err).Msg("could not update VASP verification status")
		return status.Error(codes.Aborted, "could not add new entry to VASP audit log")
	}

	// Create verification token for admin and update database
	// TODO: replace with actual authentication
	if err = models.SetAdminVerificationToken(vasp, secrets.CreateToken(48)); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not create admin verification token")
		return status.Error(codes.FailedPrecondition, "there was a problem submitting your registration review request, please contact the admins")
	}
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not save admin verification token")
		return status.Error(codes.FailedPrecondition, "there was a problem submitting your registration review request, please contact the admins")
	}

	// Step 2: send review request email to the TRISA admins.
	if _, err = s.svc.email.SendReviewRequest(vasp); err != nil {
		// TODO: When the Admin UI is up, downgrade FATAL to ERROR because the admins
		// can just check the UI for any pending reviews at that point (it is FATAL now
		// because without the email, the admins won't know there is a review).
		// Don't stop processing if review request email could not be sent.
		// NOTE: using WithLevel and Fatal does not Exit the program like log.Fatal()
		// this ensures that we issue a CRITICAL severity without stopping the server.
		sentry.Fatal(ctx).Err(err).Msg("could not send verification review email")
	} else {
		log.Info().Msg("verification review email sent to admins")
	}

	// Step 3: if the review email has been successfully sent, mark as pending review.
	if err := models.UpdateVerificationStatus(vasp, pb.VerificationState_PENDING_REVIEW, "review email sent", contactEmail); err != nil {
		sentry.Warn(ctx).Err(err).Msg("could not update VASP verification status")
		return status.Error(codes.Aborted, "could not add new entry to VASP audit log")
	}
	return nil
}

// Lookup a VASP entity by name or ID to get full details including the TRISA certification
// if it exists and the entity has been verified.
func (s *GDS) Lookup(ctx context.Context, in *api.LookupRequest) (out *api.LookupReply, err error) {
	var vasp *pb.VASP
	switch {
	case in.Id != "":
		// TODO: add registered directory to lookup
		if vasp, err = s.db.RetrieveVASP(ctx, in.Id); err != nil {
			log.Debug().Err(err).Str("id", in.Id).Str("registered_directory", in.RegisteredDirectory).Msg("could not find VASP by ID")
			return nil, status.Error(codes.NotFound, "could not find VASP by ID")
		}
	case in.CommonName != "":
		var vasps []*pb.VASP
		if vasps, err = s.db.SearchVASPs(ctx, map[string]interface{}{"name": in.CommonName}); err != nil {
			sentry.Warn(ctx).Err(err).Str("common_name", in.CommonName).Msg("could not search for common name")
			return nil, status.Error(codes.NotFound, "could not find VASP by common name")
		}

		if len(vasps) != 1 {
			// Don't warn when common name is not found, just when multiple results are returned
			if len(vasps) > 1 {
				log.Debug().Str("common_name", in.CommonName).Int("nresults", len(vasps)).Msg("multiple VASPs returned from common name search in lookup")
			} else {
				log.Debug().Msg("could not lookup VASP by common name")
			}
			return nil, status.Error(codes.NotFound, "could not find VASP by common name")
		}

		vasp = vasps[0]
	default:
		sentry.Warn(ctx).Str("rpc", "lookup").Msg("no arguments supplied")
		return nil, status.Error(codes.InvalidArgument, "please supply ID and registered directory or common name for lookup")
	}

	// Lookup should only return verified peers
	if vasp.VerificationStatus != pb.VerificationState_VERIFIED {
		return nil, status.Error(codes.NotFound, "no VASP record available")
	}

	// send lookup request activity to network activity handler
	activity.Lookup().VASP(vasp.Id).Add()

	// Prepare the response for the lookup
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
func (s *GDS) Search(ctx context.Context, in *api.SearchRequest) (out *api.SearchReply, err error) {
	// send search request activity to network activity handler
	activity.Search().Add()
	if _, ok := ctx.Deadline(); !ok {
		ctx, _ = utils.WithDeadline(ctx)
	}

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
	if vasps, err = s.db.SearchVASPs(ctx, query); err != nil {
		sentry.Error(ctx).Err(err).Msg("vasp search failed")
		return nil, status.Error(codes.Aborted, err.Error())
	}

	// Build search results to return
	out = &api.SearchReply{
		Results: make([]*api.SearchReply_Result, 0, len(vasps)),
	}
	for _, vasp := range vasps {
		// Filter out any VASPs that are not verified
		if vasp.VerificationStatus != pb.VerificationState_VERIFIED {
			continue
		}

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

// Verification returns the status of a VASP including its verification and service
// status if the directory service is performing health check monitoring.
func (s *GDS) Verification(ctx context.Context, in *api.VerificationRequest) (out *api.VerificationReply, err error) {
	var vasp *pb.VASP
	switch {
	case in.Id != "":
		// TODO: add registered directory to retrieve
		if vasp, err = s.db.RetrieveVASP(ctx, in.Id); err != nil {
			log.Debug().Err(err).Str("id", in.Id).Str("registered_directory", in.RegisteredDirectory).Msg("could not find VASP by ID")
			return nil, status.Error(codes.NotFound, "could not find VASP by ID")
		}
	case in.CommonName != "":
		var vasps []*pb.VASP
		if vasps, err = s.db.SearchVASPs(ctx, map[string]interface{}{"name": in.CommonName}); err != nil {
			sentry.Warn(ctx).Err(err).Str("common_name", in.CommonName).Msg("could not search for common name")
			return nil, status.Error(codes.NotFound, "could not find VASP by common name")
		}

		if len(vasps) != 1 {
			if len(vasps) > 1 {
				// Don't warn when common name is not found, just when multiple results are returned
				log.Debug().Str("common_name", in.CommonName).Int("nresults", len(vasps)).Msg("multiple VASPs returned from common name search in verification")
			} else {
				log.Debug().Msg("could not lookup VASP by common name")
			}
			return nil, status.Error(codes.NotFound, "could not find VASP by common name")
		}

		vasp = vasps[0]
	default:
		sentry.Warn(ctx).Str("rpc", "verification").Msg("no arguments supplied")
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
func (s *GDS) VerifyContact(ctx context.Context, in *api.VerifyContactRequest) (out *api.VerifyContactReply, err error) {
	if in.Token == "" {
		sentry.Warn(ctx).Msg("no verification token supplied")
		return nil, status.Error(codes.InvalidArgument, "could not verify contact: verification token missing from request")
	}

	if _, ok := ctx.Deadline(); !ok {
		ctx, _ = utils.WithDeadline(ctx)
	}

	// Retrieve VASP associated with contact from the database.
	var vasp *pb.VASP
	if vasp, err = s.db.RetrieveVASP(ctx, in.Id); err != nil {
		sentry.Warn(ctx).Err(err).Str("id", in.Id).Msg("could not retrieve vasp")
		return nil, status.Error(codes.NotFound, "could not find associated VASP record by ID")
	}

	found := false
	prevVerified := 0
	contactEmail := ""

	// Search through the contacts to determine the contacts verified by the supplied token.
	var contacts *models.Contacts
	if contacts, err = s.db.VASPContacts(ctx, vasp); err != nil {
		sentry.Error(ctx).Err(err).Str("id", in.Id).Msg("could not retrieve vasp contacts")
		return nil, status.Error(codes.Internal, "unable to complete contact verification")
	}

	iter := contacts.NewIterator(models.SkipDuplicates())
	for iter.Next() {
		contact := iter.Contact()

		// If the model contact is verified make sure the vasp contact is verified
		// BUG: if a previous contact is verified this will not return a 404 because
		// found is being set to true; this check has to happen after token == in.Token
		if contact.IsVerified() {
			found = true
			prevVerified++
			continue
		}

		// Perform token check and if token matches, mark contact as verified
		if contact.Email.Token == in.Token {
			found = true

			// Verify and update the models contact
			contact.Email.Verified = true
			contact.Email.VerifiedOn = time.Now().Format(time.RFC3339Nano)
			contact.Email.Token = ""

			if err = s.db.UpdateEmail(ctx, contact.Email); err != nil {
				sentry.Error(ctx).Err(err).Str("contact", contact.Email.Email).Msg("could not update contact verification status")
				return nil, status.Error(codes.Aborted, "could not update contact record")
			}

			// Record the contact as verified in the audit log
			if err := models.UpdateVerificationStatus(vasp, vasp.VerificationStatus, "contact verified", contact.Email.Email); err != nil {
				sentry.Warn(ctx).Err(err).Msg("could not append contact verification to VASP audit log")
				return nil, status.Error(codes.Aborted, "could not add new entry to VASP audit log")
			}
		}
	}

	// Check if we haven't managed to verify the contact
	if !found {
		sentry.Warn(ctx).Bool("found", found).Str("vasp", vasp.Id).Msg("could not find contact with token")
		return nil, status.Error(codes.NotFound, "could not find contact with the specified token")
	}

	// Ensures that we only send the verification email to the admins once.
	// If we have previously verified contacts, assume that we've already sent the
	// registration review email and do nothing.
	if prevVerified > 0 && vasp.VerificationStatus > pb.VerificationState_SUBMITTED {
		// Save the updated contact
		if err = s.db.UpdateVASP(ctx, vasp); err != nil {
			sentry.Error(ctx).Err(err).Msg("could not update VASP record after contact verification")
			return nil, status.Error(codes.Internal, "could not update contact after verification")
		}

		return &api.VerifyContactReply{
			Status:  vasp.VerificationStatus,
			Message: "email successfully verified",
		}, nil
	}

	if err = s.beginReview(ctx, vasp, contactEmail); err != nil {
		sentry.Error(ctx).Err(err).Str("vasp", vasp.Id).Msg("could not begin registration for contact with verified email")
		return nil, err
	}

	// Save the VASP and newly created certificate request
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not update vasp status to pending review")
		return nil, status.Error(codes.Internal, "there was a problem submitting your registration review request, please contact the admins")
	}

	return &api.VerifyContactReply{
		Status:  vasp.VerificationStatus,
		Message: "email successfully verified and verification review sent to TRISA admins",
	}, nil
}

func (s *GDS) Status(ctx context.Context, in *api.HealthCheck) (out *api.ServiceState, err error) {
	log.Debug().
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
	if s.svc.conf.Maintenance {
		out.Status = api.ServiceState_MAINTENANCE
	}

	return out, nil
}

//===========================================================================
// Helper Functions
//===========================================================================

// Get a valid email address and name from the contacts on a VASP.
func GetContactEmail(vasp *pb.VASP) string {
	if vasp.Contacts.Technical != nil && vasp.Contacts.Technical.Email != "" {
		return vasp.Contacts.Technical.Email
	}

	if vasp.Contacts.Administrative != nil && vasp.Contacts.Administrative.Email != "" {
		return vasp.Contacts.Administrative.Email
	}

	if vasp.Contacts.Legal != nil && vasp.Contacts.Legal.Email != "" {
		return vasp.Contacts.Legal.Email
	}

	if vasp.Contacts.Billing != nil && vasp.Contacts.Billing.Email != "" {
		return vasp.Contacts.Billing.Email
	}
	return ""
}

// Validate a gRPC endpoint string.
func validateEndpoint(endpoint string) (err error) {
	var host, port string
	if host, port, err = net.SplitHostPort(endpoint); err != nil {
		return errors.New("unable to parse endpoint string")
	}

	if host == "" {
		return errors.New("missing host in endpoint string")
	}

	if port == "" {
		return errors.New("missing port in endpoint string")
	}

	if _, err = strconv.Atoi(port); err != nil {
		return errors.New("endpoint port is not an integer")
	}
	return nil
}
