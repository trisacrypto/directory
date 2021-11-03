package admin

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
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/go-querystring/query"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

// New creates a new admin.v2 API client that implements the Service interface.
func New(endpoint string, creds Credentials) (_ DirectoryAdministrationClient, err error) {
	c := &APIv2{
		creds: creds,
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Timeout:       30 * time.Second,
		},
	}

	// Create cookie jar
	if c.client.Jar, err = cookiejar.New(nil); err != nil {
		return nil, fmt.Errorf("could not create cookiejar: %s", err)
	}

	if c.endpoint, err = url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("could not parse endpoint: %s", err)
	}
	return c, nil
}

// APIv2 implements the Service interface.
type APIv2 struct {
	endpoint     *url.URL
	creds        Credentials
	client       *http.Client
	accessToken  string
	refreshToken string
	csrfToken    string
}

// Ensure the API implments the Service interface.
var _ DirectoryAdministrationClient = &APIv2{}

//===========================================================================
// User Methods
//===========================================================================

// Login prepares the client for authorized requests. It uses the internal credentials
// to set the access and refresh tokens but does not return an error if the internal
// credentials is nil. Login can be called manually by the user, or it is called
// automatically on requests that must be authenticated.
func (s *APIv2) Login(ctx context.Context) (err error) {
	if s.creds != nil {
		if s.accessToken, s.refreshToken, err = s.creds.Login(s); err != nil {
			return err
		}
	}
	return nil
}

// Refresh prepares the client to continue making authorized requests. It uses the
// internal credentials to update the access and refresh tokens but does not return an
// error if the internal credentials is nil. Refresh can be called manually by the user,
// or it is called automatically on requests that must be authenticated.
func (s *APIv2) Refresh(ctx context.Context) (err error) {
	if s.creds != nil {
		if s.accessToken, s.refreshToken, err = s.creds.Refresh(s); err != nil {
			return err
		}
	}
	return nil
}

// Logout removes the cached access and refresh tokens requiring the API to login again
func (s *APIv2) Logout(ctx context.Context) (err error) {
	// Remove tokens from the cache
	s.accessToken = ""
	s.refreshToken = ""
	s.csrfToken = ""

	if s.creds != nil {
		return s.creds.Logout(s)
	}
	return nil
}

// Tokens returns the access and refresh tokens. It is not part of the
// DirectoryAdministrationClient interface, so callers will have to type check *APIv2 to
// get access to this method. It's intent is to provide the cached tokens to the
// internal credentials on refresh or for auditing purposes and use in testing.
func (s *APIv2) Tokens() (accessToken, refreshToken string) {
	return s.accessToken, s.refreshToken
}

//===========================================================================
// Client Methods
//===========================================================================

// ProtectAuthenticate sets cookies for CSRF protection (not necessary in the API but good for testing)
func (s *APIv2) ProtectAuthenticate(ctx context.Context) (err error) {
	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v2/authenticate", nil, nil); err != nil {
		return err
	}

	// Execute the request and get a response
	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}

	// NOTE: this will only work over HTTPS, not for local debugging
	cookies := s.client.Jar.Cookies(s.endpoint)
	for _, cookie := range cookies {
		if cookie.Name == CSRFCookie {
			s.csrfToken = cookie.Value
		}
	}
	return nil
}

// Authenticate the the client to the Server using the supplied credentials.
// This method calls ProtectAuthenticate before performing the authentication.
func (s *APIv2) Authenticate(ctx context.Context, in *AuthRequest) (out *AuthReply, err error) {
	if err = s.ProtectAuthenticate(ctx); err != nil {
		return nil, fmt.Errorf("could not protect authenticate from CSRF: %s", err)
	}

	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v2/authenticate", in, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &AuthReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

// Reauthenticate the the client to the Server using the supplied credentials.
// This method calls ProtectAuthenticate before performing the reauthentication.
func (s *APIv2) Reauthenticate(ctx context.Context, in *AuthRequest) (out *AuthReply, err error) {
	if err = s.ProtectAuthenticate(ctx); err != nil {
		return nil, fmt.Errorf("could not protect reauthenticate from CSRF: %s", err)
	}

	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v2/reauthenticate", in, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &AuthReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv2) Status(ctx context.Context) (out *StatusReply, err error) {
	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v2/status", nil, nil); err != nil {
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

func (s *APIv2) Summary(ctx context.Context) (out *SummaryReply, err error) {
	// Must be authenticated
	if err = s.checkAuthentication(ctx); err != nil {
		return nil, err
	}

	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v2/summary", nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &SummaryReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv2) Autocomplete(ctx context.Context) (out *AutocompleteReply, err error) {
	// Must be authenticated
	if err = s.checkAuthentication(ctx); err != nil {
		return nil, err
	}

	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v2/autocomplete", nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &AutocompleteReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv2) ListVASPs(ctx context.Context, in *ListVASPsParams) (out *ListVASPsReply, err error) {
	// Create the query params from the input
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %s", err)
	}

	// Must be authenticated
	if err = s.checkAuthentication(ctx); err != nil {
		return nil, err
	}

	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v2/vasps", nil, &params); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &ListVASPsReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv2) RetrieveVASP(ctx context.Context, id string) (out *RetrieveVASPReply, err error) {
	// Compute the path based on the id
	if id == "" {
		return nil, errors.New("id is required to compute the URL for the VASP")
	}
	path := fmt.Sprintf("/v2/vasps/%s", id)

	// Must be authenticated
	if err = s.checkAuthentication(ctx); err != nil {
		return nil, err
	}

	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &RetrieveVASPReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv2) CreateReviewNote(ctx context.Context, in *ModifyReviewNoteRequest) (out *ReviewNote, err error) {
	// vaspID is required for the endpoint
	if in.VASP == "" {
		return nil, ErrIDRequred
	}

	// Determine the path from the request
	path := fmt.Sprintf("/v2/vasps/%s/notes", in.VASP)

	// Must be authenticated
	if err = s.checkAuthentication(ctx); err != nil {
		return nil, err
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, path, in, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &ReviewNote{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv2) ListReviewNotes(ctx context.Context, id string) (out *ListReviewNotesReply, err error) {
	// vaspID is required for the endpoint
	if id == "" {
		return nil, ErrIDRequred
	}

	// Determine the path from the request
	path := fmt.Sprintf("/v2/vasps/%s/notes", id)

	// Must be authenticated
	if err = s.checkAuthentication(ctx); err != nil {
		return nil, err
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &ListReviewNotesReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv2) UpdateReviewNote(ctx context.Context, in *ModifyReviewNoteRequest) (out *Reply, err error) {
	// vaspID and noteID are required for the endpoint
	if in.VASP == "" || in.NoteID == "" {
		return nil, ErrIDRequred
	}

	// Determine the path from the request
	path := fmt.Sprintf("/v2/vasps/%s/notes/%s", in.VASP, in.NoteID)

	// Must be authenticated
	if err = s.checkAuthentication(ctx); err != nil {
		return nil, err
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPut, path, in, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &Reply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv2) DeleteReviewNote(ctx context.Context, vaspID string, noteID string) (out *Reply, err error) {
	// vaspID and noteID are required for the endpoint
	if vaspID == "" || noteID == "" {
		return nil, ErrIDRequred
	}

	// Determine the path from the request
	path := fmt.Sprintf("/v2/vasps/%s/notes/%s", vaspID, noteID)

	// Must be authenticated
	if err = s.checkAuthentication(ctx); err != nil {
		return nil, err
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &Reply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv2) Review(ctx context.Context, in *ReviewRequest) (out *ReviewReply, err error) {
	// The ID is required for the review request to determine the endpoint
	if in.ID == "" {
		return nil, ErrIDRequred
	}

	// Determine the path from the request
	path := fmt.Sprintf("/v2/vasps/%s/review", in.ID)

	// Must be authenticated
	if err = s.checkAuthentication(ctx); err != nil {
		return nil, err
	}

	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, path, in, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &ReviewReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv2) Resend(ctx context.Context, in *ResendRequest) (out *ResendReply, err error) {
	// The ID is required for the review request to determine the endpoint
	if in.ID == "" {
		return nil, ErrIDRequred
	}

	// Determine the path from the request
	path := fmt.Sprintf("/v2/vasps/%s/resend", in.ID)

	// Must be authenticated
	if err = s.checkAuthentication(ctx); err != nil {
		return nil, err
	}

	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, path, in, nil); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &ResendReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv2) ReviewTimeline(ctx context.Context, in *ReviewTimelineParams) (out *ReviewTimelineReply, err error) {
	// Create the query params from the input
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %s", err)
	}

	// Must be authenticated
	if err = s.checkAuthentication(ctx); err != nil {
		return nil, err
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v2/reviews", nil, &params); err != nil {
		return nil, err
	}

	// Execute the request and get a response
	out = &ReviewTimelineReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

//===========================================================================
// Helper Methods
//===========================================================================

const (
	userAgent    = "GDS Admin API Client/v2"
	accept       = "application/json"
	acceptLang   = "en-US,en"
	acceptEncode = "gzip, deflate, br"
	contentType  = "application/json; charset=utf-8"
)

// NewRequest creates an http.Request with the specified context and method, resolving
// the path to the root endpoint of the API (e.g. /v2) and serializes the data to JSON.
// This method also sets the default headers of all GDS Admin API v2 client requests.
func (s *APIv2) NewRequest(ctx context.Context, method, path string, data interface{}, params *url.Values) (req *http.Request, err error) {
	// Resolve the URL reference from the path
	endpoint := s.endpoint.ResolveReference(&url.URL{Path: path})
	if params != nil && len(*params) > 0 {
		endpoint.RawQuery = params.Encode()
	}

	var body io.ReadWriter
	if data != nil {
		body = &bytes.Buffer{}
		if err = json.NewEncoder(body).Encode(data); err != nil {
			return nil, fmt.Errorf("could not serialize request data: %s", err)
		}
	} else {
		body = nil
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

	// Set authorization and csrf protection if available
	if s.accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+s.accessToken)
	}
	if s.csrfToken != "" {
		req.Header.Add(CSRFHeader, s.csrfToken)
	}

	return req, nil
}

// Do executes an http request against the server, performs error checking, and
// deserializes the response data into the specified struct if requested.
func (s *APIv2) Do(req *http.Request, data interface{}, checkStatus bool) (rep *http.Response, err error) {
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

	// Check the content type to ensure data deserialization is possible
	if ct := rep.Header.Get("Content-Type"); ct != contentType {
		return rep, fmt.Errorf("unexpected content type: %q", ct)
	}

	// Deserialize the JSON data from the body
	if data != nil && rep.StatusCode >= 200 && rep.StatusCode < 300 {
		if err = json.NewDecoder(rep.Body).Decode(data); err != nil {
			return nil, fmt.Errorf("could not deserialize response data: %s", err)
		}
	}

	return rep, nil
}

// checkAuthentication ensures that the client is prepared for authentication and should
// be called by all client methods that require authentication. The check will call
// Login() or Refresh() as needed depending on the state of the client.
func (s *APIv2) checkAuthentication(ctx context.Context) (err error) {
	// If no access token is available, call Login.
	if s.accessToken == "" {
		return s.Login(ctx)
	}

	// Ignore parsing error since we'll get ValidationErrorUnverifiable but ensure that
	// a token is returned in case it was a parsing error. See the following for more:
	// https://github.com/dgrijalva/jwt-go/issues/37#issuecomment-58764625
	accessClaims := new(tokens.Claims)
	if token, err := jwt.ParseWithClaims(s.accessToken, accessClaims, nil); token == nil {
		return fmt.Errorf("could not parse access token: %s", err)
	}

	// Manually check if the access token has not expired
	now := time.Now().Unix()
	if accessClaims.ExpiresAt != 0 && now > accessClaims.ExpiresAt {
		// access token is expired, check if refresh is not expired
		if s.refreshToken != "" {
			refreshClaims := new(tokens.Claims)
			if token, _ := jwt.ParseWithClaims(s.accessToken, accessClaims, nil); token == nil {
				return fmt.Errorf("could not parse refresh token")
			}

			if now <= refreshClaims.ExpiresAt {
				// refresh token is not expired, attempt a refresh
				return s.Refresh(ctx)
			}
		}

		// access and refresh tokens are both expired, login again
		return s.Login(ctx)
	}

	// access token is not expired, try using it
	// (doesn't mean it's not valid, but a 401 will be returned from server)
	return nil
}
