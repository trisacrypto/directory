package trtl_test

import (
	"context"
	"fmt"

	"github.com/trisacrypto/directory/pkg/store/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *trtlStoreTestSuite) TestBatchIterator() {
	require := s.Require()

	// Connect to bufconn and get a trtl gRPC client
	require.NoError(s.grpc.Connect(context.Background()))
	client := pb.NewTrtlClient(s.grpc.Conn)

	// Iterate over an empty namespace
	iter := trtl.NewTrtlBatchIterator(client, "empty")
	require.False(iter.Seek([]byte("foo")), "Seek to non-existent key should return false")
	require.Nil(iter.Key(), "Key should be nil when iterator is exhausted")
	require.Nil(iter.Value(), "Value should be nil when iterator is exhausted")
	require.False(iter.Next(), "Next should return false when iterator is exhausted")
	require.Nil(iter.Value(), "Value should still be nil after calling Next")
	require.False(iter.Prev(), "Prev should return false when iterator is exhausted")
	require.Nil(iter.Value(), "Value should still be nil after calling Prev")
	require.Nil(iter.Error())
	iter.Release()

	// Namespace that contains one entry
	req := &pb.PutRequest{
		Key:       []byte("foo"),
		Value:     []byte("bar"),
		Namespace: "batch-single",
	}
	_, err := client.Put(context.Background(), req)
	require.NoError(err)

	// Calling Prev() first
	iter = trtl.NewTrtlBatchIterator(client, req.Namespace)
	require.False(iter.Prev(), "Initial Prev should return false")
	require.Nil(iter.Value(), "Value should be nil when iterator is behind the first key")
	require.True(iter.Next(), "Next should return true when iterator is behind the first key")
	require.Equal(req.Key, iter.Key())
	require.Equal(req.Value, iter.Value())
	require.Nil(iter.Error())
	iter.Release()

	// Calling Next() first
	iter = trtl.NewTrtlBatchIterator(client, req.Namespace)
	require.True(iter.Next(), "Initial Next should return true")
	require.Equal(req.Key, iter.Key())
	require.Equal(req.Value, iter.Value())
	require.False(iter.Next(), "Next should return false when iterator is exhausted")
	require.Nil(iter.Value(), "Value should be nil when iterator is exhausted")
	require.False(iter.Seek(req.Key), "Seek should return false if Next has already been called")
	require.NotNil(iter.Error())
	iter.Release()

	// Calling Seek() first
	iter = trtl.NewTrtlBatchIterator(client, req.Namespace)
	require.True(iter.Seek(req.Key), "Seek to existing key should return true")
	require.Equal(req.Key, iter.Key())
	require.Equal(req.Value, iter.Value())
	require.False(iter.Prev(), "Prev should return false when iterator is at first key")
	require.Nil(iter.Value(), "Value should be nil when iterator is behind first key")
	require.True(iter.Next(), "Next should return true when iterator is behind first key")
	require.Equal(req.Key, iter.Key())
	require.Equal(req.Value, iter.Value())
	require.False(iter.Next(), "Next should return false when iterator is at last key")
	require.Nil(iter.Value(), "Value should be nil when iterator is ahead of last key")
	require.True(iter.Prev(), "Prev should return true when iterator is ahead of last key")
	require.Equal(iter.Key(), req.Key)
	require.Equal(iter.Value(), req.Value)
	require.Nil(iter.Error())
	iter.Release()

	// Namespace that contains a page of entries
	page := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
		"x": "4",
		"y": "5",
	}
	for k, v := range page {
		req = &pb.PutRequest{
			Key:       []byte(k),
			Value:     []byte(v),
			Namespace: "page",
		}
		_, err = client.Put(context.Background(), req)
		require.NoError(err)
	}
	// Cannot call Seek after Next
	iter = trtl.NewTrtlBatchIterator(client, req.Namespace)
	require.True(iter.Next(), "First Next should return true")
	require.Equal([]byte("a"), iter.Key())
	require.Equal([]byte(page["a"]), iter.Value())
	require.False(iter.Seek([]byte("c")), "Cannot call Seek after Next")
	require.NotNil(iter.Error())
	iter.Release()

	// Seek to the end
	iter = trtl.NewTrtlBatchIterator(client, req.Namespace)
	require.False(iter.Seek([]byte("z")), "Seek to the end of the page should return false")
	require.Nil(iter.Key(), "Key should be nil when iterator is exhausted")
	require.Nil(iter.Value(), "Value should be nil when iterator is exhausted")
	require.False(iter.Next(), "Next should return false when iterator is exhausted")
	require.Nil(iter.Error())
	iter.Release()

	// Seek to the middle
	iter = trtl.NewTrtlBatchIterator(client, req.Namespace)
	require.True(iter.Seek([]byte("c")), "Seek to existing key should return true")
	require.Equal([]byte("c"), iter.Key())
	require.Equal([]byte(page["c"]), iter.Value())
	require.True(iter.Next(), "Next should return true when iterator is at first key")
	require.Equal([]byte("x"), iter.Key())
	require.Equal([]byte(page["x"]), iter.Value())
	require.True(iter.Next(), "Next should return true when iterator is at second key")
	require.Equal([]byte("y"), iter.Key())
	require.Equal([]byte(page["y"]), iter.Value())
	require.True(iter.Prev(), "Prev should return true when iterator is at third key")
	require.Equal([]byte("x"), iter.Key())
	require.Equal([]byte(page["x"]), iter.Value())
	require.True(iter.Prev(), "Prev should return true when iterator is at second key")
	require.Equal([]byte("c"), iter.Key())
	require.Equal([]byte(page["c"]), iter.Value())
	require.Nil(iter.Error())
	iter.Release()

	// Seek in between keys
	iter = trtl.NewTrtlBatchIterator(client, req.Namespace)
	require.True(iter.Seek([]byte("d")), "Seek in between keys should return true")
	require.Equal([]byte("x"), iter.Key())
	require.Equal([]byte(page["x"]), iter.Value())
	require.False(iter.Prev(), "Prev should return false when iterator is at first key")
	require.Nil(iter.Key(), "Key should be nil when iterator is behind first key")
	require.Nil(iter.Value(), "Value should be nil when iterator is behind first key")
	require.True(iter.Next(), "Next should return true when iterator behind first key")
	require.Equal([]byte("x"), iter.Key())
	require.Equal([]byte(page["x"]), iter.Value())
	require.True(iter.Next(), "Next should return true when iterator is at first key")
	require.Equal([]byte("y"), iter.Key())
	require.Equal([]byte(page["y"]), iter.Value())
	require.Nil(iter.Error())
	iter.Release()

	// Namespace that contains multiple pages of entries
	// 110 exceeds the default page size of 100
	for i := 0; i < 110; i++ {
		req = &pb.PutRequest{
			Key:       []byte(fmt.Sprintf("%d-key", i+1)),
			Value:     []byte(fmt.Sprintf("%d-value", i+1)),
			Namespace: "multi-page",
		}
		_, err := client.Put(context.Background(), req)
		require.NoError(err)
	}

	// Seek to the end of the first page - due to byte ordering the last key in the
	// first page is "9-key"
	iter = trtl.NewTrtlBatchIterator(client, req.Namespace)
	require.True(iter.Seek([]byte("9-key")), "Seek to the end of first page should return true")
	require.Equal([]byte("9-key"), iter.Key())
	require.Equal([]byte("9-value"), iter.Value())
	require.False(iter.Prev(), "Prev should return false when iterator is at first key")
	require.Nil(iter.Key(), "Key should be nil when iterator is behind first key")
	require.Nil(iter.Value(), "Value should be nil when iterator is behind first key")
	require.True(iter.Next(), "Next should return true when iterator is behind first key")
	require.Equal([]byte("9-key"), iter.Key())
	require.Equal([]byte("9-value"), iter.Value())
	require.True(iter.Next(), "Next should return true when iterator is at end of first page")
	require.Equal([]byte("90-key"), iter.Key())
	require.Equal([]byte("90-value"), iter.Value())
	require.NoError(iter.Error())
	iter.Release()

	// Seek to beginning of second page
	iter = trtl.NewTrtlBatchIterator(client, req.Namespace)
	require.True(iter.Seek([]byte("90-key")), "Seek to the beginning of second page should return true")
	require.Equal([]byte("90-key"), iter.Key())
	require.Equal([]byte("90-value"), iter.Value())
	require.False(iter.Prev(), "Prev should return false when iterator is at first key")
	require.Nil(iter.Key(), "Key should be nil when iterator is behind first key")
	require.Nil(iter.Value(), "Value should be nil when iterator is behind first key")
	require.True(iter.Next(), "Next should return true when iterator is behind first key")
	require.Equal([]byte("90-key"), iter.Key())
	require.Equal([]byte("90-value"), iter.Value())
	require.True(iter.Next(), "Next should return true when iterator is at first key")
	require.Equal([]byte("91-key"), iter.Key())
	require.Equal([]byte("91-value"), iter.Value())
	require.NoError(iter.Error())
	iter.Release()

	// Seek to the end of the second page
	iter = trtl.NewTrtlBatchIterator(client, req.Namespace)
	require.True(iter.Seek([]byte("99-key")), "Seek to the end of second page should return true")
	require.Equal([]byte("99-key"), iter.Key())
	require.Equal([]byte("99-value"), iter.Value())
	require.False(iter.Next(), "Next should return false when iterator is at very last key")
	require.Nil(iter.Key(), "Key should be nil when iterator is at very last key")
	require.Nil(iter.Value(), "Value should be nil when iterator is at very last key")
	require.NoError(iter.Error())
	iter.Release()
}

func (s *trtlStoreTestSuite) TestStreamingIterator() {
	require := s.Require()

	// Connect to bufconn and get a trtl gRPC client
	require.NoError(s.grpc.Connect(context.Background()))
	client := pb.NewTrtlClient(s.grpc.Conn)

	// Iterate over an empty namespace
	iter := trtl.NewTrtlStreamingIterator(client, "empty")
	require.False(iter.Seek([]byte("foo")), "Seek to non-existent key should return false")
	require.Nil(iter.Key(), "Key should be nil when iterator is exhausted")
	require.Nil(iter.Value(), "Value should be nil when iterator is exhausted")
	require.False(iter.Next(), "Next should return false when iterator is exhausted")
	require.Nil(iter.Value(), "Value should still be nil after calling Next")
	require.False(iter.Prev(), "Prev should return false when iterator is exhausted")
	require.Nil(iter.Value(), "Value should still be nil after calling Prev")
	require.Nil(iter.Error())
	iter.Release()

	// Namespace that contains one entry
	req := &pb.PutRequest{
		Key:       []byte("foo"),
		Value:     []byte("bar"),
		Namespace: "streaming-single",
	}
	_, err := client.Put(context.Background(), req)
	require.NoError(err)

	// Calling Prev() first
	iter = trtl.NewTrtlStreamingIterator(client, req.Namespace)
	require.False(iter.Prev(), "Initial Prev should return false")
	require.Nil(iter.Value(), "Value should be nil when iterator is behind the first key")
	require.True(iter.Next(), "Next should return true when iterator is behind the first key")
	require.Equal(req.Key, iter.Key())
	require.Equal(req.Value, iter.Value())
	require.Nil(iter.Error())
	iter.Release()

	// Calling Next() first
	iter = trtl.NewTrtlStreamingIterator(client, req.Namespace)
	require.True(iter.Next(), "Initial Next should return true")
	require.Equal(req.Key, iter.Key())
	require.Equal(req.Value, iter.Value())
	require.False(iter.Next(), "Next should return false when iterator is exhausted")
	require.Nil(iter.Value(), "Value should be nil when iterator is exhausted")
	require.False(iter.Seek(req.Key), "Seek should return false if Next has already been called")
	require.NotNil(iter.Error())
	iter.Release()

	// Calling Seek() first
	iter = trtl.NewTrtlStreamingIterator(client, req.Namespace)
	require.True(iter.Seek(req.Key), "Seek to existing key should return true")
	require.Equal(req.Key, iter.Key())
	require.Equal(req.Value, iter.Value())
	require.False(iter.Prev(), "Prev should return false when iterator is at first key")
	require.Nil(iter.Value(), "Value should be nil when iterator is behind first key")
	require.True(iter.Next(), "Next should return true when iterator is behind first key")
	require.Equal(req.Key, iter.Key())
	require.Equal(req.Value, iter.Value())
	require.False(iter.Next(), "Next should return false when iterator is at last key")
	require.Nil(iter.Value(), "Value should be nil when iterator is ahead of last key")
	require.True(iter.Prev(), "Prev should return true when iterator is ahead of last key")
	require.Equal(iter.Key(), req.Key)
	require.Equal(iter.Value(), req.Value)
	require.Nil(iter.Error())
	iter.Release()

	// Namespace that contains a page of entries
	page := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
		"x": "4",
		"y": "5",
	}
	for k, v := range page {
		req = &pb.PutRequest{
			Key:       []byte(k),
			Value:     []byte(v),
			Namespace: "streaming-page",
		}
		_, err = client.Put(context.Background(), req)
		require.NoError(err)
	}
	// Cannot call Seek after Next
	iter = trtl.NewTrtlStreamingIterator(client, req.Namespace)
	require.True(iter.Next(), "First Next should return true")
	require.Equal([]byte("a"), iter.Key())
	require.Equal([]byte(page["a"]), iter.Value())
	require.False(iter.Seek([]byte("c")), "Cannot call Seek after Next")
	require.NotNil(iter.Error())
	iter.Release()

	// Seek to the end
	iter = trtl.NewTrtlStreamingIterator(client, req.Namespace)
	require.False(iter.Seek([]byte("z")), "Seek to the end of the page should return false")
	require.Nil(iter.Key(), "Key should be nil when iterator is exhausted")
	require.Nil(iter.Value(), "Value should be nil when iterator is exhausted")
	require.False(iter.Next(), "Next should return false when iterator is exhausted")
	require.Nil(iter.Error())
	iter.Release()

	// Seek to the middle
	iter = trtl.NewTrtlStreamingIterator(client, req.Namespace)
	require.True(iter.Seek([]byte("c")), "Seek to existing key should return true")
	require.Equal([]byte("c"), iter.Key())
	require.Equal([]byte(page["c"]), iter.Value())
	require.True(iter.Next(), "Next should return true when iterator is at first key")
	require.Equal([]byte("x"), iter.Key())
	require.Equal([]byte(page["x"]), iter.Value())
	require.True(iter.Next(), "Next should return true when iterator is at second key")
	require.Equal([]byte("y"), iter.Key())
	require.Equal([]byte(page["y"]), iter.Value())
	require.True(iter.Prev(), "Prev should return true when iterator is at third key")
	require.Equal([]byte("x"), iter.Key())
	require.Equal([]byte(page["x"]), iter.Value())
	require.False(iter.Prev(), "Cannot call Prev twice without calling Next")
	require.Nil(iter.Key())
	require.Nil(iter.Value())
	require.Nil(iter.Error())
	iter.Release()

	// Seek in between keys
	iter = trtl.NewTrtlStreamingIterator(client, req.Namespace)
	require.True(iter.Seek([]byte("d")), "Seek in between keys should return true")
	require.Equal([]byte("x"), iter.Key())
	require.Equal([]byte(page["x"]), iter.Value())
	require.False(iter.Prev(), "Prev should return false when iterator is at first key")
	require.Nil(iter.Key(), "Key should be nil when iterator is behind first key")
	require.Nil(iter.Value(), "Value should be nil when iterator is behind first key")
	require.True(iter.Next(), "Next should return true when iterator behind first key")
	require.Equal([]byte("x"), iter.Key())
	require.Equal([]byte(page["x"]), iter.Value())
	require.True(iter.Next(), "Next should return true when iterator is at first key")
	require.Equal([]byte("y"), iter.Key())
	require.Equal([]byte(page["y"]), iter.Value())
	require.Nil(iter.Error())
	iter.Release()
}

func (s *trtlStoreTestSuite) TestStreamingIteratorError() {
	require := s.Require()
	iter := trtl.NewTrtlStreamingIterator(&trtlErrorClient{}, "")

	// Call interactive methods on iter should do nothing but not return an error.
	for iter.Next() {
		require.Nil(iter.Key())
		require.Nil(iter.Value())
	}

	require.False(iter.Prev())
	require.False(iter.Seek([]byte("foo")))

	// Calling release should not panic
	iter.Release()
	require.Error(iter.Error(), "iter should be in an error state without panic")

}

// Implements pb.TrtlClient but returns an error
type trtlErrorClient struct{}

func (s *trtlErrorClient) Get(context.Context, *pb.GetRequest, ...grpc.CallOption) (*pb.GetReply, error) {
	return nil, status.Error(codes.Unavailable, "trtl is down")
}

func (s *trtlErrorClient) Put(context.Context, *pb.PutRequest, ...grpc.CallOption) (*pb.PutReply, error) {
	return nil, status.Error(codes.Unavailable, "trtl is down")
}

func (s *trtlErrorClient) Delete(context.Context, *pb.DeleteRequest, ...grpc.CallOption) (*pb.DeleteReply, error) {
	return nil, status.Error(codes.Unavailable, "trtl is down")
}

func (s *trtlErrorClient) Iter(context.Context, *pb.IterRequest, ...grpc.CallOption) (*pb.IterReply, error) {
	return nil, status.Error(codes.Unavailable, "trtl is down")
}

func (s *trtlErrorClient) Batch(context.Context, ...grpc.CallOption) (pb.Trtl_BatchClient, error) {
	return nil, status.Error(codes.Unavailable, "trtl is down")
}

func (s *trtlErrorClient) Cursor(context.Context, *pb.CursorRequest, ...grpc.CallOption) (pb.Trtl_CursorClient, error) {
	return nil, status.Error(codes.Unavailable, "trtl is down")
}

func (s *trtlErrorClient) Sync(context.Context, ...grpc.CallOption) (pb.Trtl_SyncClient, error) {
	return nil, status.Error(codes.Unavailable, "trtl is down")
}

func (s *trtlErrorClient) Status(context.Context, *pb.HealthCheck, ...grpc.CallOption) (*pb.ServerStatus, error) {
	return nil, status.Error(codes.Unavailable, "trtl is down")
}
