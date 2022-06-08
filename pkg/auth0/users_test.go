package auth0_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/auth0"
)

func TestGetUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the endpoint and method
		if r.Method != http.MethodGet {
			WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("expected GET got %q", r.Method))
			return
		}

		if r.URL.Path != "/api/v2/users/auth0|62a014c5881f6b006f97ed30" {
			WriteError(w, http.StatusNotFound, fmt.Errorf("unexpected path %q", r.URL.Path))
			return
		}

		// Collect the response to write
		f, err := os.Open("testdata/example_user.json")
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
	client, err := MakeAuthenticatedClient(srv.URL)
	require.NoError(t, err, "could not create auth0 client connecting to test server")

	user, err := client.GetUser(context.TODO(), "auth0|62a014c5881f6b006f97ed30")
	require.NoError(t, err, "could not execute get user request")
	require.NotEmpty(t, user, "user not unmarshaled correctly")
	require.Equal(t, user.Email, "leopold.wentzel@gmail.com")
}

func TestUpdateUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the endpoint and method
		if r.Method != http.MethodPatch {
			WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("expected PATCH got %q", r.Method))
			return
		}

		if r.URL.Path != "/api/v2/users/auth0|62a014c5881f6b006f97ed30" {
			WriteError(w, http.StatusNotFound, fmt.Errorf("unexpected path %q", r.URL.Path))
			return
		}

		// Collect the response to write (no update actually happens for testing)
		f, err := os.Open("testdata/example_user.json")
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
	client, err := MakeAuthenticatedClient(srv.URL)
	require.NoError(t, err, "could not create auth0 client connecting to test server")

	// Create updates to the user name field, but without an ID to start
	updates := &auth0.User{
		ID:   "auth0|62a014c5881f6b006f97ed30",
		Name: "Leopold Wentzel",
	}

	_, err = client.UpdateUser(context.TODO(), "", updates)
	require.EqualError(t, err, "[400] Invalid Request: A user ID is required to make an update request", "expected user ID required")

	_, err = client.UpdateUser(context.TODO(), "auth0|62a014c5881f6b006f97ed30", updates)
	require.EqualError(t, err, "[400] Invalid Update: A user ID cannot be specified on the updates sent to the server", "expected lightweight struct validation")

	// Execute update without user ID so that it should work
	updates.ID = ""
	user, err := client.UpdateUser(context.TODO(), "auth0|62a014c5881f6b006f97ed30", updates)
	require.NoError(t, err, "could not execute get user request")
	require.NotEmpty(t, user, "user not unmarshaled correctly")
	require.Equal(t, user.Email, "leopold.wentzel@gmail.com")
}
