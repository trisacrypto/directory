package gds

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v1"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewAdmin creates a new GDS admin server derived from a parent Service.
func NewAdmin(svc *Service) (a *Admin, err error) {
	a = &Admin{
		svc:  svc,
		conf: &svc.conf.Admin,
		db:   svc.db,
	}

	// Initialize the gRPC server
	a.srv = grpc.NewServer(grpc.UnaryInterceptor(svc.serverInterceptor))
	admin.RegisterDirectoryAdministrationServer(a.srv, a)
	return a, nil
}

// Admin implements the DirectoryAdministrationServer as defined by the v1 or later GDS
// protocol buffers. This service is the primary interaction point with authorized TRISA
// users that are performing secure commands with authentication.
type Admin struct {
	admin.UnimplementedDirectoryAdministrationServer
	svc  *Service            // The parent Service the admin server uses to interact with other components
	srv  *grpc.Server        // The gRPC server that listens on its own independent port
	conf *config.AdminConfig // The admin server specific configuration (alias to s.svc.conf.Admin)
	db   store.Store         // Database connection for loading objects (alias to s.svc.db)
}

// Serve GRPC requests on the specified address.
func (s *Admin) Serve() (err error) {
	// This service should not start in maintenance mode.
	if s.svc.conf.Maintenance {
		return errors.New("could not start directory administration service in maintenance mode")
	}

	if !s.conf.Enabled {
		log.Warn().Msg("directory administration service is not enabled")
		return nil
	}

	// Listen for TCP requests on the specified address and port
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		return fmt.Errorf("could not listen on %q", s.conf.BindAddr)
	}

	// Run the server
	go func() {
		defer sock.Close()
		log.Info().
			Str("listen", s.conf.BindAddr).
			Str("version", pkg.Version()).
			Msg("directory administration server started")

		if err := s.srv.Serve(sock); err != nil {
			s.svc.echan <- err
		}
	}()

	// Now that the go routine is started return nil, meaning the service has started
	// successfully with no problems.
	return nil
}

// Shutdown the Directory Administration Service gracefully
func (s *Admin) Shutdown() (err error) {
	log.Debug().Msg("gracefully shutting down admin server")
	s.srv.GracefulStop()
	log.Debug().Msg("successful shutdown of admin server")
	return nil
}

// Review a registration request and either accept or reject it. On accept, the
// certificate request that was created on verify is used to send a Sectigo request and
// the certificate manager process watches it until the certificate has been issued. On
// reject, the VASP and certificate request records are deleted and the reject reason is
// sent to the technical contact.
func (s *Admin) Review(ctx context.Context, in *admin.ReviewRequest) (out *admin.ReviewReply, err error) {
	// Validate review request
	if in.Id == "" || in.AdminVerificationToken == "" {
		log.Error().Err(out.Error).Msg("no ID or verification token")
		return nil, status.Error(codes.InvalidArgument, "provide both the VASP ID and the admin verification token")
	}

	if !in.Accept && in.RejectReason == "" {
		log.Error().Err(out.Error).Msg("missing reject reason")
		return nil, status.Error(codes.InvalidArgument, "if rejecting the request, a reason must be supplied")
	}

	// Lookup the VASP record associated with the request
	var vasp *pb.VASP
	if vasp, err = s.db.RetrieveVASP(in.Id); err != nil {
		log.Warn().Err(err).Str("id", in.Id).Msg("could not retrieve vasp")
		return nil, status.Error(codes.NotFound, "could not retrieve VASP record by ID")
	}

	// Check that the administration verification token is correct
	var adminVerificationToken string
	if adminVerificationToken, err = models.GetAdminVerificationToken(vasp); err != nil {
		log.Error().Err(err).Msg("could not retrieve admin token from extra data field on VASP")
		return nil, status.Error(codes.Internal, "could not retrieve admin token from data")
	}
	if in.AdminVerificationToken != adminVerificationToken {
		log.Warn().Err(err).Str("vaps", in.Id).Msg("incorrect admin verification token")
		return nil, status.Error(codes.Unauthenticated, "admin verification token not accepted")
	}

	// Accept or reject the request
	out = &admin.ReviewReply{}
	if in.Accept {
		if out.Message, err = s.acceptRegistration(vasp); err != nil {
			log.Error().Err(err).Msg("could not accept VASP registration")
			return nil, status.Error(codes.FailedPrecondition, "unable to accept VASP registration request")
		}
	} else {
		if out.Message, err = s.rejectRegistration(vasp, in.RejectReason); err != nil {
			log.Error().Err(err).Msg("could not reject VASP registration")
			return nil, status.Error(codes.FailedPrecondition, "unable to reject VASP registration request")
		}
	}

	name, _ := vasp.Name()
	out.Status = vasp.VerificationStatus
	log.Info().Str("vasp", vasp.Id).Str("name", name).Bool("accepted", in.Accept).Msg("registration reviewed")
	return out, nil
}

// Accept the VASP registration and begin the certificate issuance process.
func (s *Admin) acceptRegistration(vasp *pb.VASP) (msg string, err error) {
	// Change the VASP verification status
	if err = models.SetAdminVerificationToken(vasp, ""); err != nil {
		return "", err
	}
	vasp.VerifiedOn = time.Now().Format(time.RFC3339)
	vasp.VerificationStatus = pb.VerificationState_REVIEWED
	if err = s.db.UpdateVASP(vasp); err != nil {
		return "", err
	}

	// Mark any initialized certificate requests for this VASP as ready to submit
	// NOTE: there should only be one certificate request per VASP, but no errors occur
	// if there are more than one (other than a logged warning).
	var ncertreqs int
	var careqs []*models.CertificateRequest
	if careqs, err = s.db.ListCertReqs(); err != nil {
		return "", err
	}

	for _, req := range careqs {
		if req.Vasp == vasp.Id && req.Status == models.CertificateRequestState_INITIALIZED {
			req.Status = models.CertificateRequestState_READY_TO_SUBMIT
			if err = s.db.UpdateCertReq(req); err != nil {
				return "", err
			}
			ncertreqs++
		}
	}

	switch ncertreqs {
	case 0:
		return "", errors.New("no certificate requests found for VASP registration")
	case 1:
		log.Debug().Str("vasp", vasp.Id).Msg("certificate request marked as ready to submit")
	default:
		log.Warn().Str("vasp", vasp.Id).Int("requests", ncertreqs).Msg("multiple certificate requests marked as ready to submit")
	}

	// Send successful response
	var name string
	if name, err = vasp.Name(); err != nil {
		name = vasp.Id
	}
	return fmt.Sprintf("registration request for %s has been approved and a Sectigo certificate will be requested", name), nil
}

// Reject the VASP registration and notify the contacts of the result.
func (s *Admin) rejectRegistration(vasp *pb.VASP, reason string) (msg string, err error) {
	// Change the VASP verification status
	if err = models.SetAdminVerificationToken(vasp, ""); err != nil {
		return "", err
	}
	vasp.VerificationStatus = pb.VerificationState_REJECTED
	if err = s.db.UpdateVASP(vasp); err != nil {
		return "", err
	}

	// Delete all pending certificate requests
	var ncertreqs int
	var careqs []*models.CertificateRequest
	if careqs, err = s.db.ListCertReqs(); err != nil {
		return "", err
	}

	for _, req := range careqs {
		if req.Vasp == vasp.Id {
			if err = s.db.DeleteCertReq(req.Id); err != nil {
				log.Error().Err(err).Str("id", req.Id).Msg("could not delete certificate request")
			}
			ncertreqs++
		}
	}

	// Log deletion of certificate requests
	switch ncertreqs {
	case 0:
		log.Warn().Str("vasp", vasp.Id).Msg("no certificate requests deleted")
	case 1:
		log.Debug().Str("vasp", vasp.Id).Msg("certificate request deleted")
	default:
		log.Warn().Str("vasp", vasp.Id).Msg("multiple certificate requests deleted")
	}

	// Notify the VASP contacts that the registration request has been rejected.
	if _, err = s.svc.email.SendRejectRegistration(vasp, reason); err != nil {
		return "", err
	}

	// Send successful response
	var name string
	if name, err = vasp.Name(); err != nil {
		name = vasp.Id
	}
	return fmt.Sprintf("registration request for %s has been rejected and its contacts notified", name), nil
}

// Resend emails in case they went to spam or the initial email send failed.
func (s *Admin) Resend(ctx context.Context, in *admin.ResendRequest) (out *admin.ResendReply, err error) {
	if in.Id == "" {
		log.Warn().Msg("invalid resend request: missing ID")
		return nil, status.Error(codes.InvalidArgument, "VASP record ID is required")
	}

	// Lookup the VASP record associated with the resend request
	var vasp *pb.VASP
	if vasp, err = s.db.RetrieveVASP(in.Id); err != nil {
		log.Warn().Err(err).Str("id", in.Id).Msg("could not retrieve vasp")
		return nil, status.Error(codes.NotFound, "could not retrieve VASP record by ID")
	}

	var sent int
	out = &admin.ResendReply{}

	// Handle different resend request types
	switch in.Type {
	case admin.ResendRequest_UNKNOWN:
		log.Warn().Msg("invalid resend request: unknown type")
		return nil, status.Error(codes.InvalidArgument, "specify a resend emails type")

	case admin.ResendRequest_VERIFY_CONTACT:
		if sent, err = s.svc.email.SendVerifyContacts(vasp); err != nil {
			log.Warn().Err(err).Int("sent", sent).Msg("could not resend verify contacts emails")
			return nil, status.Errorf(codes.FailedPrecondition, "could not resend contact verification emails: %s", err)
		}
		out.Message = "contact verification emails resent to all unverified contacts"

	case admin.ResendRequest_REVIEW:
		if sent, err = s.svc.email.SendReviewRequest(vasp); err != nil {
			log.Warn().Err(err).Int("sent", sent).Msg("could not resend review request")
			return nil, status.Errorf(codes.FailedPrecondition, "could not resend review request: %s", err)
		}
		out.Message = "review request resent to TRISA admins"

	case admin.ResendRequest_DELIVER_CERTS:
		// TODO: check verification state and cert request state
		// TODO: in order to implement this, we'd have to fetch the certs from Google Secrets
		// TODO: if implemented, log which contact was sent the certs (e.g. technical, admin, etc.)
		// TODO: when above implemented, also log which contact was sent certs in acceptRegistration
		return nil, status.Error(codes.Unimplemented, "resend cert delivery not yet implemented")

	case admin.ResendRequest_REJECTION:
		// Only send a rejection email if we're in the rejected state
		if vasp.VerificationStatus != pb.VerificationState_REJECTED {
			log.Warn().Err(err).Str("status", vasp.VerificationStatus.String()).Msg("cannot resend rejection emails in current state")
			return nil, status.Error(codes.FailedPrecondition, "VASP record verification status cannot send rejection email")
		}

		// A reason must be specified to send a rejection email (it's not stored)
		if in.Reason == "" {
			log.Warn().Str("resend_type", in.Type.String()).Msg("invalid resend request: missing reason argument")
			return nil, status.Error(codes.InvalidArgument, "must specify reason for rejection to resend email")
		}
		if sent, err = s.svc.email.SendRejectRegistration(vasp, in.Reason); err != nil {
			log.Warn().Err(err).Int("sent", sent).Msg("could not resend rejection emails")
			return nil, status.Errorf(codes.FailedPrecondition, "could not resend rejection emails: %s", err)
		}
		out.Message = "rejection emails resent to all verified contacts"

	default:
		log.Warn().Str("resend_type", in.Type.String()).Msg("invalid resend request: unhandled resend request type")
		return nil, status.Errorf(codes.FailedPrecondition, "unknown resend request type %q", in.Type)
	}

	out.Sent = int64(sent)
	log.Info().Str("id", vasp.Id).Int64("sent", out.Sent).Str("resend_type", in.Type.String()).Msg("resend request complete")
	return out, nil
}

// Get current counts of registration statuses and certificate requests.
func (s *Admin) Status(ctx context.Context, in *admin.StatusRequest) (out *admin.StatusReply, err error) {
	var nvasps, ncertreqs int64

	out = &admin.StatusReply{}
	if !in.NoRegistrations {
		out.Registrations = make(map[string]int64)
		var vasps []*pb.VASP
		if vasps, err = s.db.ListVASPs(); err != nil {
			log.Error().Err(err).Msg("could not list vasps from database")
			return nil, status.Error(codes.Internal, "a database error occurred")
		}
		for _, vasp := range vasps {
			nvasps++
			out.Registrations[vasp.VerificationStatus.String()]++
		}
	}

	if !in.NoCertificateRequests {
		out.CertificateRequests = make(map[string]int64)
		var certreqs []*models.CertificateRequest
		if certreqs, err = s.db.ListCertReqs(); err != nil {
			log.Error().Err(err).Msg("could not list certificate requests from database")
			return nil, status.Error(codes.Internal, "a database error occurred")
		}
		for _, certreq := range certreqs {
			ncertreqs++
			out.CertificateRequests[certreq.Status.String()]++
		}
	}

	log.Info().Int64("nvasps", nvasps).Int64("ncertreqs", ncertreqs).Msg("status counts complete")
	return out, nil
}
