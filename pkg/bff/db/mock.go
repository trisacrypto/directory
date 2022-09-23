package db

import (
	"github.com/trisacrypto/directory/pkg/store/trtl"
	"google.golang.org/grpc"
)

// NewMock creates a new mock DB object containing a store client interface. This
// method bypasses the store.Open() method to allow a bufconn trtl connection to be
// directly injected into the mock for testing purposes.
func NewMock(conn *grpc.ClientConn) (db *DB, err error) {
	db = &DB{}
	if db.db, err = trtl.NewMock(conn); err != nil {
		return nil, err
	}

	return db, nil
}
