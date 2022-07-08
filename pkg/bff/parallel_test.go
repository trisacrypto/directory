package bff_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func (s *bffTestSuite) TestParallelAdminRequests() {
	var (
		results []interface{}
		errs    []error
	)

	require := s.Require()

	// RPC that returns a status reply for both networks
	rpc := func(ctx context.Context, client admin.DirectoryAdministrationClient, network string) (rep interface{}, err error) {
		if rep, err = client.Status(ctx); err != nil {
			return nil, err
		}
		return rep, nil
	}

	// Handler that returns a healthy status reply
	healthy := func(c *gin.Context) {
		c.JSON(http.StatusOK, &admin.StatusReply{
			Status: "healthy",
		})
	}

	// Test the case where the RPC returns two errors and flatten is true
	require.NoError(s.testnet.admin.UseError(mock.StatusEP, http.StatusInternalServerError, "internal error"))
	require.NoError(s.mainnet.admin.UseError(mock.StatusEP, http.StatusInternalServerError, "internal error"))
	results, errs = s.bff.ParallelAdminRequests(context.TODO(), rpc, true)
	require.Len(results, 0, "results was not flattened")
	require.Len(errs, 2, "errors were not returned")
	require.NotNil(errs[0], "expected testnet error to be not nil")
	require.NotNil(errs[1], "expected mainnet error to be not nil")

	// Test the case where the RPC returns 2 results and flatten is true
	require.NoError(s.testnet.admin.UseHandler(mock.StatusEP, healthy))
	require.NoError(s.mainnet.admin.UseHandler(mock.StatusEP, healthy))
	results, errs = s.bff.ParallelAdminRequests(context.TODO(), rpc, true)
	require.Len(results, 2, "results was not flattened")
	require.Len(errs, 0, "errors were not flattened")
	require.NotNil(results[0], "testnet result expected")
	require.NotNil(results[1], "mainnet result expected")

	// Test the case where the RPC returns 2 results and flatten is false
	results, errs = s.bff.ParallelAdminRequests(context.TODO(), rpc, false)
	require.Len(results, 2, "results was flattened")
	require.Len(errs, 2, "errors were flattened")
	require.NotNil(results[0], "expected testnet result to be not nil")
	require.Nil(errs[0], "expected testnet error to be nil")
	require.NotNil(results[1], "expected mainnet result to be not nil")
	require.Nil(errs[1], "expected mainnet error to be nil")
}

func (s *bffTestSuite) TestParallelGDSRequests() {
	var (
		results []interface{}
		errs    []error
	)

	// Setup the test to execute requests against the status endpoint
	require := s.Require()
	rpc := func(ctx context.Context, client GlobalDirectoryClient, network string) (rep proto.Message, err error) {
		// NOTE: for the tests to pass, this must return nil, err and rep, nil instead
		// of directly returning the results from client.Status(). That's because the
		// rep nil will be (*api.ServiceState)(nil) not (protoreflect.Message)(nil) and
		// the non-interface type will not be flattened without a more extensive type
		// check or the use of reflection.
		if rep, err = client.Status(ctx, &gds.HealthCheck{}); err != nil {
			return nil, err
		}
		return rep, nil
	}

	// Test the case where the RPC returns two errors and flatten is true
	s.testnet.gds.UseError(mock.StatusRPC, codes.Unavailable, "nobody is home")
	s.mainnet.gds.UseError(mock.StatusRPC, codes.Unavailable, "nobody is home")

	results, errs = s.bff.ParallelGDSRequests(context.TODO(), rpc, true)
	require.Len(results, 0, "results was not flattened")
	require.Len(errs, 2, "errors were not returned")
	require.NotNil(errs[0], "expected testnet error to be not nil")
	require.NotNil(errs[1], "expected mainnet error to be not nil")

	// Test the case where the RPC returns two errors and flatten is false
	results, errs = s.bff.ParallelGDSRequests(context.TODO(), rpc, false)
	require.Len(results, 2, "results were flattened")
	require.Len(errs, 2, "errors were not returned")
	require.Nil(results[0], "expected testnet result to be nil")
	require.Nil(results[1], "expected mainnet result to be nil")
	require.NotNil(errs[0], "expected testnet error to be not nil")
	require.NotNil(errs[1], "expected mainnet error to be not nil")

	// Test the case where the RPC returns 1 error 1 result and flatten is true
	s.mainnet.gds.OnStatus = func(ctx context.Context, in *gds.HealthCheck) (out *gds.ServiceState, err error) {
		return &gds.ServiceState{
			Status: gds.ServiceState_HEALTHY,
		}, nil
	}

	results, errs = s.bff.ParallelGDSRequests(context.TODO(), rpc, true)
	require.Len(results, 1, "results was not flattened")
	require.Len(errs, 1, "errors were not flattened")
	require.NotNil(results[0], "result expected")
	require.NotNil(errs[0], "err also expected")

	// Test the case where the RPC returns 1 error 1 result and flatten is false
	results, errs = s.bff.ParallelGDSRequests(context.TODO(), rpc, false)
	require.Len(results, 2, "results was flattened")
	require.Len(errs, 2, "errors were flattened")
	require.Nil(results[0], "expected testnet result to be nil")
	require.NotNil(errs[0], "expected testnet error to be not nil")
	require.NotNil(results[1], "expected mainnet result to be not nil")
	require.Nil(errs[1], "expected mainnet error to be nil")

	// Test the case where the RPC returns 2 results and flatten is true
	s.testnet.gds.OnStatus = func(ctx context.Context, in *gds.HealthCheck) (out *gds.ServiceState, err error) {
		return &gds.ServiceState{
			Status: gds.ServiceState_DANGER,
		}, nil
	}

	results, errs = s.bff.ParallelGDSRequests(context.TODO(), rpc, true)
	require.Len(results, 2, "results was not flattened")
	require.Len(errs, 0, "errors were not flattened")
	require.NotNil(results[0], "testnet result expected")
	require.NotNil(results[1], "mainnet result expected")

	// Test the case where the RPC returns 2 results and flatten is false
	results, errs = s.bff.ParallelGDSRequests(context.TODO(), rpc, false)
	require.Len(results, 2, "results was flattened")
	require.Len(errs, 2, "errors were flattened")
	require.NotNil(results[0], "expected testnet result to be not nil")
	require.Nil(errs[0], "expected testnet error to be nil")
	require.NotNil(results[1], "expected mainnet result to be not nil")
	require.Nil(errs[1], "expected mainnet error to be nil")

	// Test the case where the RPC returns 1 error 1 result (but testnet this time) and flatten is true
	s.mainnet.gds.UseError(mock.StatusRPC, codes.Unavailable, "nobody is home")

	results, errs = s.bff.ParallelGDSRequests(context.TODO(), rpc, true)
	require.Len(results, 1, "results was not flattened")
	require.Len(errs, 1, "errors were not flattened")
	require.NotNil(results[0], "result expected")
	require.NotNil(errs[0], "err also expected")

	// Test the case where the RPC returns 1 error 1 result (but testnet this time) and flatten is false
	results, errs = s.bff.ParallelGDSRequests(context.TODO(), rpc, false)
	require.Len(results, 2, "results was flattened")
	require.Len(errs, 2, "errors were flattened")
	require.NotNil(results[0], "expected testnet result to be not nil")
	require.Nil(errs[0], "expected testnet error to be nil")
	require.Nil(results[1], "expected mainnet result to be nil")
	require.NotNil(errs[1], "expected mainnet error to be not nil")
}

func TestFlatten(t *testing.T) {
	repProto := &gds.ServiceState{Status: gds.ServiceState_HEALTHY}
	err := errors.New("something bad happened")

	testCases := []struct {
		results     []interface{}
		errs        []error
		expectedLen int
	}{
		{
			results:     []interface{}{nil, nil},
			errs:        []error{nil, nil},
			expectedLen: 0,
		},
		{
			results:     []interface{}{repProto, nil},
			errs:        []error{err, nil},
			expectedLen: 1,
		},
		{
			results:     []interface{}{nil, repProto},
			errs:        []error{nil, err},
			expectedLen: 1,
		},
		{
			results:     []interface{}{repProto, repProto},
			errs:        []error{err, err},
			expectedLen: 2,
		},
	}

	for idx, tc := range testCases {
		protoResults := FlattenResults(tc.results)
		require.Len(t, protoResults, tc.expectedLen, "unexpected length in test case %d", idx)

		errs := FlattenErrs(tc.errs)
		require.Len(t, errs, tc.expectedLen, "unexpected length in test case %d", idx)
	}
}
