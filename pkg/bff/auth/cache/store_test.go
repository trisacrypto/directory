package cache_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/auth/cache"
)

func TestStructStore(t *testing.T) {
	// The store is initially empty
	store := cache.NewStructStore()
	_, err := store.Get("foo")
	require.ErrorIs(t, err, cache.ErrKeyNotFound)

	// Put a value in the store and retrieve it
	require.NoError(t, store.Put("foo", "bar"))
	value, err := store.Get("foo")
	require.NoError(t, err, "could not get value from store")
	require.Equal(t, "bar", value)

	// Overwrite the value in the store and retrieve it
	require.NoError(t, store.Put("foo", "baz"))
	value, err = store.Get("foo")
	require.NoError(t, err, "could not get value from store")
	require.Equal(t, "baz", value)

	// Test concurrent access to the store
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			require.NoError(t, store.Put("foo", "bar"))
			value, err := store.Get("foo")
			require.NoError(t, err, "could not get value from store")
			require.Equal(t, "bar", value)
		}()
	}
	wg.Wait()
}
