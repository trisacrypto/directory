package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
)

func TestClient(t *testing.T) {
	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			require.Equal(t, int64(0), r.ContentLength)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "{\"hello\":\"world\"}")
			return
		}

		require.Equal(t, int64(18), r.ContentLength)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "{\"error\":\"bad request\"}")
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Ensure that the latest version of the client is returned
	bffv1, ok := client.(*api.APIv1)
	require.True(t, ok)

	// Create a new GET request to a basic path
	req, err := bffv1.NewRequest(context.TODO(), http.MethodGet, "/foo", nil, nil)
	require.NoError(t, err)

	require.Equal(t, "/foo", req.URL.Path)
	require.Equal(t, "", req.URL.RawQuery)
	require.Equal(t, http.MethodGet, req.Method)
	require.Equal(t, "GDS BFF API Client/v1", req.Header.Get("User-Agent"))
	require.Equal(t, "application/json", req.Header.Get("Accept"))
	require.Equal(t, "application/json; charset=utf-8", req.Header.Get("Content-Type"))

	// Create a new GET request with query params
	params := url.Values{}
	params.Add("q", "searching")
	params.Add("key", "open says me")
	req, err = bffv1.NewRequest(context.TODO(), http.MethodGet, "/foo", nil, &params)
	require.NoError(t, err)
	require.Equal(t, "key=open+says+me&q=searching", req.URL.RawQuery)

	data := make(map[string]string)
	rep, err := bffv1.Do(req, &data, true)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rep.StatusCode)
	require.Contains(t, data, "hello")
	require.Equal(t, "world", data["hello"])

	// Create a new POST request and check error handling
	req, err = bffv1.NewRequest(context.TODO(), http.MethodPost, "/bar", data, nil)
	require.NoError(t, err)
	rep, err = bffv1.Do(req, nil, false)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rep.StatusCode)

	req, err = bffv1.NewRequest(context.TODO(), http.MethodPost, "/bar", data, nil)
	require.NoError(t, err)
	_, err = bffv1.Do(req, nil, true)
	require.EqualError(t, err, "[400] bad request")
}

func TestStatus(t *testing.T) {
	fixture := &api.StatusReply{
		Status:  "ok",
		Uptime:  (2 * time.Second).String(),
		Version: "1.0.test",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/status", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Test with nil params
	out, err := client.Status(context.TODO(), nil)
	require.NoError(t, err)
	require.Equal(t, fixture.Status, out.Status)
	require.Equal(t, fixture.Uptime, out.Uptime)
	require.Equal(t, fixture.Version, out.Version)

	// Test with params
	out, err = client.Status(context.TODO(), &api.StatusParams{NoGDS: true})
	require.NoError(t, err)
	require.Equal(t, fixture.Status, out.Status)
	require.Equal(t, fixture.Uptime, out.Uptime)
	require.Equal(t, fixture.Version, out.Version)
}

func TestLookup(t *testing.T) {
	fixture := &api.LookupReply{
		Results: []map[string]interface{}{
			{"foo": "1", "color": "red"},
			{"foo": "2", "color": "blue"},
		},
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/lookup", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.Lookup(context.TODO(), &api.LookupParams{CommonName: "example.com"})
	require.NoError(t, err)
	require.Equal(t, fixture.Results, out.Results)
}

func TestRegister(t *testing.T) {
	fixture := &api.RegisterReply{
		Id:                  "8b2e9e78-baca-4c34-a382-8b285503c901",
		RegisteredDirectory: "vaspdirectory.net",
		CommonName:          "trisa.example.com",
		Status:              "PENDING_REVIEW",
		Message:             "Thank you for registering",
		PKCS12Password:      "supersecret squirrel",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/register/mainnet", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	req := &api.RegisterRequest{
		Network:          "MainNet",
		TRISAEndpoint:    "trisa.example.com:443",
		CommonName:       "trisa.example.com",
		Website:          "https://example.com",
		BusinessCategory: "PRIVATE_ORGANIZATION",
		VASPCategories:   []string{"ATM", "Other"},
		EstablishedOn:    "2019-01-14",
	}

	out, err := client.Register(context.TODO(), req)
	require.NoError(t, err)
	require.Equal(t, fixture.Id, out.Id)
	require.Equal(t, fixture.RegisteredDirectory, out.RegisteredDirectory)
	require.Equal(t, fixture.CommonName, out.CommonName)
	require.Equal(t, fixture.Status, out.Status)
	require.Equal(t, fixture.Message, out.Message)
	require.Equal(t, fixture.PKCS12Password, out.PKCS12Password)
}

func TestVerifyEmail(t *testing.T) {
	fixture := &api.VerifyContactReply{
		Status:  "PENDING_REVIEW",
		Message: "thank you for verifying your email",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/verify", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.VerifyContact(context.TODO(), &api.VerifyContactParams{Directory: "trisatest.net", ID: "foo", Token: "bar"})
	require.NoError(t, err)
	require.Equal(t, fixture.Status, out.Status)
	require.Equal(t, fixture.Message, out.Message)
}
