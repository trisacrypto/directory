package gds

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/utils"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	update = flag.Bool("update", false, "update the gzipped test database")
)

type gdsTestSuite struct {
	suite.Suite
	fakes   string
	db      string
	golden  string
	dbVASPs map[string]*pb.VASP
}

// Overwrite the test database with a subset of VASPs (for small-scale unit testing).
func WriteVASPs(s *gdsTestSuite, ids []string) {
	require := s.Require()
	err := os.RemoveAll(s.db)
	require.NoError(err)
	db, err := leveldb.OpenFile(s.db, nil)
	require.NoError(err)
	defer db.Close()
	for _, id := range ids {
		vasp, ok := s.dbVASPs[id]
		require.True(ok)
		data, err := proto.Marshal(vasp)
		require.NoError(err)
		err = db.Put([]byte("vasps::"+id), data, nil)
		require.NoError(err)
	}
}

// loadFixtures loads the JSON test fixtures from disk and stores them in the dbFixtures map.
func loadFixtures(s *gdsTestSuite) {
	require := s.Require()
	// Extract the gzipped archive.
	root, err := utils.ExtractGzip(s.fakes, "testdata")
	require.NoError(err)

	s.dbVASPs = make(map[string]*pb.VASP)

	// Load the JSON fixtures from disk.
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		require.NoError(err)
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		// Unmarshal the JSON into the global fixtures map.
		data, err := os.ReadFile(path)
		require.NoError(err)
		key := strings.TrimSuffix(info.Name(), ".json")
		if strings.HasPrefix(key, "vasps::") {
			vasp := &pb.VASP{}
			err = protojson.Unmarshal(data, vasp)
			require.NoError(err)
			s.dbVASPs[strings.TrimPrefix(key, "vasps::")] = vasp
			return nil
		}

		return fmt.Errorf("unrecognized prefix for file: %s", info.Name())
	})
	require.NoError(err)
	os.RemoveAll(root)
	require.NoError(err)
}

// generateDB generates an updated database and compresses it to a gzip file.
// Note: This also generates a temporary directory which the suite teardown
// should clean up.
func generateDB(s *gdsTestSuite) {
	require := s.Require()
	db, err := leveldb.OpenFile(s.db, nil)
	require.NoError(err)
	defer db.Close()

	loadFixtures(s)

	// Write all the test fixtures to the database.
	for id, vasp := range s.dbVASPs {
		data, err := proto.Marshal(vasp)
		require.NoError(err)
		err = db.Put([]byte("vasps::"+id), data, nil)
		require.NoError(err)
	}

	err = utils.WriteGzip(s.db, s.golden)
	require.NoError(err)
	log.Info().Str("db", s.golden).Msg("successfully regenerated test database")
}

func (s *gdsTestSuite) SetupSuite() {
	var err error
	require := s.Require()
	gin.SetMode(gin.TestMode)

	s.fakes = filepath.Join("testdata", "fakes.tgz")
	s.golden = filepath.Join("testdata", "db.tgz")
	s.db, err = ioutil.TempDir("testdata", "db-*")
	require.NoError(err)

	// Regenerate the test database if requested or it doesn't exist.
	// Note: generateDB calls loadFixtures under the hood in order to populate the
	// database. The difference here is whether or not the gzipped file should be
	// regenerated, which we need to do every time the JSON fixtures are updated.
	if _, err = os.Stat(s.golden); *update || os.IsNotExist(err) {
		generateDB(s)
	} else {
		loadFixtures(s)
	}
}

func (s *gdsTestSuite) BeforeTest(suite, test string) {
	// Extract the test database to a temporary directory.
	if _, err := utils.ExtractGzip(s.golden, s.db); err != nil {
		log.Warn().Err(err).Msg("unable to extract test fixtures")
		generateDB(s)
	}
}

func (s *gdsTestSuite) AfterTest(suite, test string) {
	require := s.Require()
	err := os.RemoveAll(s.db)
	require.NoError(err)
}

func TestGds(t *testing.T) {
	suite.Run(t, new(gdsTestSuite))
}

func (s *gdsTestSuite) TestFixtures() {
	require := s.Require()
	db, err := leveldb.OpenFile(s.db, nil)
	require.NoError(err)
	defer db.Close()

	require.NotEmpty(s.dbVASPs)
	for name, vasp := range s.dbVASPs {
		data, err := db.Get([]byte("vasps::"+name), nil)
		require.NoError(err)
		actual := &pb.VASP{}
		err = proto.Unmarshal(data, actual)
		require.NoError(err)
		require.True(proto.Equal(vasp, actual))
	}
}
