package cache

import (
	"errors"
	"sync"

	"github.com/trisacrypto/directory/pkg/bff/config"
)

type ResourceFetcher interface {
	Get(key string) (interface{}, error)
}

// TTLCache is a time-to-live cache that relies on a key value store for persistence,
// automatically refreshing entries when they expire and uses a FIFO eviction policy
// to manage the cache size.
type TTLCache struct {
	sync.RWMutex
	fetcher ResourceFetcher
	store   KeyValueStore
	conf    config.CacheConfig
	keys    []string
}

// NewUserCache creates a new user cache.
func NewTTLCache(conf config.CacheConfig) *TTLCache {
	return &TTLCache{
		conf: conf,
		keys: make([]string, 0),
	}
}

func (c *TTLCache) SetFetcher(fetcher ResourceFetcher) {
	c.fetcher = fetcher
}

func (c *TTLCache) SetStore(store KeyValueStore) {
	c.store = store
}

// Get returns a value by key from the cache. If the key is not found or the value has
// expired, we have to get a new value from the fetcher and handle evictions to abide
// by the cache size limit.
func (c *TTLCache) Get(key string) (value interface{}, err error) {
	if !c.conf.Enabled {
		return c.fetcher.Get(key)
	}

	c.RLock()
	if value, err = c.store.Get(key); err != nil {
		if errors.Is(err, ErrKeyNotFound) || errors.Is(err, ErrValueExpired) {
			// If the key is not found or the value is expired, fetch the new value
			if value, err = c.fetcher.Get(key); err != nil {
				c.RUnlock()
				return nil, err
			}
		} else {
			c.RUnlock()
			return nil, err
		}
	} else {
		c.RUnlock()
		return value, nil
	}
	c.RUnlock()

	// We have a new value so update the cache
	c.Lock()
	defer c.Unlock()

	// Maintain the maximum cache size by evicting the oldest entries
	if len(c.keys) >= c.conf.MaxEntries {
		evictions := int(float64(len(c.keys)) * c.conf.EvictionFraction)
		for i := 0; i < evictions; i++ {
			if err = c.store.Delete(c.keys[0]); err != nil {
				return nil, err
			}
			c.keys = c.keys[1:]
		}
	}

	if err = c.store.Put(key, value); err != nil {
		return nil, err
	}
	c.keys = append(c.keys, key)

	return value, nil
}
