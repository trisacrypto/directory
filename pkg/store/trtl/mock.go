package trtl

import (
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc"
)

// NewMock creates a mocked trtl store from a custom grpc client connection to enable
// testing with bufconn. This method avoids the Open method, so the mock will not
// reindex nor run the checkpointing go routine without an explicit call after mock.
// However it is important to ensure that all data structures are created in the mock
// function to avoid panics during testing.
func NewMock(conn *grpc.ClientConn) (store *Store, err error) {
	store = &Store{
		conn: conn,
	}
	store.client = pb.NewTrtlClient(store.conn)

	if err = store.sync(); err != nil {
		return nil, err
	}
	return store, nil
}
