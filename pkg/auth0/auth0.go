/*
Package auth0 implements a lightweight Go SDK to the Auth0 Management API.
Unfortunately, there is no supported Go SDK for Auth0 at the time of this writing. The
package is configured to connect to the Auth0 Management API with a client ID and client
secret, which wraps a basic http client for authorization and specific endpoint requests.
*/
package auth0

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Auth0 provides authenticated https requests to the Auth0 v2 Management API.
// See documentation at: https://auth0.com/docs/api/management/v2
//
// All endpoints require authentication via an access token. Access tokens are fetched
// using the client id and secret configured on the client. Because access tokens are
// machine to machine tokens and part of the paid subscription, the token is stored in
// memory for the duration of the token before a new one is fetched to ensure that as
// few tokens as possible are issued.
type Auth0 struct {
	sync.RWMutex
	client  http.Client
	creds   *Credentials
	conf    Config
	baseURL *url.URL
}

// Create a new Auth0 client with the specified configuration. If a zero-valued config
// is supplied then the config will be loaded from the environment. If there are valid
// credentials cached, then those credentials will be loaded.
func New(conf Config) (client *Auth0, err error) {
	if conf.IsZero() {
		if conf, err = NewConfig(); err != nil {
			return nil, err
		}
	} else {
		if err = conf.Validate(); err != nil {
			return nil, err
		}
	}

	client = &Auth0{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		conf:    conf,
		creds:   &Credentials{},
		baseURL: conf.BaseURL(),
	}

	// Load token from cache
	if conf.TokenCache != "" {
		if err = client.creds.LoadCache(conf.TokenCache); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// Authenticate the client with the configured client ID and client secret, getting an
// access token that can be used to make subsequent requests. In normal operations the
// access token is stored in memory and used until it expires (by default, 24 hours).
// The access token can also be cached to disk for testing or CLI programs because the
// number of M2M tokens issued per month is limited by the paid plan.
//
// For more on accessing the management API with production access tokens, see:
// https://auth0.com/docs/secure/tokens/access-tokens/get-management-api-access-tokens-for-production
func (a *Auth0) Authenticate(ctx context.Context) (err error) {
	// Create the form url-encoded payload
	payload := url.Values{}
	payload.Set(authGrantTypeKey, authGrantTypeVal)
	payload.Set(authClientIdKey, a.conf.ClientID)
	payload.Set(authClientSecretKey, a.conf.ClientSecret)
	payload.Set(authAudienceKey, a.Endpoint("/api/v2/", nil))

	var req *http.Request
	endpoint := a.Endpoint("/oauth/token", nil)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(payload.Encode())); err != nil {
		return fmt.Errorf("could not build http request: %s", err)
	}

	req.Header.Set("Content-Type", authContentType)
	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", userAgent)

	var rep *http.Response
	if rep, err = a.client.Do(req); err != nil {
		return fmt.Errorf("could not execute authentication request: %s", err)
	}
	defer rep.Body.Close()

	// Check status code
	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		return fmt.Errorf("could not complete authentication received %s", rep.Status)
	}

	// Create new credentials and parse incoming credentials from request
	a.creds = &Credentials{}
	if err = a.creds.Load(rep.Body); err != nil {
		return err
	}

	// If we have a cache for the access tokens, write it
	if a.conf.TokenCache != "" {
		a.creds.DumpCache(a.conf.TokenCache)
	}
	return nil
}

// Endpoint creates a valid URL from the specified path and query params using the base
// url specified via the configuration. String formatting can be used in the path and is
// processed before finalizing the URL.
func (a *Auth0) Endpoint(path string, query map[string]string, sfa ...interface{}) string {
	ref := &url.URL{Path: fmt.Sprintf(path, sfa...)}
	ep := a.baseURL.ResolveReference(ref)

	if len(query) > 0 {
		values := url.Values{}
		for key, val := range query {
			values.Set(key, val)
		}
		ep.RawQuery = values.Encode()
	}

	return ep.String()
}

// NewRequest creates an http request setting the default headers along with the
// authentication header. If the client has not been authenticated, then an error is
// returned. Ensure that preflight checks are run for authenticated requests.
func (a *Auth0) NewRequest(ctx context.Context, method, url string, data interface{}) (req *http.Request, err error) {
	a.RLock()
	defer a.RUnlock()
	if !a.creds.Valid() {
		return nil, ErrNotAuthenticated
	}

	if data != nil {
		// Serialize the JSON request data into the request
		body := new(bytes.Buffer)
		if err = json.NewEncoder(body).Encode(data); err != nil {
			return nil, err
		}
		if req, err = http.NewRequestWithContext(ctx, method, url, body); err != nil {
			return nil, err
		}
	} else {
		// Create a request with an empty body
		if req, err = http.NewRequestWithContext(ctx, method, url, nil); err != nil {
			return nil, err
		}
	}

	// Set Headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.creds.AccessToken))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

// Preflight prepares to send a request that needs to be authenticated by checking the
// credentials and sending any authentication requests that are required.
func (a *Auth0) Preflight(ctx context.Context) error {
	a.Lock()
	defer a.Unlock()
	if !a.creds.Valid() {
		if err := a.Authenticate(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Do performs an auth0 client request and returns the response. If a non-200 status is
// returned then the error response is parsed and an auth0 specific error is returned.
func (a *Auth0) Do(req *http.Request) (rep *http.Response, err error) {
	// Ensure that in testing mode requests are sent to a testing server
	if a.conf.Testing {
		host := req.URL.Hostname()
		if host != "localhost" && host != "127.0.0.1" {
			return nil, fmt.Errorf("hostname %q is not valid in testing mode", host)
		}
	}

	if rep, err = a.client.Do(req); err != nil {
		return nil, err
	}

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		// An error status code was returned, attempt to parse the error
		defer rep.Body.Close()
		var e *APIError

		// TODO: we could check the content-type to determine if its JSON before parsing
		if err = json.NewDecoder(rep.Body).Decode(&e); err != nil {
			// If we cannot decode the body return a generic error
			e = &APIError{StatusCode: rep.StatusCode, Status: rep.Status}
		}
		return nil, e
	}

	return rep, nil
}

// Creds returns a pointer to the credentials for testing purposes.
func (a *Auth0) Creds() *Credentials {
	return a.creds
}
