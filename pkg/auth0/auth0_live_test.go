package auth0_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/auth0"
)

// The Live Auth0 Test Suite executes read requests directly to the Auth0 API so long as
// the environment contains a configuration that allows it to connect. The purpose of
// these tests is to ensure that the SDK works with the actual Auth0 Management API and
// to facilitate development or change management. However these tests will be very
// rarely run since they require manual intervention and will never be run in CI.
type liveAuth0TestSuite struct {
	suite.Suite
	conf   auth0.Config
	client *auth0.Auth0
}

func TestLive(t *testing.T) {
	// These tests will only run if there is a valid configuration in the environment
	conf, err := auth0.NewConfig()
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

	s.conf, err = auth0.NewConfig()
	require.NoError(err, "could not initialize suite from environment configuration")

	s.client, err = auth0.New(s.conf)
	require.NoError(err, "could not initialize suite with auth0 client")
}
