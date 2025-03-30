package bff

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const ContextVerificationStatus = "verification_status"

// CheckVerification is a middleware that discovers whether or not the VASP organization
// making the request is verified in the GDS. It queries both the mainnet and testnet
// GDS in parallel and adds context keys about the verification status of the VASP.
// NOTE: the user must be authenticated before this middleware is executed.
func (s *Server) CheckVerification(c *gin.Context) {
	var (
		err    error
		claims *auth.Claims
	)

	if claims, err = auth.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("cannot check verification status on unauthenticated user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse(auth.ErrNoAuthUser))
		return
	}

	verification := &VerificationStatus{
		MainNet: &VASPVerificationStatus{Status: "NO_VERIFICATION", Verified: false},
		TestNet: &VASPVerificationStatus{Status: "NO_VERIFICATION", Verified: false},
	}

	// Create an RPC closure for making parallel GDS requests
	verify := func(ctx context.Context, client GlobalDirectoryClient, network string) (_ proto.Message, err error) {
		// Create the request from the closure
		req := &gds.VerificationRequest{}
		switch network {
		case config.TestNet:
			req.Id = claims.VASPs.TestNet
			req.RegisteredDirectory = "testnet.directory"
		case config.MainNet:
			req.Id = claims.VASPs.MainNet
			req.RegisteredDirectory = "trisa.directory"
		default:
			return nil, fmt.Errorf("unknown network %q", network)
		}

		// If there is no ID then return an empty verification reply.
		// E.g. if the user has not yet submitted a registration for the specified network.
		if req.Id == "" {
			return &gds.VerificationReply{VerificationStatus: models.VerificationState_NO_VERIFICATION}, nil
		}

		var rep *gds.VerificationReply
		if rep, err = client.Verification(ctx, req); err != nil {
			// If the code is not found or unavailable, do not return an error; just no result
			serr, _ := status.FromError(err)
			switch serr.Code() {
			case codes.NotFound:
				sentry.Warn(ctx).Err(err).Str("network", network).Msg("verification lookup for unknown VASP claims are out of date")
				return nil, nil
			default:
				sentry.Error(ctx).Err(err).Str("network", network).Msg("unsuccessful verification lookup")
				return nil, err
			}
		}
		return rep, nil
	}

	// Execute the parallel GDS verification requests, ensuring flatten is false with
	// the expectation that TestNet will be 0 index and MainNet the 1 index.
	results, _ := s.ParallelGDSRequests(sentry.RequestContext(c), verify, false)

	// Handle TestNet Results
	if results[0] != nil {
		state := results[0].(*gds.VerificationReply)
		verification.TestNet.Status = state.VerificationStatus.String()
		verification.TestNet.Verified = state.VerificationStatus == models.VerificationState_VERIFIED
	}

	// Handle MainNet Results
	if results[1] != nil {
		state := results[1].(*gds.VerificationReply)
		verification.MainNet.Status = state.VerificationStatus.String()
		verification.MainNet.Verified = state.VerificationStatus == models.VerificationState_VERIFIED
	}

	// Add verification status to the context
	c.Set(ContextVerificationStatus, verification)
	c.Next()
}

type VerificationStatus struct {
	MainNet *VASPVerificationStatus
	TestNet *VASPVerificationStatus
}

type VASPVerificationStatus struct {
	Status   string
	Verified bool
}

// A helper function to require verification to the specified network.
func RequireVerification(c *gin.Context, network string) (verified bool, err error) {
	var status *VerificationStatus
	if status, err = GetVerificationStatus(c); err != nil {
		return false, err
	}

	switch registeredDirectoryType(network) {
	case config.TestNet:
		if status.TestNet != nil {
			return status.TestNet.Verified, nil
		}
	case config.MainNet:
		if status.MainNet != nil {
			return status.MainNet.Verified, nil
		}
	default:
		return false, fmt.Errorf("unhandled directory type %q", network)
	}

	return false, ErrNoVerificationStatus
}

// A helper function to quickly retrieve the verification status from the context;
// return an error if the verification status does not exist. Panics if the status is
// not the correct type, e.g. not set by the CheckVerification middleware.
func GetVerificationStatus(c *gin.Context) (*VerificationStatus, error) {
	status, exists := c.Get(ContextVerificationStatus)
	if !exists {
		return nil, ErrNoVerificationStatus
	}
	val := status.(*VerificationStatus)
	return val, nil
}
