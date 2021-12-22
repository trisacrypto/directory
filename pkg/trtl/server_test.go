package trtl_test

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/options"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
	"github.com/trisacrypto/directory/pkg/utils/logger"

	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	metaRegion = "tauceti"
	metaOwner  = "taurian"
	metaPID    = 8
	bufSize    = 1024 * 1024
)

var (
	// CLI flag, specify go test -update to regenerate static test fixtures
	update = flag.Bool("update", false, "update the static test fixtures")

	// Path to the honu database fixture archive
	dbtgz = filepath.Join("testdata", "db.tgz")

	// dbFixtures is a test fixtures map used to both generate the static test database
	// and verify the results of the trtl DB calls. It maps a test fixture name to an entry
	// containing a namespace, key, and value stored in the database.
	// Tests should not modify these fixtures!
	dbFixtures = map[string]*dbEntry{}
)

type dbEntry struct {
	Namespace string                 `json:"namespace"`
	Key       string                 `json:"key"`
	Value     map[string]interface{} `json:"value"`
}

//===========================================================================
// Test Suite
//===========================================================================

type trtlTestSuite struct {
	suite.Suite
	trtl  *trtl.Server
	conf  *config.Config
	grpc  *bufconn.GRPCListener
	tmpdb string
}

func (s *trtlTestSuite) SetupSuite() {
	require := s.Require()

	// Load the fixtures then regenerate the test database if requested or required.
	require.NoError(s.loadFixtures())
	if _, err := os.Stat(dbtgz); *update || os.IsNotExist(err) {
		require.NoError(s.generateDB())
	}

	// Create the initial configuration, database fixture, and servers for the tests
	require.NoError(s.setupConfig())
	require.NoError(s.extractDB())
	require.NoError(s.setupServers())
}

func (s *trtlTestSuite) TearDownSuite() {
	require := s.Require()
	require.NoError(s.cleanup())
}

//===========================================================================
// Server Tests
//===========================================================================

func TestTrtl(t *testing.T) {
	suite.Run(t, new(trtlTestSuite))
}

func (s *trtlTestSuite) TestMaintenance() {
	// Becasue we're modifying the configuration, ensure we reset the test environment
	defer s.reset()
	require := s.Require()

	s.conf.Maintenance = true
	server, err := trtl.New(*s.conf)
	require.NotEmpty(server, "no maintenance mode server was returned")
	require.NoError(err, "starting the server in maintenance mode caused an error")
	require.Nil(server.GetDB(), "maintenance mode database was not nil")
}

//===========================================================================
// Test Assertions
//===========================================================================

// StatusError is a helper assertion function that checks a gRPC status error
func (s *trtlTestSuite) StatusError(err error, code codes.Code, theError string) {
	require := s.Require()
	require.Error(err, "no status error returned")

	var serr *status.Status
	serr, ok := status.FromError(err)
	require.True(ok, "error is not a grpc status error")
	require.Equal(code, serr.Code(), "status code does not match")
	require.Equal(theError, serr.Message(), "status error message does not match")
}

// EqualMeta is a helper assertion function that checks if the actual metadata matches
// expectations about what the version should be. This helper function relies on some
// test fixture information to minimize what test users must supply.
func (s *trtlTestSuite) EqualMeta(expectedKey []byte, expectedNamespace string, expectedVersion, expectedParent *pb.Version, actualMeta *pb.Meta) {
	require := s.Require()
	require.NotNil(actualMeta, "cannot compare actual meta to expectations since actual meta is nil")

	expectedMeta := &pb.Meta{
		Key:       expectedKey,
		Namespace: expectedNamespace,
		Region:    metaRegion,
		Owner:     fmt.Sprintf("%d:%s", metaPID, metaOwner),
		Version:   expectedVersion,
		Parent:    expectedParent,
	}

	// Protocol buffer sanity check, this protects us from the case where the protocol buffers
	// have changed but the tests haven't been updated.
	if proto.Equal(actualMeta, expectedMeta) {
		// If the protocol buffers are equal, then we expect that everything is ok
		return
	}

	// At this point it's up to us to determine what the difference is between the versions
	require.Equal(expectedMeta.Key, actualMeta.Key, "meta keys do not match")
	require.Equal(expectedMeta.Namespace, actualMeta.Namespace, "meta namespace does not match")
	require.Equal(expectedMeta.Region, actualMeta.Region, "meta regions do not match")
	require.Equal(expectedMeta.Owner, actualMeta.Owner, "meta owners do not match")
	s.EqualVersion(expectedMeta.Version, actualMeta.Version, "version")
	s.EqualVersion(expectedMeta.Parent, actualMeta.Parent, "parent")

	// Could not determine the difference in the protocol buffers, so error generically
	require.True(proto.Equal(actualMeta, expectedMeta), "actual and expected protocol buffer metadata does not match")
}

// EqualVersion compares to Versions to see if they are the same
func (s *trtlTestSuite) EqualVersion(expectedVersion, actualVersion *pb.Version, versionType string) {
	require := s.Require()
	if expectedVersion == nil {
		require.Nil(actualVersion, "expected %s is nil but actual %s is not", versionType, versionType)
		return
	}

	require.Equal(expectedVersion.Pid, actualVersion.Pid, "expected %s PID does not match actual %s PID", versionType, versionType)
	require.Equal(expectedVersion.Version, actualVersion.Version, "expected %s scalar does not match actual %s scalar", versionType, versionType)
	require.Equal(expectedVersion.Region, actualVersion.Region, "expected %s region does not match actual %s region", versionType, versionType)
}

//===========================================================================
// Test setup helpers
//===========================================================================

// creates a valid config for the tests so long as the current config is empty
func (s *trtlTestSuite) setupConfig() (err error) {
	if s.conf != nil || s.tmpdb != "" {
		return errors.New("cannot create configuration, run test suite cleanup first")
	}

	// Create a tmp directory for the database
	if s.tmpdb, err = ioutil.TempDir("testdata", "trtldb-*"); err != nil {
		return fmt.Errorf("could not create tmpdb: %s", err)
	}

	// Create the configuration without loading it from the environment
	conf := config.Config{
		Maintenance: false,
		BindAddr:    ":4436",
		LogLevel:    logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:  true,
		Database: config.DatabaseConfig{
			URL:           fmt.Sprintf("leveldb:///%s", s.tmpdb),
			ReindexOnBoot: false,
		},
		Replica: config.ReplicaConfig{
			Enabled:        true,
			PID:            metaPID,
			Region:         metaRegion,
			Name:           metaOwner,
			GossipInterval: 1 * time.Second,
			GossipSigma:    100 * time.Millisecond,
		},
	}

	// Mark as processed since the config wasn't loaded from the environment
	if conf, err = conf.Mark(); err != nil {
		return fmt.Errorf("could not validate test configuration: %s", err)
	}

	// Set the configuration as a pointer so individual tests can modify the config as needed
	s.conf = &conf
	return nil
}

// creates and runs all of the trtl services in preparation for testing
func (s *trtlTestSuite) setupServers() (err error) {
	if s.conf == nil {
		return errors.New("no configuration, run s.setupConfig first")
	}

	// Create the trtl server
	if s.trtl, err = trtl.New(*s.conf); err != nil {
		return fmt.Errorf("could not create trtl service: %s", err)
	}

	// Create a bufconn listener(s) so that there are no actual network requests
	s.grpc = bufconn.New(bufSize)

	// Run the test server without signals, background routines or maintenance mode checks
	// TODO: do we need to check if there was an error when starting run?
	go s.trtl.Run(s.grpc.Listener)
	return nil
}

//===========================================================================
// Test cleanup helpers
//===========================================================================

// cleanup the current temporary directory, configuration, and running services.
func (s *trtlTestSuite) cleanup() (err error) {
	// Shutdown the trtl server if it is running
	// This should shutdown all the running services and close the database
	// Note that Shutdown should be graceful and not shutdown anything not running.
	if s.trtl != nil {
		if err = s.trtl.Shutdown(); err != nil {
			return err
		}
	}

	// Shutdown the gRPC connection if it's running
	if s.grpc != nil {
		s.grpc.Release()
	}

	// Cleanup the tmpdb and delete any stray files
	if s.tmpdb != "" {
		os.RemoveAll(s.tmpdb)
	}

	// Reset all of the test suite variables
	s.tmpdb = ""
	s.conf = nil
	s.grpc = nil
	s.trtl = nil
	return nil
}

// reset the test environment, refreshing the honu database fixture and all of the
// services. This is useful if the test makes changes to the database, though it is
// somewhat heavyweight since it blows away the prior configuration, running servers,
// and open database connection.
func (s *trtlTestSuite) reset() {
	require := s.Require()

	// Cleanup previous configuration and shutdown servers, deleting the tmp database.
	require.NoError(s.cleanup(), "could not cleanup before reset")

	// Setup a new configuration and tmpdb
	require.NoError(s.setupConfig(), "could not reset configuration")

	// Extract the honu db fixture into tmpdb
	require.NoError(s.extractDB(), "could not reset db")

	// Run the trtl server on the new configuration
	require.NoError(s.setupServers(), "could not reset servers")
}

//===========================================================================
// Test fixtures management
//===========================================================================

// loads the test fixtures from a JSON file and stores them in the dbFixtures map -
// these fixtures are used both for generating the test database and for making
// comparative assertions in test code.
func (s *trtlTestSuite) loadFixtures() (err error) {
	var fixtures []byte
	if fixtures, err = ioutil.ReadFile("testdata/db.json"); err != nil {
		return fmt.Errorf("could not read fixtures at testdata/db.json: %s", err)
	}

	if err = json.Unmarshal(fixtures, &dbFixtures); err != nil {
		return fmt.Errorf("could not unmarshal fixtures at testdata/db.json: %s", err)
	}
	return nil
}

// extract honu database from generated fixtures - the honu db fixture is interacted
// with directly by the code and we expect that the data in the honu db fixture mates
// the in-memory fixtures defined by dbFixtures. If there is a mismatch, delete
// testdata/db.tgz or run the tests with the -update flag to regenerate them.
func (s *trtlTestSuite) extractDB() (err error) {
	if s.tmpdb == "" {
		return errors.New("no temporary database, run s.setupConfig first")
	}

	// Always extract the test database to a temporary directory.
	if _, err = utils.ExtractGzip(dbtgz, s.tmpdb, true); err != nil {
		return fmt.Errorf("unable to extract honu db fixture: %s", err)
	}
	return nil
}

// generates an updated database and compresses it to a gzip file.
// NOTE: loadFixtures must have been called before this method.
func (s *trtlTestSuite) generateDB() (err error) {
	// Create a new and temporary db to write the fixtures into
	if err = s.setupConfig(); err != nil {
		return err
	}

	// Ensure we cleanup so that subsequent tests can generate tmpdb directories
	defer s.cleanup()

	// Open a honu database, all fixtures will be written by Honu, which means that Honu
	// will be performing all version management, we expect that everything is at the
	// first version when the fixtures database is created.
	var db *honu.DB
	if db, err = honu.Open(s.conf.Database.URL, s.conf.GetHonuConfig()); err != nil {
		return fmt.Errorf("could not open tmp honu database in %s: %s", s.tmpdb, err)
	}
	defer db.Close()

	// Write all the test fixtures to the database.
	for _, fixture := range dbFixtures {
		var value []byte
		if value, err = json.Marshal(fixture.Value); err != nil {
			return fmt.Errorf("could not marshal %s: %s", fixture.Key, err)
		}

		// Put the data into the proper namespace, Honu takes care of versioning
		if _, err = db.Put([]byte(fixture.Key), value, options.WithNamespace(fixture.Namespace)); err != nil {
			return fmt.Errorf("could not put fixture %s to db: %s", fixture.Key, err)
		}
	}

	if err = utils.WriteGzip(s.tmpdb, dbtgz); err != nil {
		return fmt.Errorf("could not create %s from %s: %s", dbtgz, s.tmpdb, err)
	}
	return nil
}
