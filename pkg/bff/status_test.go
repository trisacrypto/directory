package bff_test

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc/codes"
)

func (s *bffTestSuite) TestStatus() {
	var err error
	require := s.Require()

	// Set the Status RPC for the mocks
	healthy := func(ctx context.Context, in *gds.HealthCheck) (*gds.ServiceState, error) {
		return &gds.ServiceState{
			Status:    gds.ServiceState_HEALTHY,
			NotBefore: time.Now().Add(1 * time.Hour).Format(time.RFC3339Nano),
			NotAfter:  time.Now().Format(time.RFC3339Nano),
		}, nil
	}
	s.testnet.OnStatus = healthy
	s.mainnet.OnStatus = healthy

	// Test Status with nil params (should default to getting GDS status)
	rep, err := s.client.Status(context.TODO(), nil)
	require.NoError(err, "could not make status request")
	require.Equal("ok", rep.Status, "server did not return an ok status")
	require.NotEmpty(rep.Uptime, "server did not return an uptime")
	require.NotEmpty(rep.Version, "server did not return a version")
	require.Equal("healthy", rep.TestNet)
	require.Equal("healthy", rep.MainNet)
	require.Equal(s.testnet.Calls["Status"], 1)
	require.Equal(s.mainnet.Calls["Status"], 1)

	// Test Status with default params
	rep, err = s.client.Status(context.TODO(), &api.StatusParams{NoGDS: false})
	require.NoError(err, "could not make status request")
	require.Equal("ok", rep.Status, "server did not return an ok status")
	require.NotEmpty(rep.Uptime, "server did not return an uptime")
	require.NotEmpty(rep.Version, "server did not return a version")
	require.Equal("healthy", rep.TestNet)
	require.Equal("healthy", rep.MainNet)
	require.Equal(s.testnet.Calls["Status"], 2)
	require.Equal(s.mainnet.Calls["Status"], 2)

	// Test Status with NoGDS
	rep, err = s.client.Status(context.TODO(), &api.StatusParams{NoGDS: true})
	require.NoError(err, "could not make status request")
	require.Equal("ok", rep.Status, "server did not return an ok status")
	require.NotEmpty(rep.Uptime, "server did not return an uptime")
	require.NotEmpty(rep.Version, "server did not return a version")
	require.Empty(rep.TestNet)
	require.Empty(rep.MainNet)
	require.Equal(s.testnet.Calls["Status"], 2)
	require.Equal(s.mainnet.Calls["Status"], 2)

	// Test Status with TestNet Error
	s.testnet.UseError(mock.StatusRPC, codes.Unavailable, "unreachable host")
	rep, err = s.client.Status(context.TODO(), nil)
	require.NoError(err, "could not make status request")
	require.Equal("ok", rep.Status, "server did not return an ok status")
	require.NotEmpty(rep.Uptime, "server did not return an uptime")
	require.NotEmpty(rep.Version, "server did not return a version")
	require.Equal("unavailable", rep.TestNet)
	require.Equal("healthy", rep.MainNet)
	require.Equal(s.testnet.Calls["Status"], 3)
	require.Equal(s.mainnet.Calls["Status"], 3)

	// Test Status with MainNet bad state
	s.mainnet.OnStatus = func(ctx context.Context, in *gds.HealthCheck) (*gds.ServiceState, error) {
		return &gds.ServiceState{
			Status:    gds.ServiceState_DANGER,
			NotBefore: time.Now().Add(1 * time.Hour).Format(time.RFC3339Nano),
			NotAfter:  time.Now().Format(time.RFC3339Nano),
		}, nil
	}
	rep, err = s.client.Status(context.TODO(), nil)
	require.NoError(err, "could not make status request")
	require.Equal("ok", rep.Status, "server did not return an ok status")
	require.NotEmpty(rep.Uptime, "server did not return an uptime")
	require.NotEmpty(rep.Version, "server did not return a version")
	require.Equal("unavailable", rep.TestNet)
	require.Equal("danger", rep.MainNet)
	require.Equal(s.testnet.Calls["Status"], 4)
	require.Equal(s.mainnet.Calls["Status"], 4)

}

func (s *bffTestSuite) TestAvailableMiddleware() {
	var err error
	require := s.Require()

	// Set the health to false to mimic the server stopping
	s.bff.SetHealth(false)
	defer s.bff.SetHealth(true)

	// Perform a status request, the Available middleware should intercept and return "stopping"
	rep, err := s.client.Status(context.TODO(), nil)
	require.NoError(err, "an error was returned from the status endpoint")
	require.Equal("stopping", rep.Status, "incorrect response from unhealthy server")
}

func (s *bffTestSuite) TestMaintenanceMode() {
	require := s.Require()

	// To test maintenance mode, we need to create a BFF server that is in maintenance
	// mode, which means we cannot use the BFF fixture created during SetupSuite.
	// A minimal maintenance mode configuration.
	conf, err := config.Config{
		Maintenance:  true,
		Mode:         gin.TestMode,
		AllowOrigins: []string{"http://localhost"},
		CookieDomain: "localhost",
		TestNet:      config.DirectoryConfig{Endpoint: "bufcon"},
		MainNet:      config.DirectoryConfig{Endpoint: "bufcon"},
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
