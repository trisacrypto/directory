package admin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

func TestAuthorization(t *testing.T) {
	// Test Authorization middleware
	router := gin.New()
	tm, err := tokens.MockTokenManager()
	require.NoError(t, err)

	// Create secure endpoint for testing purposes
	router.GET("/", admin.Authorization(tm), func(c *gin.Context) {
		claims, exists := c.Get(admin.UserClaims)
		require.True(t, exists)
		require.NotEmpty(t, claims)
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Create the test server with the authorized router
	server := httptest.NewServer(router)
	defer server.Close()

	// Create a request to the server that is not authorized
	req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)
	require.NoError(t, err)

	// Execute an unauthorized request and check a 401 response
	rep, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rep.StatusCode)

	// Read the body of the response
	data, err := readJSON(rep)
	require.NoError(t, err)
	require.Contains(t, data, "error")
	require.Equal(t, "a valid authorization is required to access this endpoint", data["error"].(string))

	// Create an access token to authorize the request
	creds := map[string]interface{}{
		"hd":      "rotational.io",
		"email":   "kate@rotational.io",
		"name":    "Kate Holland",
		"picture": "https://foo.googleusercontent.com/test!/Aoh14gJceTrUA",
	}

	accessToken, err := tm.CreateAccessToken(creds)
	require.NoError(t, err)

	// Create signed token
	tks, err := tm.Sign(accessToken)
	require.NoError(t, err)

	// Create a request to the server that is authorized
	req, err = http.NewRequest(http.MethodGet, server.URL+"/", nil)
	req.Header.Set("Authorization", "Bearer "+tks)
	require.NoError(t, err)

	// Execute an unauthorized request and check a 200 response
	rep, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rep.StatusCode)

	// Read the body of the response
	data, err = readJSON(rep)
	require.NoError(t, err)
	require.NotContains(t, data, "error")
	require.Contains(t, data, "success")

	// Create a request to the server with an invalid access token by messing with its signature
	req, err = http.NewRequest(http.MethodGet, server.URL+"/", nil)
	req.Header.Set("Authorization", "Bearer "+tks+"scrambled")
	require.NoError(t, err)

	// Execute an unauthorized request and check a 401 response
	rep, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rep.StatusCode)

	// Read the body of the response
	data, err = readJSON(rep)
	require.NoError(t, err)
	require.Contains(t, data, "error")
	require.Contains(t, data, "success")
	require.Equal(t, "a valid authorization is required to access this endpoint", data["error"].(string))
}

func TestGetAccessToken(t *testing.T) {
	// Create a gin router that constructs a client from the context
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		// Return whatever the GetAccessToken function returns
		token, err := admin.GetAccessToken(c)
		var errs string
		if err != nil {
			errs = err.Error()
		}

		c.JSON(http.StatusOK, gin.H{"token": token, "err": errs})
	})

	// Create a test server from the router
	server := httptest.NewServer(router)
	defer server.Close()

	// Test a request with no authorization header
	req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)
	require.NoError(t, err)
	rep, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	data, err := readJSON(rep)
	require.NoError(t, err)

	require.Empty(t, data["token"])
	require.Equal(t, "no access token found in request", data["err"])

	// Test a request with no bearer token header
	req, err = http.NewRequest(http.MethodGet, server.URL+"/", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "notabearertoken")

	rep, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	data, err = readJSON(rep)
	require.NoError(t, err)

	require.Empty(t, data["token"])
	require.Equal(t, "could not parser Bearer token from Authorization header", data["err"])

	// Test a request with a good bearer token header
	req, err = http.NewRequest(http.MethodGet, server.URL+"/", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer thisisyourtoken")

	rep, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	data, err = readJSON(rep)
	require.NoError(t, err)

	require.Empty(t, data["err"])
	require.Equal(t, "thisisyourtoken", data["token"])

}
