package cache

// Cache is a generic interface that can be implemented by most types of caches.
type Cache interface {
	Get(key interface{}) (data interface{}, ok bool)
	Add(key interface{}, data interface{})
	Remove(key interface{})
}
