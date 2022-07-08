package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
)

func TestCredentials(t *testing.T) {
	// Ensure that credentials can be passed to a client and used to make an
	// authenticated http request with the Authorization header.
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mimic an authenticated status endpoint
		token := r.Header.Get("Authorization")
		if matched, err := regexp.MatchString(`^Bearer eyJ.*$`, token); !matched || err != nil {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&api.Reply{Success: false, Error: fmt.Sprintf("%s", err)})
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&api.StatusReply{
			Status:  "ok",
			Uptime:  (2 * time.Second).String(),
			Version: "1.0.test",
		})
	}))
	defer ts.Close()

	var creds api.Credentials
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test a string-based access token
	creds = api.Token("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
	client, err := api.New(ts.URL, api.WithClient(ts.Client()), api.WithCredentials(creds))
	require.NoError(t, err, "unable to create an APIv1 client with token credentials")

	_, err = client.Status(ctx, &api.StatusParams{})
	require.NoError(t, err, "expected to be able to make an authenticated request with token credentials")

	// Test invalid string-based access token
	creds = api.Token("")
	client, err = api.New(ts.URL, api.WithClient(ts.Client()), api.WithCredentials(creds))
	require.NoError(t, err, "unable to create an APIv1 client with token credentials")

	// Expect that when the credentials are used they'll be validated
	_, err = client.Status(ctx, &api.StatusParams{})
	require.ErrorIs(t, err, api.ErrInvalidCredentials, "expected empty string token to be rejected")

	// Test expired credentials loaded from disk
	// NOTE: this will also test expired Auth0Token since it is used by LocalCredentials
	creds = &api.LocalCredentials{Path: "testdata/token.json"}
	client, err = api.New(ts.URL, api.WithClient(ts.Client()), api.WithCredentials(creds))
	require.NoError(t, err, "unable to create an APIv1 client with local credentials")

	// Expect that when the credentials are used they'll be loaded from disk and validated
	_, err = client.Status(ctx, &api.StatusParams{})
	require.ErrorIs(t, err, api.ErrExpiredCredentials, "expected expired credentials to be rejected")

	// Extract Auth0Token from local credentials set expires at to future time
	local, ok := creds.(*api.LocalCredentials)
	require.True(t, ok, "could not type assert local credentials")
	local.Token.ExpiresAt = time.Now().Add(10 * time.Minute)

	// Test valid credentials from local credentials
	// NOTE: this will also test Auth0Token since it is used by LocalCredentials
	client, err = api.New(ts.URL, api.WithClient(ts.Client()), api.WithCredentials(creds))
	require.NoError(t, err, "unable to create an APIv1 client with local credentials")

	_, err = client.Status(ctx, &api.StatusParams{})
	require.NoError(t, err, "expected to be able to make an authenticated request with local credentials")

	// Test dumping local credentials back to disk
	local.Path = filepath.Join(t.TempDir(), "token.json")
	require.NoError(t, local.Dump(), "could not dump local credentials back to tmp directory")
}
