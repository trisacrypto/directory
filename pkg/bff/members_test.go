package bff_test

import (
	"context"
	"errors"

	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc/codes"
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
	testnet, mainnet, testnetErr, mainnetErr := s.bff.GetSummaries(context.TODO(), expectTestnet.MemberInfo.Id, expectMainnet.MemberInfo.Id)
	require.NoError(testnetErr, "could not get testnet summary")
	require.NoError(mainnetErr, "could not get mainnet summary")
	require.True(proto.Equal(expectTestnet, testnet), "testnet summaries did not match")
	require.True(proto.Equal(expectMainnet, mainnet), "mainnet summaries did not match")

	// Test only testnet summary was returned
	s.mainnet.members.OnSummary = errorSummary
	testnet, mainnet, testnetErr, mainnetErr = s.bff.GetSummaries(context.TODO(), expectTestnet.MemberInfo.Id, expectMainnet.MemberInfo.Id)
	require.NoError(testnetErr, "could not get testnet summary")
	require.True(proto.Equal(expectTestnet, testnet), "testnet summaries did not match")
	require.Nil(mainnet, "mainnet summary should be nil")

	// Test only mainnet summary was returned
	s.testnet.members.OnSummary = errorSummary
	s.mainnet.members.OnSummary = mainnetSummary
	testnet, mainnet, testnetErr, mainnetErr = s.bff.GetSummaries(context.TODO(), expectTestnet.MemberInfo.Id, expectMainnet.MemberInfo.Id)
	require.NoError(mainnetErr, "could not get mainnet summary")
	require.True(proto.Equal(expectMainnet, mainnet), "mainnet summaries did not match")
	require.Nil(testnet, "testnet summary should be nil")

	// Test both summaries were not returned
	s.mainnet.members.OnSummary = errorSummary
	testnet, mainnet, testnetErr, mainnetErr = s.bff.GetSummaries(context.TODO(), expectTestnet.MemberInfo.Id, expectMainnet.MemberInfo.Id)
	require.Error(testnetErr, "testnet error should have been returned")
	require.Error(mainnetErr, "mainnet error should have been returned")
	require.Nil(testnet, "testnet summary should be nil")
	require.Nil(mainnet, "mainnet summary should be nil")
}

func (s *bffTestSuite) TestOverview() {
	require := s.Require()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
		OrgID:       "a2c4f8f0-f8f8-4f8f-8f8f-8f8f8f8f8f8f",
		VASPs:       map[string]string{},
	}

	// Endpoint must be authenticated
	_, err := s.client.Overview(context.TODO())
	require.EqualError(err, "[401] this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.Certificates(context.TODO())
	require.EqualError(err, "[401] user does not have permission to perform this operation", "expected error when user is not authorized")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{"read:vasp"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with correct permissions")

	// If all endpoints return an error, we should still return a response
	require.NoError(s.testnet.gds.UseError(mock.StatusRPC, codes.Unavailable, "testnet is unavailable"))
	require.NoError(s.mainnet.gds.UseError(mock.StatusRPC, codes.Unavailable, "mainnet is unavailable"))
	require.NoError(s.testnet.members.UseError(mock.SummaryRPC, codes.Unavailable, "testnet is unavailable"))
	require.NoError(s.mainnet.members.UseError(mock.SummaryRPC, codes.Unavailable, "mainnet is unavailable"))
	expected := &api.OverviewReply{
		Error: api.NetworkError{
			TestNet: "rpc error: code = Unavailable desc = testnet is unavailable",
			MainNet: "rpc error: code = Unavailable desc = mainnet is unavailable",
		},
		OrgID: claims.OrgID,
		TestNet: api.NetworkOverview{
			Status:             gds.ServiceState_UNKNOWN.String(),
			Vasps:              0,
			CertificatesIssued: 0,
			NewMembers:         0,
		},
		MainNet: api.NetworkOverview{
			Status:             gds.ServiceState_UNKNOWN.String(),
			Vasps:              0,
			CertificatesIssued: 0,
			NewMembers:         0,
		},
	}
	reply, err := s.client.Overview(context.TODO())
	require.NoError(err, "could not get overview")
	require.Equal(expected, reply, "overview reply did not match")

	// Test with a valid status response from one of the endpoints
	s.testnet.gds.OnStatus = func(ctx context.Context, in *gds.HealthCheck) (*gds.ServiceState, error) {
		return &gds.ServiceState{
			Status: gds.ServiceState_HEALTHY,
		}, nil
	}
	expected.TestNet.Status = gds.ServiceState_HEALTHY.String()
	reply, err = s.client.Overview(context.TODO())
	require.NoError(err, "could not get overview")
	require.Equal(expected, reply, "overview reply did not match")

	// Test with a valid summary response from one of the endpoints
	s.testnet.members.OnSummary = func(ctx context.Context, in *members.SummaryRequest) (*members.SummaryReply, error) {
		return &members.SummaryReply{
			Vasps:              5,
			CertificatesIssued: 6,
			NewMembers:         3,
		}, nil
	}
	expected.TestNet.Vasps = 5
	expected.TestNet.CertificatesIssued = 6
	expected.TestNet.NewMembers = 3
	expected.Error.TestNet = ""
	reply, err = s.client.Overview(context.TODO())
	require.NoError(err, "could not get overview")
	require.Equal(expected, reply, "overview reply did not match")

	// Test with both valid responses, one endpoint returns VASP details
	claims.VASPs["mainnet"] = "b2c4f8f0-f8f8-4f8f-8f8f-8f8f8f8f8f8f"
	require.NoError(s.SetClientCredentials(claims), "could not create token with VASP ID")
	s.mainnet.gds.OnStatus = func(ctx context.Context, in *gds.HealthCheck) (*gds.ServiceState, error) {
		return &gds.ServiceState{
			Status: gds.ServiceState_HEALTHY,
		}, nil
	}
	s.mainnet.members.OnSummary = func(ctx context.Context, in *members.SummaryRequest) (*members.SummaryReply, error) {
		return &members.SummaryReply{
			Vasps:              20,
			CertificatesIssued: 23,
			NewMembers:         5,
			MemberInfo: &members.VASPMember{
				Id:      claims.VASPs["mainnet"],
				Status:  models.VerificationState_SUBMITTED,
				Country: "US",
			},
		}, nil
	}
	expected.MainNet.Status = gds.ServiceState_HEALTHY.String()
	expected.MainNet.Vasps = 20
	expected.MainNet.CertificatesIssued = 23
	expected.MainNet.NewMembers = 5
	expected.MainNet.MemberDetails = api.MemberDetails{
		ID:          claims.VASPs["mainnet"],
		Status:      models.VerificationState_SUBMITTED.String(),
		CountryCode: "US",
	}
	expected.Error.MainNet = ""
	reply, err = s.client.Overview(context.TODO())
	require.NoError(err, "could not get overview")
	require.Equal(expected, reply, "overview reply did not match")

	// Test with both endpoints returning VASP details
	claims.VASPs["testnet"] = "c4f8f8f0-f8f8-4f8f-8f8f-8f8f8f8f8f8f"
	require.NoError(s.SetClientCredentials(claims), "could not create token with VASP ID")
	s.testnet.members.OnSummary = func(ctx context.Context, in *members.SummaryRequest) (*members.SummaryReply, error) {
		return &members.SummaryReply{
			Vasps:              5,
			CertificatesIssued: 6,
			NewMembers:         3,
			MemberInfo: &members.VASPMember{
				Id:      claims.VASPs["testnet"],
				Status:  models.VerificationState_VERIFIED,
				Country: "FR",
			},
		}, nil
	}
	expected.TestNet.MemberDetails = api.MemberDetails{
		ID:          claims.VASPs["testnet"],
		Status:      models.VerificationState_VERIFIED.String(),
		CountryCode: "FR",
	}
	reply, err = s.client.Overview(context.TODO())
	require.NoError(err, "could not get overview")
	require.Equal(expected, reply, "overview reply did not match")
}
