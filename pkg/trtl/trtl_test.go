package trtl_test

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/rotationalio/honu/object"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"

	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func extractGzip(s *trtlTestSuite) {
	// Read the gzip file.
	require := s.Require()
	f, err := os.Open(s.gzip)
	require.NoError(err)
	defer f.Close()
	gr, err := gzip.NewReader(f)
	require.NoError(err)
	defer gr.Close()

	// Write the contents to a directory.
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		require.NoError(err)
		switch hdr.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filepath.Join(s.db, hdr.Name), os.FileMode(hdr.Mode))
			require.NoError(err)
		case tar.TypeReg:
			f, err := os.Create(filepath.Join(s.db, hdr.Name))
			require.NoError(err)
			_, err = io.Copy(f, tr)
			require.NoError(err)
			f.Close()
		default:
			require.Fail(fmt.Sprintf("extracting %s: unexpected type: %v", hdr.Name, hdr.Typeflag))
		}
	}
}

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

var peerFoo = peers.Peer{
	Id:       1,
	Addr:     "localhost:4435",
	Name:     "foo",
	Region:   "foo",
	Created:  "foo",
	Modified: "foo",
	Extra: map[string]string{
		"foo": "bar",
	},
}

// generateDB generates an updated database and compresses it to a gzip file.
// Note: This also generates a temporary directory which the suite teardown
// should clean up.
func generateDB(s *trtlTestSuite) {
	require := s.Require()
	db, err := leveldb.OpenFile(s.db, nil)
	require.NoError(err)
	defer db.Close()

	meta := &object.Object{
		Key:       []byte("foo"),
		Namespace: "",
		Version: &object.Version{
			Pid:     1,
			Version: 2,
			Region:  "foo",
			Parent: &object.Version{
				Pid:     1,
				Version: 1,
			},
		},
	}
	meta.Data, err = proto.Marshal(&peerFoo)
	require.NoError(err)
	data, err := proto.Marshal(meta)
	require.NoError(err)
	err = db.Put(meta.Key, data, nil)
	require.NoError(err)

	val, err := db.Get(meta.Key, nil)
	require.NoError(err)
	require.Equal(data, val)

	writeGzip(s)
}

type trtlTestSuite struct {
	suite.Suite
	gzip string
	db   string
	conf config.Config
}

func (s *trtlTestSuite) SetupSuite() {
	var err error
	require := s.Require()
	s.gzip = filepath.Join("testdata", "db.tar.gz")
	s.db, err = ioutil.TempDir("testdata", "db*")
	require.NoError(err)
	// TODO: Implement --update flag for generating a new gzipped database?
	if _, err := os.Stat(s.gzip); os.IsNotExist(err) {
		generateDB(s)
	} else {
		extractGzip(s)
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
	require := s.Require()
	// Should the --update flag have specific handling here?

	server, err := trtl.New(s.conf)
	require.NoError(err)

	go server.Serve()
	defer server.Shutdown()

	// Test that we can get a response from a gRPC request.
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "localhost"+s.conf.BindAddr, grpc.WithInsecure())
	require.NoError(err)
	client := pb.NewTrtlClient(conn)
	reply, err := client.Get(ctx, &pb.GetRequest{
		Key: []byte("foo"),
		Options: &pb.Options{
			ReturnMeta: false,
		},
	})
	require.NoError(err)

	// unmarshal the reply into a Peer
	var actual peers.Peer
	err = proto.Unmarshal(reply.Value, &actual)
	require.NoError(err)
	require.True(proto.Equal(&peerFoo, &actual))
}
