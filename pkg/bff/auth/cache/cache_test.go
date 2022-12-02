package cache_test

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/auth/cache"
	"github.com/trisacrypto/directory/pkg/bff/config"
)

type TimeFetcher struct{}

func (f *TimeFetcher) Get(id string) (data interface{}, err error) {
	return time.Now(), nil
}

func TestTTLCache(t *testing.T) {
	// Disabled cache returns different values
	conf := config.CacheConfig{}
	items := cache.NewTTLCache(conf)
	items.SetFetcher(&TimeFetcher{})
	t1, err := items.Get("1")
	require.NoError(t, err, "could not get item from cache")
	t2, err := items.Get("1")
	require.NoError(t, err, "could not get item from cache")
	require.NotEqual(t, t1, t2, "if caching is disabled then different values should be returned")

	conf = config.CacheConfig{
		Enabled:          true,
		TTLMean:          time.Millisecond * 100,
		TTLSigma:         time.Millisecond * 10,
		MaxEntries:       100,
		EvictionFraction: 0.1,
	}
	items = cache.NewTTLCache(conf)
	items.SetFetcher(&TimeFetcher{})
	items.SetStore(cache.NewTTLStore(conf.TTLMean, conf.TTLSigma))

	// Fetch a key from the cache
	t1, err = items.Get("time")
	require.NoError(t, err, "could not get key from cache")
	require.NotNil(t, t1, "cache should return a value")

	// Wait for the cache entry to expire
	time.Sleep(conf.TTLMean + conf.TTLSigma*10)

	// Fetch the key again, should be a new value
	t2, err = items.Get("time")
	require.NoError(t, err, "could not get key from cache")
	require.NotNil(t, t2, "cache should return a value")
	require.NotEqual(t, t1, t2, "cache should return a new value after expiration")

	// Warm up the cache
	conf = config.CacheConfig{
		Enabled:          true,
		TTLMean:          time.Minute,
		TTLSigma:         time.Second,
		MaxEntries:       100,
		EvictionFraction: 0.1,
	}
	items = cache.NewTTLCache(conf)
	items.SetFetcher(&TimeFetcher{})
	items.SetStore(cache.NewTTLStore(conf.TTLMean, conf.TTLSigma))
	values := make([]interface{}, 0)
	for i := 0; i < conf.MaxEntries; i++ {
		val, err := items.Get(strconv.Itoa(i))
		require.NoError(t, err, "could not get key from cache")
		require.NotNil(t, val, "cache should return a value")
		values = append(values, val)
	}

	// Fetching an existing key should not trigger eviction
	val, err := items.Get("0")
	require.NoError(t, err, "could not get key from cache")
	require.NotNil(t, val, "cache should return a value")
	require.Equal(t, values[0], val, "cache should return the same value")

	// Fetching a new key should trigger eviction
	val, err = items.Get(strconv.Itoa(conf.MaxEntries))
	require.NoError(t, err, "could not get key from cache")
	require.NotNil(t, val, "cache should return a value")

	// Eviction should remove some of the keys
	for i := 0; i < 10; i++ {
		val, err := items.Get(strconv.Itoa(i))
		require.NoError(t, err, "could not get key from cache")
		require.NotNil(t, val, "cache should return a value")
		require.NotEqual(t, values[i], val, "cache should return a new value after eviction")
	}

	// Test concurrent access to the cache
	conf = config.CacheConfig{
		Enabled:          true,
		TTLMean:          time.Millisecond * 100,
		TTLSigma:         time.Millisecond * 10,
		MaxEntries:       1000,
		EvictionFraction: 0.1,
	}
	items = cache.NewTTLCache(conf)
	items.SetFetcher(&TimeFetcher{})
	items.SetStore(cache.NewTTLStore(conf.TTLMean, conf.TTLSigma))
	var wg sync.WaitGroup
	for i := 0; i < conf.MaxEntries*2; i++ {
		i := i
		go func() {
			val, err := items.Get(strconv.Itoa(i))
			require.NoError(t, err, "could not get key from cache")
			require.NotNil(t, val, "cache should return a value")
			wg.Done()
		}()
		wg.Add(1)

		go func() {
			val, err := items.Get(strconv.Itoa(i))
			require.NoError(t, err, "could not get key from cache")
			require.NotNil(t, val, "cache should return a value")
			wg.Done()
		}()
		wg.Add(1)
	}
	wg.Wait()
}
