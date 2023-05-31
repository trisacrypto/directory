package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
)

// New creates a new api.v1 API client that implements the BFF interface.
func New(endpoint string, opts ...ClientOption) (_ BFFClient, err error) {
	// Create a client with the parsed endpoint.
	c := &APIv1{}
	if c.endpoint, err = url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("could not parse endpoint: %s", err)
	}

	// Apply options
	for _, opt := range opts {
		if err = opt(c); err != nil {
			return nil, err
		}
	}

	// If a client hasn't been specified, create the default client.
	if c.client == nil {
		c.client = &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Timeout:       30 * time.Second,
		}

		// Create cookie jar for CSRF
		if c.client.Jar, err = cookiejar.New(nil); err != nil {
			return nil, fmt.Errorf("could not create cookiejar: %s", err)
		}
	}
	return c, nil
}

// APIv1 implements the BFFClient interface.
type APIv1 struct {
	endpoint *url.URL
	client   *http.Client
	creds    Credentials
}

// Ensure the API implments the BFFClient interface.
var _ BFFClient = &APIv1{}

//===========================================================================
// Client Methods
//===========================================================================

// Status performs a health check request to the BFF.
func (s *APIv1) Status(ctx context.Context, in *StatusParams) (out *StatusReply, err error) {
	// Create the query params from the input
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %s", err)
	}

	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/status", nil, &params); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	// NOTE: cannot use s.Do because we want to parse 503 Unavailable errors
	var rep *http.Response
	if rep, err = s.client.Do(req); err != nil {
		return nil, fmt.Errorf("could not execute request: %s", err)
	}
	defer rep.Body.Close()

	// Detect other errors
	if rep.StatusCode != http.StatusOK && rep.StatusCode != http.StatusServiceUnavailable {
		return nil, fmt.Errorf("[%d] %s", rep.StatusCode, rep.Status)
	}

	// Deserialize the JSON data from the response
	out = &StatusReply{}
	if err = json.NewDecoder(rep.Body).Decode(out); err != nil {
		return nil, fmt.Errorf("could not deserialize StatusReply: %s", err)
	}
	return out, nil
}

// Lookup a VASP record in both the TestNet and the MainNet.
func (s *APIv1) Lookup(ctx context.Context, in *LookupParams) (out *LookupReply, err error) {
	// Create the query params from the input
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %s", err)
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/lookup", nil, &params); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &LookupReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Verify a contact with the token sent to their email address.
func (s *APIv1) VerifyContact(ctx context.Context, in *VerifyContactParams) (out *VerifyContactReply, err error) {
	// Create the query params from the input
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %s", err)
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/verify", nil, &params); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &VerifyContactReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Login post-processes an Auth0 login or registration and sets CSRF cookies.
func (s *APIv1) Login(ctx context.Context, in *LoginParams) (err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/users/login", in, nil); err != nil {
		return err
	}

	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}
	return nil
}

// Return the set of assignable user roles.
func (s *APIv1) ListUserRoles(ctx context.Context) (out []string, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/users/roles", nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = make([]string, 0)
	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Update the user's profile.
func (s *APIv1) UpdateUser(ctx context.Context, in *UpdateUserParams) (err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPatch, "/v1/users", in, nil); err != nil {
		return err
	}

	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}
	return nil
}

// Get the user's current organization.
func (s *APIv1) UserOrganization(ctx context.Context) (out *OrganizationReply, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/users/organization", nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &OrganizationReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Create a new organization.
func (s *APIv1) CreateOrganization(ctx context.Context, in *OrganizationParams) (out *OrganizationReply, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/organizations", in, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &OrganizationReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// List available organizations.
func (s *APIv1) ListOrganizations(ctx context.Context, in *ListOrganizationsParams) (out *ListOrganizationsReply, err error) {
	// Create the query params from the input
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %s", err)
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/organizations", nil, &params); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &ListOrganizationsReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Patch an organization.
func (s *APIv1) PatchOrganization(ctx context.Context, id string, in *OrganizationParams) (out *OrganizationReply, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPatch, "/v1/organizations/"+id, in, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &OrganizationReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete an organization by ID.
func (s *APIv1) DeleteOrganization(ctx context.Context, id string) (err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, "/v1/organizations/"+id, nil, nil); err != nil {
		return err
	}

	// Execute the request and get a response
	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}
	return nil
}

// Add a collaborator to an organization.
func (s *APIv1) AddCollaborator(ctx context.Context, request *models.Collaborator) (collaborator *models.Collaborator, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/collaborators", request, nil); err != nil {
		return nil, err
	}

	collaborator = &models.Collaborator{}
	if _, err = s.Do(req, collaborator, true); err != nil {
		return nil, err
	}
	return collaborator, nil
}

// List all collaborators on an organization.
func (s *APIv1) ListCollaborators(ctx context.Context) (out *ListCollaboratorsReply, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/collaborators", nil, nil); err != nil {
		return nil, err
	}

	out = &ListCollaboratorsReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Update a collaborator's roles in an organization.
func (s *APIv1) UpdateCollaboratorRoles(ctx context.Context, id string, request *UpdateRolesParams) (collaborator *models.Collaborator, err error) {
	// ID is required for the endpoint
	if id == "" {
		return nil, ErrIDRequired
	}

	// Construct the path from the request
	path := fmt.Sprintf("/v1/collaborators/%s", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, path, request, nil); err != nil {
		return nil, err
	}

	collaborator = &models.Collaborator{}
	if _, err = s.Do(req, collaborator, true); err != nil {
		return nil, err
	}
	return collaborator, nil
}

// Delete a collaborator from an organization.
func (s *APIv1) DeleteCollaborator(ctx context.Context, id string) (err error) {
	// ID is required for the endpoint
	if id == "" {
		return ErrIDRequired
	}

	// Construct the path from the request
	path := fmt.Sprintf("/v1/collaborators/%s", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return err
	}

	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}
	return nil
}

// Load registration form data from the server to populate the front-end form.
func (s *APIv1) LoadRegistrationForm(ctx context.Context, in *RegistrationFormParams) (form *RegistrationForm, err error) {
	// Create the query params from the input
	var params url.Values
	if in != nil {
		if params, err = query.Values(in); err != nil {
			return nil, fmt.Errorf("could not encode query params: %s", err)
		}
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/register", nil, &params); err != nil {
		return nil, err
	}

	form = &RegistrationForm{}
	if _, err = s.Do(req, form, true); err != nil {
		return nil, err
	}
	return form, nil
}

// Save registration form data to the server in preparation for submitting it.
func (s *APIv1) SaveRegistrationForm(ctx context.Context, form *RegistrationForm) (out *RegistrationForm, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPut, "/v1/register", form, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	var rep *http.Response
	out = &RegistrationForm{}
	if rep, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	// Make sure we are not returning data if a 204 No Content was received
	if rep.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	return out, nil
}

// Reset the registration form on the server to its default values.
func (s *APIv1) ResetRegistrationForm(ctx context.Context, in *RegistrationFormParams) (form *RegistrationForm, err error) {
	// Create the query params from the input
	var params url.Values
	if in != nil {
		if params, err = query.Values(in); err != nil {
			return nil, fmt.Errorf("could not encode query params: %s", err)
		}
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, "/v1/register", nil, &params); err != nil {
		return nil, err
	}

	form = &RegistrationForm{}
	if _, err = s.Do(req, form, true); err != nil {
		return nil, err
	}
	return form, nil
}

// Submit the registration form to the specified network (testnet or mainnet).
func (s *APIv1) SubmitRegistration(ctx context.Context, network string) (out *RegisterReply, err error) {
	// network is required for the endpoint
	if network == "" {
		return nil, ErrNetworkRequired
	}

	// Determine the path for the request
	network = strings.ToLower(strings.TrimSpace(network))
	path := fmt.Sprintf("/v1/register/%s", network)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, path, nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &RegisterReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// RegistrationStatus returns the status of the VASP registrations for the organization.
func (s *APIv1) RegistrationStatus(ctx context.Context) (out *RegistrationStatus, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/registration", nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &RegistrationStatus{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Overview returns a high-level summary of the organization account and networks.
func (s *APIv1) Overview(ctx context.Context) (out *OverviewReply, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/overview", nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &OverviewReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Announcements returns a list of network announcments made by the admins.
func (s *APIv1) Announcements(ctx context.Context) (out *AnnouncementsReply, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/announcements", nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &AnnouncementsReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// MakeAnnouncement allows administrators to post new network announcements.
func (s *APIv1) MakeAnnouncement(ctx context.Context, in *models.Announcement) (err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/announcements", in, nil); err != nil {
		return err
	}

	// Execute the request and get a response, ensuring to check that a 200 response is returned.
	var rep *http.Response
	if rep, err = s.Do(req, nil, true); err != nil {
		return err
	}

	// If this endpoint does not return a 204 then return an error since data is unhandled.
	if rep.StatusCode != http.StatusNoContent {
		return fmt.Errorf("expected no content, received %s", rep.Status)
	}
	return nil
}

// Certificates returns the list of certificates associated with the organization.
func (s *APIv1) Certificates(ctx context.Context) (out *CertificatesReply, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/certificates", nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &CertificatesReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Details returns the sensitive details for a VASP member.
func (s *APIv1) MemberDetails(ctx context.Context, in *MemberDetailsParams) (out *MemberDetailsReply, err error) {
	// Create the query params from the input
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %s", err)
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/details", nil, &params); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &MemberDetailsReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Attention returns the set of current attention messages for the organization.
func (s *APIv1) Attention(ctx context.Context) (out *AttentionReply, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/attention", nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &AttentionReply{}
	var rep *http.Response
	if rep, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	// Make sure no data is returned if the status code is 204 (no content)
	if rep.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	return out, nil
}

//===========================================================================
// Helper Methods
//===========================================================================

const (
	userAgent    = "GDS BFF API Client/v1"
	accept       = "application/json"
	acceptLang   = "en-US,en"
	acceptEncode = "gzip, deflate, br"
	contentType  = "application/json; charset=utf-8"
)

// NewRequest creates an http.Request with the specified context and method, resolving
// the path to the root endpoint of the API (e.g. /v2) and serializes the data to JSON.
// This method also sets the default headers of all GDS Admin API v2 client requests.
func (s *APIv1) NewRequest(ctx context.Context, method, path string, data interface{}, params *url.Values) (req *http.Request, err error) {
	// Resolve the URL reference from the path
	endpoint := s.endpoint.ResolveReference(&url.URL{Path: path})
	if params != nil && len(*params) > 0 {
		endpoint.RawQuery = params.Encode()
	}

	var body io.ReadWriter
	switch {
	case data == nil:
		body = nil
	default:
		body = &bytes.Buffer{}
		if err = json.NewEncoder(body).Encode(data); err != nil {
			return nil, fmt.Errorf("could not serialize request data: %s", err)
		}
	}

	// Create the http request
	if req, err = http.NewRequestWithContext(ctx, method, endpoint.String(), body); err != nil {
		return nil, fmt.Errorf("could not create request: %s", err)
	}

	// Set the headers on the request
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", accept)
	req.Header.Add("Accept-Language", acceptLang)
	req.Header.Add("Accept-Encoding", acceptEncode)
	req.Header.Add("Content-Type", contentType)

	// Add authentication if it is available
	if s.creds != nil {
		var token string
		if token, err = s.creds.AccessToken(); err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", "Bearer "+token)
	}

	// Add CSRF protection if it is available
	if s.client.Jar != nil {
		cookies := s.client.Jar.Cookies(endpoint)
		for _, cookie := range cookies {
			if cookie.Name == "csrf_token" {
				req.Header.Add("X-CSRF-TOKEN", cookie.Value)
			}
		}
	}

	return req, nil
}

// Do executes an http request against the server, performs error checking, and
// deserializes the response data into the specified struct if requested.
func (s *APIv1) Do(req *http.Request, data interface{}, checkStatus bool) (rep *http.Response, err error) {
	if rep, err = s.client.Do(req); err != nil {
		return rep, fmt.Errorf("could not execute request: %s", err)
	}
	defer rep.Body.Close()

	// Detect errors if they've occurred
	if checkStatus {
		if rep.StatusCode < 200 || rep.StatusCode >= 300 {
			// Attempt to read the error response from the JSON, ignore body
			// deserialization or read errors and simply return the status error.
			var reply Reply
			if err = json.NewDecoder(rep.Body).Decode(&reply); err == nil {
				if reply.Error != "" {
					return rep, fmt.Errorf("[%d] %s", rep.StatusCode, reply.Error)
				}
			}
			return rep, errors.New(rep.Status)
		}
	}

	// Deserialize the JSON data from the body
	// TODO: what if this is protocol buffer JSON?
	if data != nil && rep.StatusCode >= 200 && rep.StatusCode < 300 && rep.StatusCode != http.StatusNoContent {
		// Check the content type to ensure data deserialization is possible
		if ct := rep.Header.Get("Content-Type"); ct != contentType {
			return rep, fmt.Errorf("unexpected content type: %q", ct)
		}

		if err = json.NewDecoder(rep.Body).Decode(data); err != nil {
			return nil, fmt.Errorf("could not deserialize response data: %s", err)
		}
	}

	return rep, nil
}

// SetCredentials is a helper function for external users to override credentials at
// runtime and is used extensively in testing the BFF server.
func (c *APIv1) SetCredentials(creds Credentials) {
	c.creds = creds
}

// SetCSRFProtect is a helper function to set CSRF cookies on the client. This is not
// possible in a browser because of the HttpOnly flag. This method should only be used
// for testing purposes and an error is returned if the URL is not localhost. For live
// clients - the server should set these cookies. If protect is false, then the cookies
// are removed from the client by setting the cookies to an empty slice.
func (c *APIv1) SetCSRFProtect(protect bool) error {
	if c.client.Jar == nil {
		return errors.New("client does not have a cookie jar, cannot set cookies")
	}

	if c.endpoint.Hostname() != "127.0.0.1" && c.endpoint.Hostname() != "localhost" {
		return fmt.Errorf("csrf protect is for local testing only, cannot set cookies for %s", c.endpoint.Hostname())
	}

	// The URL for the cookies
	u := c.endpoint.ResolveReference(&url.URL{Path: "/"})

	var cookies []*http.Cookie
	if protect {
		cookies = []*http.Cookie{
			{
				Name:     "csrf_token",
				Value:    "testingcsrftoken",
				Expires:  time.Now().Add(10 * time.Minute),
				HttpOnly: false,
			},
			{
				Name:     "csrf_reference_token",
				Value:    "testingcsrftoken",
				Expires:  time.Now().Add(10 * time.Minute),
				HttpOnly: true,
			},
		}
	} else {
		cookies = c.client.Jar.Cookies(u)
		for _, cookie := range cookies {
			cookie.MaxAge = -1
		}
	}

	c.client.Jar.SetCookies(u, cookies)
	return nil
}
