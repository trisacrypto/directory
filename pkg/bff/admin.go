package bff

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/config"
	records "github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// GetCertificates makes parallel calls to the admin services to get the certificate
// information for both testnet and mainnet. If testnetID or mainnetID are empty
// strings, this will simply return a nil response for the corresponding network so
// the caller can distinguish between a non registration and an error.
func (s *Server) GetCertificates(ctx context.Context, testnetID, mainnetID string) (testnetCerts, mainnetCerts *admin.ListCertificatesReply, testnetErr, mainnetErr error) {
	// Create the RPC which can do both testnet and mainnet calls
	rpc := func(ctx context.Context, client admin.DirectoryAdministrationClient, network string) (rep interface{}, err error) {
		var vaspID string
		switch network {
		case config.TestNet:
			vaspID = testnetID
		case config.MainNet:
			vaspID = mainnetID
		default:
			return nil, fmt.Errorf("unknown network: %s", network)
		}

		if vaspID == "" {
			// The VASP is not registered for this network, so do not error and return
			// nil
			return nil, nil
		}
		return client.ListCertificates(ctx, vaspID)
	}

	// Perform the parallel requests
	results, errs := s.ParallelAdminRequests(ctx, rpc, false)
	if len(errs) != 2 || len(results) != 2 {
		err := fmt.Errorf("unexpected number of results from parallel requests: %d", len(results))
		return nil, nil, err, err
	}

	// Parse the results
	var ok bool
	if errs[0] != nil {
		testnetErr = errs[0]
	} else if results[0] != nil {
		if testnetCerts, ok = results[0].(*admin.ListCertificatesReply); !ok {
			testnetErr = fmt.Errorf("unexpected testnet result type returned from parallel certificate requests: %T", results[0])
		}
	}

	if errs[1] != nil {
		mainnetErr = errs[1]
	} else if results[1] != nil {
		if mainnetCerts, ok = results[1].(*admin.ListCertificatesReply); !ok {
			mainnetErr = fmt.Errorf("unexpected mainnet result type returned from parallel certificate requests: %T", results[1])
		}
	}

	return testnetCerts, mainnetCerts, testnetErr, mainnetErr
}

// Certificates returns the list of certificates for the authenticated user.
func (s *Server) Certificates(c *gin.Context) {
	var err error

	// Get the bff claims from the context
	var claims *auth.Claims
	if claims, err = auth.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("unable to retrieve bff claims from context")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Extract the VASP IDs from the claims
	// Note that if testnet or mainnet are absent from the VASPs struct, the ID will
	// default to an empty string, and GetCertificates will return nil for that network
	// instead of an error.
	testnetID := claims.VASPs.TestNet
	mainnetID := claims.VASPs.MainNet

	// Get the certificate replies from the admin APIs
	testnet, mainnet, testnetErr, mainnetErr := s.GetCertificates(c.Request.Context(), testnetID, mainnetID)

	// Construct the response
	out := &api.CertificatesReply{
		Error:   api.NetworkError{},
		TestNet: make([]api.Certificate, 0),
		MainNet: make([]api.Certificate, 0),
	}

	// Populate the testnet response
	if testnetErr != nil {
		out.Error.TestNet = testnetErr.Error()
	} else if testnet != nil {
		for _, cert := range testnet.Certificates {
			out.TestNet = append(out.TestNet, api.Certificate{
				SerialNumber: cert.SerialNumber,
				IssuedAt:     cert.IssuedAt,
				ExpiresAt:    cert.ExpiresAt,
				Revoked:      cert.Status == models.CertificateState_REVOKED.String(),
				Details:      cert.Details,
			})
		}
	}

	// Populate the mainnet response
	if mainnetErr != nil {
		out.Error.MainNet = mainnetErr.Error()
	} else if mainnet != nil {
		for _, cert := range mainnet.Certificates {
			out.MainNet = append(out.MainNet, api.Certificate{
				SerialNumber: cert.SerialNumber,
				IssuedAt:     cert.IssuedAt,
				ExpiresAt:    cert.ExpiresAt,
				Revoked:      cert.Status == models.CertificateState_REVOKED.String(),
				Details:      cert.Details,
			})
		}
	}

	c.JSON(http.StatusOK, out)
}

const (
	testnetName          = "TestNet"
	mainnetName          = "MainNet"
	supportEmail         = "support@rotational.io"
	StartRegistration    = "Start the registration and verification process for your organization to receive an X.509 Identity Certificate and become a trusted member of the TRISA network."
	CompleteRegistration = "Complete the registration process and verification process for your organization to receive an X.509 Identity Certificate and become a trusted member of the TRISA network."
	SubmitTestnet        = "Review and submit your " + testnetName + " registration."
	SubmitMainnet        = "Review and submit your " + mainnetName + " registration."
	VerifyEmails         = "Your organization's %s registration has been submitted and verification emails have been sent to the contacts specified in the form. Contacts and email addresses must be verified as the first step in the approval process. Please request that contacts verify their email addresses promptly so that the TRISA Validation Team can proceed with the validation process. Please contact TRISA support at " + supportEmail + " if contacts have not received the verification email and link."
	RegistrationPending  = "Your organization's %s registration has been received and is pending approval. The TRISA Validation Team will notify you about the outcome."
	RegistrationRejected = "Your organization's %s registration has been rejected by the TRISA Validation Team. This means your organization is not a verified member of the TRISA network and cannot communicate with other members. Please contact TRISA support at " + supportEmail + " for additional details and next steps."
	RegistrationApproved = "Your organization's %s registration has been approved by the TRISA Validation Team. Take the next steps to integrate, test, and begin sending complicance messages with TRISA-verified counterparties."
	RenewCertificate     = "Your organization's %s X.509 Identity Certificate will expire on %s. Start the renewal process to receive a new X.509 Identity Certificate and remain a trusted member of the TRISA network."
	CertificateRevoked   = "Your organization's %s X.509 Identity Certificate has been revoked by TRISA. This means your organization is no longer a verified member of the TRISA network and can no longer communicate with other members. Please contact TRISA support at " + supportEmail + " for additional details and next steps."
)

// GetVASPs makes parallel calls to the admin APIs to retrieve VASP records from
// testnet and mainnet. If testnet or mainnet are empty strings, this will simply
// return a nil response for the corresponding network so the caller can distinguish
// between a non registration and an error.
func (s *Server) GetVASPs(ctx context.Context, testnetID, mainnetID string) (testnetVASP, mainnetVASP *pb.VASP, testnetErr, mainnetErr error) {
	// Create the RPC which can do both testnet and mainnet calls
	rpc := func(ctx context.Context, client admin.DirectoryAdministrationClient, network string) (rep interface{}, err error) {
		var vaspID string
		switch network {
		case config.TestNet:
			vaspID = testnetID
		case config.MainNet:
			vaspID = mainnetID
		default:
			return nil, fmt.Errorf("unknown network: %s", network)
		}

		if vaspID == "" {
			// The VASP is not registered for this network, so do not error and return
			// nil
			return nil, nil
		}
		return client.RetrieveVASP(ctx, vaspID)
	}

	// Perform the parallel requests
	results, errs := s.ParallelAdminRequests(ctx, rpc, false)
	if len(errs) != 2 || len(results) != 2 {
		err := fmt.Errorf("unexpected number of results from parallel requests: %d", len(results))
		return nil, nil, err, err
	}

	// Parse the results
	if errs[0] != nil {
		testnetErr = errs[0]
	} else if results[0] != nil {
		if testnetReply, ok := results[0].(*admin.RetrieveVASPReply); ok {
			testnetVASP = &pb.VASP{}
			if err := wire.Unwire(testnetReply.VASP, testnetVASP); err != nil {
				testnetErr = fmt.Errorf("could not unwire testnet result to VASP: %s", err)
				testnetVASP = nil
			}
		} else {
			testnetErr = fmt.Errorf("unexpected testnet result type returned from parallel certificate requests: %T", results[0])
		}
	}

	if errs[1] != nil {
		mainnetErr = errs[1]
	} else if results[1] != nil {
		if mainnetReply, ok := results[1].(*admin.RetrieveVASPReply); ok {
			mainnetVASP = &pb.VASP{}
			if err := wire.Unwire(mainnetReply.VASP, mainnetVASP); err != nil {
				mainnetErr = fmt.Errorf("could not unwire mainnet result to VASP: %s", err)
				mainnetVASP = nil
			}
		} else {
			mainnetErr = fmt.Errorf("unexpected mainnet result type returned from parallel certificate requests: %T", results[1])
		}
	}

	return testnetVASP, mainnetVASP, testnetErr, mainnetErr
}

// registrationMessage returns a message corresponding to the registration state of the
// VASP. These states are distinct, so only one message is returned.
func registrationMessage(vasp *pb.VASP, network string) (msg *api.AttentionMessage, err error) {
	const expireLayout = "January 2, 2006"

	if vasp == nil {
		return nil, nil
	}

	switch {
	case vasp.VerificationStatus == pb.VerificationState_SUBMITTED:
		// Verify contact emails have been sent and are pending verification
		return &api.AttentionMessage{
			Message:  fmt.Sprintf(VerifyEmails, network),
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_VERIFY_EMAILS.String(),
		}, nil
	case vasp.VerificationStatus > pb.VerificationState_SUBMITTED && vasp.VerificationStatus < pb.VerificationState_VERIFIED:
		// The VASP is pending review and certificate issuance
		return &api.AttentionMessage{
			Message:  fmt.Sprintf(RegistrationPending, network),
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_NO_ACTION.String(),
		}, nil
	case vasp.IdentityCertificate != nil && vasp.IdentityCertificate.Revoked:
		// The VASP's certificate has been revoked
		return &api.AttentionMessage{
			Message:  fmt.Sprintf(CertificateRevoked, network),
			Severity: records.AttentionSeverity_ALERT.String(),
			Action:   records.AttentionAction_CONTACT_SUPPORT.String(),
		}, nil
	case vasp.IdentityCertificate != nil:
		// Certificate has been issued, check if it is about to expire
		var expiresAt time.Time
		if expiresAt, err = time.Parse(time.RFC3339, vasp.IdentityCertificate.NotAfter); err != nil {
			return nil, err
		}

		// Warn if less than 30 days before expiration
		if time.Until(expiresAt) < 30*24*time.Hour {
			return &api.AttentionMessage{
				Message:  fmt.Sprintf(RenewCertificate, network, expiresAt.Format(expireLayout)),
				Severity: records.AttentionSeverity_WARNING.String(),
				Action:   records.AttentionAction_RENEW_CERTIFICATE.String(),
			}, nil
		}
	case vasp.VerificationStatus == pb.VerificationState_VERIFIED:
		// The VASP is verified and the certificate has been issued
		return &api.AttentionMessage{
			Message:  fmt.Sprintf(RegistrationApproved, network),
			Severity: records.AttentionSeverity_SUCCESS.String(),
			Action:   records.AttentionAction_NO_ACTION.String(),
		}, nil
	case vasp.VerificationStatus == pb.VerificationState_REJECTED:
		// The VASP has been rejected, so no certificate was issued
		return &api.AttentionMessage{
			Message:  fmt.Sprintf(RegistrationRejected, network),
			Severity: records.AttentionSeverity_ALERT.String(),
			Action:   records.AttentionAction_CONTACT_SUPPORT.String(),
		}, nil
	default:
	}
	return nil, nil
}

// Attention returns the current attention messages for the authenticated user.
func (s *Server) Attention(c *gin.Context) {
	var err error

	// Get the bff claims from the context
	var claims *auth.Claims
	if claims, err = auth.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("unable to retrieve bff claims from context")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Retrieve the organization from the claims
	// NOTE: This method handles the error logging and response.
	var org *records.Organization
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Attention messages to return
	messages := make([]*api.AttentionMessage, 0)

	// Check the registration state, at most one of these messages will be returned
	testnetSubmitted := (org.Testnet != nil && org.Testnet.Submitted != "")
	mainnetSubmitted := (org.Mainnet != nil && org.Mainnet.Submitted != "")
	switch {
	case org.Registration == nil || org.Registration.State == nil || org.Registration.State.Started == "":
		// Registration has not started
		messages = append(messages, &api.AttentionMessage{
			Message:  StartRegistration,
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_START_REGISTRATION.String(),
		})
	case !testnetSubmitted && !mainnetSubmitted:
		// Registration has started but has not been completed
		messages = append(messages, &api.AttentionMessage{
			Message:  CompleteRegistration,
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_COMPLETE_REGISTRATION.String(),
		})
	case testnetSubmitted && !mainnetSubmitted:
		// Registration is submitted for testnet but not for mainnet
		messages = append(messages, &api.AttentionMessage{
			Message:  SubmitMainnet,
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_SUBMIT_MAINNET.String(),
		})
	case !testnetSubmitted && mainnetSubmitted:
		// Registration is submitted for mainnet but not for testnet
		messages = append(messages, &api.AttentionMessage{
			Message:  SubmitTestnet,
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_SUBMIT_TESTNET.String(),
		})
	default:
	}

	// Get the VASP records from the admin APIs
	// NOTE: This will not attempt to retrieve the VASP records if the VASP ID is not
	// set in the claims for a network and a nil result will be returned instead for
	// that network.
	testnetVASP, mainnetVASP, testnetErr, mainnetErr := s.GetVASPs(c.Request.Context(), claims.VASPs.TestNet, claims.VASPs.MainNet)

	if testnetErr != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(testnetErr))
		return
	}

	if mainnetErr != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(mainnetErr))
		return
	}

	// Get attention messages relating to certificates
	var testnetMsg *api.AttentionMessage
	if testnetMsg, err = registrationMessage(testnetVASP, testnetName); err != nil {
		log.Error().Err(err).Msg("could not get testnet certificate attention message")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}
	if testnetMsg != nil {
		messages = append(messages, testnetMsg)
	}

	var mainnetMsg *api.AttentionMessage
	if mainnetMsg, err = registrationMessage(mainnetVASP, mainnetName); err != nil {
		log.Error().Err(err).Msg("could not get mainnet certificate attention message")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}
	if mainnetMsg != nil {
		messages = append(messages, mainnetMsg)
	}

	// Build the response
	if len(messages) == 0 {
		c.JSON(http.StatusNoContent, nil)
	} else {
		c.JSON(http.StatusOK, &api.AttentionReply{
			Messages: messages,
		})
	}
}

// RegistrationStatus returns the registration status for both testnet and mainnet for
// the user.
func (s *Server) RegistrationStatus(c *gin.Context) {
	var err error

	// Retrieve the organization from the claims
	// NOTE: This method handles the error logging and response.
	var org *records.Organization
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Build the response
	// TODO: We should be querying the VASP record instead to allow for re-registration
	out := &api.RegistrationStatus{}
	if org.Testnet != nil && org.Testnet.Submitted != "" {
		out.TestNetSubmitted = org.Testnet.Submitted
	}
	if org.Mainnet != nil && org.Mainnet.Submitted != "" {
		out.MainNetSubmitted = org.Mainnet.Submitted
	}
	c.JSON(http.StatusOK, out)
}
