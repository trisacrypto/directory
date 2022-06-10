package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"gopkg.in/square/go-jose.v2"
)

func TestClaims(t *testing.T) {
	claims := &auth.Claims{
		Scope:       "openid profile email",
		Permissions: []string{"read:foo", "write:foo", "delete:foo", "read:bar"},
	}

	// Test Validate
	require.NoError(t, claims.Validate(context.TODO()), "claims should be valid")

	// Test HasScope
	require.False(t, claims.HasScope("machine"), "unexpected scope returned")
	require.True(t, claims.HasScope("profile"), "expected profile to be in scope")

	// Test Permissions
	require.False(t, claims.HasPermission("write:bar"), "unexpected permission returned")
	require.True(t, claims.HasPermission("write:foo"), "expected permission to be true")
	require.False(t, claims.HasAllPermissions("write:foo", "write:bar"), "only has one permission")
	require.False(t, claims.HasAllPermissions("delete:bar", "write:bar"), "has no permissions")
	require.True(t, claims.HasAllPermissions("delete:foo", "write:foo", "read:foo"), "has all permissions")

	// Test Claim creation function
	customClaims := auth.NewClaims()
	require.IsType(t, claims, customClaims, "new claims did not return the expected type")
}

func TestClaimsContext(t *testing.T) {
	// Load claims fixture
	data, err := ioutil.ReadFile("testdata/validated_claims.json")
	require.NoError(t, err, "could not read validated claims fixture")

	vclaims := &validator.ValidatedClaims{CustomClaims: &auth.Claims{}, RegisteredClaims: validator.RegisteredClaims{}}
	err = json.Unmarshal(data, vclaims)
	require.NoError(t, err, "could not unmarshal validated claims fixture")
	require.Equal(t, "rebecca@example.com", vclaims.CustomClaims.(*auth.Claims).Email, "could not unmarshal to custom claims interface")

	// Create gin context fixture
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Test empty context errors
	_, err = auth.GetClaims(c)
	require.ErrorIs(t, err, auth.ErrNoClaims, "empty context returned bff claims?")
	_, err = auth.GetRegisteredClaims(c)
	require.ErrorIs(t, err, auth.ErrNoClaims, "empty context returned registered claims?")

	// Test context with values
	c.Set(auth.ContextBFFClaims, vclaims.CustomClaims.(*auth.Claims))
	c.Set(auth.ContextRegisteredClaims, vclaims.RegisteredClaims)

	bclaims, err := auth.GetClaims(c)
	require.NoError(t, err, "could not fetch bff claims")
	require.Equal(t, "6f0d943d-6cd7-4745-bc9d-6d65e32c70e9", bclaims.OrgID)
	require.Equal(t, "eee784b5-49b3-452e-97d5-1b01e79f5e62", bclaims.VASP["testnet"])
	require.True(t, bclaims.HasAllPermissions("add:collaborators", "read:certificates"))

	rclaims, err := auth.GetRegisteredClaims(c)
	require.NoError(t, err, "could not fetch registered claims")
	require.Equal(t, "auth0|9d56yPBwjRv1tp3gEIA", rclaims.Subject)
	require.Equal(t, "https://example.auth0.com/", rclaims.Issuer)
}

func TestAuthenticate(t *testing.T) {
	// Can only test bad paths of Authenticate middleware since a live JWT token is
	// required to obtain the happy path. It is possible to get one, but it will expire.
	// See TestAuthenticatePublicKeys for a mock of the Auth0 known keys endpoint.

	// A valid issuer url is required to create the middleware.
	conf := config.AuthConfig{}
	_, err := auth.Authenticate(conf)
	require.Error(t, err, "expected invalid issuer url error")

	conf.Domain = "example.auth0.com"
	conf.Audience = "http://localhost:3000"
	authenticate, err := auth.Authenticate(conf)
	require.NoError(t, err, "could not create valid authenticate middleware")

	// Create default handler
	success := func(c *gin.Context) {
		c.JSON(http.StatusOK, api.Reply{Success: true})
	}

	// Test anonymous user (no authorization header in request)
	c, s, w := createTestContext(http.MethodGet, "/", nil, authenticate, success)
	rep, code, err := doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusOK, code)
	require.Contains(t, rep, "success", "response does not contain a success field")
	require.True(t, rep["success"].(bool), "success is not true")

	claims, err := auth.GetClaims(c)
	require.NoError(t, err, "expected anonymous claims on context")
	require.True(t, claims.IsAnonymous(), "expected anonymous claims on context")

	// Test non-bearer token
	c, s, w = createTestContext(http.MethodGet, "/", nil, authenticate, success)
	c.Request.Header.Set("authorization", "token foo")
	_, code, err = doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusOK, code)

	claims, err = auth.GetClaims(c)
	require.NoError(t, err, "expected anonymous claims on context")
	require.True(t, claims.IsAnonymous(), "expected anonymous claims on context")

	// Test forbidden error with incorrectly signed token
	token, err := ioutil.ReadFile("testdata/invalid_token.txt")
	require.NoError(t, err, "could not read invalid token fixture")

	c, s, w = createTestContext(http.MethodGet, "/", nil, authenticate, success)
	c.Request.Header.Set("authorization", "Bearer "+string(token))
	rep, code, err = doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusForbidden, code)
	require.Contains(t, rep, "error", "expected error on JSON response")
	require.Equal(t, "invalid authorization token", rep["error"])
}

func TestAuthenticatePublicKeys(t *testing.T) {
	// Creates a test server that serves well known jwks keys instead of the Auth0
	// tenant - used to mock Auth0 (not a live test) but checks the happy path.
	// NOTE: this test is fragile, e.g. if Auth0 changes its implementation.
	require.NoError(t, createTokenFixtures(), "could not create required token key fixtures")

	// The test server returns the well known token to authenticate the token like Auth0 does.
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("testdata/openid-configuration.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err.Error())
		}

		config := make(map[string]interface{})
		if err = json.Unmarshal(data, &config); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err.Error())
		}

		// Add the server URL to the openid paths
		endpoints := []string{"issuer", "authorization_endpoint", "token_endpoint", "device_authorization_endpoint", "userinfo_endpoint", "mfa_challenge_endpoint", "jwks_uri", "registration_endpoint", "revocation_endpoint"}
		for _, endpoint := range endpoints {
			u := &url.URL{Scheme: "https", Host: r.Host, Path: config[endpoint].(string)}
			config[endpoint] = u.String()
		}

		w.Header().Add("content-type", "application/json")
		json.NewEncoder(w).Encode(&config)
	})
	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/jwks.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err.Error())
		}
		defer f.Close()

		w.Header().Add("content-type", "application/json")
		io.Copy(w, f)
	})

	srv := httptest.NewTLSServer(mux)
	defer srv.Close()

	testURL, _ := url.Parse(srv.URL)
	conf := config.AuthConfig{
		Domain:   testURL.Host,
		Audience: testURL.String(),
	}

	// Setup authentication middleware
	authenticate, err := auth.Authenticate(conf, auth.WithHTTPClient(srv.Client()))
	require.NoError(t, err, "expected valid authenticate middleware")

	// Create default handler
	success := func(c *gin.Context) {
		c.JSON(http.StatusOK, api.Reply{Success: true})
	}

	// Create a valid token to authenticate
	tks, err := createValidToken(srv.URL)
	require.NoError(t, err, "could not create valid token")

	// Execute request expecting success
	c, s, w := createTestContext(http.MethodGet, "/", nil, authenticate, success)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tks))
	rep, code, err := doRequest(s, w, c)

	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusOK, code)
	require.Contains(t, rep, "success", "response does not contain a success field")
	require.True(t, rep["success"].(bool), "success is not true")
}

func TestAuthorize(t *testing.T) {
	// Setup authorization middleware
	authorize := auth.Authorize("read:foo", "read:bar")

	// Create default handler
	success := func(c *gin.Context) {
		c.JSON(http.StatusOK, api.Reply{Success: true})
	}

	// Test unauthorized no claims on context
	c, s, w := createTestContext(http.MethodGet, "/", nil, authorize, success)
	rep, code, err := doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "could not authorize request", rep["error"], "unexpected error returned from authorize")

	// Test anonymous user on context
	authenticate := func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.AnonymousClaims)
		c.Next()
	}

	c, s, w = createTestContext(http.MethodGet, "/", nil, authenticate, authorize, success)
	rep, code, err = doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "this endpoint requires authentication", rep["error"], "unexpected error returned from authorize")

	// Test user does not have permissions
	authenticate = func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.Claims{Permissions: []string{"write:foo"}})
		c.Next()
	}

	c, s, w = createTestContext(http.MethodGet, "/", nil, authenticate, authorize, success)
	rep, code, err = doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "user does not have permission to perform this operation", rep["error"], "unexpected error returned from authorize")

	// Test user does not have all permissions
	authenticate = func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.Claims{Permissions: []string{"read:foo"}})
		c.Next()
	}

	c, s, w = createTestContext(http.MethodGet, "/", nil, authenticate, authorize, success)
	rep, code, err = doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "user does not have permission to perform this operation", rep["error"], "unexpected error returned from authorize")

	// Test user does have permissions
	authenticate = func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.Claims{Permissions: []string{"read:foo", "read:bar"}})
		c.Next()
	}

	c, s, w = createTestContext(http.MethodGet, "/", nil, authenticate, authorize, success)
	rep, code, err = doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusOK, code)
	require.Contains(t, rep, "success", "response does not contain a success field")
	require.True(t, rep["success"].(bool), "success is not true")
}

func TestUserInfo(t *testing.T) {
	t.Skip("not implemented fully yet")

	// Setup UserInfo middleware
	middleware, err := auth.UserInfo(config.AuthConfig{Domain: "example.auth0.com", ClientID: "example", ClientSecret: "supersecretsquirrel"})
	require.NoError(t, err, "could not create user info middleware")

	// Create default handler
	success := func(c *gin.Context) {
		c.JSON(http.StatusOK, api.Reply{Success: true})
	}

	// Test userinfo no claims on context
	c, s, w := createTestContext(http.MethodGet, "/", nil, middleware, success)
	rep, code, err := doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "could not authorize request", rep["error"], "unexpected error returned from authorize")

	// Test anonymous user on context
	authenticate := func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.AnonymousClaims)
		c.Next()
	}

	c, s, w = createTestContext(http.MethodGet, "/", nil, authenticate, middleware, success)
	rep, code, err = doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "this endpoint requires authentication", rep["error"], "unexpected error returned from authorize")

	// Test user does not have permissions
	authenticate = func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.Claims{Permissions: []string{"write:foo"}})
		c.Next()
	}

	c, s, w = createTestContext(http.MethodGet, "/", nil, authenticate, middleware, success)
	rep, code, err = doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "user does not have permission to perform this operation", rep["error"], "unexpected error returned from authorize")

	// Test user does not have all permissions
	authenticate = func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.Claims{Permissions: []string{"read:foo"}})
		c.Next()
	}

	c, s, w = createTestContext(http.MethodGet, "/", nil, authenticate, middleware, success)
	rep, code, err = doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "user does not have permission to perform this operation", rep["error"], "unexpected error returned from authorize")

	// Test user does have permissions
	authenticate = func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.Claims{Permissions: []string{"read:foo", "read:bar"}})
		c.Next()
	}

	c, s, w = createTestContext(http.MethodGet, "/", nil, authenticate, middleware, success)
	rep, code, err = doRequest(s, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusOK, code)
	require.Contains(t, rep, "success", "response does not contain a success field")
	require.True(t, rep["success"].(bool), "success is not true")
}

func createTestContext(method, target string, body io.Reader, handlers ...gin.HandlerFunc) (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
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

func doRequest(s *gin.Engine, w *httptest.ResponseRecorder, c *gin.Context) (data map[string]interface{}, code int, err error) {
	s.HandleContext(c)

	rep := w.Result()
	defer rep.Body.Close()

	data = make(map[string]interface{})
	var raw []byte
	if raw, err = ioutil.ReadAll(rep.Body); err != nil {
		return nil, 0, err
	}

	if err = json.Unmarshal(raw, &data); err != nil {
		return nil, 0, err
	}
	return data, rep.StatusCode, nil
}

// Create RSA key fixtures and well known keys testdata if they don't exist yet.
func createTokenFixtures() (err error) {
	// check if fixtures already exist.
	if err = checkFixtures(); err == nil {
		return nil
	}

	if err = generateKeyPairFixture(); err != nil {
		return err
	}

	if err = generateJWKSFixture(); err != nil {
		return err
	}

	return nil
}

const kid = "StyqeY8Kl4Eam28KsUs"

var fixturePaths = []string{
	"testdata/token_keys.pem",
	"testdata/jwks.json",
}

// Returns an error if the expected fixtures are missing.
func checkFixtures() error {
	for _, path := range fixturePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("required fixture %s does not exist", path)
		}
	}
	return nil
}

func generateKeyPairFixture() (err error) {
	// create an rsa token key pair for use with tokens
	var keypair *rsa.PrivateKey
	if keypair, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return err
	}

	var f *os.File
	if f, err = os.OpenFile("testdata/token_keys.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return err
	}
	defer f.Close()

	block := &pem.Block{Type: "PRIVATE KEY"}
	if block.Bytes, err = x509.MarshalPKCS8PrivateKey(keypair); err != nil {
		return err
	}

	if err = pem.Encode(f, block); err != nil {
		return err
	}

	block = &pem.Block{Type: "PUBLIC KEY"}
	if block.Bytes, err = x509.MarshalPKIXPublicKey(&keypair.PublicKey); err != nil {
		return err
	}

	if err = pem.Encode(f, block); err != nil {
		return err
	}
	return nil
}

func generateJWKSFixture() (err error) {
	var key *rsa.PrivateKey
	if key, err = loadKeyFixture(); err != nil {
		return err
	}

	webkeys := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Key:       &key.PublicKey,
				KeyID:     kid,
				Algorithm: jwt.SigningMethodRS256.Alg(),
			},
		},
	}

	var data []byte
	if data, err = json.Marshal(webkeys); err != nil {
		return err
	}

	return ioutil.WriteFile("testdata/jwks.json", data, 0644)
}

func loadKeyFixture() (_ *rsa.PrivateKey, err error) {
	var data []byte
	if data, err = ioutil.ReadFile("testdata/token_keys.pem"); err != nil {
		return nil, err
	}

	var block *pem.Block
	for {
		block, data = pem.Decode(data)
		if block == nil {
			break
		}

		if block.Type == "PRIVATE KEY" {
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			return key.(*rsa.PrivateKey), nil
		}
	}
	return nil, errors.New("could not find private key in pem data")
}

func createValidToken(iss string) (tks string, err error) {
	var key *rsa.PrivateKey
	if key, err = loadKeyFixture(); err != nil {
		return "", err
	}

	var data []byte
	if data, err = ioutil.ReadFile("testdata/token_claims_template.json"); err != nil {
		return "", err
	}

	claims := make(jwt.MapClaims)
	if err = json.Unmarshal(data, &claims); err != nil {
		return "", err
	}
	claims["iss"] = iss + "/"
	claims["aud"] = []string{iss}
	claims["iat"] = jwt.NewNumericDate(time.Now())
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(10 * time.Minute))

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	return token.SignedString(key)
}
