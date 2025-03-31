package bff

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/config"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// The default query parameters against the TRISAMembers gRPC API
const (
	DefaultMembersTimeout   = 25 * time.Second
	DefaultMembersPageSize  = 200
	DefaultMembersDirectory = "trisa.directory"
)

// GetSummaries makes parallel calls to the members service to get the summary
// information for both testnet and mainnet. If an endpoint returned an error, then a
// nil value is returned from this function for that endpoint instead of an error.
func (s *Server) GetSummaries(ctx context.Context, testnetID, mainnetID string) (testnetSummary, mainnetSummary *members.SummaryReply, testnetErr, mainnetErr error) {
	// Create the RPC which can do both testnet and mainnet calls
	rpc := func(ctx context.Context, client GlobalDirectoryClient, network string) (rep proto.Message, err error) {
		req := &members.SummaryRequest{}
		switch network {
		case config.TestNet:
			req.MemberId = testnetID
		case config.MainNet:
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
//
// @Summary Get summary information for the overview dashboard [read:vasp]
// @Description Returns a high level summary representing the state of each directory service and VASP registrations.
// @Tags overview
// @Produce json
// @Success 200 {object} api.OverviewReply
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /overview [get]
func (s *Server) Overview(c *gin.Context) {
	var err error

	// Get the bff claims from the context
	var claims *auth.Claims
	if claims, err = auth.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("unable to retrieve bff claims from context")
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
		sentry.Error(c).Err(err).Msg("unable to retrieve status information")
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
				sentry.Error(c).Msg("expected VASP details from testnet Summary RPC")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse(fmt.Errorf("could not retrieve testnet VASP details")))
				return
			}

			out.TestNet.MemberDetails = api.MemberDetails{
				ID:          testnet.MemberInfo.Id,
				Status:      testnet.MemberInfo.Status.String(),
				CountryCode: testnet.MemberInfo.Country,
				FirstListed: testnet.MemberInfo.FirstListed,
				VerifiedOn:  testnet.MemberInfo.VerifiedOn,
				LastUpdated: testnet.MemberInfo.LastUpdated,
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
				sentry.Error(c).Msg("expected VASP details from mainnet Summary RPC")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse(fmt.Errorf("could not retrieve mainnet VASP details")))
				return
			}

			out.MainNet.MemberDetails = api.MemberDetails{
				ID:          mainnet.MemberInfo.Id,
				Status:      mainnet.MemberInfo.Status.String(),
				CountryCode: mainnet.MemberInfo.Country,
				FirstListed: mainnet.MemberInfo.FirstListed,
				VerifiedOn:  mainnet.MemberInfo.VerifiedOn,
				LastUpdated: mainnet.MemberInfo.LastUpdated,
			}
		}
	}

	c.JSON(http.StatusOK, out)
}

// MemberList is an authenticated endpoint that returns a list of all verified VASPs in
// the requested directory (e.g. either TestNet or MainNet). This endpoint requires the
// read:vasp permission and is only available to organizations that have themselves been
// verified through the TRISA directory that they are querying.
//
// @Summary List verified VASPs in the specified directory [read:vasp].
// @Description Returns a list of verified VASPs in the specified directory so long as the organization is a verified member of that directory.
// @Tags members
// @Accept json
// @Produce json
// @Param params body api.MemberPageInfo true "Directory and Pagination"
// @Success 200 {object} object "VASP List"
// @Failure 400 {object} api.Reply "VASP ID and directory are required"
// @Failure 401 {object} api.Reply
// @Failure 404 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /members [get]
func (s *Server) MemberList(c *gin.Context) {
	var (
		err    error
		params *api.MemberPageInfo
		rep    *members.ListReply
	)

	// Bind the query parameters and add reasonable defaults
	params = &api.MemberPageInfo{}
	if err = c.ShouldBindQuery(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request with query params")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Add reasonable defaults to the params if they do not exist
	if params.PageSize <= 0 {
		params.PageSize = DefaultMembersPageSize
	}

	// By default query the mainnet if a directory is not specified
	if params.Directory == "" {
		params.Directory = DefaultMembersDirectory
	}

	// Validate the registered directory
	params.Directory = strings.ToLower(params.Directory)
	if !validRegisteredDirectory(params.Directory) {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("unknown registered directory"))
		return
	}

	// Check that the requester is verified with the directory they're trying to access.
	if verified, err := RequireVerification(c, params.Directory); err != nil {
		sentry.Error(c).Err(err).Msg("could not require verification for members endpoint")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve members list"))
		return
	} else if !verified {
		c.JSON(http.StatusUnavailableForLegalReasons, api.ErrorResponse("listing GDS members is only available to verified TRISA members"))
		return
	}

	// Execute the members list request against the specified GDS
	log.Debug().Str("registered_directory", params.Directory).Int32("page_size", params.PageSize).Str("page_token", params.PageToken).Msg("members list request")
	req := &members.ListRequest{
		PageSize:  params.PageSize,
		PageToken: params.PageToken,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), DefaultMembersTimeout)
	defer cancel()

	switch registeredDirectoryType(params.Directory) {
	case config.TestNet:
		rep, err = s.testnetGDS.List(ctx, req)
	case config.MainNet:
		rep, err = s.mainnetGDS.List(ctx, req)
	default:
		sentry.Error(c).Str("registered_directory", params.Directory).Msg("unhandled directory")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve member list"))
		return
	}

	// Handle gRPC errors
	if err != nil {
		if serr, ok := status.FromError(err); ok {
			sentry.Error(c).Err(err).Str("code", serr.Code().String()).Str("grpc_error", serr.Message()).Msg("members list rpc error")
			if serr.Code() == codes.Unavailable {
				c.JSON(http.StatusServiceUnavailable, api.ErrorResponse("specified directory is currently unavailable, please try again later"))
				return
			}

			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve member list"))
			return
		}

		sentry.Error(c).Err(err).Msg("unhandled error from gRPC service")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve member list"))
		return
	}

	// Create the list response
	c.JSON(http.StatusOK, &api.MemberListReply{
		VASPs:         rep.Vasps,
		NextPageToken: rep.NextPageToken,
	})
}

// MemberDetail endpoint is an authenticated endpoint that returns more detailed
// information about a verified VASP member from the specified directory (either
// TestNet or MainNet). This endpoint requires the read:vasp permission and is only
// available to organizations that have themselves been verified through the TRISA
// directory that they are querying.
//
// @Summary Get details for a VASP in the specified directory [read:vasp]
// @Description Returns details for a VASP by ID and directory so long as the organization is a verified member of that directory.
// @Tags members
// @Accept json
// @Produce json
// @Param params body api.MemberDetailsParams true "VASP ID and directory"
// @Success 200 {object} object "VASP Details"
// @Failure 400 {object} api.Reply "VASP ID and directory are required"
// @Failure 401 {object} api.Reply
// @Failure 404 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /members/{id} [get]
func (s *Server) MemberDetail(c *gin.Context) {
	// Bind the parameters associated with the MemberDetails request
	params := &api.MemberDetailsParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request with query params")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Overwrite the params ID with the vaspID from the URL
	params.ID = c.Param("vaspID")

	// By default, query the mainnet if a directory is not specified
	if params.Directory == "" {
		params.Directory = DefaultMembersDirectory
	}

	// Validate the registered directory
	params.Directory = strings.ToLower(params.Directory)
	if !validRegisteredDirectory(params.Directory) {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("unknown registered directory"))
		return
	}

	// Check that the requester is verified with the directory they're trying to access.
	if verified, err := RequireVerification(c, params.Directory); err != nil {
		sentry.Error(c).Err(err).Msg("could not require verification for members endpoint")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve members list"))
		return
	} else if !verified {
		c.JSON(http.StatusUnavailableForLegalReasons, api.ErrorResponse("listing GDS members is only available to verified TRISA members"))
		return
	}

	// Do the members request
	log.Debug().Str("registered_directory", params.Directory).Msg("issuing members detail request")
	req := &members.DetailsRequest{
		MemberId: params.ID,
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), DefaultMembersTimeout)
	defer cancel()

	var (
		err error
		rep *members.MemberDetails
	)

	switch registeredDirectoryType(params.Directory) {
	case config.TestNet:
		rep, err = s.testnetGDS.Details(ctx, req)
	case config.MainNet:
		rep, err = s.mainnetGDS.Details(ctx, req)
	default:
		sentry.Error(c).Str("registered_directory", params.Directory).Msg("unknown directory")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve member details"))
		return
	}

	// Handle errors from the members endpoint
	if err != nil {
		serr, _ := status.FromError(err)
		switch serr.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, api.ErrorResponse(serr.Message()))
		default:
			sentry.Error(c).Err(err).Str("code", serr.Code().String()).Str("registered_directory", params.Directory).Msg("could not retrieve member details")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(serr.Message()))
		}
		return
	}

	// Create the member details response, rewiring the protobuf types for a better
	// JSON representation.
	out := api.MemberDetailsReply{}

	// Marshal the VASP summary details
	if rep.MemberSummary == nil {
		sentry.Error(c).Msg("did not receive summary details from members detail RPC")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve member details"))
		return
	}

	if out.Summary, err = wire.Rewire(rep.MemberSummary); err != nil {
		sentry.Error(c).Err(err).Msg("could not serialize summary details")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Marshal the legal person details
	if rep.LegalPerson == nil {
		sentry.Error(c).Msg("did not receive legal person details from members detail RPC")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve member details"))
		return
	}

	if out.LegalPerson, err = wire.Rewire(rep.LegalPerson); err != nil {
		sentry.Error(c).Err(err).Msg("could not serialize legal person details")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Marshal the contacts details
	if rep.Contacts == nil {
		sentry.Error(c).Msg("did not receive contacts details from members detail RPC")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve member details"))
		return
	}

	if out.Contacts, err = wire.Rewire(rep.Contacts); err != nil {
		sentry.Error(c).Err(err).Msg("could not serialize contacts details")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Marshal the Trixo form details
	if rep.Trixo == nil {
		sentry.Error(c).Msg("did not receive trixo form details from members detail RPC")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve member details"))
		return
	}

	if out.Trixo, err = wire.Rewire(rep.Trixo); err != nil {
		sentry.Error(c).Err(err).Msg("could not serialize trixo form details")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, out)
}
