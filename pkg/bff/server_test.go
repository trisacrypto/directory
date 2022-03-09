package bff_test

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

// The BFF Test Suite provides mock functionality and fixtures for running BFF tests
// that expect to interact with two GDS services: TestNet and MainNet.
type bffTestSuite struct {
	suite.Suite
	bff     *bff.Server
	client  api.BFFClient
	testnet *mock.GDS
	mainnet *mock.GDS
}

func (s *bffTestSuite) SetupSuite() {
	var err error
	require := s.Require()

	// This configuration will run the BFF as a fully functional server on an open port
	// on the system for local loop-back only. It is also in test mode, so a Gin context
	// can also be used to test endpoints. The BFF server runs for the duration of the
	// tests and must be shutdown when the test suite terminates.
	conf, err := config.Config{
		Maintenance: false,
		BindAddr:    "127.0.0.1:0",
		Mode:        gin.TestMode,
		LogLevel:    logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:  true,
		TestNet: config.DirectoryConfig{
			Insecure: true,
			Endpoint: "bufnet",
			Timeout:  1 * time.Second,
		},
		MainNet: config.DirectoryConfig{
			Insecure: true,
			Endpoint: "bufnet",
			Timeout:  1 * time.Second,
		},
	}.Mark()
	require.NoError(err, "could not mark configuration")

	// Create the GDS mocks for testnet and mainnet
	s.testnet, err = mock.NewGDS(conf.TestNet)
	require.NoError(err, "could not create testnet mock")

	s.mainnet, err = mock.NewGDS(conf.MainNet)
	require.NoError(err, "could not create mainnet mock")

	s.bff, err = bff.New(conf)
	require.NoError(err, "could not create the bff")

	// Add the GDS mock clients to the BFF server
	tnClient, err := s.testnet.Client()
	require.NoError(err, "could not create testnet GDS client")
	mnClient, err := s.mainnet.Client()
	require.NoError(err, "could not create mainnet GDS client")
	s.bff.SetClients(tnClient, mnClient)

	// Start the BFF server - the goal of the BFF tests is to have the server run for
	// the entire duration of the tests. Implement reset methods to ensure the server
	// state doesn't change between tests.
	go s.bff.Serve()

	// Create the BFF client for making requests to the server
	s.client, err = api.New(s.bff.GetURL())
	require.NoError(err, "could not initialize BFF client")
}

func (s *bffTestSuite) TearDownSuite() {
	require := s.Require()
	err := s.bff.Shutdown()
	require.NoError(err, "could not shutdown the BFF server after tests")
}

func TestBFF(t *testing.T) {
	suite.Run(t, new(bffTestSuite))
}
