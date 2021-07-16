package leveldb

import (
	"errors"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/trisacrypto/directory/pkg/gds/global/v1"
	"github.com/trisacrypto/directory/pkg/gds/store/wire"
)

// Iter returns a wrapper for the leveldb iterator that produces global.Object
// references for replication, unmarshaling the object representation as needed. If the
// namespace is provided then a bytes prefix is used to scan over a portion of the db.
func (s *Store) Iter(namespace string) global.ObjectIterator {
	var prefix *util.Range
	if namespace != "" {
		prefix = util.BytesPrefix([]byte(namespace))
	}
	return &objectIterator{ldb: s.db.NewIterator(prefix, nil)}
}

// Get retrieves an object by the specified key. If withData is specified then the
// object bytes are loaded into the Data Any field; otherwise just metadata is returned.
func (s *Store) Get(namespace, key string, withData bool) (obj *global.Object, err error) {
	var data []byte
	if data, err = s.db.Get([]byte(key), nil); err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, wire.ErrObjectNotFound
		}
		return nil, err
	}

	if obj, err = wire.UnmarshalObject(namespace, data, withData); err != nil {
		return nil, err
	}
	return obj, nil
}

// Put puts an object into the specified key, the Data Any field should be present.
// NOTE: there is no delete method because replication requires the Put of a tombstone.
// TODO: implementation required!
func (s *Store) Put(namespace, key string, obj *global.Object) error {
	return errors.New("not implemented yet")
}

// ObjectIterator mirrors the leveldb iterator interface but requires the Store to
// directly manipulate the keys and values into the expected object interface.
type objectIterator struct {
	ldb iterator.Iterator
}

// Next moves the iterator to the next key/value pair. It returns false if the iterator is exhausted.
func (it *objectIterator) Next() bool {
	return it.ldb.Next()
}

// Error returns any accumulated error. Exhausting all the key/value pairs is not
// considered to be an error.
func (it *objectIterator) Error() error {
	return it.ldb.Error()
}

// Key returns the key of the current key/value pair, or nil if done. The caller should
// not modify the contents of the returned slice, and its contents may change on the
// next call to any 'seeks method'.
func (it *objectIterator) Key() string {
	return string(it.ldb.Key())
}

// Object returns the parsed object using wire.UnmarshalObject. If withData is specified
// then the Data Any field is populated, otherwise just metadata is returned.
func (it *objectIterator) Object(withData bool) (*global.Object, error) {
	data := it.ldb.Value()
	prefix := strings.Split(it.Key(), "::")[0]
	return wire.UnmarshalObject(prefix, data, withData)
}

// Release releases associated resources. Release should always success and can be
// called multiple times without causing error.
func (it *objectIterator) Release() {
	it.ldb.Release()
}
