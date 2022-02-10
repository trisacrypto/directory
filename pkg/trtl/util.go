package trtl

import "github.com/trisacrypto/directory/pkg/trtl/internal"

// SeekCursor returns a next page token that will cause the iter method to
// start pagination at the specified key. This is a hack to support the Seek() method in
// the trtl store iterator interface. In the future, the Iter method should be able to
// seek via an RPC.
func SeekCursor(pageSize int32, nextKey []byte, namespace string) (string, error) {
	cursor := &internal.PageCursor{
		PageSize:  pageSize,
		NextKey:   nextKey,
		Namespace: namespace,
	}
	return cursor.Dump()
}
