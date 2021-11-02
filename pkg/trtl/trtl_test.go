package trtl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rotationalio/honu/object"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"

	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func generateDB(s *suite.Suite) (db *leveldb.DB) {
	// TODO: Load from a gzipped file.
	db, err := leveldb.OpenFile("testdata/tmp", nil)
	defer db.Close()
	s.NoError(err)

	peer := &peers.Peer{
		Id:       1,
		Addr:     "localhost:1313",
		Name:     "foo",
		Region:   "foo",
		Created:  "foo",
		Modified: "foo",
		Extra: map[string]string{
			"foo": "bar",
		},
	}
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
	meta.Data, err = proto.Marshal(peer)
	s.NoError(err)
	data, err := proto.Marshal(meta)
	s.NoError(err)
	err = db.Put(meta.Key, data, nil)
	s.NoError(err)

	val, err := db.Get(meta.Key, nil)
	s.NoError(err)
	s.Equal(data, val)

	return db
}

type trtlTestSuite struct {
	suite.Suite
	path string
	db   *leveldb.DB
}

func (s *trtlTestSuite) SetupSuite() {
	// TODO: Swap the path for a gzipped database.
	// TODO: Implement --update flag for generating a new gzipped database?
	s.db = generateDB(&s.Suite)
}

func (s *trtlTestSuite) TearDownSuite() {
	// TODO: Delete the extracted version of the database.
}

func TestTrtl(t *testing.T) {
	suite.Run(t, new(trtlTestSuite))
}

// Test that we can call the Get RPC and get the correct response.
func (s *trtlTestSuite) TestGet() {
	// Should the --update flag have specific handling here?

	config := config.Config{
		Enabled:  true,
		BindAddr: "localhost:1313",
	}

	server, err := trtl.New(config)
	s.NoError(err)

	go server.Serve()
	defer server.Shutdown()

	fmt.Println("started server")

	// Test that we can get a response from a gRPC request.
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "localhost:4435", grpc.WithInsecure())
	s.NoError(err)
	client := pb.NewTrtlClient(conn)
	reply, err := client.Get(ctx, &pb.GetRequest{
		Key: []byte("foo"),
		Options: &pb.Options{
			ReturnMeta: false,
		},
	})
	s.NoError(err)

	// Compare the reply to the expected value.
	peer := &peers.Peer{
		Id:       1,
		Addr:     "localhost:1313",
		Name:     "foo",
		Region:   "foo",
		Created:  "foo",
		Modified: "foo",
		Extra: map[string]string{
			"foo": "bar",
		},
	}

	// unmrshal the reply into a Peer
	var actual peers.Peer
	err = proto.Unmarshal(reply.Value, &actual)
	s.NoError(err)
	s.True(proto.Equal(peer, &actual))
}

// Test that we can start and stop a trtl server.
func TestServer(t *testing.T) {
	t.Skip()
	// TODO: For the real tests we probably want to avoid binding a real address.
	config := config.Config{
		Enabled:  true,
		BindAddr: "localhost:1313",
	}

	server, err := trtl.New(config)
	require.NoError(t, err)

	err = server.Serve()
	require.NoError(t, err)

	// Test that we can get a response from a gRPC request.
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "localhost:1313", grpc.WithInsecure())
	require.NoError(t, err)
	client := pb.NewTrtlClient(conn)
	_, err = client.Get(ctx, &pb.GetRequest{Key: []byte("foo")})
	require.EqualError(t, err, "rpc error: code = Unimplemented desc = not implemented")

	err = server.Shutdown()
	require.NoError(t, err)
}
