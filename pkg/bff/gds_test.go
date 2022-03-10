package bff_test

import (
	"context"

	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	"google.golang.org/grpc/codes"
)

func (s *bffTestSuite) TestLookup() {
	require := s.Require()
	params := &api.LookupParams{}

	// Test Bad Request (no parameters)
	_, err := s.client.Lookup(context.TODO(), params)
	require.EqualError(err, "[400] must provide either uuid or common_name in query params", "expected a 400 error with no params")

	// Provide some params
	params.CommonName = "api.alice.vaspbot.net"

	// Test NotFound returns a 404
	require.NoError(s.testnet.UseError(mock.LookupRPC, codes.NotFound, "not found"))
	require.NoError(s.mainnet.UseError(mock.LookupRPC, codes.NotFound, "not found"))
	_, err = s.client.Lookup(context.TODO(), params)
	require.EqualError(err, "[404] no results returned for query", "expected a 404 error when both GDS return not found")

	// Test InternalError when both GDS return Unavailable
	require.NoError(s.testnet.UseError(mock.LookupRPC, codes.Unavailable, "cannot connect"))
	require.NoError(s.mainnet.UseError(mock.LookupRPC, codes.Unavailable, "cannot connect"))
	_, err = s.client.Lookup(context.TODO(), params)
	require.EqualError(err, "[500] unable to execute Lookup request", "expected a 500 error when both GDS return unavailable")

	// Test one result from TestNet
	require.NoError(s.testnet.UseFixture(mock.LookupRPC, "testdata/testnet/lookup_reply.json"))
	require.NoError(s.mainnet.UseError(mock.LookupRPC, codes.NotFound, "not found"))
	rep, err := s.client.Lookup(context.TODO(), params)
	require.NoError(err, "could not fetch expected result from testnet")
	require.Len(rep.Results, 1, "expected one result back from server")
	require.Equal("6a57fea4-8fb7-42f3-bf0c-55fecccd2e53", rep.Results[0]["id"])

	// Test one result from MainNet
	require.NoError(s.testnet.UseError(mock.LookupRPC, codes.NotFound, "not found"))
	require.NoError(s.mainnet.UseFixture(mock.LookupRPC, "testdata/mainnet/lookup_reply.json"))
	rep, err = s.client.Lookup(context.TODO(), params)
	require.NoError(err, "could not fetch expected result from mainnet")
	require.Len(rep.Results, 1, "expected one result back from server")
	require.Equal("ca0cff66-719f-4a62-8086-be953699b27d", rep.Results[0]["id"])

	// Test results from both TestNet and MainNet
	require.NoError(s.testnet.UseFixture(mock.LookupRPC, "testdata/testnet/lookup_reply.json"))
	require.NoError(s.mainnet.UseFixture(mock.LookupRPC, "testdata/mainnet/lookup_reply.json"))
	rep, err = s.client.Lookup(context.TODO(), params)
	require.NoError(err, "could not fetch expected result from mainnet and testnet")
	require.Len(rep.Results, 2, "expected two results back from server")
}
