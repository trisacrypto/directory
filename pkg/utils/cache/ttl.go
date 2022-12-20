package cache

import (
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/trisacrypto/directory/pkg/bff/config"
)

// TTL is a wrapper around a thread-safe LRU cache that adds a time-to-live value to
// each entry.
type TTL struct {
	conf config.CacheConfig
	lru  *lru.Cache
}

func New(conf config.CacheConfig) (Cache, error) {
	if !conf.Enabled {
		return &Disabled{}, nil
	}
	return NewTTL(conf)
}

func NewTTL(conf config.CacheConfig) (cache *TTL, err error) {
	cache = &TTL{
		conf: conf,
	}

	if cache.lru, err = lru.New(int(conf.Size)); err != nil {
		return nil, err
	}

	return cache, nil
}

// TTL implements the Cache interface
var _ Cache = &TTL{}

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
func (c *TTL) Get(key interface{}) (data interface{}, ok bool) {
	// This will panic if the value is not a ttlEntry so we must prevent external
	// packages from accessing the underlying cache directly.
	var value interface{}
	if value, ok = c.lru.Get(key); !ok || value.(*ttlEntry).Expired() {
		return nil, false
	}

	return value.(*ttlEntry).data, true
}

// Add stores data in the cache by key.
func (c *TTL) Add(key interface{}, data interface{}) {
	// Create a new entry and store it in the cache
	c.lru.Add(key, &ttlEntry{
		data:    data,
		expires: time.Now().Add(c.conf.Expiration),
	})
}

// Remove removes data from the cache by key.
func (c *TTL) Remove(key interface{}) {
	c.lru.Remove(key)
}
