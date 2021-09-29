package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

// New creates a new admin.v2 API client that implements the Service interface.
func New(endpoint string) (_ DirectoryAdministrationClient, err error) {
	c := &APIv2{
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       30 * time.Second,
		},
	}
	if c.endpoint, err = url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("could not parse endpoint: %s", err)
	}
	return c, nil
}

// APIv2 implements the Service interface.
type APIv2 struct {
	endpoint *url.URL
	client   *http.Client
}

// Ensure the API implments the Service interface.
var _ DirectoryAdministrationClient = &APIv2{}

// Authenticate the the client to the Server using the supplied credentials.
// TODO: this method should prepare the APIv2 client to send and recv authenticated requests
// TODO: does the client need to call the ProtectAuthenticate endpoint before auth?
func (s APIv2) Authenticate(ctx context.Context, in *AuthRequest) (out *AuthReply, err error) {
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
// TODO: this method should prepare the APIv2 client to send and recv authenticated requests
func (s APIv2) Reauthenticate(ctx context.Context, in *AuthRequest) (out *AuthReply, err error) {
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

func (s APIv2) Status(ctx context.Context) (out *StatusReply, err error) {
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

func (s APIv2) Summary(ctx context.Context) (out *SummaryReply, err error) {
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

func (s APIv2) Autocomplete(ctx context.Context) (out *AutocompleteReply, err error) {
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

func (s APIv2) ListVASPs(ctx context.Context, in *ListVASPsParams) (out *ListVASPsReply, err error) {
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

func (s APIv2) RetrieveVASP(ctx context.Context, id string) (out *RetrieveVASPReply, err error) {
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

func (s APIv2) Review(ctx context.Context, in *ReviewRequest) (out *ReviewReply, err error) {
	// The ID is required for the review request to determine the endpoint
	if in.ID == "" {
		return nil, ErrIDRequred
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

func (s APIv2) Resend(ctx context.Context, in *ResendRequest) (out *ResendReply, err error) {
	// The ID is required for the review request to determine the endpoint
	if in.ID == "" {
		return nil, ErrIDRequred
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
func (s APIv2) NewRequest(ctx context.Context, method, path string, data interface{}, params *url.Values) (req *http.Request, err error) {
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

	return req, nil
}

// Do executes an http request against the server, performs error checking, and
// deserializes the response data into the specified struct if requested.
func (s APIv2) Do(req *http.Request, data interface{}, checkStatus bool) (rep *http.Response, err error) {
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
