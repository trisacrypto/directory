package bff_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *bffTestSuite) TestCheckVerification() {
	require := s.Require()

	// Check verification should error without any claims
	c, r := createEmptyGinContext()
	s.bff.CheckVerification(c)
	require.Equal(http.StatusInternalServerError, r.Code)
	_, err := bff.GetVerificationStatus(c)
	require.ErrorIs(err, bff.ErrNoVerificationStatus)

	// Create claims with no verification information
	claims := &auth.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
		OrgID:       "a2c4f8f0-f8f8-4f8f-8f8f-8f8f8f8f8f8f",
		VASPs:       auth.VASPs{},
	}

	c, _ = createEmptyGinContext()
	c.Set(auth.ContextBFFClaims, claims)
	s.bff.CheckVerification(c)

	status, err := bff.GetVerificationStatus(c)
	require.NoError(err, "could not fetch verification status")
	require.Equal("NO_VERIFICATION", status.MainNet.Status)
	require.Equal("NO_VERIFICATION", status.TestNet.Status)
	require.False(status.MainNet.Verified)
	require.False(status.TestNet.Verified)
	require.Equal(0, s.mainnet.gds.Calls[mock.VerificationRPC], "expected no calls to GDS with no VASP IDs")
	require.Equal(0, s.testnet.gds.Calls[mock.VerificationRPC], "expected no calls to GDS with no VASP IDs")

	// Setup Verification Mocks
	s.mainnet.gds.OnVerification = verificationMock(config.MainNet)
	s.testnet.gds.OnVerification = verificationMock(config.TestNet)

	testCases := []struct {
		testnet       string
		mainnet       string
		testnetStatus string
		mainnetStatus string
		testName      string
	}{
		{"a246f9ff-094a-4fa8-b151-1c8d76e02e86", "", "VERIFIED", "NO_VERIFICATION", "testnet only, verified"},
		{"9cbcd158-9b37-4200-803a-17fbc188f677", "", "SUBMITTED", "NO_VERIFICATION", "testnet only, not verified"},
		{"86088e9d-b8b3-4d1d-996c-261d204b0f1a", "", "NO_VERIFICATION", "NO_VERIFICATION", "testnet only, not found"},
		{"", "0846137c-fd14-474e-99b6-f4f33f7f3a86", "NO_VERIFICATION", "VERIFIED", "mainnet only, verified"},
		{"", "fee76a8c-684b-4d79-bc2f-4439e42a597a", "NO_VERIFICATION", "REJECTED", "mainnet only, not verified"},
		{"", "86088e9d-b8b3-4d1d-996c-261d204b0f1a", "NO_VERIFICATION", "NO_VERIFICATION", "mainnet only, not found"},
		{"a246f9ff-094a-4fa8-b151-1c8d76e02e86", "0846137c-fd14-474e-99b6-f4f33f7f3a86", "VERIFIED", "VERIFIED", "both verified"},
		{"9cbcd158-9b37-4200-803a-17fbc188f677", "fee76a8c-684b-4d79-bc2f-4439e42a597a", "SUBMITTED", "REJECTED", "both not verified"},
		{"86088e9d-b8b3-4d1d-996c-261d204b0f1a", "86088e9d-b8b3-4d1d-996c-261d204b0f1a", "NO_VERIFICATION", "NO_VERIFICATION", "both not found"},
	}

	testnetCalls, mainnetCalls := 0, 0

	for i, tc := range testCases {
		claims.VASPs.MainNet = tc.mainnet
		claims.VASPs.TestNet = tc.testnet

		c, _ = createEmptyGinContext()
		c.Set(auth.ContextBFFClaims, claims)
		s.bff.CheckVerification(c)

		if tc.testnet != "" {
			testnetCalls++
		}
		if tc.mainnet != "" {
			mainnetCalls++
		}

		status, err = bff.GetVerificationStatus(c)
		require.NoError(err, "could not fetch verification status in test case %d (%s)", i, tc.testName)
		require.Equal(tc.mainnetStatus, status.MainNet.Status, "unexpected mainnet status in test case %d (%s)", i, tc.testName)
		require.Equal(tc.testnetStatus, status.TestNet.Status, "unexpected testnet status in test case %d (%s)", i, tc.testName)
		require.Equal(tc.mainnetStatus == "VERIFIED", status.MainNet.Verified, "unexpected mainnet verified state in test case %d (%s)", i, tc.testName)
		require.Equal(tc.testnetStatus == "VERIFIED", status.TestNet.Verified, "unexpected testnet verified state in test case %d (%s)", i, tc.testName)

		require.Equal(mainnetCalls, s.mainnet.gds.Calls[mock.VerificationRPC], "unexpected number of mainnet calls in test case %d (%s)", i, tc.testName)
		require.Equal(testnetCalls, s.testnet.gds.Calls[mock.VerificationRPC], "unexpected number of testnet calls in test case %d (%s)", i, tc.testName)
	}

	// Update Verification Mocks to return errors
	s.testnet.gds.UseError(mock.VerificationRPC, codes.DataLoss, "something bad happened")
	s.mainnet.gds.UseError(mock.VerificationRPC, codes.DataLoss, "something bad happened")

	claims.VASPs.TestNet = "0846137c-fd14-474e-99b6-f4f33f7f3a86"
	claims.VASPs.MainNet = "a246f9ff-094a-4fa8-b151-1c8d76e02e86"
	c, _ = createEmptyGinContext()
	c.Set(auth.ContextBFFClaims, claims)
	s.bff.CheckVerification(c)

	status, err = bff.GetVerificationStatus(c)
	require.NoError(err, "could not fetch verification status")
	require.Equal("NO_VERIFICATION", status.MainNet.Status)
	require.Equal("NO_VERIFICATION", status.TestNet.Status)
	require.False(status.MainNet.Verified)
	require.False(status.TestNet.Verified)
	require.Equal(mainnetCalls+1, s.mainnet.gds.Calls[mock.VerificationRPC], "expected error call")
	require.Equal(testnetCalls+1, s.testnet.gds.Calls[mock.VerificationRPC], "expected error call")
}

func verificationMock(network string) func(context.Context, *gds.VerificationRequest) (*gds.VerificationReply, error) {
	vasps := make(map[string]models.VerificationState)
	switch network {
	case config.MainNet:
		vasps["0846137c-fd14-474e-99b6-f4f33f7f3a86"] = models.VerificationState_VERIFIED
		vasps["fee76a8c-684b-4d79-bc2f-4439e42a597a"] = models.VerificationState_REJECTED
	case config.TestNet:
		vasps["a246f9ff-094a-4fa8-b151-1c8d76e02e86"] = models.VerificationState_VERIFIED
		vasps["9cbcd158-9b37-4200-803a-17fbc188f677"] = models.VerificationState_SUBMITTED
	}

	return func(_ context.Context, in *gds.VerificationRequest) (out *gds.VerificationReply, err error) {
		var ok bool
		out = &gds.VerificationReply{}

		if out.VerificationStatus, ok = vasps[in.Id]; !ok {
			return nil, status.Error(codes.NotFound, "vasp not found")
		}

		return out, nil
	}
}

// Create an empty test context for testing middleware on
func createEmptyGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	router := gin.New()
	router.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	recorder := &httptest.ResponseRecorder{}

	// Create a test context for testing the middleware on
	c := gin.CreateTestContextOnly(recorder, router)
	c.Request, _ = http.NewRequest(http.MethodGet, "http://localhost/v1/members", nil)
	return c, recorder
}
