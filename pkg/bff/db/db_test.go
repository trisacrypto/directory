package db_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/bff/db"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/mock"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
)

const (
	bufsize = 1024 * 1024
)

// The DB Test Suite creates a running trtl process in memory and a database connection
// to it through bufconn, allowing us to test the BFF interactions with a live trtl db.
// This suite is intended to specifically test the database code in this package.
type dbTestSuite struct {
	suite.Suite
	db   *db.DB
	trtl *trtl.Server
	conn *bufconn.GRPCListener
}

func TestDatabaseSuite(t *testing.T) {
	s := new(dbTestSuite)
	suite.Run(t, s)
}

func (s *dbTestSuite) SetupSuite() {
	var err error
	require := s.Require()

	// Create a trtl config that points trtl to a temporary directory for storage
	conf := mock.Config()
	conf.Database.URL = "leveldb:///" + s.T().TempDir()
	conf, err = conf.Mark()
	require.NoError(err, "could not mark config as processed")

	// Create a bufconn to connect to the trtl server on
	s.conn = bufconn.New(bufsize, "")
	require.NoError(s.conn.Connect(context.Background()), "could not dial the bufconn")

	// Start the Trtl server
	s.trtl, err = trtl.New(conf)
	require.NoError(err, "could not create the trtl server")
	go s.trtl.Run(s.conn.Listener)

	// Connect our database to the running Trtl server
	s.db, err = db.DirectConnect(s.conn.Conn)
	require.NoError(err, "could not connect the database to the trtl server")
}

func (s *dbTestSuite) TearDownSuite() {
	var err error
	require := s.Require()

	err = s.db.Close()
	require.NoError(err, "could not close database")

	err = s.trtl.Shutdown()
	require.NoError(err, "could not shutdown trtl in-memory process")

	s.conn.Release()
}

func (s *dbTestSuite) TestBasicOperations() {
	// Test Put, Get, and Delete against the database
	var err error
	require := s.Require()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := []byte("thisisthetestkey")
	value := []byte("thisisthetestvaluehawtness")
	namespace := "thisisthetestnamespace"

	// Expect a not found error when the key is not there
	_, err = s.db.Get(ctx, key, namespace)
	require.ErrorIs(err, db.ErrNotFound, "expected not found error before Put")

	// Should be able to successfully Put the key-value pair
	err = s.db.Put(ctx, key, value, namespace)
	require.NoError(err, "could not Put to the database")

	retrieved, err := s.db.Get(ctx, key, namespace)
	require.NoError(err, "could not fetch key just put to db")
	require.True(bytes.Equal(value, retrieved), "retrieved value not identical to original")

	// Should have been put to the correct namespace
	_, err = s.db.Get(ctx, key, "thisisnotthetestnamespace")
	require.ErrorIs(err, db.ErrNotFound, "expected not found error on wrong namespace")

	// Should be able to Delete the key-value pair
	err = s.db.Delete(ctx, key, namespace)
	require.NoError(err, "could not Delete key from the database")

	_, err = s.db.Get(ctx, key, namespace)
	require.ErrorIs(err, db.ErrNotFound, "expected not found error after Delete")
}
