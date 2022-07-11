package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff/admin/mock"
	apiv2 "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

const (
	StatusEP           = "Status"
	ListCertificatesEP = "ListCertificates"
)

func NewAdmin() (a *Admin, err error) {
	a = &Admin{}

	// Create the mock token manager for authentication
	if a.tokens, err = tokens.MockTokenManager(); err != nil {
		return nil, err
	}

	// Configure the router for the mock
	gin.SetMode(gin.TestMode)
	a.router = gin.New()
	a.authorize = apiv2.Authorization(a.tokens)
	a.csrf = apiv2.DoubleCookie()
	a.setupHandlers()

	a.srv = httptest.NewTLSServer(a.router)
	if a.client, err = mock.New(a.tokens, a.srv); err != nil {
		return nil, err
	}

	return a, nil
}

// Admin is a mock implementation the DirectoryAdministrationServer using an httptest
// server. The Client() method provides access to the DirectoryAdministrationClient
// that can be used by tests to make calls to the Admin API using the BFF client. This
// allows the BFF tests to set specific responses for individual endpoints. To
// configure the behavior for specific endpoints, use the UseFixture, UseError, or
// UseHandler methods.
type Admin struct {
	srv       *httptest.Server
	router    *gin.Engine
	client    apiv2.DirectoryAdministrationClient
	tokens    *tokens.TokenManager
	authorize gin.HandlerFunc
	csrf      gin.HandlerFunc
	handlers  map[string]gin.HandlerFunc
	Calls     map[string]int
}

// UseFixture allows tests to use a JSON fixture from disk as the response.
func (a *Admin) UseFixture(endpoint, path string) (err error) {
	// Read the fixture data from disk
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return fmt.Errorf("could not read fixture data: %s", err)
	}

	switch endpoint {
	case StatusEP:
		out := &apiv2.StatusReply{}
		if err = json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal fixture data: %s", err)
		}
		a.handle(endpoint, func(c *gin.Context) {
			c.JSON(http.StatusOK, out)
		})
	case ListCertificatesEP:
		out := &apiv2.ListCertificatesReply{}
		if err = json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal fixture data: %s", err)
		}
		a.handle(endpoint, func(c *gin.Context) {
			c.JSON(http.StatusOK, out)
		})
	default:
		return fmt.Errorf("unsupported endpoint: %s", endpoint)
	}

	return nil
}

// UseError allows tests to specify an HTTP error for the response.
func (a *Admin) UseError(endpoint string, code int, msg string) {
	a.handle(endpoint, func(c *gin.Context) {
		c.JSON(code, msg)
	})
}

// UseHandler allows tests to specify a custom handler for the endpoint.
func (a *Admin) UseHandler(endpoint string, handler gin.HandlerFunc) {
	a.handle(endpoint, handler)
}

// handle is a helper function that adds a handler to the handlers map and
// returns that handler function when the endpoint is called.
func (a *Admin) handle(endpoint string, handler gin.HandlerFunc) {
	if _, ok := a.handlers[endpoint]; !ok {
		panic(fmt.Sprintf("endpoint %s not initialized", endpoint))
	}
	a.handlers[endpoint] = handler
}

// initHandler initializes a handler for the given endpoint on the gin router with the
// provided method, path, and middleware. After calling this function the configured
// handler for the endpoint will return a 200 and an empty response. This should only
// be called once for each endpoint to avoid a gin panic. UseFixture, UseError, or
// UseHandler should be used to configure different endpoint responses from the tests.
func (a *Admin) initHandler(endpoint, method, path string, middleware ...gin.HandlerFunc) {
	// Copy any provided middleware into the handlers slice
	handlers := make([]gin.HandlerFunc, len(middleware), len(middleware)+1)
	copy(handlers, middleware)

	a.handlers[endpoint] = func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	}

	handlerFunc := func(c *gin.Context) {
		if handler, ok := a.handlers[endpoint]; ok {
			a.Calls[endpoint]++
			handler(c)
		} else {
			c.JSON(http.StatusNotFound, fmt.Sprintf("unhandled endpoint %s", endpoint))
		}
	}
	handlers = append(handlers, handlerFunc)

	a.router.Handle(method, path, handlers...)
}

func (a *Admin) setupHandlers() {
	a.handlers = make(map[string]gin.HandlerFunc)
	a.Calls = make(map[string]int)

	a.initHandler(StatusEP, http.MethodGet, "/v2/status")
	a.initHandler(ListCertificatesEP, http.MethodGet, "/v2/vasps/:vaspID/certificates", a.authorize)
}

func (a *Admin) Shutdown() {
	a.srv.Close()
}

func (a *Admin) Client() apiv2.DirectoryAdministrationClient {
	return a.client
}
