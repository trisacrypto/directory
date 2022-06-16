package auth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/config"
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
	c.Set(auth.ContextRegisteredClaims, &vclaims.RegisteredClaims)

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
	// Create context, gin.Engine, and http test writer to execute tests
	c, srv, w := createTestContext(http.MethodGet, "/", nil, authenticate, success)
	rep, code, err := doRequest(srv, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusOK, code)
	require.Contains(t, rep, "success", "response does not contain a success field")
	require.True(t, rep["success"].(bool), "success is not true")

	claims, err := auth.GetClaims(c)
	require.NoError(t, err, "expected anonymous claims on context")
	require.True(t, claims.IsAnonymous(), "expected anonymous claims on context")

	// Test non-bearer token
	// Create context, gin.Engine, and http test writer to execute tests
	c, srv, w = createTestContext(http.MethodGet, "/", nil, authenticate, success)
	c.Request.Header.Set("authorization", "token foo")
	_, code, err = doRequest(srv, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusOK, code)

	claims, err = auth.GetClaims(c)
	require.NoError(t, err, "expected anonymous claims on context")
	require.True(t, claims.IsAnonymous(), "expected anonymous claims on context")

	// Test forbidden error with incorrectly signed token
	token, err := ioutil.ReadFile("testdata/invalid_token.txt")
	require.NoError(t, err, "could not read invalid token fixture")

	// Create context, gin.Engine, and http test writer to execute tests
	c, srv, w = createTestContext(http.MethodGet, "/", nil, authenticate, success)
	c.Request.Header.Set("authorization", "Bearer "+string(token))
	rep, code, err = doRequest(srv, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusForbidden, code)
	require.Contains(t, rep, "error", "expected error on JSON response")
	require.Equal(t, "invalid authorization token", rep["error"])
}

func TestAuthenticatePublicKeys(t *testing.T) {
	// Creates a test server that serves well known jwks keys instead of the Auth0
	// tenant - used to mock Auth0 (not a live test) but checks the happy path.
	// NOTE: this test is fragile, e.g. if Auth0 changes its implementation.
	srv, err := authtest.New()
	require.NoError(t, err, "could not create authtest server")
	defer srv.Close()

	// Setup authentication middleware
	authenticate, err := auth.Authenticate(srv.Config(), auth.WithHTTPClient(srv.Client()))
	require.NoError(t, err, "expected valid authenticate middleware")

	// Create default handler
	success := func(c *gin.Context) {
		c.JSON(http.StatusOK, api.Reply{Success: true})
	}

	// Create a valid token to authenticate
	tks, err := srv.NewToken()
	require.NoError(t, err, "could not create valid token")

	// Execute request expecting success
	c, mux, w := createTestContext(http.MethodGet, "/", nil, authenticate, success)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tks))
	rep, code, err := doRequest(mux, w, c)

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
	// Create context, gin.Engine, and http test writer to execute tests
	c, srv, w := createTestContext(http.MethodGet, "/", nil, authorize, success)
	rep, code, err := doRequest(srv, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "could not authorize request", rep["error"], "unexpected error returned from authorize")

	// Test anonymous user on context
	authenticate := func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.AnonymousClaims)
		c.Next()
	}

	// Create context, gin.Engine, and http test writer to execute tests
	c, srv, w = createTestContext(http.MethodGet, "/", nil, authenticate, authorize, success)
	rep, code, err = doRequest(srv, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "this endpoint requires authentication", rep["error"], "unexpected error returned from authorize")

	// Test user does not have permissions
	authenticate = func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.Claims{Permissions: []string{"write:foo"}})
		c.Next()
	}

	// Create context, gin.Engine, and http test writer to execute tests
	c, srv, w = createTestContext(http.MethodGet, "/", nil, authenticate, authorize, success)
	rep, code, err = doRequest(srv, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "user does not have permission to perform this operation", rep["error"], "unexpected error returned from authorize")

	// Test user does not have all permissions
	authenticate = func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.Claims{Permissions: []string{"read:foo"}})
		c.Next()
	}

	// Create context, gin.Engine, and http test writer to execute tests
	c, srv, w = createTestContext(http.MethodGet, "/", nil, authenticate, authorize, success)
	rep, code, err = doRequest(srv, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "user does not have permission to perform this operation", rep["error"], "unexpected error returned from authorize")

	// Test user does have permissions
	authenticate = func(c *gin.Context) {
		c.Set(auth.ContextBFFClaims, &auth.Claims{Permissions: []string{"read:foo", "read:bar"}})
		c.Next()
	}

	// Create context, gin.Engine, and http test writer to execute tests
	c, srv, w = createTestContext(http.MethodGet, "/", nil, authenticate, authorize, success)
	rep, code, err = doRequest(srv, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusOK, code)
	require.Contains(t, rep, "success", "response does not contain a success field")
	require.True(t, rep["success"].(bool), "success is not true")
}

func TestUserInfo(t *testing.T) {
	t.Skip("cannot implement these tests unless auth0 is mocked, which will likely be very difficult")

	// Setup UserInfo middleware
	middleware, err := auth.UserInfo(config.AuthConfig{Domain: "example.auth0.com", ClientID: "example", ClientSecret: "supersecretsquirrel"})
	require.NoError(t, err, "could not create user info middleware")

	// Create default handler
	success := func(c *gin.Context) {
		c.JSON(http.StatusOK, api.Reply{Success: true})
	}

	// Test userinfo no claims on context
	// Create context, gin.Engine, and http test writer to execute tests
	c, srv, w := createTestContext(http.MethodGet, "/", nil, middleware, success)
	rep, code, err := doRequest(srv, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "could not authorize request", rep["error"], "unexpected error returned from middleware")

	// Test claims without subject on the context
	authenticate := func(c *gin.Context) {
		c.Set(auth.ContextRegisteredClaims, &validator.RegisteredClaims{Subject: ""})
		c.Next()
	}

	// Create context, gin.Engine, and http test writer to execute tests
	c, srv, w = createTestContext(http.MethodGet, "/", nil, authenticate, middleware, success)
	rep, code, err = doRequest(srv, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusUnauthorized, code)
	require.Contains(t, rep, "error", "response does not contain json error")
	require.Equal(t, "this endpoint requires authentication", rep["error"], "unexpected error returned from authorize")

	// Test user does have registered claims on context
	authenticate = func(c *gin.Context) {
		c.Set(auth.ContextRegisteredClaims, &validator.RegisteredClaims{Subject: "test|1234567890abcdefg"})
		c.Next()
	}

	// Create context, gin.Engine, and http test writer to execute tests
	c, srv, w = createTestContext(http.MethodGet, "/", nil, authenticate, middleware, success)
	rep, code, err = doRequest(srv, w, c)
	require.NoError(t, err, "could not handle test request")
	require.Equal(t, http.StatusOK, code)
	require.Contains(t, rep, "success", "response does not contain a success field")
	require.True(t, rep["success"].(bool), "success is not true")

	// TODO: check that the user info on the context is as expected.
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
		return nil, 0, err
	}
	return data, rep.StatusCode, nil
}
