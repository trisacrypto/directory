package bff

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
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
		case testnet:
			vaspID = testnetID
		case mainnet:
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
