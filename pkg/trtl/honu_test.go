package trtl_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	codes "google.golang.org/grpc/codes"
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
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := pb.NewTrtlClient(s.grpc.Conn)

	// Retrieve a value from a reserved namespace - should fail.
	_, err := client.Get(ctx, &pb.GetRequest{
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
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := pb.NewTrtlClient(s.grpc.Conn)

	// Put a value to a reserved namespace - should fail.
	_, err := client.Put(ctx, &pb.PutRequest{
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
	// TODO update os.Getenv to use test fixtures in sc-2098
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
	// TODO this test modifies the DB state so could cause subsequent tests to have unexpected results
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

// Test that we can call the Delete RPC and get the correct response.
func (s *trtlTestSuite) TestDelete() {
	require := s.Require()
	ctx := context.Background()
	tempNS := "temp"

	// Start the gRPC client.
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := pb.NewTrtlClient(s.grpc.Conn)

	// Delete from a reserved namespace - should fail.
	_, err := client.Delete(ctx, &pb.DeleteRequest{
		Namespace: "default",
		Key:       []byte("jerky"),
	})
	require.Error(err)

	// Delete a value without the key - should fail.
	_, err = client.Delete(ctx, &pb.DeleteRequest{
		Namespace: tempNS,
	})
	require.Error(err)

	// Doing a valid Put, then valid Delete should not error & should return success
	tempKey := []byte("tempkey")
	tempVal := []byte("deleteme")

	_, err = client.Put(ctx, &pb.PutRequest{
		Key:       tempKey,
		Value:     tempVal,
		Namespace: tempNS,
	})
	require.NoError(err)
	noMeta, err := client.Delete(ctx, &pb.DeleteRequest{
		Key:       tempKey,
		Namespace: tempNS,
	})
	require.NoError(err)
	require.True(noMeta.Success)
	require.Nil(noMeta.Meta)

	// Doing a valid Put, then Delete with meta gives back success and metadata
	_, err = client.Put(ctx, &pb.PutRequest{
		Key:       tempKey,
		Value:     tempVal,
		Namespace: tempNS,
	})
	require.NoError(err)
	withMeta, err := client.Delete(ctx, &pb.DeleteRequest{
		Key:       tempKey,
		Namespace: tempNS,
		Options: &pb.Options{
			ReturnMeta: true,
		},
	})
	require.NoError(err)
	// TODO update os.Getenv to use test fixtures in sc-2098
	pid := os.Getenv("TRTL_REPLICA_PID")
	expectedPID, err := strconv.Atoi(pid)
	require.NoError(err)
	expectedRegion := os.Getenv("TRTL_REPLICA_REGION")
	owner := bytes.Join([][]byte{[]byte(pid), []byte(expectedRegion)}, []byte(":"))
	expectedMeta := &pb.Meta{
		Key:       tempKey,
		Namespace: tempNS,
		Region:    expectedRegion,
		Owner:     string(owner),
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
	require.NoError(err)
	require.True(withMeta.Success)
	require.NotNil(withMeta.Meta)
	require.Equal(tempKey, withMeta.Meta.Key)
	require.Equal(tempNS, withMeta.Meta.Namespace)
	require.True(proto.Equal(expectedMeta, withMeta.Meta))
}

// Test that we can call the Batch RPC and get the correct response.
func (s *trtlTestSuite) TestBatch() {
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := pb.NewTrtlClient(s.grpc.Conn)

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

func (s *trtlTestSuite) TestIter() {
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := pb.NewTrtlClient(s.grpc.Conn)

	// Test cannot use reserved namespace
	_, err := client.Iter(ctx, &pb.IterRequest{Namespace: "index"})
	s.StatusError(err, codes.PermissionDenied, "cannot used reserved namespace")

	// Test Invalid Options
	_, err = client.Iter(ctx, &pb.IterRequest{Namespace: "people", Options: &pb.Options{IterNoKeys: true, IterNoValues: true}})
	s.StatusError(err, codes.InvalidArgument, "cannot specify no keys, values, and no return meta: no data would be returned")

	// Test iter default prefix, default options, expecting 1 response from default namespace
	rep, err := client.Iter(ctx, &pb.IterRequest{})
	require.NoError(err, "could not iterate default prefix with default options")
	require.Len(rep.Values, 1, "too many responses returned, did the fixtures change?")
	require.Empty(rep.NextPageToken, "a next page token was returned for a one page response")
	require.NotEmpty(rep.Values[0].Key, "key not supplied in iter by default")
	require.NotEmpty(rep.Values[0].Value, "value not supplied in iter by default")
	require.Equal("default", rep.Values[0].Namespace, "non-default namespace fetched")
	require.Empty(rep.Values[0].Meta, "meta returned by default")

	// Test invalid page token
	token := "CAISLHBlb3BsZTo6NDZlNzg5MTctOGQyMC00N2MwLWIwZDEtZTUyMDQxNDlhOTM2"
	_, err = client.Iter(ctx, &pb.IterRequest{Options: &pb.Options{PageToken: "foo"}})
	s.StatusError(err, codes.InvalidArgument, "invalid page token")
	_, err = client.Iter(ctx, &pb.IterRequest{Namespace: "people", Options: &pb.Options{PageToken: token, PageSize: 27}})
	s.StatusError(err, codes.InvalidArgument, "page size cannot change between requests")

	// Test ordered non-paginated request with prefix
	rep, err = client.Iter(ctx, &pb.IterRequest{Namespace: "people", Prefix: []byte("215")})
	require.NoError(err, "could not fetch complete iteration")
	require.Empty(rep.NextPageToken, "a next page token was returned for a one page response")
	require.Len(rep.Values, 5, "incorrect responses returned did the fixtures change?")

	expectedOrder := []string{
		"alice", "bob", "charlie", "darlene", "erica",
		"franklin", "gregor", "helen", "ivan", "juliet",
	}

	for i := 0; i < 5; i++ {
		expected := dbFixtures[expectedOrder[i]]
		require.NotNil(expected)
		pair := rep.Values[i]
		require.Equal("people", pair.Namespace)
		require.Empty(pair.Meta)
		require.Equal(expected.Key, string(pair.Key))

		var value map[string]interface{}
		require.NoError(json.Unmarshal(pair.Value, &value), "could not unmarshal value from db")
		require.Equal(expected.Value, value)
	}

	// Test No Keys
	rep, err = client.Iter(ctx, &pb.IterRequest{Namespace: "people", Prefix: []byte("215"), Options: &pb.Options{IterNoKeys: true}})
	require.NoError(err, "could not fetch complete iteration")
	require.NotEmpty(rep.Values, "no values returned, expected more than 1")

	for _, pair := range rep.Values {
		require.Empty(pair.Key, "key returned on no keys")
		require.NotEmpty(pair.Value, "value not returned")
		require.Empty(pair.Meta, "meta returned without request")
	}

	// Test No Values
	rep, err = client.Iter(ctx, &pb.IterRequest{Namespace: "people", Prefix: []byte("215"), Options: &pb.Options{IterNoValues: true}})
	require.NoError(err, "could not fetch complete iteration")
	require.NotEmpty(rep.Values, "no values returned, expected more than 1")

	for _, pair := range rep.Values {
		require.NotEmpty(pair.Key, "key not returned")
		require.Empty(pair.Value, "value returned on no values")
		require.Empty(pair.Meta, "meta returned without request")
	}

	// Test Return Meta
	rep, err = client.Iter(ctx, &pb.IterRequest{Namespace: "people", Prefix: []byte("215"), Options: &pb.Options{ReturnMeta: true}})
	require.NoError(err, "could not fetch complete iteration")
	require.NotEmpty(rep.Values, "no values returned, expected more than 1")
	for _, pair := range rep.Values {
		require.NotEmpty(pair.Key, "key not returned")
		require.NotEmpty(pair.Value, "value not returned")
		require.NotEmpty(pair.Meta, "meta not returned on request")
	}

	// Test paginated request with odd numbers of pages
	var (
		pages, people int
		pageToken     string
	)

	for queries := 0; queries < 6; queries++ {
		req := &pb.IterRequest{
			Namespace: "people",
			Options: &pb.Options{
				PageSize:  3,
				PageToken: pageToken,
			},
		}

		rep, err = client.Iter(ctx, req)
		require.NoError(err, "could not make paginated request")
		require.LessOrEqual(len(rep.Values), 3, "invalid page size returned")

		pages++
		people += len(rep.Values)

		pageToken = rep.NextPageToken
		if rep.NextPageToken == "" {
			break
		}
	}

	require.Equal(4, pages, "number of people pages changed, have fixtures been modified?")
	require.Equal(10, people, "number of people has changed, have fixtures been modified?")

	// Test paginated request with even numbers of pages
	pageToken = ""
	pages = 0
	people = 0

	for queries := 0; queries < 4; queries++ {
		req := &pb.IterRequest{
			Namespace: "people",
			Options: &pb.Options{
				PageSize:  5,
				PageToken: pageToken,
			},
		}

		rep, err = client.Iter(ctx, req)
		require.NoError(err, "could not make paginated request")
		require.Equal(len(rep.Values), 5, "invalid page size returned")

		pages++
		people += len(rep.Values)

		pageToken = rep.NextPageToken
		if rep.NextPageToken == "" {
			break
		}
	}

	require.Equal(2, pages, "number of people pages changed, have fixtures been modified?")
	require.Equal(10, people, "number of people has changed, have fixtures been modified?")
}

func (s *trtlTestSuite) TestCursor() {
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := pb.NewTrtlClient(s.grpc.Conn)

	// Test cannot use reserved namespace
	stream, err := client.Cursor(ctx, &pb.CursorRequest{Namespace: "index"})
	require.NoError(err, "could not create cursor stream")
	_, err = stream.Recv()
	s.StatusError(err, codes.PermissionDenied, "cannot used reserved namespace")

	// Test Invalid Options
	stream, err = client.Cursor(ctx, &pb.CursorRequest{Namespace: "people", Options: &pb.Options{IterNoKeys: true, IterNoValues: true}})
	require.NoError(err, "could not create cursor stream")
	_, err = stream.Recv()
	s.StatusError(err, codes.InvalidArgument, "cannot specify no keys, values, and no return meta: no data would be returned")

	// Test iter default prefix, default options, expecting 1 response from default namespace
	results := make([]*pb.KVPair, 0, 1)
	stream, err = client.Cursor(ctx, &pb.CursorRequest{})
	require.NoError(err, "could not create cursor stream")
	for {
		rep, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(err, "received non-EOF error from recv")
		require.NotEmpty(rep.Key, "key not supplied in iter by default")
		require.NotEmpty(rep.Value, "value not supplied in iter by default")
		require.Equal("default", rep.Namespace, "non-default namespace fetched")
		require.Empty(rep.Meta, "meta returned by default")

		results = append(results, rep)
	}
	require.Len(results, 1, "too many responses returned, did the fixtures change?")

	// Test ordered request with prefix
	expectedOrder := []string{
		"alice", "bob", "charlie", "darlene", "erica",
		"franklin", "gregor", "helen", "ivan", "juliet",
	}

	stream, err = client.Cursor(ctx, &pb.CursorRequest{Namespace: "people", Prefix: []byte("215")})
	require.NoError(err, "could not create cursor stream")

	i := 0
	for {
		pair, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(err, "received non-EOF error from recv")

		expected := dbFixtures[expectedOrder[i]]
		require.NotNil(expected)

		require.Equal("people", pair.Namespace)
		require.Empty(pair.Meta)
		require.Equal(expected.Key, string(pair.Key), "incorrect key on index %d", i)

		var value map[string]interface{}
		require.NoError(json.Unmarshal(pair.Value, &value), "could not unmarshal value from db")
		require.Equal(expected.Value, value)
		i++
	}
	require.Equal(5, i, "expected 5 results returned, have fixtures changed?")

	// Test No Keys
	stream, err = client.Cursor(ctx, &pb.CursorRequest{Namespace: "people", Prefix: []byte("216"), Options: &pb.Options{IterNoKeys: true}})
	require.NoError(err, "could not create cursor stream")

	for {
		pair, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(err, "received non-EOF error from recv")

		require.Empty(pair.Key, "key returned on no keys")
		require.NotEmpty(pair.Value, "value not returned")
		require.Empty(pair.Meta, "meta returned without request")
	}

	// Test No Values
	stream, err = client.Cursor(ctx, &pb.CursorRequest{Namespace: "people", Prefix: []byte("216"), Options: &pb.Options{IterNoValues: true}})
	require.NoError(err, "could not create cursor stream")

	for {
		pair, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(err, "received non-EOF error from recv")

		require.NotEmpty(pair.Key, "key not returned")
		require.Empty(pair.Value, "value returned on no values")
		require.Empty(pair.Meta, "meta returned without request")
	}

	// Test Return Meta
	stream, err = client.Cursor(ctx, &pb.CursorRequest{Namespace: "people", Prefix: []byte("216"), Options: &pb.Options{ReturnMeta: true}})
	require.NoError(err, "could not create cursor stream")

	for {
		pair, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(err, "received non-EOF error from recv")

		require.NotEmpty(pair.Key, "key not returned")
		require.NotEmpty(pair.Value, "value not returned")
		require.NotEmpty(pair.Meta, "meta not returned on request")
	}
}
