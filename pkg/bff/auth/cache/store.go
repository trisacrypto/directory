package cache

import (
	"math/rand"
	"time"
)

// KeyValueStore allows keys to be stored, retrieved, and deleted.
type KeyValueStore interface {
	Get(key string) (interface{}, error)
	Put(key string, value interface{}) error
	Delete(key string) error
}

// TTLStore is an in-memory key value store with values that expire.
type TTLStore struct {
	ttlMean  time.Duration
	ttlSigma time.Duration
	entries  map[string]*ttlEntry
}

func NewTTLStore(ttlMean, ttlSigma time.Duration) *TTLStore {
	return &TTLStore{
		ttlMean:  ttlMean,
		ttlSigma: ttlSigma,
		entries:  make(map[string]*ttlEntry),
	}
}

type ttlEntry struct {
	value   interface{}
	expires time.Time
}

// TTLStore implements the KeyValueStore interface
var _ KeyValueStore = &TTLStore{}

// Get returns the value by key or an error if the key is missing or the value has
// expired.
func (s *TTLStore) Get(key string) (value interface{}, err error) {
	var (
		entry *ttlEntry
		ok    bool
	)
	if entry, ok = s.entries[key]; !ok {
		return nil, ErrKeyNotFound
	}

	if time.Now().After(entry.expires) {
		return nil, ErrValueExpired
	}

	return entry.value, nil
}

func (s *TTLStore) Put(key string, value interface{}) error {
	s.entries[key] = &ttlEntry{
		value:   value,
		expires: expiration(s.ttlMean, s.ttlSigma),
	}
	return nil
}

func (s *TTLStore) Delete(key string) error {
	delete(s.entries, key)
	return nil
}

// Helper to compute expiration times for TTL entries.
func expiration(mean, sigma time.Duration) time.Time {
	// Use a normal distribution of jitter to avoid cache stampedes
	ns := rand.NormFloat64() * float64(sigma)
	return time.Now().Add(mean + time.Duration(ns))
}
