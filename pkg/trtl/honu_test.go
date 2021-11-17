package trtl_test

import (
	"context"
	"encoding/json"

	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/protobuf/proto"
)

// Test that we can call the Get RPC and get the correct response.
func (s *trtlTestSuite) TestGet() {
	var actual interface{}

	require := s.Require()
	alice := dbFixtures["alice"]
	object := dbFixtures["object"]
	ctx := context.Background()

	// Start the gRPC client.
	conn, err := s.connect()
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

// Test that we can call the Batch RPC and get the correct response.
func (s *trtlTestSuite) TestBatch() {
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	conn, err := s.connect()
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
	for _, e := range reply.Errors {
		require.Contains(requests, e.Id)
		require.Equal(requests[e.Id].Id, e.Id)
		require.Contains(e.Error, "not implemented")
	}
}
