package bff

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/config"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

// ConnectMembers creates a gRPC client to the TRISA Members Service specified in the
// configuration. This method is used to connect to both the TestNet and the MainNet and
// to connect to mock GDS services in testing using buffconn.
func ConnectMembers(conf config.DirectoryConfig) (_ members.TRISAMembersClient, err error) {
	// Create the Dial options with required credentials
	var opts []grpc.DialOption
	if conf.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	// Connect the directory client (non-blocking)
	var cc *grpc.ClientConn
	if cc, err = grpc.DialContext(ctx, conf.Endpoint, opts...); err != nil {
		return nil, err
	}
	return members.NewTRISAMembersClient(cc), nil
}

// GetSummaries makes parallel calls to the members service to get the summary
// information for both testnet and mainnet. If an endpoint returned an error, then a
// nil value is returned from this function for that endpoint instead of an error.
func (s *Server) GetSummaries(ctx context.Context, testnetID, mainnetID string) (testsum *members.SummaryReply, mainsum *members.SummaryReply, err error) {
	// Create the RPC which can do both testnet and mainnet calls
	rpc := func(ctx context.Context, client members.TRISAMembersClient, network string) (rep proto.Message, err error) {
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
	results, errs := s.ParallelMembersRequests(ctx, rpc, false)
	if len(errs) != 2 || len(results) != 2 {
		return nil, nil, fmt.Errorf("unexpected number of results from parallel requests: %d", len(results))
	}

	// Parse the results
	var ok bool
	if errs[0] != nil {
		testsum = nil
	} else if results[0] == nil {
		return nil, nil, fmt.Errorf("nil testnet result returned from parallel requests")
	} else if testsum, ok = results[0].(*members.SummaryReply); !ok {
		return nil, nil, fmt.Errorf("unexpected testnet status result type returned from parallel requests: %T", results[0])
	}

	if errs[1] != nil {
		mainsum = nil
	} else if results[1] == nil {
		return nil, nil, fmt.Errorf("nil mainnet status result returned from parallel requests")
	} else if mainsum, ok = results[1].(*members.SummaryReply); !ok {
		return nil, nil, fmt.Errorf("unexpected mainnet status result type returned from parallel requests: %T", results[1])
	}

	return testsum, mainsum, nil
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
	testnetID := claims.VASP[testnet]
	mainnetID := claims.VASP[mainnet]

	out := &api.OverviewReply{
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
	var testnet, mainnet *members.SummaryReply
	testnet, mainnet, err = s.GetSummaries(c.Request.Context(), testnetID, mainnetID)
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve summary information")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Populate the summary responses
	if testnet != nil {
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

	if mainnet != nil {
		out.MainNet.Vasps = int(mainnet.Vasps)
		out.MainNet.CertificatesIssued = int(mainnet.CertificatesIssued)
		out.MainNet.NewMembers = int(mainnet.NewMembers)

		if mainnetID != "" {
			// Check if we received the VASP details from the mainnet
			if mainnet.MemberInfo == nil {
				log.Error().Msg("could not retrieve mainnet VASP details")
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
