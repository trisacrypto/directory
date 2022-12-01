package cache

import (
	"math/rand"
	"time"

	"github.com/trisacrypto/directory/pkg/bff/config"
)

type ResourceFetcher interface {
	Get(key string) (interface{}, error)
}

// TTLStore is a key-value store that manages expiration for keys.
type TTLStore interface {
	Get(key string) (interface{}, error)
	Put(key string, value interface{}) error
	Delete(key string) error
}

type StructStore struct {
	conf    config.CacheConfig
	fetcher ResourceFetcher
	entries map[string]interface{}
}

func NewStructStore(conf config.CacheConfig, fetcher ResourceFetcher) *StructStore {
	return &StructStore{
		conf:    conf,
		fetcher: fetcher,
		entries: make(map[string]interface{}),
	}
}

type ttlEntry struct {
	value   interface{}
	expires time.Time
}

func expiration(mean, sigma time.Duration) time.Time {
	// Add a normal distribution of jitter to the TTL to avoid cache stampedes
	ns := rand.NormFloat64() * float64(sigma)
	return time.Now().Add(mean + time.Duration(ns))
}

func NewTTLEntry(conf config.CacheConfig) *ttlEntry {
	return &ttlEntry{
		expires: expiration(conf.TTLMean, conf.TTLSigma),
	}
}

// StructStore implements the TTLStore interface
var _ TTLStore = &StructStore{}

func (s *StructStore) Get(key string) (value interface{}, err error) {
	var ok bool
	if data, ok = s.entries[key]; !ok || time.Now().After(value.(*ttlEntry).expires) {
		if value, err = s.fetcher.Get(key); err != nil {
			return nil, err
		}
	}

	return value, nil
}

func (s *StructStore) Put(key string, value interface{}) error {
	s.entries[key] = value
	return nil
}

func (s *StructStore) Delete(key string) error {
	delete(s.entries, key)
	return nil
}
