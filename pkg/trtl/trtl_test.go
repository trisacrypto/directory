package trtl_test

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/rotationalio/honu/object"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/utils"

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

	err = utils.WriteGzip(s.db, s.gzip)
	require.NoError(err)
	log.Info().Msg("successfully regenerated test fixtures")
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
	if _, err = utils.ExtractGzip(s.gzip, s.db); err != nil {
		// Regenerate the test database if the extraction failed.
		log.Warn().Err(err).Msg("unable to extract test fixtures")
		generateDB(s)
	}

	// Load default config and add database path.
	os.Setenv("TRTL_DATABASE_URL", "leveldb:///"+s.db)
	os.Setenv("TRTL_REPLICA_PID", "8")
	os.Setenv("TRTL_REPLICA_REGION", "minneapolis")
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

// Test that we can call the Put RPC and get the correct response.
func (s *trtlTestSuite) TestPut() {
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

	// Put a value from a reserved namespace - should fail.
	_, err = client.Put(ctx, &pb.PutRequest{
		Namespace: "default",
		Key:       []byte(object.Key),
		Value:     []byte("foo"),
	})
	require.Error(err)

	// Put a value without the key - should fail.
	_, err = client.Put(ctx, &pb.PutRequest{
		Namespace: object.Namespace,
		Value:     []byte("foo"),
	})
	require.Error(err)

	// Put without a value - should fail.
	_, err = client.Put(ctx, &pb.PutRequest{
		Namespace: object.Namespace,
		Key:       []byte(object.Key),
	})
	require.Error(err)

	// Put a value to the default namespace.
	reply, err := client.Put(ctx, &pb.PutRequest{
		Key:   []byte("testKey"),
		Value: []byte("testVal"),
	})
	require.NoError(err)
	require.True(reply.Success)
	require.Empty((reply.Meta))

	// Put a value to a valid namespace with return_meta=false.
	reply, err = client.Put(ctx, &pb.PutRequest{
		Namespace: alice.Namespace,
		Key:       []byte(alice.Key),
		Value:     []byte("arlo guthrie"),
		Options: &pb.Options{
			ReturnMeta: false,
		},
	})
	require.NoError(err)
	require.True(reply.Success)
	require.Empty((reply.Meta))

	// Put a value with return_meta=true.
	expectedPID, err := strconv.Atoi((os.Getenv("TRTL_REPLICA_PID")))
	require.NoError(err)
	expectedRegion := os.Getenv("TRTL_REPLICA_REGION")
	expectedMeta := &pb.Meta{
		Key:       []byte(alice.Key),
		Namespace: alice.Namespace,
		Region:    metaRegion,
		Owner:     metaOwner,
		Version: &pb.Version{
			Pid:     uint64(expectedPID),
			Version: 4,
			Region:  expectedRegion,
		},
		Parent: &pb.Version{
			Pid:     uint64(expectedPID),
			Version: 3,
			Region:  expectedRegion,
		},
	}
	reply, err = client.Put(ctx, &pb.PutRequest{
		Namespace: alice.Namespace,
		Key:       []byte(alice.Key),
		Value:     []byte("cheshire cat"),
		Options: &pb.Options{
			ReturnMeta: true,
		},
	})
	require.NoError(err)
	require.True(reply.Success)
	require.NotNil(reply.Meta)
	require.Equal([]byte(alice.Key), reply.Meta.Key)
	require.Equal(alice.Namespace, reply.Meta.Namespace)
	require.True(proto.Equal(expectedMeta, reply.Meta))
}

// Test that we can call the Batch RPC and get the correct response.
func (s *trtlTestSuite) TestBatch() {
	require := s.Require()

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

	requests := map[int64]*pb.BatchRequest{
		1: {
			Id: 1,
			Request: &pb.BatchRequest_Put{
				Put: &pb.PutRequest{
					Key:       []byte("foo"),
					Namespace: "default",
					Value:     []byte("bar"),
				},
			},
		},
		2: {
			Id: 2,
			Request: &pb.BatchRequest_Delete{
				Delete: &pb.DeleteRequest{
					Key:       []byte("foo"),
					Namespace: "default",
				},
			},
		},
	}
	stream, err := client.Batch(ctx)
	require.NoError(err)
	for _, r := range requests {
		err = stream.Send(r)
		require.NoError(err)
	}
	reply, err := stream.CloseAndRecv()
	require.NoError(err)
	require.Equal(int64(len(requests)), reply.Operations)
	require.Equal(int64(len(requests)), reply.Failed)
	require.Equal(int64(0), reply.Successful)
	require.Len(reply.Errors, len(requests))
	require.Contains(requests, reply.Errors[1].Id)
	require.Equal(requests[reply.Errors[1].Id].Id, reply.Errors[1].Id)
}
