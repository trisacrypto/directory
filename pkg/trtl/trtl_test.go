package trtl_test

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/rotationalio/honu/object"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
)

const (
	metaRegion = "us-east-1"
	metaOwner  = "foo"
	bufSize    = 1024 * 1024
)

var (
	// CLI flag, specify go test -update to regenerate static test fixtures
	update = flag.Bool("update", false, "update the static test fixtures")

	// dbFixtures is a test fixtures map used to both generate the static test database
	// and verify the results of the trtl DB calls. It maps a test fixture name to an entry
	// containing a namespace, key, and value stored in the database.
	dbFixtures = map[string]*dbEntry{}

	// the version of all objects in the fixtures
	metaVersion = &object.Version{
		Pid:     1,
		Version: 2,
		Region:  metaRegion,
		Parent: &object.Version{
			Pid:     1,
			Version: 1,
			Region:  metaRegion,
		},
	}
)

type dbEntry struct {
	Namespace string                 `json:"namespace"`
	Key       string                 `json:"key"`
	Value     map[string]interface{} `json:"value"`
}

type trtlTestSuite struct {
	suite.Suite
	gzip string
	db   string
	trtl *trtl.Server
	conf config.Config
	conn *bufconn.Listener
}

// loads the test fixtures from a JSON file and stores them in the dbFixtures map
func (s *trtlTestSuite) loadFixtures() {
	require := s.Require()
	fixtures, err := ioutil.ReadFile("testdata/db.json")
	require.NoError(err)
	err = json.Unmarshal(fixtures, &dbFixtures)
	require.NoError(err)
}

// generates an updated database and compresses it to a gzip file.
// Note: This also generates a temporary directory which the suite teardown
// should clean up.
func (s *trtlTestSuite) generateDB() {
	require := s.Require()
	db, err := leveldb.OpenFile(s.db, nil)
	require.NoError(err)
	defer db.Close()

	s.loadFixtures()

	// Write all the test fixtures to the database.
	for _, fixture := range dbFixtures {
		// Must be wrapped in a honu object to be retrievable with honu.
		meta := &object.Object{
			Key:       []byte(fixture.Key),
			Namespace: fixture.Namespace,
			Region:    metaRegion,
			Owner:     metaOwner,
			Version:   metaVersion,
		}
		meta.Data, err = json.Marshal(fixture.Value)
		require.NoError(err)
		data, err := proto.Marshal(meta)
		require.NoError(err)
		err = db.Put([]byte(fixture.Namespace+"::"+fixture.Key), data, nil)
		require.NoError(err)
	}

	err = utils.WriteGzip(s.db, s.gzip)
	require.NoError(err)
	log.Info().Msg("successfully regenerated test fixtures")
}

func (s *trtlTestSuite) connect() (*grpc.ClientConn, error) {
	ctx := context.Background()
	return grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(s.dialer), grpc.WithInsecure())
}

func (s *trtlTestSuite) dialer(context.Context, string) (net.Conn, error) {
	return s.conn.Dial()
}

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

func (s *trtlTestSuite) SetupSuite() {
	var err error
	require := s.Require()
	s.gzip = filepath.Join("testdata", "db.tar.gz")
	s.db, err = ioutil.TempDir("testdata", "db*")
	require.NoError(err)

	// Regenerate the test database if requested or it doesn't exist.
	// Note: generateDB calls loadFixtures under the hood in order to populate the
	// database. The difference here is whether or not the gzipped file should be
	// regenerated, which we need to do every time db.json is updated.
	if _, err = os.Stat(s.gzip); *update || os.IsNotExist(err) {
		s.generateDB()
	} else {
		s.loadFixtures()
	}

	// Always extract the test database to a temporary directory.
	if _, err = utils.ExtractGzip(s.gzip, s.db); err != nil {
		// Regenerate the test database if the extraction failed.
		log.Warn().Err(err).Msg("unable to extract test fixtures")
		s.generateDB()
	}

	// Load default config and add database path.
	os.Setenv("TRTL_DATABASE_URL", "leveldb:///"+s.db)
	os.Setenv("TRTL_REPLICA_PID", "8")
	os.Setenv("TRTL_REPLICA_REGION", "minneapolis")
	s.conf, err = config.New()
	require.NoError(err)

	// Create the trtl server
	s.trtl, err = trtl.New(s.conf)
	require.NoError(err)

	// Create a bufcon listener so that there are no actual network requests
	s.conn = bufconn.Listen(bufSize)

	// Run the test server without signals, background routines or maintenance mode checks
	// TODO: do we need to check if there was an error when starting run?
	go s.trtl.Run(s.conn)
}

func (s *trtlTestSuite) TearDownSuite() {
	require := s.Require()
	err := os.RemoveAll(s.db)
	require.NoError(err)
}

func TestTrtl(t *testing.T) {
	suite.Run(t, new(trtlTestSuite))
}
