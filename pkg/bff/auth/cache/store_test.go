package cache_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/auth/cache"
)

func TestStructStore(t *testing.T) {
	// The store is initially empty
	store := cache.NewTTLStore(time.Minute, time.Second)
	_, err := store.Get("foo")
	require.ErrorIs(t, err, cache.ErrKeyNotFound)

	// Should not error deleting a missing key
	require.NoError(t, store.Delete("missing"))

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

	// Add a second value to the store and retrieve it
	require.NoError(t, store.Put("bar", "baz2"))
	value, err = store.Get("bar")
	require.NoError(t, err, "could not get value from store")
	require.Equal(t, "baz2", value)

	// Delete the first value from the store, should not be able to retrieve it
	require.NoError(t, store.Delete("foo"))
	_, err = store.Get("foo")
	require.ErrorIs(t, err, cache.ErrKeyNotFound)

	// After the expiration time values cannot be retrieved
	ttlMean := time.Millisecond * 100
	ttlSigma := time.Millisecond * 10
	store = cache.NewTTLStore(ttlMean, ttlSigma)
	require.NoError(t, store.Put("foo", "bar"))
	time.Sleep(ttlMean + ttlSigma*10)
	_, err = store.Get("foo")
	require.ErrorIs(t, err, cache.ErrValueExpired)
}
