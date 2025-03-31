package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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

func TestLookupAutocomplete(t *testing.T) {
	fixture := map[string]string{
		"Alice": "alice",
		"Bob":   "bob",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/lookup/autocomplete", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a new Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.LookupAutocomplete(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture, out)
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

	out, err := client.VerifyContact(context.TODO(), &api.VerifyContactParams{Directory: "testnet.directory", ID: "foo", Token: "bar"})
	require.NoError(t, err)
	require.Equal(t, fixture.Status, out.Status)
	require.Equal(t, fixture.Message, out.Message)
}

func TestLogin(t *testing.T) {
	// Test login with credentials and a TLS connection sets csrf protection
	var err error
	token := api.Token("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

	// Create a Test TLS Server
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/v1/users/login" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Header.Get("Authorization") != "Bearer "+string(token) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Set double cookies
		cookie := http.Cookie{
			Name:     "csrf_token",
			Value:    "thisisanexamplecookietoken",
			MaxAge:   600,
			Secure:   true,
			HttpOnly: false,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)

		cookie.Name = "csrf_reference_token"
		cookie.HttpOnly = true
		http.SetCookie(w, &cookie)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Fetch the https client and add a cookie jar
	client := ts.Client()
	client.Jar, err = cookiejar.New(nil)
	require.NoError(t, err, "couldn't add a cookie jar to the https client")

	// Create the BFF api client
	bff, err := api.New(ts.URL, api.WithClient(client), api.WithCredentials(token))
	require.NoError(t, err, "couldn't create BFF client with https and credentials")

	// Execute the Login request
	err = bff.Login(context.TODO(), &api.LoginParams{})
	require.NoError(t, err, "could not login using the bff client")

	// Check to ensure double cookies are set. This doesn't test our code, but ensures
	// that tests that depend on double cookies will work in the future.
	u, err := url.Parse(ts.URL)
	require.NoError(t, err, "could not parse test server url")
	require.Len(t, client.Jar.Cookies(u), 2, "expected two cookies set in the cookie jar")
}

func TestListUserRoles(t *testing.T) {
	fixture := []string{"leader", "collaborator"}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/users/roles", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.ListUserRoles(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture, out)
}

func TestUpdateUser(t *testing.T) {
	fixture := &api.UpdateUserParams{
		Name: "John Doe",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/v1/users", r.URL.Path)

		w.WriteHeader(http.StatusNoContent)
	}))

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	err = client.UpdateUser(context.TODO(), fixture)
	require.NoError(t, err)
}

func TestCreateOrganization(t *testing.T) {
	fixture := &api.OrganizationReply{
		ID:     "8b2e9e78-baca-4c34-a382-8b285503c901",
		Name:   "Alice VASP",
		Domain: "alicevasp.io",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/organizations", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.CreateOrganization(context.TODO(), &api.OrganizationParams{
		Name:   fixture.Name,
		Domain: fixture.Domain,
	})
	require.NoError(t, err)
	require.Equal(t, fixture, out)
}

func TestUserOrganization(t *testing.T) {
	fixture := &api.OrganizationReply{
		ID:     "8b2e9e78-baca-4c34-a382-8b285503c901",
		Name:   "Alice VASP",
		Domain: "alicevasp.io",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/users/organization", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.UserOrganization(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture, out)
}

func TestListOrganizations(t *testing.T) {
	fixture := &api.ListOrganizationsReply{
		Organizations: []*api.OrganizationReply{
			{
				ID:     "8b2e9e78-baca-4c34-a382-8b285503c901",
				Name:   "Alice VASP",
				Domain: "alicevasp.io",
			},
			{
				ID:     "c22ef329-2b8a-4b0c-9c1f-2b8a4b0c9c1f",
				Name:   "Bob VASP",
				Domain: "bobvasp.io",
			},
		},
		Count:    2,
		Page:     2,
		PageSize: 10,
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/organizations", r.URL.Path)
		require.Equal(t, "page=2&page_size=10", r.URL.RawQuery)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.ListOrganizations(context.TODO(), &api.ListOrganizationsParams{
		Page:     2,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.Equal(t, fixture, out)
}

func TestPatchOrganization(t *testing.T) {
	fixture := &api.OrganizationReply{
		Name:   "Alice VASP",
		Domain: "alicevasp.io",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/v1/organizations/8b2e9e78-baca-4c34-a382-8b285503c901", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	req := &api.OrganizationParams{
		Name:   fixture.Name,
		Domain: fixture.Domain,
	}
	out, err := client.PatchOrganization(context.TODO(), "8b2e9e78-baca-4c34-a382-8b285503c901", req)
	require.NoError(t, err)
	require.Equal(t, fixture, out)
}

func TestDeleteOrganization(t *testing.T) {
	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/organizations/8b2e9e78-baca-4c34-a382-8b285503c901", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Do the request
	require.NoError(t, client.DeleteOrganization(context.TODO(), "8b2e9e78-baca-4c34-a382-8b285503c901"))
}

func TestAddCollaborator(t *testing.T) {
	fixture := &models.Collaborator{
		Email:     "alice@example.com",
		CreatedAt: time.Now().Format(time.RFC3339Nano),
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/collaborators", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	request := &models.Collaborator{
		Email: "alice@example.com",
	}
	collab, err := client.AddCollaborator(context.TODO(), request)
	require.NoError(t, err)
	require.Equal(t, fixture.Email, collab.Email)
}

func TestListCollaborators(t *testing.T) {
	fixture := &api.ListCollaboratorsReply{
		Collaborators: []*models.Collaborator{
			{
				Id:        "wWD4zGmk8L8rA2J1I1PQYA",
				Email:     "alice@example.com",
				Roles:     []string{"Organization Collaborator"},
				CreatedAt: time.Now().Format(time.RFC3339Nano),
			},
			{
				Id:        "cD4FGmk8L8rA2J1I1PQYA",
				Email:     "bob@example.com",
				Roles:     []string{"Organization Leader"},
				CreatedAt: time.Now().Format(time.RFC3339Nano),
			},
		},
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/collaborators", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	collabs, err := client.ListCollaborators(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture, collabs)
}

func TestUpdateCollaboratorRoles(t *testing.T) {
	fixture := &models.Collaborator{
		Id:        "wWD4zGmk8L8rA2J1I1PQYA",
		Email:     "alice@example.com",
		Roles:     []string{"Organization Collaborator"},
		CreatedAt: time.Now().Format(time.RFC3339Nano),
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/collaborators/wWD4zGmk8L8rA2J1I1PQYA", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Should return an error if no ID is present
	_, err = client.UpdateCollaboratorRoles(context.TODO(), "", &api.UpdateRolesParams{})
	require.EqualError(t, err, api.ErrIDRequired.Error())

	// Valid request should succeed and return the fixture
	request := &api.UpdateRolesParams{
		Roles: []string{"Organization Collaborator"},
	}
	collab, err := client.UpdateCollaboratorRoles(context.TODO(), fixture.Id, request)
	require.NoError(t, err)
	require.Equal(t, fixture.Email, collab.Email)
}

func TestDeleteCollaborator(t *testing.T) {
	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/collaborators/wWD4zGmk8L8rA2J1I1PQYA", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Should return an error if no ID is present
	err = client.DeleteCollaborator(context.TODO(), "")
	require.EqualError(t, err, api.ErrIDRequired.Error())

	// Valid request should succeed
	require.NoError(t, client.DeleteCollaborator(context.TODO(), "wWD4zGmk8L8rA2J1I1PQYA"))
}

func TestLoadRegistrationForm(t *testing.T) {
	// Load a fixture from testdata
	fixture := &models.RegistrationForm{}
	err := loadFixture("testdata/registration.pb.json", fixture)
	require.NoError(t, err, "could not load registration fixture")

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/register", r.URL.Path)

		form := &api.RegistrationForm{
			Form: fixture,
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(form)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.LoadRegistrationForm(context.TODO(), nil)
	require.NoError(t, err)
	require.Equal(t, fixture, out.Form)
}

func TestSaveRegistrationForm(t *testing.T) {
	// Load a fixture from testdata
	fixture := &api.RegistrationForm{
		Form: &models.RegistrationForm{
			Website: "https://example.com",
		},
	}
	err := loadFixture("testdata/registration.pb.json", fixture)
	require.NoError(t, err, "could not load registration fixture")

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/register", r.URL.Path)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Should return the form fixture
	out, err := client.SaveRegistrationForm(context.TODO(), fixture)
	require.NoError(t, err)
	require.Equal(t, fixture, out)
	ts.Close()

	// Create a Test Server that returns 204 No Content
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/register", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err = api.New(ts.URL)
	require.NoError(t, err)

	// Should return no content
	out, err = client.SaveRegistrationForm(context.TODO(), fixture)
	require.NoError(t, err)
	require.Nil(t, out)
}

func TestResetRegistrationForm(t *testing.T) {
	// Load a defualt form
	fixture := models.NewRegisterForm()

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/register", r.URL.Path)

		form := &api.RegistrationForm{
			Form: fixture,
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(form)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Should return no content
	out, err := client.ResetRegistrationForm(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, fixture, out.Form)
}

func TestSubmitRegistration(t *testing.T) {
	fixture := &api.RegisterReply{
		Id:                  "8b2e9e78-baca-4c34-a382-8b285503c901",
		RegisteredDirectory: "trisa.directory",
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

	out, err := client.SubmitRegistration(context.TODO(), "MainNet")
	require.NoError(t, err)
	require.Equal(t, fixture.Id, out.Id)
	require.Equal(t, fixture.RegisteredDirectory, out.RegisteredDirectory)
	require.Equal(t, fixture.CommonName, out.CommonName)
	require.Equal(t, fixture.Status, out.Status)
	require.Equal(t, fixture.Message, out.Message)
	require.Equal(t, fixture.PKCS12Password, out.PKCS12Password)
}

func TestRegistrationStatus(t *testing.T) {
	fixture := &api.RegistrationStatus{
		TestNetSubmitted: time.Now().Format(time.RFC3339),
		MainNetSubmitted: time.Now().Add(time.Hour).Format(time.RFC3339),
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/registration", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))

	// Create a client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.RegistrationStatus(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture.TestNetSubmitted, out.TestNetSubmitted)
	require.Equal(t, fixture.MainNetSubmitted, out.MainNetSubmitted)
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

func TestAnnoucements(t *testing.T) {
	fixture := &api.AnnouncementsReply{
		Announcements: []*models.Announcement{
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

		in := &models.Announcement{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode register request")

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	req := &models.Announcement{
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

		in := &models.Announcement{}
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

	req := &models.Announcement{Title: "200"}
	err = client.MakeAnnouncement(context.TODO(), req)
	require.EqualError(t, err, "expected no content, received 200 OK")

	req = &models.Announcement{Title: "400"}
	err = client.MakeAnnouncement(context.TODO(), req)
	require.EqualError(t, err, "400 Bad Request")
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

func TestMemberList(t *testing.T) {
	fixture := &api.MemberListReply{
		VASPs: []*members.VASPMember{
			{
				Id:                  "7b8e1638-cf44-4b72-a4ae-08ff0352563b",
				RegisteredDirectory: "testnet",
				CommonName:          "example.com",
				Endpoint:            "trisa.example.com",
			},
			{
				Id:                  "03f47724-4751-40d4-8dda-fa5468f3b4a7",
				RegisteredDirectory: "testnet",
				CommonName:          "foobear.io",
				Endpoint:            "trisa-test.foobear.io.com",
			},
		},
		NextPageToken: "thenextpage",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/members", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	req := &api.MemberPageInfo{
		Directory: "testnet",
		PageSize:  100,
		PageToken: "theprevpage",
	}

	out, err := client.MemberList(context.TODO(), req)
	require.NoError(t, err)
	require.Equal(t, fixture, out)
}

func TestMemberDetails(t *testing.T) {
	fixture := &api.MemberDetailsReply{
		Summary: map[string]interface{}{
			"id":                   "8b2e9e78-baca-4c34-a382-8b285503c901",
			"registered_directory": "testnet.directory",
			"common_name":          "trisa.example.com",
			"endpoint":             "trisa.example.com:443",
			"name":                 "Trisa TestNet",
			"website":              "https://trisa.example.com",
			"country":              "US",
			"vasp_categories":      []interface{}{"P2P"},
			"verified_on":          "2022-04-21T12:05:23Z",
		},
		LegalPerson: map[string]interface{}{
			"country_of_registration": "US",
			"customer_number":         "123456789",
		},
		Contacts: map[string]interface{}{
			"technical": map[string]interface{}{
				"name":  "Rhyne Howard",
				"email": "rhoward@testing.net",
				"phone": "000-123-4567",
			},
		},
		Trixo: map[string]interface{}{
			"compliance_threshold":          0.0,
			"compliance_threshold_currency": "USD",
			"conducts_customer_kyc":         false,
		},
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/members/8b2e9e78-baca-4c34-a382-8b285503c901", r.URL.Path)
		require.Contains(t, r.URL.Query(), "registered_directory")
		require.Equal(t, "testnet.directory", r.URL.Query().Get("registered_directory"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	req := &api.MemberDetailsParams{
		ID:        "8b2e9e78-baca-4c34-a382-8b285503c901",
		Directory: "testnet.directory",
	}

	out, err := client.MemberDetails(context.TODO(), req)
	require.NoError(t, err)
	require.Equal(t, fixture, out)

	// Ensure member ID is required
	_, err = client.MemberDetails(context.TODO(), &api.MemberDetailsParams{})
	require.ErrorIs(t, err, api.ErrMissingMemberID)
}

func TestAttention(t *testing.T) {
	fixture := &api.AttentionReply{
		Messages: []*api.AttentionMessage{
			{
				Message:  bff.SubmitMainnet,
				Severity: models.AttentionSeverity_INFO.String(),
				Action:   models.AttentionAction_SUBMIT_MAINNET.String(),
			},
			{
				Message:  fmt.Sprintf(bff.CertificateRevoked, "testnet"),
				Severity: models.AttentionSeverity_ALERT.String(),
				Action:   models.AttentionAction_CONTACT_SUPPORT.String(),
			},
		},
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/attention", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.Attention(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture, out)
}

func TestNoAttention(t *testing.T) {
	// Create a Test Server that returns no content
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/attention", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Should return no error and nil response for 204 no content
	out, err := client.Attention(context.TODO())
	require.NoError(t, err)
	require.Nil(t, out)
}

func TestNetworkActivity(t *testing.T) {
	TestNetActivity := []api.Activity{
		{
			Date:   "Aug 21",
			Events: 10,
		},
		{
			Date:   "Aug 22",
			Events: 20,
		},
	}

	MainNetActivity := []api.Activity{
		{
			Date:   "Aug 21",
			Events: 15,
		},
		{
			Date:   "Aug 22",
			Events: 25,
		},
	}

	fixture := &api.NetworkActivityReply{
		TestNet: TestNetActivity,
		MainNet: MainNetActivity,
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/network/activity", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.NetworkActivity(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture, out)
	require.Equal(t, fixture.TestNet, out.TestNet)
	require.Equal(t, fixture.MainNet, out.MainNet)
}

func loadFixture(path string, v interface{}) (err error) {
	switch t := v.(type) {
	case proto.Message:
		return loadPBFixture(path, t)
	default:
		return loadJSONFixture(path, t)
	}
}

func loadPBFixture(path string, v proto.Message) (err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return err
	}

	pbjson := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	if err = pbjson.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

func loadJSONFixture(path string, v interface{}) (err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}
