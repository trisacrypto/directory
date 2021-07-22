package global

// ObjectStore allows Replicas to directly interact with the underlying objects managed
// by the interfaces described above, rather than by going through their versioned
// change mechanisms. This is a very low-level interface for object management.
type ObjectStore interface {
	Iter(namespace string) ObjectIterator
	Get(namespace, key string, withData bool) (*Object, error)
	Put(obj *Object) error
}

// ObjectIterator mirrors the leveldb iterator interface but requires the Store to
// directly manipulate the keys and values into the expected object interface.
type ObjectIterator interface {
	Next() bool
	Error() error
	Key() string
	Object(withData bool) (*Object, error)
	Release()
}
