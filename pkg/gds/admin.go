package gds

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v1"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Review a registration request and either accept or reject it. On accept, the
// certificate request that was created on verify is used to send a Sectigo request and
// the certificate manager process watches it until the certificate has been issued. On
// reject, the VASP and certificate request records are deleted and the reject reason is
// sent to the technical contact.
func (s *Server) Review(ctx context.Context, in *admin.ReviewRequest) (out *admin.ReviewReply, err error) {
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
	if vasp, err = s.db.Retrieve(in.Id); err != nil {
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
func (s *Server) acceptRegistration(vasp *pb.VASP) (msg string, err error) {
	// Change the VASP verification status
	if err = models.SetAdminVerificationToken(vasp, ""); err != nil {
		return "", err
	}
	vasp.VerifiedOn = time.Now().Format(time.RFC3339)
	vasp.VerificationStatus = pb.VerificationState_REVIEWED
	if err = s.db.Update(vasp); err != nil {
		return "", err
	}

	// Mark any initialized certificate requests for this VASP as ready to submit
	// NOTE: there should only be one certificate request per VASP, but no errors occur
	// if there are more than one (other than a logged warning).
	var ncertreqs int
	var careqs []*models.CertificateRequest
	if careqs, err = s.db.ListCertRequests(); err != nil {
		return "", err
	}

	for _, req := range careqs {
		if req.Vasp == vasp.Id && req.Status == models.CertificateRequestState_INITIALIZED {
			req.Status = models.CertificateRequestState_READY_TO_SUBMIT
			if err = s.db.SaveCertRequest(req); err != nil {
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
func (s *Server) rejectRegistration(vasp *pb.VASP, reason string) (msg string, err error) {
	// Change the VASP verification status
	if err = models.SetAdminVerificationToken(vasp, ""); err != nil {
		return "", err
	}
	vasp.VerificationStatus = pb.VerificationState_REJECTED
	if err = s.db.Update(vasp); err != nil {
		return "", err
	}

	// Delete all pending certificate requests
	var ncertreqs int
	var careqs []*models.CertificateRequest
	if careqs, err = s.db.ListCertRequests(); err != nil {
		return "", err
	}

	for _, req := range careqs {
		if req.Vasp == vasp.Id {
			if err = s.db.DeleteCertRequest(req.Id); err != nil {
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
	if err = s.email.SendRejectRegistration(vasp, reason); err != nil {
		return "", err
	}

	// Send successful response
	var name string
	if name, err = vasp.Name(); err != nil {
		name = vasp.Id
	}
	return fmt.Sprintf("registration request for %s has been rejected and its contacts notified", name), nil
}
