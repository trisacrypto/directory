package server_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/sectigo/server"
)

func (s *serverTestSuite) TestLogin() {
	// Create a new client to ensure the client is not logged in
	require := s.Require()
	profile := sectigo.Config{
		Username: "badusername",
		Password: "incorrectpassword",
		Profile:  sectigo.ProfileCipherTraceEE,
		Testing:  true,
		Endpoint: s.srv.URL(),
	}
	client, err := sectigo.New(profile)
	require.NoError(err)

	err = client.Authenticate()
	require.ErrorIs(err, sectigo.ErrNotAuthorized)

	profile.Username = sectigo.MockUsername
	profile.Password = sectigo.MockPassword
	client, err = sectigo.New(profile)
	require.NoError(err)

	err = client.Authenticate()
	require.NoError(err)
}

func (s *serverTestSuite) TestRefresh() {
	require := s.Require()
	err := s.client.Refresh()
	require.NoError(err)
}

func TestTokens(t *testing.T) {
	conf := server.AuthConfig{
		Issuer:  "http://localhost:8831",
		Subject: "testuser",
		Scopes:  []string{"ROLE_USER"},
	}

	tokens, err := server.NewTokens(conf)
	require.NoError(t, err, "could not create token manager")

	ats, rts, err := tokens.SignedTokenPair()
	require.NoError(t, err, "could not create signed token pair")
	require.NotEmpty(t, rts, "no refresh token returned")

	claims, err := tokens.Verify(ats)
	require.NoError(t, err, "could not verify access token")
	require.Equal(t, conf.Issuer, claims.Issuer)
	require.Equal(t, conf.Subject, claims.Subject)
	require.Equal(t, conf.Scopes, claims.Scopes)
}
