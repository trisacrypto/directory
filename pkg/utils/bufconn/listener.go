package bufconn

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// GRPCListener handles gRPC connections using a bufconn listener. This is useful for
// testing when it's unnecessary to have a live gRPC server running. The normal
// workflow is to call New() to start the listener, Connect() to start a gRPC
// connection under which to send client calls, Close() to close the connection, and
// Release() to close the listener.
type GRPCListener struct {
	Listener *bufconn.Listener
	Conn     *grpc.ClientConn
}

func New(bufSize int) *GRPCListener {
	return &GRPCListener{
		Listener: bufconn.Listen(bufSize),
	}
}

func (g *GRPCListener) Connect() (err error) {
	ctx := context.Background()
	if g.Conn, err = grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(g.dialer), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return err
	}
	return nil
}

func (g *GRPCListener) dialer(context.Context, string) (net.Conn, error) {
	return g.Listener.Dial()
}

func (g *GRPCListener) Close() {
	g.Conn.Close()
}

func (s *GRPCListener) Release() {
	s.Listener.Close()
}
