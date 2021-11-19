package gds_test

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
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
	"github.com/trisacrypto/directory/pkg/utils"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	vaspPrefix = "vasps::"
	certPrefix = "certreqs::"
)

var (
	update          = flag.Bool("update", false, "update the gzipped test databases")
	smallDBFixtures = []string{
		vaspPrefix + "d9da630e-41aa-11ec-9d29-acde48001122",
		vaspPrefix + "d9efca14-41aa-11ec-9d29-acde48001122",
	}
)

type gdsTestSuite struct {
	suite.Suite
	fakes       string
	dbPath      string
	smallDBPath string
	dbGzip      string
	smallDBGzip string
	fixtures    map[string]interface{}
	conf        *config.Config
	db          store.Store
	tokens      *tokens.TokenManager
}

func getVASPIDFromKey(key string) string {
	return strings.TrimPrefix(key, vaspPrefix)
}

// loadFixtures loads the JSON test fixtures from disk and stores them in the dbFixtures map.
func (s *gdsTestSuite) loadFixtures() {
	require := s.Require()
	// Extract the gzipped archive.
	root, err := utils.ExtractGzip(s.fakes, "testdata", true)
	require.NoError(err)

	s.fixtures = make(map[string]interface{})

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
		switch {
		case strings.HasPrefix(key, vaspPrefix):
			vasp := &pb.VASP{}
			err = protojson.Unmarshal(data, vasp)
			require.NoError(err)
			s.fixtures[key] = vasp
			return nil
		case strings.HasPrefix(key, certPrefix):
			cert := &models.CertificateRequest{}
			err = protojson.Unmarshal(data, cert)
			require.NoError(err)
			s.fixtures[key] = cert
			return nil
		}
		return fmt.Errorf("unrecognized prefix for file: %s", info.Name())
	})
	require.NoError(err)
	require.NoError(os.RemoveAll(root))
}

// writeObj writes a protobuf object to the given database.
func (s *gdsTestSuite) writeObj(db *leveldb.DB, key string, obj interface{}) {
	var (
		data []byte
		err  error
	)
	require := s.Require()
	switch obj.(type) {
	case *pb.VASP:
		vasp := obj.(*pb.VASP)
		data, err = proto.Marshal(vasp)
		require.NoError(err)
	case *models.CertificateRequest:
		cert := obj.(*models.CertificateRequest)
		data, err = proto.Marshal(cert)
		require.NoError(err)
	default:
		require.Fail(fmt.Sprintf("unrecognized object for key: %s", key))
	}
	require.NoError(db.Put([]byte(key), data, nil))
}

// generateDB generates an updated database and compresses it to a gzip file.
// Note: This also generates a temporary directory which the suite teardown
// should clean up.
func (s *gdsTestSuite) generateDB() {
	require := s.Require()
	db, err := leveldb.OpenFile(s.dbPath, nil)
	require.NoError(err)
	defer db.Close()
	small, err := leveldb.OpenFile(s.smallDBPath, nil)
	require.NoError(err)
	defer small.Close()
	s.loadFixtures()

	// Write all the test fixtures to the database.
	for id, obj := range s.fixtures {
		s.writeObj(db, id, obj)
	}
	require.NoError(utils.WriteGzip(s.dbPath, s.dbGzip))
	log.Info().Str("db", s.dbGzip).Msg("successfully regenerated test database")

	// Write test fixtures to the small database.
	for _, id := range smallDBFixtures {
		obj, ok := s.fixtures[id]
		require.True(ok)
		s.writeObj(small, id, obj)
	}
	require.NoError(utils.WriteGzip(s.smallDBPath, s.smallDBGzip))
	log.Info().Str("db", s.smallDBGzip).Msg("successfully regenerated test database")
}

func (s *gdsTestSuite) SetupSuite() {
	var err error
	require := s.Require()
	gin.SetMode(gin.TestMode)

	s.fakes = filepath.Join("testdata", "fakes.tgz")
	s.dbGzip = filepath.Join("testdata", "db.tgz")
	s.smallDBGzip = filepath.Join("testdata", "smalldb.tgz")
	s.dbPath, err = ioutil.TempDir("testdata", "db-*")
	require.NoError(err)
	s.smallDBPath, err = ioutil.TempDir("testdata", "smalldb-*")
	require.NoError(err)

	s.conf = &config.Config{
		Admin: config.AdminConfig{
			CookieDomain: "example.com",
			Oauth: config.OauthConfig{
				GoogleAudience:         "http://localhost",
				AuthorizedEmailDomains: []string{"example.com"},
			},
		},
		Email: config.EmailConfig{
			ServiceEmail: "service@example.com",
			AdminEmail:   "admin@example.com",
			Testing:      true,
		},
	}

	// Regenerate the test database if requested or it doesn't exist.
	// Note: generateDB calls loadFixtures under the hood in order to populate the
	// database. The difference here is whether or not the gzipped file should be
	// regenerated, which we need to do every time the JSON fixtures are updated.
	_, dbErr := os.Stat(s.dbGzip)
	_, smallDBErr := os.Stat(s.smallDBGzip)
	if *update || os.IsNotExist(dbErr) || os.IsNotExist(smallDBErr) {
		s.generateDB()
	} else {
		s.loadFixtures()
	}
}

func (s *gdsTestSuite) BeforeTest(suite, test string) {
	// Extract the test database to a temporary directory.
	if _, err := utils.ExtractGzip(s.dbGzip, s.dbPath, false); err != nil {
		log.Warn().Err(err).Str("db", s.dbGzip).Msg("unable to extract test fixtures")
		s.generateDB()
	}

	// Extract the small test database to a temporary directory.
	if _, err := utils.ExtractGzip(s.smallDBGzip, s.smallDBPath, false); err != nil {
		log.Warn().Err(err).Str("db", s.smallDBGzip).Msg("unable to extract test fixtures")
		s.generateDB()
	}
}

func (s *gdsTestSuite) AfterTest(suite, test string) {
	require := s.Require()
	require.NoError(os.RemoveAll(s.dbPath))
	require.NoError(os.RemoveAll(s.smallDBPath))
}

func TestGds(t *testing.T) {
	suite.Run(t, new(gdsTestSuite))
}

func (s *gdsTestSuite) TestFixtures() {
	require := s.Require()
	db, err := leveldb.OpenFile(s.dbPath, nil)
	require.NoError(err)
	defer db.Close()
	smallDB, err := leveldb.OpenFile(s.smallDBPath, nil)
	require.NoError(err)
	defer smallDB.Close()

	require.NotEmpty(s.fixtures)
	for name, obj := range s.fixtures {
		data, err := db.Get([]byte(name), nil)
		require.NoError(err)
		switch obj.(type) {
		case *pb.VASP:
			expected := obj.(*pb.VASP)
			actual := &pb.VASP{}
			err = proto.Unmarshal(data, actual)
			require.NoError(err)
			require.True(proto.Equal(expected, actual))
		case *models.CertificateRequest:
			expected := obj.(*models.CertificateRequest)
			actual := &models.CertificateRequest{}
			err = proto.Unmarshal(data, actual)
			require.NoError(err)
			require.True(proto.Equal(expected, actual))
		default:
			require.Fail(fmt.Sprintf("unrecognized object for key: %s", name))
		}
	}

	for _, key := range smallDBFixtures {
		data, err := smallDB.Get([]byte(key), nil)
		require.NoError(err)
		switch {
		case strings.HasPrefix(key, vaspPrefix):
			expected := s.fixtures[key].(*pb.VASP)
			actual := &pb.VASP{}
			err = proto.Unmarshal(data, actual)
			require.NoError(err)
			require.True(proto.Equal(expected, actual))
		case strings.HasPrefix(key, certPrefix):
			expected := s.fixtures[key].(*models.CertificateRequest)
			actual := &models.CertificateRequest{}
			err = proto.Unmarshal(data, actual)
			require.NoError(err)
			require.True(proto.Equal(expected, actual))
		default:
			require.Fail(fmt.Sprintf("unrecognized object for key: %s", key))
		}
	}
}
