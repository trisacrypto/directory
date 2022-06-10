package auth_test

import (
	"net/http"
	"testing"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
)

// The Live Auth0 Test Suite executes read requests directly to the Auth0 API so long as
// the environment contains a configuration that allows it to connect. The purpose of
// these tests is to ensure that the auth middleware works with the actual Auth0
// Management API and to facilitate development or change management. However these
// tests will be very rarely run since they require manual intervention and will never
// be run in CI.
type liveAuth0TestSuite struct {
	suite.Suite
	conf authtest.Config
}

func TestLive(t *testing.T) {
	// These tests will only run if there is a valid configuration in the environment
	conf, err := authtest.NewConfig()
	if err != nil {
		t.Skip("live tests require local environment configuration")
	}

	// Do not run the live tests if there is no access token cacheing
	if conf.TokenCache == "" {
		t.Skip("live tests require a token cache to prevent issuing multiple M2M tokens")
	}

	//  Log the situation for the tests and run the test suite.
	t.Logf("live tests starting with auth0 client %s, using token cache %s", conf.ClientID, conf.TokenCache)
	suite.Run(t, new(liveAuth0TestSuite))
}

func (s *liveAuth0TestSuite) SetupSuite() {
	var err error
	require := s.Require()

	s.conf, err = authtest.NewConfig()
	require.NoError(err, "could not initialize suite from environment configuration")
}

func (s *liveAuth0TestSuite) TestUserInfo() {
	// NOTE: this user ID must be in the tenant for this test to work.
	userID := "auth0|62a014c5881f6b006f97ed30"
	require := s.Require()

	middleware, err := auth.UserInfo(s.conf.AuthConfig())
	require.NoError(err, "could not create middleware with live config")

	ctx, _, _ := createTestContext(http.MethodGet, "/", nil, middleware)
	ctx.Set(auth.ContextRegisteredClaims, &validator.RegisteredClaims{Subject: userID})

	// Executing the middleware should put the user info on the context
	middleware(ctx)
	user, err := auth.GetUserInfo(ctx)
	require.NoError(err, "no user info was added to the context")
	require.NotNil(user.Email, "no email was returned for the user")
	require.Equal("leopold.wentzel@gmail.com", *user.Email)
}
