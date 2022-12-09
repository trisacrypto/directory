package cache

import (
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/trisacrypto/directory/pkg/bff/config"
)

// TTLCache is a wrapper around a thread-safe LRU cache that adds a time-to-live value
// to each entry.
type TTLCache struct {
	conf config.CacheConfig
	lru  *lru.Cache
}

func New(conf config.CacheConfig) (cache *TTLCache, err error) {
	cache = &TTLCache{
		conf: conf,
	}

	if conf.Enabled {
		if cache.lru, err = lru.New(int(conf.Size)); err != nil {
			return nil, err
		}
	}

	return cache, nil
}

// ttlEntry stores the data and expiration time
type ttlEntry struct {
	data    interface{}
	expires time.Time
}

// Expired returns true if the entry has expired.
func (e *ttlEntry) Expired() bool {
	return e.expires.Before(time.Now())
}

// Get returns data from the cache by key or false if the key does not exist or has
// expired.
func (c *TTLCache) Get(key interface{}) (data interface{}, ok bool) {
	if !c.conf.Enabled {
		return nil, false
	}

	// This will panic if the value is not a ttlEntry so we must prevent external
	// packages from accessing the underlying cache directly.
	var value interface{}
	if value, ok = c.lru.Get(key); !ok || value.(*ttlEntry).Expired() {
		return nil, false
	}

	return value.(*ttlEntry).data, true
}

// Add stores data in the cache by key.
func (c *TTLCache) Add(key interface{}, data interface{}) {
	if c.conf.Enabled {
		// Create a new entry and store it in the cache
		c.lru.Add(key, &ttlEntry{
			data:    data,
			expires: time.Now().Add(c.conf.Expiration),
		})
	}
}

// Remove removes data from the cache by key.
func (c *TTLCache) Remove(key interface{}) {
	if c.conf.Enabled {
		c.lru.Remove(key)
	}
}
