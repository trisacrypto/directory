package tokens_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

type TokenTestSuite struct {
	suite.Suite
	testdata map[string]string
}

func (s *TokenTestSuite) SetupSuite() {
	// Create the keys map from the testdata directory to create new token managers.
	s.testdata = make(map[string]string)
	s.testdata["1yAwhf28bXi3IWP6FYcGa0dcrfq"] = "testdata/1yAwhf28bXi3IWP6FYcGa0dcrfq.pem"
	s.testdata["1yAxs5vPqCrg433fPrFENevvzen"] = "testdata/1yAxs5vPqCrg433fPrFENevvzen.pem"
}

func (s *TokenTestSuite) TestTokenManager() {
	require := s.Require()
	tm, err := tokens.New(s.testdata)
	require.NoError(err, "could not initialize token manager")

	keys := tm.Keys()
	require.Len(keys, 2)
	require.Equal("1yAxs5vPqCrg433fPrFENevvzen", tm.CurrentKey().String())

	// Create an access token from simple claims
	creds := map[string]interface{}{
		"hd":      "rotational.io",
		"email":   "kate@rotational.io",
		"name":    "Kate Holland",
		"picture": "https://foo.googleusercontent.com/test!/Aoh14gJceTrUA",
	}

	accessToken, err := tm.CreateAccessToken(creds)
	require.NoError(err, "could not create access token from claims")
	require.IsType(&tokens.Claims{}, accessToken.Claims)

	time.Sleep(500 * time.Millisecond)
	now := time.Now()

	// Check access token claims
	ac := accessToken.Claims.(*tokens.Claims)
	require.NotZero(ac.Id)
	require.Equal("", ac.Audience)
	require.Equal("", ac.Issuer)
	require.Equal("", ac.Subject)
	require.True(time.Unix(ac.IssuedAt, 0).Before(now))
	require.True(time.Unix(ac.NotBefore, 0).Before(now))
	require.True(time.Unix(ac.ExpiresAt, 0).After(now))
	require.Equal(creds["hd"], ac.Domain)
	require.Equal(creds["email"], ac.Email)
	require.Equal(creds["name"], ac.Name)
	require.Equal(creds["picture"], ac.Picture)

	// Create a refresh token from the access token
	refreshToken, err := tm.CreateRefreshToken(accessToken)
	require.NoError(err, "could not create refresh token from access token")
	require.IsType(&tokens.Claims{}, refreshToken.Claims)

	// Check refresh token claims
	// Check access token claims
	rc := refreshToken.Claims.(*tokens.Claims)
	require.Equal(ac.Id, rc.Id, "access and refresh tokens must have same jid")
	require.Equal(ac.Audience, rc.Audience)
	require.Equal(ac.Issuer, rc.Issuer)
	require.Equal(ac.Subject, rc.Subject)
	require.True(time.Unix(rc.IssuedAt, 0).Equal(time.Unix(ac.IssuedAt, 0)))
	require.True(time.Unix(rc.NotBefore, 0).After(now))
	require.True(time.Unix(rc.ExpiresAt, 0).After(time.Unix(rc.NotBefore, 0)))
	require.Empty(rc.Domain)
	require.Empty(rc.Email)
	require.Empty(rc.Name)
	require.Empty(rc.Picture)

	// Sign the access token
	atks, err := tm.Sign(accessToken)
	require.NoError(err, "could not sign access token")

	// Sign the refresh token
	rtks, err := tm.Sign(refreshToken)
	require.NoError(err, "could not sign refresh token")
	require.NotEqual(atks, rtks, "identical access and refresh tokens")

	// Validate the access token
	_, err = tm.Verify(atks)
	require.NoError(err, "could not validate access token")

	// Validate the refresh token (should be invalid because of not before in the future)
	_, err = tm.Verify(rtks)
	require.Error(err, "refresh token is valid?")
}

// TODO: test validation of audience, issuer, and subject.
// TODO: test time based validation (not before, issued at, expires)
// TODO: test signed with wrong key
func (s *TokenTestSuite) TestValidTokens() {
	require := s.Require()
	tm, err := tokens.New(s.testdata)
	require.NoError(err, "could not initialize token manager")

	// Default creds
	creds := map[string]interface{}{
		"hd":      "rotational.io",
		"email":   "kate@rotational.io",
		"name":    "Kate Holland",
		"picture": "https://foo.googleusercontent.com/test!/Aoh14gJceTrUA",
	}

	// Test required creds
	for _, key := range []string{"hd", "email"} {
		// remove key
		orig := creds[key]
		delete(creds, key)

		// Fail when creating an access token
		_, err = tm.CreateAccessToken(creds)
		require.Error(err)

		// replace key
		creds[key] = orig
	}

	// Test optional creds
	for _, key := range []string{"name", "picture"} {
		// remove key
		orig := creds[key]
		delete(creds, key)

		// Fail when creating an access token
		_, err = tm.CreateAccessToken(creds)
		require.NoError(err)

		// replace key
		creds[key] = orig
	}

	// Test optional and empty
	creds["name"] = ""
	_, err = tm.CreateAccessToken(creds)
	require.NoError(err)

	// Test not optional and empty
	creds["email"] = ""
	_, err = tm.CreateAccessToken(creds)
	require.Error(err)

	// Test bad parsing
	creds["email"] = 1234
	_, err = tm.CreateAccessToken(creds)
	require.Error(err)
}

// Test that a token signed with an old cert can still be verified.
func (s *TokenTestSuite) TestKeyRotation() {
	require := s.Require()

	// Create the "old token manager"
	testdata := make(map[string]string)
	testdata["1yAwhf28bXi3IWP6FYcGa0dcrfq"] = "testdata/1yAwhf28bXi3IWP6FYcGa0dcrfq.pem"
	oldTM, err := tokens.New(testdata)
	require.NoError(err, "could not initialize old token manager")

	// Create the "new" token manager with the new key
	newTM, err := tokens.New(s.testdata)
	require.NoError(err, "could not initialize new token manager")

	// Create a valid token with the "old token manager"
	token, err := oldTM.CreateAccessToken(map[string]interface{}{
		"hd":      "rotational.io",
		"email":   "kate@rotational.io",
		"name":    "Kate Holland",
		"picture": "https://foo.googleusercontent.com/test!/Aoh14gJceTrUA",
	})
	require.NoError(err)

	tks, err := oldTM.Sign(token)
	require.NoError(err)

	// Validate token with "new token manager"
	_, err = newTM.Verify(tks)
	require.NoError(err)
}

// Execute suite as a go test.
func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}
