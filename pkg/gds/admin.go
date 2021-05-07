package trisads

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v1"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// Review a registration request and either accept or reject it. On accept, the
// certificate request that was created on verify is used to send a Sectigo request and
// the certificate manager process watches it until the certificate has been issued. On
// reject, the VASP and certificate request records are deleted and the reject reason is
// sent to the technical contact.
func (s *Server) Review(ctx context.Context, in *admin.ReviewRequest) (out *admin.ReviewReply, err error) {
	out = &admin.ReviewReply{}

	// Validate review request
	if in.Id == "" || in.AdminVerificationToken == "" {
		out.Error = &api.Error{
			Code:    400,
			Message: "provide both the VASP id to review and the admin verification token",
		}
		log.Error().Err(out.Error).Msg("bad review request")
		return out, nil
	}

	if !in.Accept && in.RejectReason == "" {
		out.Error = &api.Error{
			Code:    400,
			Message: "if rejecting the request, a reason must be supplied",
		}
		log.Error().Err(out.Error).Msg("bad review request")
		return out, nil
	}

	// Lookup the VASP record associated with the request
	var vasp *pb.VASP
	if vasp, err = s.db.Retrieve(in.Id); err != nil {
		log.Error().Err(err).Str("id", in.Id).Msg("could not retrieve vasp")
		out.Error = &api.Error{
			Code:    404,
			Message: err.Error(),
		}
		return out, nil
	}

	// Check that the administration verification token is correct
	if in.AdminVerificationToken != vasp.AdminVerificationToken {
		log.Error().Err(err).Str("token", in.AdminVerificationToken).Msg("incorrect admin verification token")
		out.Error = &api.Error{
			Code:    403,
			Message: "admin verification token not accepted",
		}
		return out, nil
	}

	// Accept or reject the request
	if in.Accept {
		if out.Message, err = s.acceptRegistration(vasp); err != nil {
			log.Error().Err(out.Error).Msg("could not accept VASP registration")
			out.Error = &api.Error{Code: 500, Message: "could not review VASP registration"}
			return out, nil
		}
	} else {
		if out.Message, err = s.rejectRegistration(vasp, in.RejectReason); err != nil {
			log.Error().Err(out.Error).Msg("could not reject VASP registration")
			out.Error = &api.Error{Code: 500, Message: "could not review VASP registration"}
			return out, nil
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
	vasp.AdminVerificationToken = ""
	vasp.VerificationStatus = pb.VerificationState_REVIEWED
	if err = s.db.Update(vasp); err != nil {
		return "", err
	}

	// Mark any initialized certificate requests for this VASP as ready to submit
	// NOTE: there should only be one certificate request per VASP, but no errors occur
	// if there are more than one (other than a logged warning).
	var ncertreqs int
	var careqs []*pb.CertificateRequest
	if careqs, err = s.db.ListCertRequests(); err != nil {
		return "", err
	}

	for _, req := range careqs {
		if req.Vasp == vasp.Id && req.Status == pb.CertificateRequestState_INITIALIZED {
			req.Status = pb.CertificateRequestState_READY_TO_SUBMIT
			if err = s.db.SaveCertRequest(req); err != nil {
				return "", err
			}
			ncertreqs++
		}
	}

	if ncertreqs == 0 {
		return "", errors.New("no certificate requests found for VASP registration")
	}

	if ncertreqs == 1 {
		log.Debug().Str("vasp", vasp.Id).Msg("certificate request marked as ready to submit")
	} else {
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
	vasp.AdminVerificationToken = ""
	vasp.VerificationStatus = pb.VerificationState_REJECTED
	if err = s.db.Update(vasp); err != nil {
		return "", err
	}

	// Delete all pending certificate requests
	var ncertreqs int
	var careqs []*pb.CertificateRequest
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
	if err = s.RejectRegistrationEmail(vasp, reason); err != nil {
		return "", err
	}

	// Send successful response
	var name string
	if name, err = vasp.Name(); err != nil {
		name = vasp.Id
	}
	return fmt.Sprintf("registration request for %s has been rejected and its contacts notified", name), nil
}
