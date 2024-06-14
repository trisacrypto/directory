package gds_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/gds"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/fixtures"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

const (
	bufSize = 1024 * 1024
)

var (
	fixturesPath = filepath.Join("testdata", "fakes.tgz")
	dbPath       = filepath.Join("testdata", "db")
)

// The GDS Test Suite provides mock functionality and database fixtures for the services
// defined in the GDS package. Most tests in this package should be methods of the Test
// Suite. On startup, the reference fixtures are loaded from the `fakes.tgz` dataset and
// a mock service is created with an empty database. If tests require fixtures to be
// loaded they should call the loadFixtures() or loadSmallFixtures() methods to point
// the mock service at a database that has those fixtures or to loadEmptyFixtures() if
// they require an empty database. If the test modifies the database they should defer
// a call to resetFixtures, resetSmallFixtures, or resetEmptyFixtures as necessary.
//
// Tests should use accessor methods such as s.svc.GetAdmin() or s.svc.GetStore() to
// access the internals of the service and services for testing purposes.
type gdsTestSuite struct {
	suite.Suite
	svc            *gds.Service
	grpc           *bufconn.GRPCListener
	conf           *config.Config
	fixtures       *fixtures.Library
	expectedEmails emails.ExpectedEmails
}

// SetConfig allows a custom config to be specified by the tests.
// Note that loadFixtures() needs to be called in order for the config to be used.
func (s *gdsTestSuite) SetConfig(conf config.Config) {
	s.conf = &conf
}

// ResetConfig back to the default.
func (s *gdsTestSuite) ResetConfig() {
	s.conf = nil
}

func (s *gdsTestSuite) SetupSuite() {
	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog MUST be false otherwise this will be overwritten
	logger.Discard()
	gin.SetMode(gin.TestMode)

	// Start with an empty fixtures service
	s.LoadEmptyFixtures()
}

// SetupGDS starts the GDS server
// Run this inside the test methods after loading the appropriate fixtures
func (s *gdsTestSuite) SetupGDS() {

	// Using a bufconn listener allows us to avoid network requests
	s.grpc = bufconn.New(bufSize, "")
	go s.svc.GetGDS().Run(s.grpc.Listener)
}

// SetupMembers starts the Members server
// Run this inside the test methods after loading the appropriate fixtures
func (s *gdsTestSuite) SetupMembers() {

	// Using a bufconn listener allows us to avoid network requests
	s.grpc = bufconn.New(bufSize, "")
	go s.svc.GetMembers().Run(s.grpc.Listener)
}

// Helper function to shutdown any previously running GDS or Members servers and release the gRPC connection
func (s *gdsTestSuite) shutdownServers() {
	// Shutdown old GDS and Members servers, if they exist
	if s.svc != nil {
		if err := s.svc.GetGDS().Shutdown(); err != nil {
			log.Warn().Err(err).Msg("could not shutdown GDS server to start new one")
		}
		if err := s.svc.GetMembers().Shutdown(); err != nil {
			log.Warn().Err(err).Msg("could not shutdown Members server to start new one")
		}
	}
	if s.grpc != nil {
		s.grpc.Release()
	}
}

func (s *gdsTestSuite) TearDownSuite() {
	if s.svc != nil && s.svc.GetStore() != nil {
		s.svc.GetStore().Close()
	}

	s.shutdownServers()
	s.fixtures.Close()
	os.RemoveAll(dbPath)
	logger.ResetLogger()
}

func TestGDSLevelDB(t *testing.T) {
	var err error
	s := new(gdsTestSuite)
	s.fixtures, err = fixtures.New(fixturesPath, dbPath, fixtures.StoreLevelDB)
	require.NoError(t, err)
	suite.Run(t, s)
}

func TestGDSTrtl(t *testing.T) {
	var err error
	s := new(gdsTestSuite)
	s.fixtures, err = fixtures.New(fixturesPath, dbPath, fixtures.StoreTrtl)
	require.NoError(t, err)
	suite.Run(t, s)
}

//===========================================================================
// Test fixtures management
//===========================================================================

// loadFixtures loads a new set of fixtures into the database. This method must respect
// the ftype variable on the test suite, which indicates which fixtures are currently
// loaded. If the ftype is different than the indicated fixture type, then this causes
// the current database to be completely overwritten with the indicated fixtures. Tests
// that require database fixtures should call the appropriate load method, either
// LoadFullFixtures, LoadSmallFixtures, or LoadEmptyFixtures to ensure that the correct
// fixtures are present before test execution.
func (s *gdsTestSuite) loadFixtures(ftype fixtures.FixtureType) {
	// If we're already at the specified fixture type and no custom config is provided,
	// do nothing
	if s.fixtures.FixtureType() == ftype && s.conf == nil {
		log.Info().Uint8("ftype", uint8(ftype)).Msg("CACHED FIXTURE")
		return
	}

	// Close the current service
	var err error
	require := s.Require()
	if s.svc != nil && s.svc.GetStore() != nil {
		if err := s.svc.GetStore().Close(); err != nil {
			log.Warn().Err(err).Msg("could not close service store to load new fixtures")
		}
	}

	s.shutdownServers()

	// Load the new fixtures
	require.NoError(s.fixtures.Load(ftype))

	// Use the custom config if specified
	var conf config.Config
	if s.conf != nil {
		conf = *s.conf
	} else {
		conf = gds.MockConfig()
	}

	// Create the new GDS with the configured store
	switch s.fixtures.StoreType() {
	case fixtures.StoreLevelDB:
		conf.Database.URL = "leveldb:///" + s.fixtures.DBPath()
		s.svc, err = gds.NewMock(conf, nil)
		require.NoError(err, "could not create mock GDS with leveldb store")
	case fixtures.StoreTrtl:
		conn, err := s.fixtures.ConnectTrtl(context.Background())
		require.NoError(err, "could not connect to trtl server")
		s.svc, err = gds.NewMock(conf, conn)
		require.NoError(err, "could not create mock GDS with trtl store")
	default:
		require.Fail("unrecognized store type %d", s.fixtures.StoreType())
	}

	// Create the expected emails factory from the configuration.
	s.expectedEmails = emails.ExpectedEmailsFactory(conf.Email.ServiceEmail)

	log.Info().Uint8("ftype", uint8(ftype)).Msg("FIXTURE LOADED")
}

func (s *gdsTestSuite) LoadEmptyFixtures() {
	s.loadFixtures(fixtures.Empty)
}

// LoadFullFixtures loads the JSON test fixtures from disk and stores them in the dbFixtures map.
func (s *gdsTestSuite) LoadFullFixtures() {
	s.loadFixtures(fixtures.Full)
}

func (s *gdsTestSuite) LoadSmallFixtures() {
	s.loadFixtures(fixtures.Small)
}

// ResetFixtures uncaches the current database which causes the next call to
// loadFixtures to generate a new database that overwrites the current one. Tests that
// modify the database should call ResetFixtures to ensure that the fixtures are reset
// for the next test.
func (s *gdsTestSuite) ResetFixtures() {
	// Set the ftype to unknown to ensure that loadFixtures loads the fixture.
	s.fixtures.Reset()
}

// SetVerificationStatus sets the verification status of a VASP fixture on the
// database. This is useful for testing VerificationState checks without having to use
// multiple fixtures.
func (s *gdsTestSuite) SetVerificationStatus(id string, status pb.VerificationState) {
	require := s.Require()

	// Retrieve the VASP from the database
	vasp, err := s.svc.GetStore().RetrieveVASP(context.Background(), id)
	require.NoError(err, "VASP not found in database")

	// Set the verification status and write back to the database
	vasp.VerificationStatus = status
	require.NoError(s.svc.GetStore().UpdateVASP(context.Background(), vasp), "could not update VASP")
}
