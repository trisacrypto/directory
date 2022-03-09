package bff_test

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/config"
)

func (s *bffTestSuite) TestStatus() {
	var err error
	require := s.Require()

	rep, err := s.client.Status(context.TODO())
	require.NoError(err, "could not make status request")
	require.Equal("ok", rep.Status, "server did not return an ok status")
	require.NotEmpty(rep.Uptime, "server did not return an uptime")
	require.NotEmpty(rep.Version, "server did not return a version")
}

func (s *bffTestSuite) TestAvailableMiddleware() {
	var err error
	require := s.Require()

	// Set the health to false to mimic the server stopping
	s.bff.SetHealth(false)
	defer s.bff.SetHealth(true)

	// Perform a status request, the Available middleware should intercept and return "stopping"
	rep, err := s.client.Status(context.TODO())
	require.NoError(err, "an error was returned from the status endpoint")
	require.Equal("stopping", rep.Status, "incorrect response from unhealthy server")
}

func (s *bffTestSuite) TestMaintenanceMode() {
	require := s.Require()

	// To test maintenance mode, we need to create a BFF server that is in maintenance
	// mode, which means we cannot use the BFF fixture created during SetupSuite.
	// A minimal maintenance mode configuration.
	conf, err := config.Config{
		Maintenance: true,
		Mode:        gin.TestMode,
		TestNet:     config.DirectoryConfig{Endpoint: "bufcon"},
		MainNet:     config.DirectoryConfig{Endpoint: "bufcon"},
	}.Mark()
	require.NoError(err, "configuration is not valid")

	// Run the maintenance mode server
	server, err := bff.New(conf)
	require.NoError(err, "could not create maintenance mode server")
	go server.Serve()
	defer server.Shutdown()

	time.Sleep(500 * time.Millisecond)

	// Create an http client that sends GET requests to random endpoints, all requests
	// should return a 503 error indicating that the service is unavailable.
	client := &http.Client{}

	for _, path := range []string{"/", "/v1", "/v1/", "/v1/status", "/v1/register", "/status"} {
		url := server.GetURL() + path
		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(err, "could not make HTTP request")
		rep, err := client.Do(req)
		require.NoError(err, "could not execute HTTP request")
		require.Equal(http.StatusServiceUnavailable, rep.StatusCode, "expected unavailable status")
	}
}
