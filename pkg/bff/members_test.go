package bff_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
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

func createTestContext(method, target string, body io.Reader, handlers ...gin.HandlerFunc) (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	fmt.Println("createTestContext", method, target)
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(method, target, body)
	req.Header.Set("content-type", "application/json")

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	c.Request = req

	if len(handlers) > 1 {
		r.Handle(method, target, handlers...)
	}
	return c, r, w
}

func doRequest(srv *gin.Engine, w *httptest.ResponseRecorder, c *gin.Context) (data map[string]interface{}, code int, err error) {
	srv.HandleContext(c)

	rep := w.Result()
	defer rep.Body.Close()

	data = make(map[string]interface{})
	var raw []byte
	if raw, err = ioutil.ReadAll(rep.Body); err != nil {
		return nil, 0, err
	}

	if err = json.Unmarshal(raw, &data); err != nil {
		fmt.Println(string(raw))
		return nil, 0, err
	}
	return data, rep.StatusCode, nil
}

func (s *bffTestSuite) TestOverview() {
	require := s.Require()

	// Test 401 with no access token
	_, err := s.client.Overview(context.TODO())
	require.ErrorContains(err, "401", "should return 401 with no token")

	// Test 401 authenticated user without read:vasp permission
	token, err := s.bff.GetMockAuth().NewToken()
	require.NoError(err, "could not create token")

	// Try to set the token in the context for the client method
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", s.bff.GetURL()+"/v1/overview", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)
	_, err = s.client.Overview(c)
	require.ErrorContains(err, "401", "should return 401 with no read:vasp permission")

	// Try to set the token in the context with an HTTP request
	c, srv, r := createTestContext("GET", s.bff.GetURL()+"/v1/overview", nil, nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)
	_, code, err := doRequest(srv, r, c)
	require.NoError(err, "could not get overview")
	require.Equal(401, code, "should return 401")

	// Test 401 authenticated user with the wrong permission scope

	// Test 200 response with authenticated user with read:vasp permission
}
