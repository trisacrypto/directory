package cache_test

import (
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
	// Get a value into the cache
	conf := config.CacheConfig{
		TTLMean:  time.Millisecond * 100,
		TTLSigma: time.Millisecond * 10,
	}
	items := cache.NewTTLCache(conf, &TimeFetcher{}, cache.NewStructStore())
	value, err := items.Get("foo")
	require.NoError(t, err, "could not get value from cache")
	require.NotEmpty(t, value, "value should not be empty")

	// Value should still be in the cache
	actual, err := items.Get("foo")
	require.NoError(t, err, "could not get value from cache")
	require.Equal(t, value, actual, "value should not have changed")

	// After the expiration time the value should be refreshed
	time.Sleep(conf.TTLMean + conf.TTLSigma)
	actual, err = items.Get("foo")
	require.NoError(t, err, "could not get value from cache")
	require.NotEmpty(t, actual, "value should not be empty")
	require.NotEqual(t, value, actual, "value should have changed")
}
