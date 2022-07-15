package bff

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

// GetSummaries makes parallel calls to the members service to get the summary
// information for both testnet and mainnet. If an endpoint returned an error, then a
// nil value is returned from this function for that endpoint instead of an error.
func (s *Server) GetSummaries(ctx context.Context, testnetID, mainnetID string) (testnetSummary, mainnetSummary *members.SummaryReply, testnetErr, mainnetErr error) {
	// Create the RPC which can do both testnet and mainnet calls
	rpc := func(ctx context.Context, client GlobalDirectoryClient, network string) (rep proto.Message, err error) {
		req := &members.SummaryRequest{}
		switch network {
		case testnet:
			req.MemberId = testnetID
		case mainnet:
			req.MemberId = mainnetID
		default:
			return nil, fmt.Errorf("unknown network: %s", network)
		}
		return client.Summary(ctx, req)
	}

	// Perform the parallel requests
	results, errs := s.ParallelGDSRequests(ctx, rpc, false)
	if len(errs) != 2 || len(results) != 2 {
		err := fmt.Errorf("unexpected number of results from parallel requests: %d", len(results))
		return nil, nil, err, err
	}

	// Parse the results
	var ok bool
	if errs[0] != nil {
		testnetErr = errs[0]
	} else if results[0] == nil {
		testnetErr = fmt.Errorf("nil result returned from parallel requests")
	} else if testnetSummary, ok = results[0].(*members.SummaryReply); !ok {
		testnetErr = fmt.Errorf("unexpected summary result type returned from parallel requests: %T", results[0])
	}

	if errs[1] != nil {
		mainnetErr = errs[1]
	} else if results[1] == nil {
		mainnetErr = fmt.Errorf("nil result returned from parallel requests")
	} else if mainnetSummary, ok = results[1].(*members.SummaryReply); !ok {
		mainnetErr = fmt.Errorf("unexpected summary result type returned from parallel requests: %T", results[1])
	}

	return testnetSummary, mainnetSummary, testnetErr, mainnetErr
}

// Overview endpoint is an authenticated endpoint that requires the read:vasp permission.
func (s *Server) Overview(c *gin.Context) {
	var err error

	// Get the bff claims from the context
	var claims *auth.Claims
	if claims, err = auth.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("unable to retrieve bff claims from context")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Extract the VASP IDs from the claims
	testnetID := claims.VASPs.TestNet
	mainnetID := claims.VASPs.MainNet

	out := api.OverviewReply{
		OrgID:   claims.OrgID,
		TestNet: api.NetworkOverview{},
		MainNet: api.NetworkOverview{},
	}

	// Get the status for both testnet and mainnet
	var testnetStatus, mainnetStatus *gds.ServiceState
	if testnetStatus, mainnetStatus, err = s.GetStatuses(c); err != nil {
		log.Error().Err(err).Msg("unable to retrieve status information")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Populate the status responses
	if testnetStatus != nil {
		out.TestNet.Status = testnetStatus.Status.String()
	} else {
		out.TestNet.Status = gds.ServiceState_UNKNOWN.String()
	}

	if mainnetStatus != nil {
		out.MainNet.Status = mainnetStatus.Status.String()
	} else {
		out.MainNet.Status = gds.ServiceState_UNKNOWN.String()
	}

	// Get the summaries for both testnet and mainnet
	testnet, mainnet, testnetErr, mainnetErr := s.GetSummaries(c.Request.Context(), testnetID, mainnetID)
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve summary information")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Populate the summary responses
	if testnetErr != nil {
		out.Error.TestNet = testnetErr.Error()
	} else if testnet != nil {
		out.TestNet.Vasps = int(testnet.Vasps)
		out.TestNet.CertificatesIssued = int(testnet.CertificatesIssued)
		out.TestNet.NewMembers = int(testnet.NewMembers)

		if testnetID != "" {
			// Check if we received the VASP details from the testnet
			if testnet.MemberInfo == nil {
				log.Error().Msg("expected VASP details from testnet Summary RPC")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse(fmt.Errorf("could not retrieve testnet VASP details")))
				return
			}

			out.TestNet.MemberDetails = api.MemberDetails{
				ID:          testnet.MemberInfo.Id,
				Status:      testnet.MemberInfo.Status.String(),
				CountryCode: testnet.MemberInfo.Country,
			}
		}
	}

	if mainnetErr != nil {
		out.Error.MainNet = mainnetErr.Error()
	} else if mainnet != nil {
		out.MainNet.Vasps = int(mainnet.Vasps)
		out.MainNet.CertificatesIssued = int(mainnet.CertificatesIssued)
		out.MainNet.NewMembers = int(mainnet.NewMembers)

		if mainnetID != "" {
			// Check if we received the VASP details from the mainnet
			if mainnet.MemberInfo == nil {
				log.Error().Msg("expected VASP details from mainnet Summary RPC")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse(fmt.Errorf("could not retrieve mainnet VASP details")))
				return
			}

			out.MainNet.MemberDetails = api.MemberDetails{
				ID:          mainnet.MemberInfo.Id,
				Status:      mainnet.MemberInfo.Status.String(),
				CountryCode: mainnet.MemberInfo.Country,
			}
		} else {
			out.MainNet.MemberDetails = api.MemberDetails{
				Status: pb.VerificationState_NO_VERIFICATION.String(),
			}
		}
	}

	c.JSON(http.StatusOK, out)
}
