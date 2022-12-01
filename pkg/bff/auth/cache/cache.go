package cache

import (
	"math/rand"
	"sync"
	"time"

	"github.com/trisacrypto/directory/pkg/bff/config"
)



// TTLCache is a time-to-live cache that relies on a key value store for persistence,
// automatically refreshing entries when they expire and uses a FIFO eviction policy
// to manage the cache size.
type TTLCache struct {
	sync.RWMutex
	conf    config.CacheConfig
	entries   TTLStore
	keys    []string
}

type ttlEntry struct {
	value   interface{}
	expires time.Time
}

func NewTTLEntry(conf config.CacheConfig) *ttlEntry {
	// Add a normal distribution of jitter to the TTL to avoid cache stampedes
	ns := rand.NormFloat64() * float64(conf.TTLSigma)
	return &ttlEntry{
		expires: time.Now().Add(conf.TTLMean + time.Duration(ns)),
	}
}

// NewUserCache creates a new user cache.
func NewTTLCache(conf config.CacheConfig, fetcher ResourceFetcher, store KeyValueStore) *TTLCache {
	return &TTLCache{
		conf:    conf,
		index:   make(map[string]*ttlEntry),
		fetcher: fetcher,
		store:   store,
		keys:    make([]string, 0),
	}
}

// Get returns the user record from the cache
func (c *TTLCache) Get(key string) (value interface{}, err error) {
	c.RLock()
	entry, ok := c.index[key]
	c.RUnlock()

	// Fetch the new or current version of the resource
	if !ok || time.Now().After(entry.expires) {
		if value, err = c.fetcher.Get(key); err != nil {
			return nil, err
		}
	} else if value, err = c.store.Get(key); err == nil {
		return value, nil
	} else {
		return nil, err
	}

	c.Lock()
	defer c.Unlock()
	entry, ok = c.index[key]
	if !ok {
		// If the cache is full, evict the oldest entry
		if len(c.index) >= c.conf.MaxSize {
			key := c.keys[0]
			c.keys = c.keys[1:]
			delete(c.index, key)
			if err = c.store.Delete(key); err != nil {
				return nil, err
			}
		}

		// Add the new entry to the cache

	}
	if entry, ok = c.index[key]; !ok || time.Now().After(entry.expires) {
		// Update the cache entry
		c.index[key] = NewTTLEntry(c.conf)
		c.keys = append(c.keys, key)


	// If the value is not in the cache, add it
	if !ok || time.Now().After(entry.expires) {
		c.Lock()
		c.index[key] = NewTTLEntry(c.conf)
		c.Unlock()

		if err = c.store.Put(key, value); err != nil {
			return nil, err
		}
	} else if value, err = c.store.Get(key); err != nil {
		return nil, err
	}

	return value, nil
}