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
	profile = &client.Profile{DatabaseURL: "leveldb:///foo"}
	db, err := profile.OpenLevelDB()
	require.NoError(t, err)
	defer db.Close()
	defer os.RemoveAll("foo")
}
