package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/store"
)

func TestDSNParsing(t *testing.T) {
	cases := []struct {
		uri string
		dsn *store.DSN
	}{
		{"leveldb:///fixtures/db", &store.DSN{"leveldb", "fixtures/db"}},
		{"sqlite3:///fixtures/db", &store.DSN{"sqlite3", "fixtures/db"}},
		{"leveldb:////data/db", &store.DSN{"leveldb", "/data/db"}},
		{"sqlite3:////data/db", &store.DSN{"sqlite3", "/data/db"}},
	}

	for _, tc := range cases {
		dsn, err := store.ParseDSN(tc.uri)
		require.NoError(t, err)
		require.Equal(t, tc.dsn, dsn)
	}

	// Test error cases
	_, err := store.ParseDSN("foo")
	require.Error(t, err)

	_, err = store.ParseDSN("foo://")
	require.Error(t, err)
}
