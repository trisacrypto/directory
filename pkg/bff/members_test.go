package bff_test

import (
	"context"
	"errors"

	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

func (s *bffTestSuite) TestGetSummaries() {
	require := s.Require()

	// Set the Summary RPC for the mocks
	expectTestnet := &members.SummaryReply{
		Vasps:              10,
		CertificatesIssued: 9,
		NewMembers:         3,
		MemberInfo: &members.VASPMember{
			Id:                  "a2c4f8f0-f8f8-4f8f-8f8f-8f8f8f8f8f8f",
			RegisteredDirectory: "testnet",
			CommonName:          "alice.vaspbot.net",
			Status:              pb.VerificationState_VERIFIED,
		},
	}
	testnetSummary := func(ctx context.Context, in *members.SummaryRequest) (*members.SummaryReply, error) {
		return expectTestnet, nil
	}

	expectMainnet := &members.SummaryReply{
		Vasps:              30,
		CertificatesIssued: 32,
		NewMembers:         5,
		MemberInfo: &members.VASPMember{
			Id:                  "b2c4f8f0-f8f8-4f8f-8f8f-8f8f8f8f8f8f",
			RegisteredDirectory: "mainnet",
			CommonName:          "alice.vaspbot.net",
			Status:              pb.VerificationState_SUBMITTED,
		},
	}
	mainnetSummary := func(ctx context.Context, in *members.SummaryRequest) (*members.SummaryReply, error) {
		return expectMainnet, nil
	}

	errorSummary := func(ctx context.Context, in *members.SummaryRequest) (*members.SummaryReply, error) {
		return nil, errors.New("unreachable host")
	}

	s.testnet.members.OnSummary = testnetSummary
	s.mainnet.members.OnSummary = mainnetSummary

	// Test both summaries were returned
	testnet, mainnet, err := s.bff.GetSummaries(context.TODO(), expectTestnet.MemberInfo.Id, expectMainnet.MemberInfo.Id)
	require.NoError(err, "could not get summaries")
	require.True(proto.Equal(expectTestnet, testnet), "testnet summaries did not match")
	require.True(proto.Equal(expectMainnet, mainnet), "mainnet summaries did not match")

	// Test only testnet summary was returned
	s.mainnet.members.OnSummary = errorSummary
	testnet, mainnet, err = s.bff.GetSummaries(context.TODO(), expectTestnet.MemberInfo.Id, expectMainnet.MemberInfo.Id)
	require.NoError(err, "could not get summaries")
	require.True(proto.Equal(expectTestnet, testnet), "testnet summaries did not match")
	require.Nil(mainnet, "mainnet summary should be nil")

	// Test only mainnet summary was returned
	s.testnet.members.OnSummary = errorSummary
	s.mainnet.members.OnSummary = mainnetSummary
	testnet, mainnet, err = s.bff.GetSummaries(context.TODO(), expectTestnet.MemberInfo.Id, expectMainnet.MemberInfo.Id)
	require.NoError(err, "could not get summaries")
	require.Nil(testnet, "testnet summary should be nil")
	require.True(proto.Equal(expectMainnet, mainnet), "mainnet summaries did not match")

	// Test both summaries were not returned
	s.mainnet.members.OnSummary = errorSummary
	testnet, mainnet, err = s.bff.GetSummaries(context.TODO(), expectTestnet.MemberInfo.Id, expectMainnet.MemberInfo.Id)
	require.NoError(err, "could not get summaries")
	require.Nil(testnet, "testnet summary should be nil")
	require.Nil(mainnet, "mainnet summary should be nil")
}

func (s *bffTestSuite) TestOverview() {
	// TODO: need to mock authentication before these tests will work.
	s.T().Skip("not implemented yet")

	// Test 401 with no access token
	// Test 401 authenticated user without read:vasp permission
	// Test 200 response with authenticated user with read:vasp permission
}
