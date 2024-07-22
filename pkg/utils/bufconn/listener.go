package bufconn

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const Endpoint = "passthrough://bufnet"

const bufsize = 1024 * 1024

// GRPCListener handles gRPC connections using a bufconn listener. This is useful for
// testing when it's unnecessary to have a live gRPC server running. The normal
// workflow is to call New() to start the listener, Connect() to start a gRPC
// connection under which to send client calls, Close() to close the connection, and
// Release() to close the listener.
type GRPCListener struct {
	Listener *bufconn.Listener
	Target   string
	Conn     *grpc.ClientConn
}

func New(target string) *GRPCListener {
	if target == "" {
		target = Endpoint
	}

	return &GRPCListener{
		Listener: bufconn.Listen(bufsize),
		Target:   target,
	}
}

func (g *GRPCListener) Connect(ctx context.Context, opts ...grpc.DialOption) (err error) {
	if len(opts) == 0 {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	opts = append([]grpc.DialOption{grpc.WithContextDialer(g.Dialer)}, opts...)
	if g.Conn, err = grpc.NewClient(g.Target, opts...); err != nil {
		return err
	}
	return err
}

func (g *GRPCListener) Dialer(context.Context, string) (net.Conn, error) {
	return g.Listener.Dial()
}

func (g *GRPCListener) Close() {
	g.Conn.Close()
}

func (s *GRPCListener) Release() {
	s.Listener.Close()
}
