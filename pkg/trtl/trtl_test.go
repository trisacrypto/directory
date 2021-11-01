package trtl_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc"
)

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
