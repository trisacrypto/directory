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
	apiv2 "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
)

// GetCertificates makes parallel calls to the admin services to get the certificate
// information for both testnet and mainnet. If testnetID or mainnetID are empty
// strings, this will simply return a nil response for the corresponding network so
// the caller can distinguish between a non registration and an error.
func (s *Server) GetCertificates(ctx context.Context, testnetID, mainnetID string) (testcerts *apiv2.ListCertificatesReply, maincerts *apiv2.ListCertificatesReply, err error) {
	// Create the RPC which can do both testnet and mainnet calls
	rpc := func(ctx context.Context, client apiv2.DirectoryAdministrationClient, network string) (rep interface{}, err error) {
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
		return nil, nil, fmt.Errorf("unexpected number of results from parallel requests: %d", len(results))
	}

	// Parse the results
	var ok bool
	if errs[0] != nil {
		return nil, nil, errs[0]
	}

	if results[0] != nil {
		if testcerts, ok = results[0].(*apiv2.ListCertificatesReply); !ok {
			return nil, nil, fmt.Errorf("unexpected testnet result type returned from parallel certificate requests: %T", results[0])
		}
	}

	if errs[1] != nil {
		return nil, nil, errs[1]
	}

	if results[1] != nil {
		if maincerts, ok = results[1].(*apiv2.ListCertificatesReply); !ok {
			return nil, nil, fmt.Errorf("unexpected mainnet result type returned from parallel certificate requests: %T", results[1])
		}
	}

	return testcerts, maincerts, nil
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
	testnetID := claims.VASP[testnet]
	mainnetID := claims.VASP[mainnet]

	// Get the certificate replies from the admin APIs
	var testnet, mainnet *admin.ListCertificatesReply
	if testnet, mainnet, err = s.GetCertificates(c.Request.Context(), testnetID, mainnetID); err != nil {
		log.Error().Err(err).Msg("unable to get certificates from the admin APIs")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Construct the response
	out := &api.CertificatesReply{
		TestNet: make([]api.Certificate, 0),
		MainNet: make([]api.Certificate, 0),
	}

	// Populate the testnet response
	if testnet != nil {
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
	if mainnet != nil {
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
