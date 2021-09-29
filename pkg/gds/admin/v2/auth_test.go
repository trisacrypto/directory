package admin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
)

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
