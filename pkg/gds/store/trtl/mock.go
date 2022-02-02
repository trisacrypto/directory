package trtl

import (
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc"
)

// NewMock creates a mocked trtl store from a custom grpc client connection to enable
// testing with bufconn.
func NewMock(conn *grpc.ClientConn) (store *Store, err error) {
	store = &Store{
		conn: conn,
	}
	store.client = pb.NewTrtlClient(store.conn)
	return store, nil
}
