package db_test

import (
	"context"
	"testing"

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
	s.db, err = db.NewMock(s.conn.Conn)
	require.NoError(err, "could not connect to the trtl server")
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
