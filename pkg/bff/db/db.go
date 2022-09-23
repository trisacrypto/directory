package db

import (
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/store/config"
)

func New(conf config.StoreConfig) (db *DB, err error) {
	db = &DB{}
	if db.db, err = store.Open(conf); err != nil {
		return nil, err
	}
	return db, nil
}

// DB is a wrapper around a store interface to provide access to the directory service
// CRUD operations as well as high level BFF-specific database interactions implemented
// as methods on the DB struct.
type DB struct {
	db store.Store
}

// Close the database.
func (store *DB) Close() error {
	return store.db.Close()
}
