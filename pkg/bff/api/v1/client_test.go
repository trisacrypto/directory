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
		TestNet: map[string]interface{}{"foo": "2", "color": "blue"},
		MainNet: map[string]interface{}{"foo": "1", "color": "red"},
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
	require.Equal(t, fixture.TestNet, out.TestNet)
	require.Equal(t, fixture.MainNet, out.MainNet)
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

		in := &api.RegisterRequest{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode register request")

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

func TestVerifyContact(t *testing.T) {
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

func TestOverview(t *testing.T) {
	fixture := &api.OverviewReply{
		OrgID: "ba2202bf-635e-414e-a7bc-86f309dc95e0",
		TestNet: api.NetworkOverview{
			Status:             "online",
			Vasps:              8,
			CertificatesIssued: 7,
			NewMembers:         3,
			MemberDetails: api.MemberDetails{
				ID:          "8b2e9e78-baca-4c34-a382-8b285503c901",
				Status:      "VERIFIED",
				CountryCode: "FK",
				Certificate: map[string]interface{}{
					"common_name": "trisa.example.com",
				},
			},
		},
		MainNet: api.NetworkOverview{
			Status:             "pending",
			Vasps:              12,
			CertificatesIssued: 21,
			NewMembers:         5,
			MemberDetails: api.MemberDetails{
				ID:          "c34c9e78-baca-4c34-a382-8b285503c901",
				Status:      "SUBMITTED",
				CountryCode: "FK",
				Certificate: map[string]interface{}{
					"common_name": "trisa.example.com",
				},
			},
		},
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/overview", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.Overview(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture, out)
	require.Equal(t, fixture.OrgID, out.OrgID)
	require.Equal(t, fixture.TestNet.Status, out.TestNet.Status)
	require.Equal(t, fixture.TestNet.MemberDetails, out.TestNet.MemberDetails)
	require.Equal(t, fixture.MainNet.CertificatesIssued, out.MainNet.CertificatesIssued)
	require.Equal(t, fixture.MainNet.MemberDetails, out.MainNet.MemberDetails)
}

func TestCertificates(t *testing.T) {
	fixture := &api.CertificatesReply{
		TestNet: []api.Certificate{
			{
				SerialNumber: "ABC83132333435363738",
				IssuedAt:     time.Now().AddDate(-1, -1, 0).Format(time.RFC3339),
				ExpiresAt:    time.Now().AddDate(0, -1, 0).Format(time.RFC3339),
				Revoked:      true,
				Details: map[string]interface{}{
					"common_name": "trisa.example.com",
				},
			},
		},
		MainNet: []api.Certificate{
			{
				SerialNumber: "DEF83132333435363738",
				IssuedAt:     time.Now().Format(time.RFC3339),
				ExpiresAt:    time.Now().AddDate(1, 0, 0).Format(time.RFC3339),
				Revoked:      false,
				Details: map[string]interface{}{
					"common_name": "trisa.example.com",
				},
			},
		},
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/certificates", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.Certificates(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture, out)
	require.Equal(t, fixture.TestNet, out.TestNet)
	require.Equal(t, fixture.MainNet, out.MainNet)
}

func TestAnnoucements(t *testing.T) {
	fixture := &api.AnnouncementsReply{
		Announcements: []*api.Announcement{
			{
				Title:    "Upcoming TRISA Working Group Call",
				Body:     "Join us on Thursday Apr 28 for the TRISA Working Group.",
				PostDate: "2022-04-20",
				Author:   "admin@trisa.io",
			},
			{
				Title:    "Routine Maintenance Scheduled",
				Body:     "The GDS will be undergoing routine maintenance on Apr 7.",
				PostDate: "2022-04-01",
				Author:   "admin@trisa.io",
			},
			{
				Title:    "Beware the Ides of March",
				Body:     "I have a bad feeling about tomorrow.",
				PostDate: "2022-03-14",
				Author:   "julius@caesar.com",
			},
		},
		LastUpdated: "2022-04-21T12:05:23Z",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/announcements", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.Announcements(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture, out)
	require.Len(t, out.Announcements, 3)
	require.Equal(t, "2022-04-21T12:05:23Z", out.LastUpdated)
}

func TestMakeAnnoucement(t *testing.T) {
	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/announcements", r.URL.Path)

		in := &api.Announcement{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode register request")

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	req := &api.Announcement{
		Title: "The Happenings",
		Body:  "Things are going on, we're all very busy, and you should join us!",
	}

	err = client.MakeAnnouncement(context.TODO(), req)
	require.NoError(t, err)
}

func TestMakeAnnoucementErrors(t *testing.T) {
	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/announcements", r.URL.Path)

		in := &api.Announcement{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode register request")

		switch in.Title {
		case "200":
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
		case "400":
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
		}

	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	req := &api.Announcement{Title: "200"}
	err = client.MakeAnnouncement(context.TODO(), req)
	require.EqualError(t, err, "expected no content, received 200 OK")

	req = &api.Announcement{Title: "400"}
	err = client.MakeAnnouncement(context.TODO(), req)
	require.EqualError(t, err, "400 Bad Request")
}
