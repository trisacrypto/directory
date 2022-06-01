package bff_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/db"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

// loadFixture is a helper function to return an unwired JSON protocol buffer for the
// the BFF client to post to the server, which will then rewire it for GDS requests.
func loadFixture(path string) (fixture map[string]interface{}, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}

	fixture = make(map[string]interface{})
	if err = json.Unmarshal(data, &fixture); err != nil {
		return nil, err
	}
	return fixture, nil
}

// The BFF Test Suite provides mock functionality and fixtures for running BFF tests
// that expect to interact with two GDS services: TestNet and MainNet.
type bffTestSuite struct {
	suite.Suite
	bff      *bff.Server
	client   api.BFFClient
	testnet  *mock.GDS
	mainnet  *mock.GDS
	dbPath   string
	trtl     *trtl.Server
	trtlsock *bufconn.GRPCListener
}

func (s *bffTestSuite) SetupSuite() {
	var err error
	require := s.Require()

	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog MUST be false otherwise this will be overriden
	logger.Discard()

	// Setup a mock trtl server for the tests
	s.SetupTrtl()

	// This configuration will run the BFF as a fully functional server on an open port
	// on the system for local loop-back only. It is also in test mode, so a Gin context
	// can also be used to test endpoints. The BFF server runs for the duration of the
	// tests and must be shutdown when the test suite terminates.
	conf, err := config.Config{
		Maintenance:  false,
		BindAddr:     "127.0.0.1:0",
		Mode:         gin.TestMode,
		LogLevel:     logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:   false,
		AllowOrigins: []string{"http://localhost"},
		CookieDomain: "localhost",
		Auth0: config.AuthConfig{
			Issuer:        "http://auth.localhost/",
			Audience:      "http://localhost",
			ProviderCache: 5 * time.Minute,
		},
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
		Database: config.DatabaseConfig{
			URL:      "trtl://bufnet/",
			Insecure: true,
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

	// Direct connect the BFF server to the database
	db, err := db.DirectConnect(s.trtlsock.Conn)
	require.NoError(err, "could not direct connect db to the BFF server")
	s.bff.SetDB(db)

	// Start the BFF server - the goal of the BFF tests is to have the server run for
	// the entire duration of the tests. Implement reset methods to ensure the server
	// state doesn't change between tests.
	go s.bff.Serve()

	// Wait for 500 ms to ensure the BFF starts serving
	time.Sleep(500 * time.Millisecond)

	// Create the BFF client for making requests to the server
	require.NotEmpty(s.bff.GetURL(), "no url to connect to client on")
	s.client, err = api.New(s.bff.GetURL())
	require.NoError(err, "could not initialize BFF client")
}

func (s *bffTestSuite) AfterTest(suiteName, testName string) {
	s.testnet.Reset()
	s.mainnet.Reset()
}

func (s *bffTestSuite) TearDownSuite() {
	require := s.Require()
	require.NoError(s.bff.Shutdown(), "could not shutdown the BFF server after tests")
	s.testnet.Shutdown()
	s.mainnet.Shutdown()

	// Shutdown and cleanup trtl
	s.trtl.Shutdown()
	s.trtlsock.Release()
	os.RemoveAll(s.dbPath)
	s.dbPath = ""

	// Cleanup logger
	logger.ResetLogger()
}

func TestBFF(t *testing.T) {
	suite.Run(t, new(bffTestSuite))
}

// SetupTrtl starts a Trtl server on a bufconn for testing with the BFF
func (s *bffTestSuite) SetupTrtl() {
	var err error
	require := s.Require()

	// Create a temporary directory for the testing database
	s.dbPath, err = ioutil.TempDir("", "trtldb-*")
	require.NoError(err, "could not create a temporary directory for trtl")

	conf := trtl.MockConfig()
	conf.Database.URL = "leveldb:///" + s.dbPath
	conf, err = conf.Mark()
	require.NoError(err, "could not validate mock config")

	// Start the Trtl server
	s.trtl, err = trtl.New(conf)
	require.NoError(err, "could not start trtl server")

	s.trtlsock = bufconn.New(1024 * 1024)
	go s.trtl.Run(s.trtlsock.Listener)

	// Connect to the running trtl server
	require.NoError(s.trtlsock.Connect(), "could not connect to trtl socket")
}
