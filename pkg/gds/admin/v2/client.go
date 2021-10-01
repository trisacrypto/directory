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
	"os"
	"path/filepath"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/go-querystring/query"
	"github.com/kelseyhightower/envconfig"
	"github.com/shibukawa/configdir"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

// New creates a new admin.v2 API client that implements the Service interface.
func New(endpoint string) (_ DirectoryAdministrationClient, err error) {
	c := &APIv2{
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
	endpoint    *url.URL
	client      *http.Client
	accessToken string
	csrfToken   string
}

// Ensure the API implments the Service interface.
var _ DirectoryAdministrationClient = &APIv2{}

// Login prepares the client for authorized requests. It first looks for a stored token
// on disk; if it finds none then it uses either a token manager or redirects the client
// to the browser to authenticate with the server. If the access token is expired, it
// uses the refresh token to reauthenticate and saves the resulting tokens back to disk.
func (s *APIv2) Login(ctx context.Context) (err error) {
	var creds *AuthReply
	if creds, err = s.Credentials(); err != nil {
		// No credentials were found begin login process
		// TODO: refactor this section to allow either local token generation or online login
		return s.GenerateCredentials()
	}
	s.accessToken = creds.AccessToken

	// If we've successfully loaded the credentials check the access token to make sure it's still valid
	parser := &jwt.Parser{}
	claims := &tokens.Claims{}
	if _, _, err = parser.ParseUnverified(creds.AccessToken, claims); err != nil {
		s.DeleteCredentials()
		return errors.New("access tokens unparseable, cached credentials have been deleted, please try again")
	}

	if err = claims.Valid(); err == nil {
		return nil
	}

	// Access token is invalid attempt to reauthenticate
	if creds, err = s.Reauthenticate(ctx, &AuthRequest{Credential: creds.RefreshToken}); err != nil {
		s.DeleteCredentials()
		return fmt.Errorf("could not reauthenticate (%s), cached credentials have been deleted, please try again", err)
	}

	// Save refreshed creds to disk
	s.accessToken = creds.AccessToken

	var data []byte
	if data, err = json.Marshal(creds); err != nil {
		return err
	}

	folders := cfgd.QueryFolders(configdir.Global)
	folders[0].WriteFile(credentials, data)

	return nil
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

func (s *APIv2) Review(ctx context.Context, in *ReviewRequest) (out *ReviewReply, err error) {
	// The ID is required for the review request to determine the endpoint
	if in.ID == "" {
		return nil, ErrIDRequred
	}

	if err = s.ProtectAuthenticate(ctx); err != nil {
		return nil, fmt.Errorf("could not protect review from CSRF: %s", err)
	}

	// Determine the path from the request
	path := fmt.Sprintf("/v2/vasps/%s/review", in.ID)

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

	if err = s.ProtectAuthenticate(ctx); err != nil {
		return nil, fmt.Errorf("could not protect resend from CSRF: %s", err)
	}

	// Determine the path from the request
	path := fmt.Sprintf("/v2/vasps/%s/resend", in.ID)

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

	// Set authorizatoin and csrf protection if available
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

var cfgd = configdir.New("rotational", "gds")

const (
	credentials = "credentials.json"
)

// Credentials returns the cached access and refresh tokens from disk.
func (s *APIv2) Credentials() (creds *AuthReply, err error) {
	folder := cfgd.QueryFolderContainsFile(credentials)
	if folder != nil {
		var data []byte
		if data, err = folder.ReadFile(credentials); err != nil {
			return nil, err
		}

		creds = &AuthReply{}
		if err = json.Unmarshal(data, creds); err != nil {
			return nil, err
		}

		return creds, nil
	}
	return nil, errors.New("no credentials are available")
}

// DeleteCredentials removed cached access and refresh tokens
func (s *APIv2) DeleteCredentials() (err error) {
	folder := cfgd.QueryFolderContainsFile(credentials)
	if folder != nil && folder.Exists(credentials) {
		return os.Remove(filepath.Join(folder.Path, credentials))
	}

	return nil
}

// TODO: do better than this when we have client profiles
type ClientConfig struct {
	TokenKeys map[string]string `envconfig:"GDS_ADMIN_TOKEN_KEYS"`
}

// GenerateCredentials creates a token manager generate and save credentials
func (s *APIv2) GenerateCredentials() (err error) {
	var conf ClientConfig
	if err = envconfig.Process("gds", &conf); err != nil {
		return err
	}

	if len(conf.TokenKeys) == 0 {
		return errors.New("invalid configuration: token keys are required for local key generation")
	}

	var tm *tokens.TokenManager
	if tm, err = tokens.New(conf.TokenKeys); err != nil {
		return err
	}

	var accessToken, refreshToken *jwt.Token

	claims := map[string]interface{}{
		"hd":      "rotational.io",
		"email":   "admin@rotational.io",
		"name":    "GDS Admin CLI",
		"picture": "",
	}

	// Create the access and refresh tokens from the claims
	if accessToken, err = tm.CreateAccessToken(claims); err != nil {
		return err
	}

	if refreshToken, err = tm.CreateRefreshToken(accessToken); err != nil {
		return err
	}

	// Sign the tokens and return the response
	creds := new(AuthReply)
	if creds.AccessToken, err = tm.Sign(accessToken); err != nil {
		return err
	}
	if creds.RefreshToken, err = tm.Sign(refreshToken); err != nil {
		return err
	}

	// Save the credentials to disk
	var data []byte
	if data, err = json.Marshal(creds); err != nil {
		return err
	}

	folders := cfgd.QueryFolders(configdir.Global)
	folders[0].WriteFile(credentials, data)

	s.accessToken = creds.AccessToken
	return nil
}
