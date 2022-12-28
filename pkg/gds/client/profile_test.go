package client_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/client"
)

// Test that OpenLevelDB opens a LevelDB database.
func TestOpenLevelDB(t *testing.T) {
	// No DB specified
	profile := &client.Profile{}
	_, err := profile.OpenLevelDB()
	require.Error(t, err)

	// Invalid DB scheme
	profile = &client.Profile{DatabaseURL: "foo://bar"}
	_, err = profile.OpenLevelDB()
	require.Error(t, err)

	// Valid DB DSN
	dbpath, err := os.MkdirTemp("", "gdsdb-*")
	require.NoError(t, err, "could not create temp directory")
	defer os.RemoveAll(dbpath)

	profile = &client.Profile{DatabaseURL: "leveldb:///" + dbpath}
	db, err := profile.OpenLevelDB()
	require.NoError(t, err)
	defer db.Close()
}
