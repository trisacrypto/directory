package trtl

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	trtlpb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

type iterWrapper struct {
	iter *trtlIterator
}

type vaspIterator struct {
	iterWrapper
}

type certReqIterator struct {
	iterWrapper
}

func (i *iterWrapper) Next() bool {
	return i.iter.Next()
}

func (i *iterWrapper) Prev() bool {
	return i.iter.Prev()
}

func (i *iterWrapper) Error() error {
	return i.iter.Error()
}

func (i *iterWrapper) Release() {
	i.iter.Release()
}

type trtlIterator struct {
	client        trtlpb.TrtlClient
	values        []*trtlpb.KVPair
	index         int
	nextPageToken string
	cursor        trtlpb.Trtl_CursorClient
	cancel        context.CancelFunc
	snapshot      bool
	namespace     string
	err           error
}

func NewTrtlIterator(client trtlpb.TrtlClient, snapshot bool, namespace string) (iter *trtlIterator) {
	iter = &trtlIterator{
		client:    client,
		snapshot:  snapshot,
		namespace: namespace,
	}

	var ctx context.Context
	ctx, iter.cancel = getContext()

	if snapshot {
		request := &trtlpb.CursorRequest{
			Namespace: namespace,
		}
		iter.cursor, iter.err = iter.client.Cursor(ctx, request)
		iter.values = make([]*trtlpb.KVPair, 0)
	} else {
		//defer iter.cancel()
		request := &trtlpb.IterRequest{
			Namespace: namespace,
		}

		var reply *trtlpb.IterReply
		var err error
		if reply, err = iter.client.Iter(ctx, request); err != nil {
			iter.err = err
			return
		}
		iter.values = reply.Values
		iter.nextPageToken = reply.NextPageToken
	}
	iter.index = -1

	return iter
}

func (i *trtlIterator) nextCursor() bool {
	if i.cursor == nil {
		i.err = errors.New("nil cursor encountered on Next()")
		return false
	}

	i.index++
	if i.index >= len(i.values) {
		// Fetch a new value from the cursor
		var value *trtlpb.KVPair
		var err error
		if value, err = i.cursor.Recv(); err != nil {
			if err != io.EOF {
				i.err = err
			}
			i.cancel()
			return false
		}
		i.values = append(i.values, value)
	}

	return true
}

func (i *trtlIterator) nextIter() bool {
	i.index++
	if i.index >= len(i.values) && i.nextPageToken != "" {
		ctx, cancel := getContext()
		defer cancel()

		// Retrieve the next page from the server
		request := &trtlpb.IterRequest{
			Namespace: i.namespace,
			Options: &trtlpb.Options{
				PageToken: i.nextPageToken,
			},
		}
		var reply *trtlpb.IterReply
		var err error
		if reply, err = i.client.Iter(ctx, request); err != nil {
			i.err = err
			return false
		}

		// Add the new values to the slice
		i.values = append(i.values, reply.Values...)
		i.nextPageToken = reply.NextPageToken
	} else if i.index >= len(i.values) {
		// No more values from the iterator
		i.index = len(i.values) - 1
		return false
	}

	return true
}

func (i *trtlIterator) Next() bool {
	if i.err != nil {
		return false
	}
	if i.snapshot {
		return i.nextCursor()
	}
	return i.nextIter()
}

func (i *trtlIterator) Prev() bool {
	i.index--
	if i.index < 0 {
		i.index = 0
		return false
	}
	return true
}

func (i *trtlIterator) Seek(key []byte) bool {
	for current := i.Value(); bytes.Compare(current, key) < 0; current = i.Value() {
		if !i.Next() {
			return false
		}
	}
	return true
}

func (i *trtlIterator) Key() []byte {
	return i.values[i.index].Key
}

func (i *trtlIterator) Value() []byte {
	return i.values[i.index].Value
}

func (i *trtlIterator) Error() error {
	return i.err
}

func (i *trtlIterator) Release() {
	i.cancel()
}

func (i *vaspIterator) VASP() (*pb.VASP, error) {
	vasp := new(pb.VASP)
	if err := proto.Unmarshal(i.iter.Value(), vasp); err != nil {
		log.Error().Err(err).Str("type", wire.NamespaceVASPs).Str("key", string(i.iter.Key())).Msg("corrupted data encountered")
		return nil, err
	}
	return vasp, nil
}

func (i *vaspIterator) All() (vasps []*pb.VASP, err error) {
	vasps = make([]*pb.VASP, 0)
	defer i.iter.Release()

	for i.iter.Next() {
		vasp := new(pb.VASP)
		if err = proto.Unmarshal(i.iter.Value(), vasp); err != nil {
			return nil, err
		}
		vasps = append(vasps, vasp)
	}

	if err = i.iter.Error(); err != nil {
		return nil, err
	}
	return vasps, nil
}

func (i *vaspIterator) Id() string {
	// The VASP ID is prefix + uuid so strip off the prefix and return the string
	return string(i.iter.Key())
}

func (i *vaspIterator) Seek(vaspID string) bool {
	return i.iter.Seek([]byte(vaspID))
}

func (i *certReqIterator) CertReq() (*models.CertificateRequest, error) {
	r := new(models.CertificateRequest)
	if err := proto.Unmarshal(i.iter.Value(), r); err != nil {
		log.Error().Err(err).Str("type", wire.NamespaceCertReqs).Str("key", string(i.iter.Key())).Msg("corrupted data encountered")
		return nil, err
	}
	return r, nil
}

func (i *certReqIterator) All() (reqs []*models.CertificateRequest, err error) {
	reqs = make([]*models.CertificateRequest, 0)
	defer i.iter.Release()
	for i.iter.Next() {
		r := new(models.CertificateRequest)
		if err = proto.Unmarshal(i.iter.Value(), r); err != nil {
			return nil, err
		}
		reqs = append(reqs, r)
	}

	if err = i.iter.Error(); err != nil {
		return nil, err
	}

	return reqs, nil
}
