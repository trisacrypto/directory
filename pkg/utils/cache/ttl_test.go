package cache_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/utils/cache"
)

func TestTTLCache(t *testing.T) {
	// Configure the cache
	conf := config.CacheConfig{
		Size:       100,
		Expiration: time.Millisecond,
	}
	items, err := cache.NewTTL(conf)
	require.NoError(t, err, "could not create cache")

	// Should return false for a non-existent key
	_, ok := items.Get("foo")
	require.False(t, ok, "cache should not return a value for a non-existent key")

	// Should be able to add and get a value
	items.Add("foo", "bar")
	val, ok := items.Get("foo")
	require.True(t, ok, "cache should return a value for an existing key")
	require.Equal(t, "bar", val, "cache should return the correct value")

	// After expiration the value should be removed
	time.Sleep(2 * conf.Expiration)
	_, ok = items.Get("foo")
	require.False(t, ok, "cache should not return an expired value")

	// Should be able to remove values
	items.Add("foo", "bar")
	items.Remove("foo")
	_, ok = items.Get("foo")
	require.False(t, ok, "cache should not return a removed value")

	// Fill up the cache
	conf = config.CacheConfig{
		Size:       100,
		Expiration: time.Minute,
	}
	items, err = cache.NewTTL(conf)
	require.NoError(t, err, "could not create cache")
	for i := 0; i < int(conf.Size); i++ {
		items.Add(i, i)
	}

	// With LRU eviction the oldest value should be removed
	items.Add("foo", "bar")
	_, ok = items.Get(0)
	require.False(t, ok, "expected first value to be evicted")

	// The second value should still be accessible
	val, ok = items.Get(1)
	require.True(t, ok, "second value was evicted")
	require.Equal(t, 1, val, "second value did not match")
}
