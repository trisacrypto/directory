package bff_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	"github.com/trisacrypto/directory/pkg/store"
	storeconfig "github.com/trisacrypto/directory/pkg/store/config"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// The BFF Test Suite provides mock functionality and fixtures for running BFF tests
// that expect to interact with two GDS services: TestNet and MainNet.
type bffTestSuite struct {
	suite.Suite
	bff     *bff.Server
	client  api.BFFClient
	testnet mockNetwork
	mainnet mockNetwork
	db      *mock.Trtl
	auth    *authtest.Server
}

type mockNetwork struct {
	db      *mock.Trtl
	gds     *mock.GDS
	members *mock.Members
}

func (s *bffTestSuite) SetupSuite() {
	var err error
	require := s.Require()

	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog MUST be false otherwise this will be overriden
	logger.Discard()

	// Setup a mock trtl database for the bff
	s.db, err = mock.NewTrtl()
	require.NoError(err, "could not create mock trtl database")

	// Start the authtest server for authentication verification
	s.auth, err = authtest.Serve()
	require.NoError(err, "could not start the authtest server")

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
		LoginURL:     "http://localhost/auth/login",
		RegisterURL:  "http://localhost/auth/register",
		Auth0:        s.auth.Config(),
		TestNet: config.NetworkConfig{
			Directory: config.DirectoryConfig{
				Insecure: true,
				Endpoint: "passthrough://bufnet",
				Timeout:  1 * time.Second,
			},
			Members: config.MembersConfig{
				Endpoint: "passthrough://bufnet",
				Timeout:  1 * time.Second,
				MTLS: config.MTLSConfig{
					Insecure: true,
				},
			},
		},
		MainNet: config.NetworkConfig{
			Directory: config.DirectoryConfig{
				Insecure: true,
				Endpoint: "passthrough://bufnet",
				Timeout:  1 * time.Second,
			},
			Members: config.MembersConfig{
				Endpoint: "passthrough://bufnet",
				Timeout:  1 * time.Second,
				MTLS: config.MTLSConfig{
					Insecure: true,
				},
			},
		},
		Database: storeconfig.StoreConfig{
			URL:      "trtl://bufnet/",
			Insecure: true,
		},
		Email: config.EmailConfig{
			ServiceEmail: "service@example.com",
			Testing:      true,
		},
		UserCache: config.CacheConfig{
			Enabled:    true,
			Expiration: 1 * time.Minute,
			Size:       100,
		},
	}.Mark()
	require.NoError(err, "could not mark configuration")

	// Create the database mocks for testnet and mainnet
	s.testnet.db, err = mock.NewTrtl()
	require.NoError(err, "could not create testnet trtl database")

	s.mainnet.db, err = mock.NewTrtl()
	require.NoError(err, "could not create mainnet trtl database")

	// Create the GDS mocks for testnet and mainnet
	s.testnet.gds, err = mock.NewGDS(conf.TestNet.Directory)
	require.NoError(err, "could not create testnet mock")

	s.mainnet.gds, err = mock.NewGDS(conf.MainNet.Directory)
	require.NoError(err, "could not create mainnet mock")

	// Create the members mocks for testnet and mainnet
	s.testnet.members, err = mock.NewMembers(conf.TestNet.Members)
	require.NoError(err, "could not create testnet mock")

	s.mainnet.members, err = mock.NewMembers(conf.MainNet.Members)
	require.NoError(err, "could not create mainnet mock")

	s.bff, err = bff.New(conf)
	require.NoError(err, "could not create the bff")

	// Create the mock testnet clients
	testnetClient := &bff.GDSClient{}
	require.NoError(testnetClient.ConnectGDS(conf.TestNet.Directory, s.testnet.gds.DialOpts()...), "could not connect to testnet GDS")
	require.NoError(testnetClient.ConnectMembers(conf.TestNet.Members, s.testnet.members.DialOpts()...), "could not connect to testnet members")

	// Create the mock mainnet client
	mainnetClient := &bff.GDSClient{}
	require.NoError(mainnetClient.ConnectGDS(conf.MainNet.Directory, s.mainnet.gds.DialOpts()...), "could not connect to mainnet GDS")
	require.NoError(mainnetClient.ConnectMembers(conf.MainNet.Members, s.mainnet.members.DialOpts()...), "could not connect to mainnet members")

	// Add the mock clients to the mock
	s.bff.SetGDSClients(testnetClient, mainnetClient)

	// Inject the database clients into the BFF
	s.bff.SetDB(s.DB())
	s.bff.SetTestNetDB(s.TestNetDB())
	s.bff.SetMainNetDB(s.MainNetDB())

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
	s.testnet.gds.Reset()
	s.mainnet.gds.Reset()
	s.testnet.members.Reset()
	s.mainnet.members.Reset()

	// Ensure any credentials set on the client are reset
	s.client.(*api.APIv1).SetCredentials(nil)
	s.client.(*api.APIv1).SetCSRFProtect(false)
}

func (s *bffTestSuite) TearDownSuite() {
	require := s.Require()
	require.NoError(s.bff.Shutdown(), "could not shutdown the BFF server after tests")
	s.testnet.gds.Shutdown()
	s.mainnet.gds.Shutdown()
	s.testnet.members.Shutdown()
	s.mainnet.members.Shutdown()

	// Shutdown the authtest server
	authtest.Close()

	// Shutdown and cleanup the trtls
	s.db.Shutdown()
	s.testnet.db.Shutdown()
	s.mainnet.db.Shutdown()

	// Cleanup logger
	logger.ResetLogger()
}

func TestBFF(t *testing.T) {
	suite.Run(t, new(bffTestSuite))
}

// Helper function to set the credentials on the test client from claims, reducing 3 or
// 4 lines of code into a single helper function call to make tests more readable.
func (s *bffTestSuite) SetClientCredentials(claims *authtest.Claims) error {
	token, err := s.auth.NewTokenWithClaims(claims)
	if err != nil {
		return err
	}

	s.client.(*api.APIv1).SetCredentials(api.Token(token))
	return nil
}

// Helper function to set cookies for CSRF protection on the BFF client
func (s *bffTestSuite) SetClientCSRFProtection() error {
	s.client.(*api.APIv1).SetCSRFProtect(true)
	return nil
}

// DB returns the configured database client for the BFF
func (s *bffTestSuite) DB() store.Store {
	client, err := s.db.Client()
	s.Require().NoError(err, "could not init BFF database client")
	return client
}

// TestnetDB returns the configured database client for the testnet
func (s *bffTestSuite) TestNetDB() store.Store {
	client, err := s.testnet.db.Client()
	s.Require().NoError(err, "could not init testnet database client")
	return client
}

// MainnetDB returns the configured database client for the mainnet
func (s *bffTestSuite) MainNetDB() store.Store {
	client, err := s.mainnet.db.Client()
	s.Require().NoError(err, "could not init mainnet database client")
	return client
}

// Resets the BFF Trtl database to an initial state by deleting the current database
// and opening a new one. Tests should defer this method if they modify the BFF
// database.
func (s *bffTestSuite) ResetDB() {
	var err error
	require := s.Require()
	s.db.Shutdown()
	s.db, err = mock.NewTrtl()
	require.NoError(err, "could not create new BFF trtl database")
	client, err := s.db.Client()
	require.NoError(err, "could not create new BFF trtl client")
	s.bff.SetDB(client)
}

// Resets the testnet Trtl database to an initial state by deleting the current
// database and opening a new one. Tests should defer this method if they modify the
// testnet database.
func (s *bffTestSuite) ResetTestNetDB() {
	var err error
	require := s.Require()
	s.testnet.db.Shutdown()
	s.testnet.db, err = mock.NewTrtl()
	require.NoError(err, "could not create new testnet trtl database")
	client, err := s.testnet.db.Client()
	require.NoError(err, "could not create new testnet trtl client")
	s.bff.SetTestNetDB(client)
}

// Resets the mainnet Trtl database to an initial state by deleting the current
// database and opening a new one. Tests should defer this method if they modify the
// mainnet database.
func (s *bffTestSuite) ResetMainNetDB() {
	var err error
	require := s.Require()
	s.mainnet.db.Shutdown()
	s.mainnet.db, err = mock.NewTrtl()
	require.NoError(err, "could not create new mainnet trtl database")
	client, err := s.mainnet.db.Client()
	require.NoError(err, "could not create new mainnet trtl client")
	s.bff.SetMainNetDB(client)
}

// Custom assertion to ensure a formatted error contains the correct status code and
// message.
func (s *bffTestSuite) requireError(err error, status int, message string, msgAndArgs ...interface{}) {
	require := s.Require()
	require.EqualError(err, fmt.Sprintf("[%d] %s", status, message), msgAndArgs...)
}

// Helper function to load test fixtures from disk. If v is a proto.Message it is loaded
// using protojson, otherwise it is loaded using encoding/json.
func loadFixture(path string, v interface{}) (err error) {
	switch t := v.(type) {
	case proto.Message:
		return loadPBFixture(path, t)
	default:
		return loadJSONFixture(path, t)
	}
}

func loadPBFixture(path string, v proto.Message) (err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return err
	}

	pbjson := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	if err = pbjson.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

func loadJSONFixture(path string, v interface{}) (err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}
