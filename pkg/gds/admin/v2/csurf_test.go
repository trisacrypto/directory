package admin_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
)

func TestDoubleCookies(t *testing.T) {
	// Test both the DoubleCookie middleware and the SetDoubleCookieTokens handler.
	router := gin.New()

	// Add login route that sets the cookies
	router.GET("/login", func(c *gin.Context) {
		err := admin.SetDoubleCookieTokens(c, time.Now().Add(time.Minute*10).Unix())
		require.NoError(t, err)
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Add request route that requires double cookie submit
	router.POST("/action", admin.DoubleCookie(), func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"success": true})
	})

	// Create an http client with a cookie jar
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := http.Client{
		Jar: jar,
	}

	// Create the test server with the CSRF protected router
	server := httptest.NewServer(router)

	// Attempt to make a request to the server that is not CSRF protected
	req, err := http.NewRequest(http.MethodPost, server.URL+"/action", nil)
	require.NoError(t, err)

	// Ensure the request is Forbidden
	rep, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, rep.StatusCode)

	// Check the data in the response
	data, err := readJSON(rep)
	require.NoError(t, err)
	require.Contains(t, data, "error")
	require.Contains(t, data, "success")
	require.Equal(t, admin.ErrCSRFVerification.Error(), data["error"].(string))
	require.False(t, data["success"].(bool))

	// Login and set the cookies
	req, err = http.NewRequest(http.MethodGet, server.URL+"/login", nil)
	require.NoError(t, err)

	// Ensure the request is Ok
	rep, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rep.StatusCode)

	// Check that we're getting two cookies back in the response
	cookies := rep.Cookies()
	require.Len(t, cookies, 2)

	// Check the data in the response
	data, err = readJSON(rep)
	require.NoError(t, err)
	require.NotContains(t, data, "error")
	require.Contains(t, data, "success")
	require.True(t, data["success"].(bool))

	// TESTING HACK: add the secure cookies to the client cookie jar even though this
	// is not TLS. We're interested in the DoubleCookie middleware, not whether or not
	// the http client does the right thing with cookies.
	ep, _ := url.Parse(server.URL)
	for _, cookie := range cookies {
		cookie.Secure = false
	}
	client.Jar.SetCookies(ep, cookies)

	// Attempt to send a request with the cookies but no X-CSRF-TOKEN header
	req, err = http.NewRequest(http.MethodPost, server.URL+"/action", nil)
	require.NoError(t, err)

	// Ensure the request is Forbidden
	rep, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, rep.StatusCode)

	// Check the data in the response
	data, err = readJSON(rep)
	require.NoError(t, err)
	require.Contains(t, data, "error")
	require.Contains(t, data, "success")
	require.Equal(t, admin.ErrCSRFVerification.Error(), data["error"].(string))
	require.False(t, data["success"].(bool))

	// Send a request with the cookies but an incorrect X-CSRF-TOKEN header
	req, err = http.NewRequest(http.MethodPost, server.URL+"/action", nil)
	req.Header.Set(admin.CSRFHeader, "foo")
	require.NoError(t, err)

	// Ensure the request is Forbidden
	rep, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, rep.StatusCode)

	// Check the data in the response
	data, err = readJSON(rep)
	require.NoError(t, err)
	require.Contains(t, data, "error")
	require.Contains(t, data, "success")
	require.Equal(t, admin.ErrCSRFVerification.Error(), data["error"].(string))
	require.False(t, data["success"].(bool))

	var cookieToken string
	for _, cookie := range cookies {
		if cookie.Name == admin.CSRFCookie {
			cookieToken = cookie.Value
		}
	}

	// Finally, send a request with the cookies but a valid X-CSRF-TOKEN header
	req, err = http.NewRequest(http.MethodPost, server.URL+"/action", nil)
	req.Header.Set(admin.CSRFHeader, cookieToken)
	require.NoError(t, err)

	// Ensure the request is Ok
	rep, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rep.StatusCode)

	// Check the data in the response
	data, err = readJSON(rep)
	require.NoError(t, err)
	require.NotContains(t, data, "error")
	require.Contains(t, data, "success")
	require.True(t, data["success"].(bool))

}

func readJSON(rep *http.Response) (map[string]interface{}, error) {
	defer rep.Body.Close()
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func TestGenerateCSRFToken(t *testing.T) {
	token1, err := admin.GenerateCSRFToken()
	require.NoError(t, err)
	require.Len(t, token1, 44)

	token2, err := admin.GenerateCSRFToken()
	require.NoError(t, err)

	require.NotEqual(t, token1, token2)
}
