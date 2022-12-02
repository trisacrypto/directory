package auth_test

import (
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/config"
)

func (s *authTestSuite) TestUserFetcher() {
	require := s.Require()
	fetcher := auth.NewUserFetcher(s.auth0.User)

	// Test that an error is returned if the user does not exist
	_, err := fetcher.Get("not-a-user")
	require.Error(err, "expected an error when fetching a non-existent user")

	// Test that user details can be fetched
	data, err := fetcher.Get(authtest.UserID)
	require.NoError(err, "could not fetch user details")
	require.NotNil(data, "user details should not be nil")

	// Test that user details can be asserted to the UserDetails type
	details, ok := data.(*auth.UserDetails)
	require.True(ok, "could not assert data to the UserDetails type")

	// Test that the user details are correct
	expected := &auth.UserDetails{
		Name:  authtest.Name,
		Roles: []string{authtest.UserRole},
	}
	require.Equal(expected, details, "user details do not match")
}

func (s *authTestSuite) TestUserCache() {
	require := s.Require()

	// Create a new user cache with a 1 second TTL
	conf := config.CacheConfig{
		Enabled:          true,
		TTLMean:          time.Minute,
		TTLSigma:         time.Second,
		MaxEntries:       100,
		EvictionFraction: 0.1,
	}
	users := auth.NewUserCache(conf, s.auth0.User)

	// Fetch a user from the cache
	u, err := users.Get(authtest.UserID)
	require.NoError(err, "could not get user from cache")

	// Test that the user details are correct
	expected := &auth.UserDetails{
		Name:  authtest.Name,
		Roles: []string{authtest.UserRole},
	}
	require.Equal(expected, u, "user details do not match")

	// Test that an error is returned if the user does not exist
	_, err = users.Get("not-a-user")
	require.Error(err, "expected an error when fetching a non-existent user")
}

type authTestSuite struct {
	suite.Suite
	srv   *authtest.Server
	auth0 *management.Management
}

func (s *authTestSuite) SetupSuite() {
	var err error
	require := s.Require()

	// Setup the authtest server and management client
	s.srv, err = authtest.Serve()
	require.NoError(err, "could not start the authtest server")

	s.auth0, err = auth.NewManagementClient(s.srv.Config())
	require.NoError(err, "could not create the auth0 management client")
}

func (s *authTestSuite) TearDownSuite() {
	// Shutdown the authtest server
	authtest.Close()
}
