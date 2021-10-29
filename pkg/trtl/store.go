package trtl

import (
	"errors"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/iterator"
	"github.com/rotationalio/honu/object"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
)

// Store is an interface for interacting with an underlying honu database.
type Store interface {
	Close() error
	KeyValueStore
	PeerStore
}

// KeyValueStore is an interface for interacting with key-value pairs in honu.
type KeyValueStore interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	Iter(prefix []byte) (iterator.Iterator, error)
	Object(key []byte) (*object.Object, error)
}

// PeerStore is an interface for interacting with Peer objects stored in honu.
type PeerStore interface {
	CreatePeer(peer *peers.Peer) (string, error)
	DeletePeer(id string) error
	ListPeers() iterator.Iterator
	AllPeers() ([]*peers.Peer, error)
}

type HonuStore struct {
	db *honu.DB
}

func NewHonuStore(db *honu.DB) *HonuStore {
	return &HonuStore{db: db}
}

func (s *HonuStore) Close() error {
	return s.db.Close()
}

func (s *HonuStore) Get(key []byte) ([]byte, error) {
	return s.db.Get(key)
}

func (s *HonuStore) Put(key []byte, value []byte) error {
	return s.db.Put(key, value)
}

func (s *HonuStore) Delete(key []byte) error {
	return s.db.Delete(key)
}

func (s *HonuStore) Iter(prefix []byte) (iterator.Iterator, error) {
	return s.db.Iter(prefix)
}

func (s *HonuStore) Object(key []byte) (*object.Object, error) {
	return s.db.Object(key)
}

func (s *HonuStore) CreatePeer(peer *peers.Peer) (string, error) {
	return "", errors.New("not implemented")
}

func (s *HonuStore) DeletePeer(id string) error {
	return errors.New("not implemented")
}

func (s *HonuStore) ListPeers() iterator.Iterator {
	return iterator.NewEmptyIterator(errors.New("not implemented"))
}

func (s *HonuStore) AllPeers() ([]*peers.Peer, error) {
	return nil, errors.New("not implemented")
}
