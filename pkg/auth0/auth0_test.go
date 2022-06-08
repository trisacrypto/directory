package auth0_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/auth0"
)

func TestAuthenticate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the endpoint and method
		if r.Method != http.MethodPost {
			WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("expected POST got %q", r.Method))
			return
		}

		if r.URL.Path != "/oauth/token" {
			WriteError(w, http.StatusNotFound, fmt.Errorf("expected /oauth/token got %q", r.URL.Path))
			return
		}

		// Confirm the header
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			WriteError(w, http.StatusUnsupportedMediaType, "unexpected content-type header")
			return
		}

		if err := r.ParseForm(); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		// Confirm the data that was sent in the request
		if !r.PostForm.Has("grant_type") || r.PostForm.Get("grant_type") != "client_credentials" {
			WriteError(w, http.StatusBadRequest, "missing or incorrect grant_type")
			return
		}

		if !r.PostForm.Has("client_id") || r.PostForm.Get("client_id") != "hello" {
			WriteError(w, http.StatusBadRequest, "missing or incorrect client_id")
			return
		}

		if !r.PostForm.Has("client_secret") || r.PostForm.Get("client_secret") != "world" {
			WriteError(w, http.StatusBadRequest, "missing or incorrect client_secret")
			return
		}

		if !r.PostForm.Has("audience") || !strings.HasSuffix(r.PostForm.Get("audience"), "/api/v2/") {
			WriteError(w, http.StatusBadRequest, "missing or incorrect audience")
			return
		}

		// Collect the response to write
		f, err := os.Open("testdata/example_token.json")
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}
		defer f.Close()

		// Everything is fine, write 200 response with token
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
	}))
	defer srv.Close()

	// Create the auth0 client for testing
	testURL, _ := url.Parse(srv.URL)
	client, err := auth0.New(auth0.Config{Domain: testURL.Host, ClientID: "hello", ClientSecret: "world", Testing: true})
	require.NoError(t, err, "could not create auth0 client connecting to test server")

	// Check that the credentials are zero-valued to start
	require.Empty(t, client.Creds(), "expected zero-valued credentials for test")

	// Execute valid authenticate request
	err = client.Authenticate()
	require.NoError(t, err, "could not authenticate client")

	// Credentials should be not empty and valued after authentication
	creds := client.Creds()
	require.NotEmpty(t, creds, "credentials zero-valued after successful authentication?")
	require.True(t, creds.Valid(), "credentials should be valid after successful authentication")
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
		WriteError(w, http.StatusInternalServerError, "Internal Error")
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

// WriteError is a helper method to write error responses in the test server.
func WriteError(w http.ResponseWriter, statusCode int, err interface{}) {
	switch e := err.(type) {
	case *auth0.APIError:
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(e)
	case error:
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, e.Error())
	case string:
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, e)
	default:
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(e)
	}
}
