package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/store/config"
	"github.com/trisacrypto/directory/pkg/store/mock"
)

func TestOpenMockStore(t *testing.T) {
	mockdb, err := store.Open(config.StoreConfig{
		URL:           "mock:///",
		ReindexOnBoot: false,
	})
	require.NoError(t, err, "should be able to open a mock store with mock:///")

	_, ok := mockdb.(*mock.Store)
	require.True(t, ok, "a mock database store was returned")
}

func TestDSNParsing(t *testing.T) {
	cases := []struct {
		uri string
		dsn *store.DSN
	}{
		{"leveldb:///fixtures/db", &store.DSN{"leveldb", "fixtures/db"}},
		{"leveldb:////data/db", &store.DSN{"leveldb", "/data/db"}},
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
