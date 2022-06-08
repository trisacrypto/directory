package auth0_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/auth0"
)

func TestLive(t *testing.T) {
	// These tests will only run if there is a valid configuration in the environment
	conf, err := auth0.NewConfig()
	if err != nil {
		t.Skip("live tests require local environment configuration")
	}

	// Do not run the live tests if there is no access token cacheing
	if conf.TokenCache == "" {
		t.Skip("live tests require a token cache to prevent issuing multiple M2M tokens")
	}

	//  Log the situation for the tests
	t.Logf("live tests starting with auth0 client %s, using token cache %s", conf.ClientID, conf.TokenCache)

	// Create the client to start the live tests
	client, err := auth0.New(conf)
	require.NoError(t, err, "could not create auth0 client for live testing")

	// TODO: don't call authenticate directly
	err = client.Authenticate()
	require.NoError(t, err, "could not authenticate the request")

}

func TestEndpoint(t *testing.T) {
	client, err := auth0.New(auth0.Config{Domain: "example.auth0.com", Testing: true})
	require.NoError(t, err, "could not create testing client")

	testCases := []struct {
		path     string
		query    map[string]string
		sfa      []interface{}
		expected string
	}{
		{"/v2/users", nil, nil, "http://example.auth0.com/v2/users"},
		{"/v2/users", map[string]string{"sort": "created"}, nil, "http://example.auth0.com/v2/users?sort=created"},
		{"/v2/users", map[string]string{"sort": "created", "impending": "doom"}, nil, "http://example.auth0.com/v2/users?impending=doom&sort=created"},
		{"/v2/users/%d", nil, []interface{}{123}, "http://example.auth0.com/v2/users/123"},
		{"/v2/users/%d", map[string]string{"sort": "created", "impending": "doom"}, []interface{}{123}, "http://example.auth0.com/v2/users/123?impending=doom&sort=created"},
	}

	for i, tc := range testCases {
		endpoint := client.Endpoint(tc.path, tc.query, tc.sfa...)
		require.Equal(t, tc.expected, endpoint, "unexpected endpoint for test case %d", i+1)
	}
}

func TestDo(t *testing.T) {
	// Create a test mux for various request paths
	mux := http.NewServeMux()
	mux.HandleFunc("/valid", func(w http.ResponseWriter, r *http.Request) {
		// Returns a valid 200 response
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"success": true}`)
	})
	mux.HandleFunc("/invalid", func(w http.ResponseWriter, r *http.Request) {
		// Returns an invalid 401 response
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, `{"statusCode": 401, "error": "Unauthorized", "message": "Missing authentication"}`)
	})
	mux.HandleFunc("/crash", func(w http.ResponseWriter, r *http.Request) {
		// Crashes the server returning text/html instead of JSON
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Error")
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	// Get the hostname of the test server
	testURL, _ := url.Parse(srv.URL)

	// Create the auth0 client for testing
	client, err := auth0.New(auth0.Config{Domain: testURL.Host, Testing: true})
	require.NoError(t, err, "could not create auth0 client")

	// A valid request should return a response with an open body
	req, err := http.NewRequest(http.MethodGet, client.Endpoint("/valid", nil), nil)
	require.NoError(t, err, "could not create request to /valid")

	// Execute valid request
	rep, err := client.Do(req)
	require.NoError(t, err, "could note execute request to /valid")
	require.Equal(t, http.StatusOK, rep.StatusCode)
	data := make(map[string]interface{})
	require.NoError(t, json.NewDecoder(rep.Body).Decode(&data), "response body not open or json not decodable")
	require.Contains(t, data, "success", "did not receive a correct response")

	// An invalid request should return an API error with a closed body
	req, err = http.NewRequest(http.MethodGet, client.Endpoint("/invalid", nil), nil)
	require.NoError(t, err, "could not create request to /invalid")

	// Execute invalid request
	_, err = client.Do(req)
	require.Error(t, err, "expected an error response from test server")
	require.EqualError(t, err, "[401] Unauthorized: Missing authentication")
	require.IsType(t, &auth0.APIError{}, err, "expected an API error back from the server")

	// An crash request should return an API error even though JSON cannot be parsed
	req, err = http.NewRequest(http.MethodGet, client.Endpoint("/crash", nil), nil)
	require.NoError(t, err, "could not create request to /crash")

	// Execute crash request
	_, err = client.Do(req)
	require.Error(t, err, "expected an error response from test server")
	require.EqualError(t, err, "[500] 500 Internal Server Error")
	require.IsType(t, &auth0.APIError{}, err, "expected an API error back from the server")

}

func TestDoProtectTesting(t *testing.T) {
	// In testing mode, the client should only allow localhost and 127.0.0.1
	client, err := auth0.New(auth0.Config{Domain: "example.auth0.com", Testing: true})
	require.NoError(t, err, "could not create auth0 client")

	req, err := http.NewRequest(http.MethodGet, client.Endpoint("/", nil), nil)
	require.NoError(t, err, "could not create http request")

	_, err = client.Do(req)
	require.EqualError(t, err, `hostname "example.auth0.com" is not valid in testing mode`)
}
