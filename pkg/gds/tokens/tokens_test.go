package tokens_test

import (
	"crypto/rand"
	"crypto/rsa"
	"strings"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

type TokenTestSuite struct {
	suite.Suite
	testdata map[string]string
}

func (s *TokenTestSuite) SetupSuite() {
	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog MUST be false otherwise this will be overriden
	logger.Discard()

	// Create the keys map from the testdata directory to create new token managers.
	s.testdata = make(map[string]string)
	s.testdata["1yAwhf28bXi3IWP6FYcGa0dcrfq"] = "testdata/1yAwhf28bXi3IWP6FYcGa0dcrfq.pem"
	s.testdata["1yAxs5vPqCrg433fPrFENevvzen"] = "testdata/1yAxs5vPqCrg433fPrFENevvzen.pem"
}

func (s *TokenTestSuite) TearDownSuite() {
	logger.ResetLogger()
}

func (s *TokenTestSuite) TestTokenManager() {
	require := s.Require()
	tm, err := tokens.New(s.testdata, "http://localhost:3000")
	require.NoError(err, "could not initialize token manager")

	keys := tm.Keys()
	require.Len(keys, 2)
	require.Equal("1yAxs5vPqCrg433fPrFENevvzen", tm.CurrentKey().String())

	// Create an access token from simple claims
	creds := map[string]interface{}{
		"sub":     "102374163855881761273",
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
	require.Equal("http://localhost:3000", ac.Audience)
	require.Empty(ac.Issuer, "issuer is a duplicate of audience")
	require.Equal("102374163855881761273", ac.Subject)
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

	// Verify relative nbf and exp claims of access and refresh tokens
	require.True(time.Unix(ac.IssuedAt, 0).Equal(time.Unix(rc.IssuedAt, 0)), "access and refresh tokens do not have same iss timestamp")
	require.Equal(45*time.Minute, time.Unix(rc.NotBefore, 0).Sub(time.Unix(ac.IssuedAt, 0)), "refresh token nbf is not 45 minutes after access token iss")
	require.Equal(15*time.Minute, time.Unix(ac.ExpiresAt, 0).Sub(time.Unix(rc.NotBefore, 0)), "refresh token active does not overlap active token active by 15 minutes")
	require.Equal(60*time.Minute, time.Unix(rc.ExpiresAt, 0).Sub(time.Unix(ac.ExpiresAt, 0)), "refresh token does not expire 1 hour after access token")

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

func (s *TokenTestSuite) TestValidTokens() {
	require := s.Require()
	tm, err := tokens.New(s.testdata, "http://localhost:3000")
	require.NoError(err, "could not initialize token manager")

	// Default creds
	creds := map[string]interface{}{
		"sub":     "102374163855881761273",
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
	for _, key := range []string{"name", "picture", "sub"} {
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

func (s *TokenTestSuite) TestInvalidTokens() {
	// Create the token manager
	require := s.Require()
	tm, err := tokens.New(s.testdata, "http://localhost:3000")
	require.NoError(err, "could not initialize token manager")

	// Manually create a token to validate with the token manager
	now := time.Now()
	claims := &tokens.Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewString(),                  // id not validated
			Audience:  "http://foo.example.com",          // wrong audience
			Subject:   "102374163855881761273",           // sub not validated
			IssuedAt:  now.Add(-1 * time.Hour).Unix(),    // iat not validated
			NotBefore: now.Add(15 * time.Minute).Unix(),  // nbf is validated and is after now
			ExpiresAt: now.Add(-30 * time.Minute).Unix(), // exp is validated and is before now
		},
		Domain:  "rotational.io",
		Email:   "kate@rotational.io",
		Name:    "Kate Holland",
		Picture: "https://foo.googleusercontent.com/test!/Aoh14gJceTrUA",
	}

	// Test validation signed with wrong kid
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "1zSQqRhO7lU1qXoOZvvF2kJdRt5"
	badkey, err := rsa.GenerateKey(rand.Reader, 1024)
	require.NoError(err, "could not generate bad rsa keys")
	tks, err := token.SignedString(badkey)
	require.NoError(err, "could not sign token with bad kid")

	_, err = tm.Verify(tks)
	require.EqualError(err, "unknown signing key")

	// Test validation signed with good kid but wrong key
	token.Header["kid"] = "1yAwhf28bXi3IWP6FYcGa0dcrfq"
	tks, err = token.SignedString(badkey)
	require.NoError(err, "could not sign token with bad keys and good kid")

	_, err = tm.Verify(tks)
	require.EqualError(err, "crypto/rsa: verification error")

	// Test time-based validation: nbf
	tks, err = tm.Sign(token)
	require.NoError(err, "could not sign token with good keys")

	_, err = tm.Verify(tks)
	require.EqualError(err, "token is not valid yet")

	// Test time-based validation: exp
	claims.NotBefore = now.Add(-1 * time.Hour).Unix()
	tks, err = tm.Sign(jwt.NewWithClaims(jwt.SigningMethodRS256, claims))
	require.NoError(err, "could not sign token with good keys")

	// NOTE: actual error message is "token is expired by 30m0s" however, if the clock
	// minute happens to tick over the message could be "token is expired by 30m1s" so
	// to prevent test failure, we're only testing the prefix.
	_, err = tm.Verify(tks)
	require.True(strings.HasPrefix(err.Error(), "token is expired"))

	// Test audience verification
	claims.ExpiresAt = now.Add(1 * time.Hour).Unix()
	tks, err = tm.Sign(jwt.NewWithClaims(jwt.SigningMethodRS256, claims))
	require.NoError(err, "could not sign token with good keys")

	_, err = tm.Verify(tks)
	require.EqualError(err, "invalid audience \"http://foo.example.com\"")

	// Token is finally valid
	claims.Audience = "http://localhost:3000"
	tks, err = tm.Sign(jwt.NewWithClaims(jwt.SigningMethodRS256, claims))
	require.NoError(err, "could not sign token with good keys")
	_, err = tm.Verify(tks)
	require.NoError(err, "claims are still not valid")
}

// Test that a token signed with an old cert can still be verified.
// This also tests that the correct signing key is required.
func (s *TokenTestSuite) TestKeyRotation() {
	require := s.Require()

	// Create the "old token manager"
	testdata := make(map[string]string)
	testdata["1yAwhf28bXi3IWP6FYcGa0dcrfq"] = "testdata/1yAwhf28bXi3IWP6FYcGa0dcrfq.pem"
	oldTM, err := tokens.New(testdata, "http://localhost:3000")
	require.NoError(err, "could not initialize old token manager")

	// Create the "new" token manager with the new key
	newTM, err := tokens.New(s.testdata, "http://localhost:3000")
	require.NoError(err, "could not initialize new token manager")

	// Create a valid token with the "old token manager"
	token, err := oldTM.CreateAccessToken(map[string]interface{}{
		"sub":     "102374163855881761273",
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

	// A token created by the "new token mangaer" should not be verified by the old one.
	tks, err = newTM.Sign(token)
	require.NoError(err)

	_, err = oldTM.Verify(tks)
	require.Error(err)
}

// Test that a token can be parsed even if it is expired. This is necessary to parse
// access tokens in order to use a refresh token to extract the claims.
func (s *TokenTestSuite) TestParseExpiredToken() {
	require := s.Require()
	tm, err := tokens.New(s.testdata, "http://localhost:3000")
	require.NoError(err, "could not initialize token manager")

	// Default creds
	creds := map[string]interface{}{
		"sub":     "102374163855881761273",
		"hd":      "rotational.io",
		"email":   "kate@rotational.io",
		"name":    "Kate Holland",
		"picture": "https://foo.googleusercontent.com/test!/Aoh14gJceTrUA",
	}

	accessToken, err := tm.CreateAccessToken(creds)
	require.NoError(err, "could not create access token from claims")
	require.IsType(&tokens.Claims{}, accessToken.Claims)

	// Modify claims to be expired
	claims := accessToken.Claims.(*tokens.Claims)
	claims.IssuedAt = time.Unix(claims.IssuedAt, 0).Add(-24 * time.Hour).Unix()
	claims.ExpiresAt = time.Unix(claims.ExpiresAt, 0).Add(-24 * time.Hour).Unix()
	claims.NotBefore = time.Unix(claims.NotBefore, 0).Add(-24 * time.Hour).Unix()
	accessToken.Claims = claims

	// Create signed token
	tks, err := tm.Sign(accessToken)
	require.NoError(err, "could not create expired access token from claims")

	// Ensure that verification fails; claims are invalid.
	pclaims, err := tm.Verify(tks)
	require.Error(err, "expired token was somehow validated?")
	require.Empty(pclaims, "verify returned claims even after error")

	// Parse token without verifying claims but verifying the signature
	pclaims, err = tm.Parse(tks)
	require.NoError(err, "claims were validated in parse")
	require.NotEmpty(pclaims, "parsing returned empty claims without error")

	// Check claims
	require.Equal(claims.Id, pclaims.Id)
	require.Equal(claims.ExpiresAt, pclaims.ExpiresAt)
	require.Equal(creds["sub"], claims.Subject)
	require.Equal(creds["hd"], claims.Domain)
	require.Equal(creds["email"], claims.Email)
	require.Equal(creds["name"], claims.Name)
	require.Equal(creds["picture"], claims.Picture)

	// Ensure signature is still validated on parse
	tks += "abcdefg"
	claims, err = tm.Parse(tks)
	require.Error(err, "claims were parsed with bad signature")
	require.Empty(claims, "bad signature token returned non-empty claims")
}

// Execute suite as a go test.
func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}
