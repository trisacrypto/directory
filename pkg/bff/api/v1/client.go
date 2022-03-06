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
	"time"
)

// New creates a new api.v1 API client that implements the BFF interface.
func New(endpoint string) (_ BFFClient, err error) {
	c := &APIv1{
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

// APIv1 implements the BFFClient interface.
type APIv1 struct {
	endpoint *url.URL
	client   *http.Client
}

// Ensure the API implments the BFFClient interface.
var _ BFFClient = &APIv1{}

//===========================================================================
// Client Methods
//===========================================================================

func (s *APIv1) Status(ctx context.Context) (out *StatusReply, err error) {
	//  Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/status", nil, nil); err != nil {
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
