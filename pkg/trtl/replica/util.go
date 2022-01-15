package replica

import (
	"encoding/base64"
	"sync"

	"github.com/rotationalio/honu/object"
)

// b64e encodes []byte keys and values as base64 encoded strings suitable for logging.
func b64e(src []byte) string {
	return base64.RawURLEncoding.EncodeToString(src)
}

// nsmap is a lightweight tool for keeping track of what objects we've seen during
// gossip. It is threadsafe and implements set methods. This should be replaced by a
// bloom filter for memory saving increases in performance
type nsmap struct {
	sync.RWMutex
	seen map[string]map[string]struct{}
}

func (m *nsmap) Add(obj *object.Object) {
	m.Lock()
	defer m.Unlock()
	if m.seen == nil {
		m.seen = make(map[string]map[string]struct{})
	}

	if _, ok := m.seen[obj.Namespace]; !ok {
		m.seen[obj.Namespace] = make(map[string]struct{})
	}

	m.seen[obj.Namespace][b64e(obj.Key)] = struct{}{}
}

func (m *nsmap) In(namespace string, key []byte) bool {
	m.RLock()
	defer m.RUnlock()
	if m.seen != nil {
		if nsm, nok := m.seen[namespace]; nok {
			_, kok := nsm[b64e(key)]
			return kok
		}
	}
	return false
}
