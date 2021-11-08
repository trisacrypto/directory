package trtl_test

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/rotationalio/honu/object"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	update = flag.Bool("update", false, "update the static test fixtures")
)

// dbFixtures is a test fixtures map used to both generate the static test database
// and verify the results of the trtl DB calls. It maps a test fixture name to an entry
// containing a namespace, key, and value stored in the database.
var dbFixtures = map[string]*dbEntry{}

const metaRegion = "us-east-1"
const metaOwner = "foo"

var metaVersion = &object.Version{
	Pid:     1,
	Version: 2,
	Region:  metaRegion,
	Parent: &object.Version{
		Pid:     1,
		Version: 1,
		Region:  metaRegion,
	},
}

type dbEntry struct {
	Namespace string                 `json:"namespace"`
	Key       string                 `json:"key"`
	Value     map[string]interface{} `json:"value"`
}

type trtlTestSuite struct {
	suite.Suite
	gzip string
	db   string
	conf config.Config
}

// loadFixtures loads the test fixtures from a JSON file and stores them in the
// dbFixtures map.
func loadFixtures(s *trtlTestSuite) {
	require := s.Require()
	fixtures, err := ioutil.ReadFile("testdata/db.json")
	require.NoError(err)
	err = json.Unmarshal(fixtures, &dbFixtures)
	require.NoError(err)
}

// generateDB generates an updated database and compresses it to a gzip file.
// Note: This also generates a temporary directory which the suite teardown
// should clean up.
func generateDB(s *trtlTestSuite) {
	require := s.Require()
	db, err := leveldb.OpenFile(s.db, nil)
	require.NoError(err)
	defer db.Close()

	loadFixtures(s)

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

	writeGzip(s)
	log.Info().Msg("successfully regenerated test fixtures")
}

// extractGzip extracts the gzipped database to the temporary directory.
func extractGzip(s *trtlTestSuite) (err error) {
	var (
		f  *os.File
		gr *gzip.Reader
	)

	// Read the gzip file.
	if f, err = os.Open(s.gzip); err != nil {
		return err
	}
	defer f.Close()
	if gr, err = gzip.NewReader(f); err != nil {
		return err
	}
	defer gr.Close()

	// Write the contents to the temporary directory.
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(filepath.Join(s.db, hdr.Name), os.FileMode(hdr.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			var reg *os.File
			if reg, err = os.Create(filepath.Join(s.db, hdr.Name)); err != nil {
				return err
			}
			if _, err = io.Copy(reg, tr); err != nil {
				reg.Close()
				return err
			}
			reg.Close()
		default:
			return fmt.Errorf("extracting %s: unknown type flag: %c", hdr.Name, hdr.Typeflag)
		}
	}
	return nil
}

// writeGzip writes the database in the temporary directory to a gzipped file.
func writeGzip(s *trtlTestSuite) {
	require := s.Require()
	// Create a gzip file.
	f, err := os.Create(s.gzip)
	require.NoError(err)
	defer f.Close()
	w := gzip.NewWriter(f)
	defer w.Close()

	// Create a tar file.
	tw := tar.NewWriter(w)
	defer tw.Close()

	// Write the DB to the tar file.
	err = filepath.Walk(s.db, func(path string, info os.FileInfo, err error) error {
		require.NoError(err)
		hdr, err := tar.FileInfoHeader(info, "")
		require.NoError(err)
		hdr.Name = path[len(s.db):]
		err = tw.WriteHeader(hdr)
		require.NoError(err)
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		require.NoError(err)
		defer f.Close()
		_, err = io.Copy(tw, f)
		require.NoError(err)
		return nil
	})
	require.NoError(err)
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
		generateDB(s)
	} else {
		loadFixtures(s)
	}

	// Always extract the test database to a temporary directory.
	if err = extractGzip(s); err != nil {
		// Regenerate the test database if the extraction failed.
		log.Warn().Err(err).Msg("unable to extract test fixtures")
		generateDB(s)
	}

	// Load default config and add database path.
	os.Setenv("TRTL_DATABASE_URL", "leveldb:///"+s.db)
	os.Setenv("TRTL_PID", "1")
	os.Setenv("TRTL_REGION", "foo")
	s.conf, err = config.New()
	require.NoError(err)
}

func (s *trtlTestSuite) TearDownSuite() {
	require := s.Require()
	err := os.RemoveAll(s.db)
	require.NoError(err)
}

func TestTrtl(t *testing.T) {
	suite.Run(t, new(trtlTestSuite))
}

// Test that we can call the Get RPC and get the correct response.
func (s *trtlTestSuite) TestGet() {
	var actual interface{}

	require := s.Require()
	alice := dbFixtures["alice"]
	object := dbFixtures["object"]

	// Start the server.
	server, err := trtl.New(s.conf)
	require.NoError(err)
	go server.Serve()
	defer server.Shutdown()

	// Start the gRPC client.
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "localhost"+s.conf.BindAddr, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()
	client := pb.NewTrtlClient(conn)

	// Retrieve a value from a reserved namespace - should fail.
	_, err = client.Get(ctx, &pb.GetRequest{
		Namespace: "default",
		Key:       []byte(object.Key),
	})
	require.Error(err)

	// Retrieve a value without the key - should fail.
	_, err = client.Get(ctx, &pb.GetRequest{
		Namespace: object.Namespace,
	})
	require.Error(err)

	// Retrieve a value from the default namespace.
	reply, err := client.Get(ctx, &pb.GetRequest{
		Key: []byte(object.Key),
	})
	require.NoError(err)
	err = json.Unmarshal(reply.Value, &actual)
	require.NoError(err)
	require.Equal(object.Value, actual)

	// Retrieve a value from a valid namespace.
	reply, err = client.Get(ctx, &pb.GetRequest{
		Namespace: alice.Namespace,
		Key:       []byte(alice.Key),
	})
	require.NoError(err)
	err = json.Unmarshal(reply.Value, &actual)
	require.NoError(err)
	require.Equal(alice.Value, actual)

	// Retrieve a value from a non-existent namespace - should fail.
	_, err = client.Get(ctx, &pb.GetRequest{
		Namespace: "invalid",
		Key:       []byte(alice.Key),
	})
	require.Error(err)

	// Retrieve a value from a non-existent key - should fail.
	_, err = client.Get(ctx, &pb.GetRequest{
		Namespace: alice.Namespace,
		Key:       []byte("invalid"),
	})
	require.Error(err)

	// Retrieve a value with return_meta=false.
	reply, err = client.Get(ctx, &pb.GetRequest{
		Namespace: alice.Namespace,
		Key:       []byte(alice.Key),
		Options: &pb.Options{
			ReturnMeta: false,
		},
	})
	require.NoError(err)
	require.Nil(reply.Meta)
	err = json.Unmarshal(reply.Value, &actual)
	require.NoError(err)
	require.Equal(alice.Value, actual)

	// Retrieve a value with return_meta=true.
	expectedMeta := &pb.Meta{
		Key:       []byte(alice.Key),
		Namespace: alice.Namespace,
		Region:    metaRegion,
		Owner:     metaOwner,
		Version: &pb.Version{
			Pid:     metaVersion.Pid,
			Version: metaVersion.Version,
			Region:  metaVersion.Region,
		},
		Parent: &pb.Version{
			Pid:     metaVersion.Parent.Pid,
			Version: metaVersion.Parent.Version,
			Region:  metaVersion.Parent.Region,
		},
	}
	reply, err = client.Get(ctx, &pb.GetRequest{
		Namespace: alice.Namespace,
		Key:       []byte(alice.Key),
		Options: &pb.Options{
			ReturnMeta: true,
		},
	})
	require.NoError(err)
	require.NotNil(reply.Meta)
	require.Equal([]byte(alice.Key), reply.Meta.Key)
	require.Equal(alice.Namespace, reply.Meta.Namespace)
	require.True(proto.Equal(expectedMeta, reply.Meta))
}
