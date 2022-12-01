package auth_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
)

func TestUserFetcher(t *testing.T) {
	// Setup the authtest server and management client
	srv, err := authtest.Serve()
	require.NoError(t, err, "could not start the authtest server")
	defer srv.Close()

	auth0, err := auth.NewManagementClient(srv.Config())
	require.NoError(t, err, "could not create the auth0 management client")
	fetcher := auth.NewUserFetcher(auth0.User)

	// Test that an error is returned if the user does not exist
	_, err = fetcher.Get("not-a-user")
	require.Error(t, err, "expected an error when fetching a non-existent user")

	// Test that user details can be fetched
	data, err := fetcher.Get(authtest.UserID)
	require.NoError(t, err, "could not fetch user details")
	require.NotNil(t, data, "user details should not be nil")

	// Test that user details can be asserted to the UserDetails type
	details, ok := data.(*auth.UserDetails)
	require.True(t, ok, "could not assert data to the UserDetails type")

	// Test that the user details are correct
	expected := &auth.UserDetails{
		Name:  authtest.Name,
		Roles: []string{authtest.UserRole},
	}
	require.Equal(t, expected, details, "user details do not match")
}
