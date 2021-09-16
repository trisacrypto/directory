package admin_test

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
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
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
	client, err := admin.New(ts.URL)
	require.NoError(t, err)

	// Ensure that the latest version of the client is returned
	apiv2, ok := client.(*admin.APIv2)
	require.True(t, ok)

	// Create a new GET request to a basic path
	req, err := apiv2.NewRequest(context.TODO(), http.MethodGet, "/foo", nil, nil)
	require.NoError(t, err)

	require.Equal(t, "/foo", req.URL.Path)
	require.Equal(t, "", req.URL.RawQuery)
	require.Equal(t, http.MethodGet, req.Method)
	require.Equal(t, "GDS Admin API Client/v2", req.Header.Get("User-Agent"))
	require.Equal(t, "application/json", req.Header.Get("Accept"))
	require.Equal(t, "application/json; charset=utf-8", req.Header.Get("Content-Type"))

	// Create a new GET request with query params
	params := url.Values{}
	params.Add("q", "searching")
	params.Add("key", "open says me")
	req, err = apiv2.NewRequest(context.TODO(), http.MethodGet, "/foo", nil, &params)
	require.NoError(t, err)
	require.Equal(t, "key=open+says+me&q=searching", req.URL.RawQuery)

	data := make(map[string]string)
	rep, err := apiv2.Do(req, &data, true)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rep.StatusCode)
	require.Contains(t, data, "hello")
	require.Equal(t, "world", data["hello"])

	// Create a new POST request and check error handling
	req, err = apiv2.NewRequest(context.TODO(), http.MethodPost, "/bar", data, nil)
	require.NoError(t, err)
	rep, err = apiv2.Do(req, nil, false)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rep.StatusCode)

	req, err = apiv2.NewRequest(context.TODO(), http.MethodPost, "/bar", data, nil)
	require.NoError(t, err)
	_, err = apiv2.Do(req, nil, true)
	require.EqualError(t, err, "[400] bad request")
}

func TestStatus(t *testing.T) {
	fixture := &admin.StatusReply{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.test",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v2/status", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := admin.New(ts.URL)
	require.NoError(t, err)

	out, err := client.Status(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture.Status, out.Status)
	require.True(t, fixture.Timestamp.Equal(out.Timestamp))
	require.Equal(t, fixture.Version, out.Version)
}

func TestListVASPs(t *testing.T) {
	fixture := &admin.ListVASPsReply{
		VASPs: []admin.VASPSnippet{
			{
				ID:                 "af367d27-b0e7-48b5-8987-e48a0712a826",
				Name:               "Alice VASP",
				CommonName:         "trisa.alice.us",
				VerificationStatus: "verified",
				LastUpdated:        "2021-08-15T12:32:41Z",
				Traveler:           false,
				VerifiedContacts:   []string{"administrative", "technical"},
			},
			{
				ID:                 "5a26150d-ac6b-4bc8-973f-9065b815286c",
				Name:               "Bob VASP",
				CommonName:         "trisa.bob.co.uk",
				VerificationStatus: "pending review",
				LastUpdated:        "2021-09-11T22:02:39Z",
				Traveler:           false,
				VerifiedContacts:   []string{"billing", "technical"},
			},
		},
		Page:     2,
		PageSize: 10,
		Count:    12,
	}

	params := &admin.ListVASPsParams{
		Page:     2,
		PageSize: 10,
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v2/vasps", r.URL.Path)
		require.Equal(t, "page=2&page_size=10", r.URL.RawQuery)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := admin.New(ts.URL)
	require.NoError(t, err)

	out, err := client.ListVASPs(context.TODO(), params)
	require.NoError(t, err)
	require.Equal(t, fixture.VASPs, out.VASPs)
	require.Equal(t, fixture.Page, out.Page)
	require.Equal(t, fixture.PageSize, out.PageSize)
	require.Equal(t, fixture.Count, out.Count)
}

func TestReview(t *testing.T) {
	fixture := &admin.ReviewReply{
		Status:  "reviewed",
		Message: "the message has been received and the registration updated",
	}

	req := &admin.ReviewRequest{
		ID:                     "1234",
		AdminVerificationToken: "foo",
		Accept:                 true,
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v2/vasps/1234/review", r.URL.Path)

		// Must be able to deserialize the request
		in := new(admin.ReviewRequest)
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err)

		require.Equal(t, req.ID, in.ID)
		require.Equal(t, req.AdminVerificationToken, in.AdminVerificationToken)
		require.Equal(t, req.Accept, in.Accept)
		require.Equal(t, req.RejectReason, in.RejectReason)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := admin.New(ts.URL)
	require.NoError(t, err)

	out, err := client.Review(context.TODO(), req)
	require.NoError(t, err)
	require.Equal(t, fixture.Status, out.Status)
	require.Equal(t, fixture.Message, out.Message)
}

func TestResend(t *testing.T) {
	fixture := &admin.ResendReply{
		Sent:    3,
		Message: "the certificates were successfully redelivered",
	}

	req := &admin.ResendRequest{
		ID:     "1234",
		Action: admin.ResendDeliverCerts,
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v2/vasps/1234/resend", r.URL.Path)

		// Must be able to deserialize the request
		in := new(admin.ResendRequest)
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err)

		require.Equal(t, req.ID, in.ID)
		require.Equal(t, req.Action, in.Action)
		require.Equal(t, req.Reason, in.Reason)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := admin.New(ts.URL)
	require.NoError(t, err)

	out, err := client.Resend(context.TODO(), req)
	require.NoError(t, err)
	require.Equal(t, fixture.Sent, out.Sent)
	require.Equal(t, fixture.Message, out.Message)
}
