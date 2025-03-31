package bff_test

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"time"

	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
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
			FirstListed:         time.Now().AddDate(0, 0, -2).Format(time.RFC3339),
			VerifiedOn:          time.Now().AddDate(0, 0, -1).Format(time.RFC3339),
			LastUpdated:         time.Now().Format(time.RFC3339),
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
			FirstListed:         time.Now().AddDate(0, 0, -3).Format(time.RFC3339),
			VerifiedOn:          time.Now().AddDate(0, 0, -2).Format(time.RFC3339),
			LastUpdated:         time.Now().Format(time.RFC3339),
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
	require.Error(mainnetErr, "expected mainnet error")
	require.True(proto.Equal(expectTestnet, testnet), "testnet summaries did not match")
	require.Nil(mainnet, "mainnet summary should be nil")

	// Test only mainnet summary was returned
	s.testnet.members.OnSummary = errorSummary
	s.mainnet.members.OnSummary = mainnetSummary
	testnet, mainnet, testnetErr, mainnetErr = s.bff.GetSummaries(context.TODO(), expectTestnet.MemberInfo.Id, expectMainnet.MemberInfo.Id)
	require.Error(testnetErr, "expected testnet error")
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
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.Overview(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{auth.ReadVASP}
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

	// Test with valid responses from both endpoints
	firstListed := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	verifiedOn := time.Now().AddDate(0, 0, -2).Format(time.RFC3339)
	lastUpdated := time.Now().Format(time.RFC3339)
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
				Id:          claims.VASPs["mainnet"],
				Status:      pb.VerificationState_SUBMITTED,
				Country:     "US",
				FirstListed: firstListed,
				VerifiedOn:  verifiedOn,
				LastUpdated: lastUpdated,
			},
		}, nil
	}
	expected.MainNet.Status = gds.ServiceState_HEALTHY.String()
	expected.MainNet.Vasps = 20
	expected.MainNet.CertificatesIssued = 23
	expected.MainNet.NewMembers = 5
	expected.Error.MainNet = ""
	reply, err = s.client.Overview(context.TODO())
	require.NoError(err, "could not get overview")
	require.Equal(expected, reply, "overview reply did not match")

	// Test that one of the endpoints returns VASP details if the ID is in the claims
	claims.VASPs["mainnet"] = "b2c4f8f0-f8f8-4f8f-8f8f-8f8f8f8f8f8f"
	require.NoError(s.SetClientCredentials(claims), "could not create token with VASP ID")
	expected.MainNet.MemberDetails = api.MemberDetails{
		ID:          claims.VASPs["mainnet"],
		Status:      pb.VerificationState_SUBMITTED.String(),
		CountryCode: "US",
		FirstListed: firstListed,
		VerifiedOn:  verifiedOn,
		LastUpdated: lastUpdated,
	}
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
				Id:          claims.VASPs["testnet"],
				Status:      pb.VerificationState_VERIFIED,
				Country:     "FR",
				FirstListed: firstListed,
				VerifiedOn:  verifiedOn,
				LastUpdated: lastUpdated,
			},
		}, nil
	}
	expected.TestNet.MemberDetails = api.MemberDetails{
		ID:          claims.VASPs["testnet"],
		Status:      pb.VerificationState_VERIFIED.String(),
		CountryCode: "FR",
		FirstListed: firstListed,
		VerifiedOn:  verifiedOn,
		LastUpdated: lastUpdated,
	}
	reply, err = s.client.Overview(context.TODO())
	require.NoError(err, "could not get overview")
	require.Equal(expected, reply, "overview reply did not match")
}

func (s *bffTestSuite) TestMemberList() {
	require := s.Require()

	// Create initial claims
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
		OrgID:       "a2c4f8f0-f8f8-4f8f-8f8f-8f8f8f8f8f8f",
		VASPs:       map[string]string{},
	}

	// Create members list params request
	req := &api.MemberPageInfo{}

	// Set initial RPC handlers to return an error
	require.NoError(s.mainnet.members.UseError(mock.ListRPC, codes.Unavailable, "members list is mock unavailable"))
	require.NoError(s.testnet.members.UseError(mock.ListRPC, codes.Unavailable, "members list is mock unavailable"))

	// Ensure that the read:vasp permission is required
	require.NoError(s.SetClientCredentials(claims), "could not create token with claims")
	_, err := s.client.MemberList(context.TODO(), req)
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error with no permissions on claims")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{auth.ReadVASP}
	require.NoError(s.SetClientCredentials(claims), "could not create token with claims")

	// Ensure a valid directory is required
	req.Directory = "unrecognized.io"
	_, err = s.client.MemberList(context.TODO(), req)
	s.requireError(err, http.StatusBadRequest, "unknown registered directory", "expected invalid directory")

	// Ensure that check verification middleware is required to access the specific directory
	req.Directory = "testnet.directory"
	_, err = s.client.MemberList(context.TODO(), req)
	s.requireError(err, http.StatusUnavailableForLegalReasons, "listing GDS members is only available to verified TRISA members")

	req.Directory = "trisa.directory"
	_, err = s.client.MemberList(context.TODO(), req)
	s.requireError(err, http.StatusUnavailableForLegalReasons, "listing GDS members is only available to verified TRISA members")

	// Enable verification mocks
	s.mainnet.gds.OnVerification = verificationMock(config.MainNet)
	s.testnet.gds.OnVerification = verificationMock(config.TestNet)

	// Ensure that revoked or pending VASPs cannot access member list
	// NOTE: these UUIDs are REJECTED/SUBMITTED states by default set by the verificationMock() function
	claims.VASPs["mainnet"] = "fee76a8c-684b-4d79-bc2f-4439e42a597a"
	claims.VASPs["testnet"] = "9cbcd158-9b37-4200-803a-17fbc188f677"
	require.NoError(s.SetClientCredentials(claims), "could not create token with claims")

	req.Directory = "testnet.directory"
	_, err = s.client.MemberList(context.TODO(), req)
	s.requireError(err, http.StatusUnavailableForLegalReasons, "listing GDS members is only available to verified TRISA members")

	req.Directory = "trisa.directory"
	_, err = s.client.MemberList(context.TODO(), req)
	s.requireError(err, http.StatusUnavailableForLegalReasons, "listing GDS members is only available to verified TRISA members")

	// Add claims for verified VASPs for the remainder of the tests
	// NOTE: these UUIDs are VERIFIED states by default set by the verificationMock() function
	claims.VASPs["mainnet"] = "0846137c-fd14-474e-99b6-f4f33f7f3a86"
	claims.VASPs["testnet"] = "a246f9ff-094a-4fa8-b151-1c8d76e02e86"
	require.NoError(s.SetClientCredentials(claims), "could not create token with claims")

	// Ensure errors are returned from the testnet and mainnet directory when the mocks
	// are set to return unavailable errors.
	req.Directory = "testnet.directory"
	_, err = s.client.MemberList(context.TODO(), req)
	s.requireError(err, http.StatusServiceUnavailable, "specified directory is currently unavailable, please try again later", "expected grpc pass through error")

	req.Directory = "trisa.directory"
	_, err = s.client.MemberList(context.TODO(), req)
	s.requireError(err, http.StatusServiceUnavailable, "specified directory is currently unavailable, please try again later", "expected grpc pass through error")

	// Test the happy path with VASPs correctly returned from both TestNet and MainNet
	s.testnet.members.UseFixture(mock.ListRPC, "testdata/testnet/list_reply.json")
	s.mainnet.members.UseFixture(mock.ListRPC, "testdata/mainnet/list_reply.json")

	req.Directory = "testnet.directory"
	out, err := s.client.MemberList(context.TODO(), req)
	require.NoError(err, "expected valid response from testnet")
	require.Len(out.VASPs, 5)
	require.Equal(out.NextPageToken, "mLB9CU8O8xQj2XEyjAtlfvTj9imoXnLv/1p8fTLchTg=")

	req.Directory = "trisa.directory"
	out, err = s.client.MemberList(context.TODO(), req)
	require.NoError(err, "expected valid response from mainnet")
	require.Len(out.VASPs, 3)
	require.Empty(out.NextPageToken, "expected mainnet next page token to be empty")

	// Test default is the MainNet
	req.Directory = ""
	other, err := s.client.MemberList(context.TODO(), req)
	require.NoError(err, "could not make request with no directory param")
	require.Equal(out, other, "expected default response to match mainnet")
}

func (s *bffTestSuite) TestMemberDetail() {
	require := s.Require()

	testnetID := "7a96ca2c-2818-4106-932e-1bcfd743b04c"
	mainnetID := "9e069e01-8515-4d57-b9a5-e249f7ab4fca"

	testnetDetails := &members.MemberDetails{}
	mainnetDetails := &members.MemberDetails{}
	testnetFixture := filepath.Join("testdata", "testnet", "details_reply.json")
	mainnetFixture := filepath.Join("testdata", "mainnet", "details_reply.json")
	require.NoError(loadFixture(testnetFixture, testnetDetails))
	require.NoError(loadFixture(mainnetFixture, mainnetDetails))

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
		OrgID:       "a2c4f8f0-f8f8-4f8f-8f8f-8f8f8f8f8f8f",
		VASPs:       map[string]string{},
	}

	// Set initial handlers to return an error
	require.NoError(s.testnet.members.UseError(mock.DetailsRPC, codes.Unavailable, "endpoint is unavailable"))
	require.NoError(s.mainnet.members.UseError(mock.DetailsRPC, codes.Unavailable, "endpoint is unavailable"))

	// Endpoint must be authenticated
	req := &api.MemberDetailsParams{ID: "foo"}
	_, err := s.client.MemberDetails(context.TODO(), req)
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.MemberDetails(context.TODO(), req)
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{auth.ReadVASP}
	require.NoError(s.SetClientCredentials(claims), "could not create token with correct permissions")

	// Test with unrecognized directory
	req.ID = "b2c4f8f0-f8f8-4f8f-8f8f-8f8f8f8f8f8f"
	req.Directory = "unrecognized.net"
	_, err = s.client.MemberDetails(context.TODO(), req)
	s.requireError(err, http.StatusBadRequest, "unknown registered directory", "expected error when directory is unrecognized")

	// Ensure that check verification middleware is required to access the specific directory
	req.Directory = "testnet.directory"
	_, err = s.client.MemberDetails(context.TODO(), req)
	s.requireError(err, http.StatusUnavailableForLegalReasons, "listing GDS members is only available to verified TRISA members")

	req.Directory = "trisa.directory"
	_, err = s.client.MemberDetails(context.TODO(), req)
	s.requireError(err, http.StatusUnavailableForLegalReasons, "listing GDS members is only available to verified TRISA members")

	// Enable verification mocks
	s.mainnet.gds.OnVerification = verificationMock(config.MainNet)
	s.testnet.gds.OnVerification = verificationMock(config.TestNet)

	// Ensure that revoked or pending VASPs cannot access member list
	// NOTE: these UUIDs are REJECTED/SUBMITTED states by default set by the verificationMock() function
	claims.VASPs["mainnet"] = "fee76a8c-684b-4d79-bc2f-4439e42a597a"
	claims.VASPs["testnet"] = "9cbcd158-9b37-4200-803a-17fbc188f677"
	require.NoError(s.SetClientCredentials(claims), "could not create token with claims")

	req.Directory = "testnet.directory"
	_, err = s.client.MemberDetails(context.TODO(), req)
	s.requireError(err, http.StatusUnavailableForLegalReasons, "listing GDS members is only available to verified TRISA members")

	req.Directory = "trisa.directory"
	_, err = s.client.MemberDetails(context.TODO(), req)
	s.requireError(err, http.StatusUnavailableForLegalReasons, "listing GDS members is only available to verified TRISA members")

	// Add claims for verified VASPs for the remainder of the tests
	// NOTE: these UUIDs are VERIFIED states by default set by the verificationMock() function
	claims.VASPs["mainnet"] = "0846137c-fd14-474e-99b6-f4f33f7f3a86"
	claims.VASPs["testnet"] = "a246f9ff-094a-4fa8-b151-1c8d76e02e86"
	require.NoError(s.SetClientCredentials(claims), "could not create token with claims")

	// Test error is returned when VASP does not exist in the requested directory
	require.NoError(s.testnet.members.UseError(mock.DetailsRPC, codes.NotFound, "member not found"))
	req.Directory = "testnet.directory"
	_, err = s.client.MemberDetails(context.TODO(), req)
	s.requireError(err, http.StatusNotFound, "member not found", "expected error when VASP does not exist")

	// Test successful response from testnet
	actualSummary := &members.VASPMember{}
	actualPerson := &ivms101.LegalPerson{}
	actualContacts := &pb.Contacts{}
	actualTrixo := &pb.TRIXOQuestionnaire{}
	require.NoError(s.testnet.members.UseFixture(mock.DetailsRPC, testnetFixture))

	req.ID = testnetID
	reply, err := s.client.MemberDetails(context.TODO(), req)
	require.NoError(err, "could not get member details")
	require.NoError(wire.Unwire(reply.Summary, actualSummary), "could not unmarshal summary in response")
	require.Equal(testnetDetails.MemberSummary, actualSummary, "response summary did not match")
	require.NoError(wire.Unwire(reply.LegalPerson, actualPerson), "could not unmarshal legal person in response")
	require.Equal(testnetDetails.LegalPerson, actualPerson, "response legal person did not match")
	require.NoError(wire.Unwire(reply.Contacts, actualContacts), "could not unmarshal contacts in response")
	require.Equal(testnetDetails.Contacts, actualContacts, "response contacts did not match")
	require.NoError(wire.Unwire(reply.Trixo, actualTrixo), "could not unmarshal trixo in response")
	require.Equal(testnetDetails.Trixo, actualTrixo, "response trixo did not match")

	// Ensure that enums are represented as strings in the response
	require.Equal(testnetDetails.MemberSummary.BusinessCategory.String(), reply.Summary["business_category"], "expected business category to be string representation of enum")

	// Test successful response from mainnet and mixed case directory
	req.Directory = "TRISA.directory"
	actualPerson = &ivms101.LegalPerson{}
	actualContacts = &pb.Contacts{}
	actualTrixo = &pb.TRIXOQuestionnaire{}
	require.NoError(s.mainnet.members.UseFixture(mock.DetailsRPC, mainnetFixture))

	req.ID = mainnetID
	reply, err = s.client.MemberDetails(context.TODO(), req)
	require.NoError(err, "could not get member details")
	require.NoError(wire.Unwire(reply.Summary, actualSummary), "could not unmarshal summary in response")
	require.Equal(mainnetDetails.MemberSummary, actualSummary, "response summary did not match")
	require.NoError(wire.Unwire(reply.LegalPerson, actualPerson), "could not unmarshal legal person in response")
	require.Equal(mainnetDetails.LegalPerson, actualPerson, "response legal person did not match")
	require.NoError(wire.Unwire(reply.Contacts, actualContacts), "could not unmarshal contacts in response")
	require.Equal(mainnetDetails.Contacts, actualContacts, "response contacts did not match")
	require.NoError(wire.Unwire(reply.Trixo, actualTrixo), "could not unmarshal trixo in response")
	require.Equal(mainnetDetails.Trixo, actualTrixo, "response trixo did not match")
}
