package trtl

import (
	"context"
	"errors"
	"io"

	"github.com/google/uuid"
	bff "github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/trtl"
	trtlpb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

type vaspIterator struct {
	trtlIterator
}

type certIterator struct {
	trtlIterator
}

type certReqIterator struct {
	trtlIterator
}

type organizationIterator struct {
	trtlIterator
}

type emailIterator struct {
	trtlIterator
}

// trtlIterator is an interface that is implemented by both the trtlBatchIterator and
// trtlStreamingIterator to iterate over values in the trtl store. The general workflow
// is to instantiate the iterator with either NewTrtlBatchIterator or
// NewTrtlStreamingIterator, call Next() to move the iterator to the next value, and if
// Next() returns True, call Key() and/or Value() to retrieve the current key and/or
// value. If Next() returns false, the iterator is already pointing to the last value.
// The caller should check Error() to verify that no errors occurred during iteration
// and Release() to clean up networking resources.
type trtlIterator interface {
	Next() bool
	Prev() bool
	Seek(key []byte) bool
	Key() []byte
	Value() []byte
	Error() error
	Release()
}

// trtlBatchIterator implements a batch iterator for the trtl store. Internally, the
// key-value pairs are loaded from the server in paginated batches. This provides the
// caller with materialized views of the data, and is best suited for operations such
// as reading all of the objects in a namespace at once.
type trtlBatchIterator struct {
	client        trtlpb.TrtlClient
	values        []*trtlpb.KVPair
	index         int
	nextPageToken string
	iterCalled    bool
	namespace     string
	err           error
}

func NewTrtlBatchIterator(client trtlpb.TrtlClient, namespace string) *trtlBatchIterator {
	return &trtlBatchIterator{
		client:    client,
		namespace: namespace,
		index:     -1,
		values:    make([]*trtlpb.KVPair, 0),
	}
}

func (i *trtlBatchIterator) Next() bool {
	if i.err != nil {
		return false
	}

	i.index++

	// If our next index will not refer to a value from the current page and there is
	// no next page, then we're done iterating. Note that we must also check the
	// starting case to fetch the first page.
	if i.iterCalled && i.index >= len(i.values) && i.nextPageToken == "" {
		// No more values from the iterator
		return false
	}

	if i.index >= 0 && i.index < len(i.values) {
		// If the index is in the range of the current set of values, we do
		// not have to fetch the next page since the next value is already loaded.
		return true
	}

	request := &trtlpb.IterRequest{
		Namespace: i.namespace,
		Options: &trtlpb.Options{
			PageToken: i.nextPageToken,
		},
	}

	var reply *trtlpb.IterReply
	var err error
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	if reply, err = i.client.Iter(ctx, request); err != nil {
		i.err = err
		return false
	}

	// We can't rely on checking the length of the values in the batch to determine
	// whether or not to fetch new values because it's possible for the Iter() call to
	// return no values. So, we set a flag here to remind us that we've already called
	// Iter() and there's no point calling it again if nextPageToken is empty.
	i.iterCalled = true

	// Add the new values to the slice
	i.values = append(i.values, reply.Values...)
	i.nextPageToken = reply.NextPageToken

	// If no new values are returned, then we're done iterating
	return i.index < len(i.values)
}

func (i *trtlBatchIterator) Prev() bool {
	i.index--
	if i.index < 0 {
		i.index = -1
		return false
	}
	return i.index < len(i.values)
}

func (i *trtlBatchIterator) Seek(key []byte) bool {
	var err error

	if len(i.values) > 0 {
		i.err = errors.New("cannot call Seek() on a batch iterator after Next() has been called")
		return false
	}

	// Using a page size of 100 to match the default page size in the members list RPC
	if i.nextPageToken, err = trtl.SeekCursor(100, key, i.namespace); err != nil {
		i.err = err
		return false
	}

	return i.Next()
}

func (i *trtlBatchIterator) Key() []byte {
	if i.index < 0 || i.index > len(i.values)-1 {
		return nil
	}
	return i.values[i.index].Key
}

func (i *trtlBatchIterator) Value() []byte {
	if i.index < 0 || i.index > len(i.values)-1 {
		return nil
	}
	return i.values[i.index].Value
}

func (i *trtlBatchIterator) Error() error {
	return i.err
}

func (i *trtlBatchIterator) Release() {
	// Cleanup values stored in memory and mark for garbage collection
	i.values = nil
}

// trtlStreamingIterator implements a streaming iterator for the trtl store.
// The iterator fetches the next value from the server when the caller calls Next(),
// and only the previous and current values are stored in memory. This is best suited
// for use cases where snapshot isolation or one-at-a-time processing is required.
type trtlStreamingIterator struct {
	client    trtlpb.TrtlClient
	cursor    trtlpb.Trtl_CursorClient
	cancel    context.CancelFunc
	prev      *trtlpb.KVPair
	current   *trtlpb.KVPair
	next      *trtlpb.KVPair
	eof       bool
	namespace string
	err       error
}

func NewTrtlStreamingIterator(client trtlpb.TrtlClient, namespace string) *trtlStreamingIterator {
	return &trtlStreamingIterator{
		client:    client,
		namespace: namespace,
	}
}

func (i *trtlStreamingIterator) Next() bool {
	if i.cursor == nil {
		var ctx context.Context
		ctx, i.cancel = utils.WithDeadline(context.Background())
		request := &trtlpb.CursorRequest{
			Namespace: i.namespace,
		}
		i.cursor, i.err = i.client.Cursor(ctx, request)
	}

	if i.err != nil {
		return false
	}

	if i.next != nil {
		// We have already loaded the next value
		i.prev = i.current
		i.current = i.next
		i.next = nil
		return true
	}

	if i.eof {
		return false
	}

	// Fetch the next value from the cursor
	var val *trtlpb.KVPair
	if val, i.err = i.cursor.Recv(); i.err != nil {
		if i.err == io.EOF {
			i.err = nil
			i.eof = true
		}
		i.prev = i.current
		i.current = nil
		return false
	}

	// Enable one Prev() call
	i.prev = i.current
	i.current = val
	return true
}

// The streaming iterator only stores the current and previous values, so successive
// calls to Prev() are only valid if there is at least one Next() call in between them.
func (i *trtlStreamingIterator) Prev() bool {
	i.next = i.current
	i.current = i.prev
	i.prev = nil

	return i.current != nil
}

// Seek() can only be called once before Next() and sets the initial position of the
// iterator to the specified key if it exists.
func (i *trtlStreamingIterator) Seek(key []byte) bool {
	if i.cursor != nil {
		i.err = errors.New("cursor already initialized, cannot call seek")
		return false
	}

	var ctx context.Context
	ctx, i.cancel = utils.WithDeadline(context.Background())
	request := &trtlpb.CursorRequest{
		Namespace: i.namespace,
		SeekKey:   key,
	}
	i.cursor, i.err = i.client.Cursor(ctx, request)

	return i.Next()
}

func (i *trtlStreamingIterator) Key() []byte {
	if i.current == nil {
		return nil
	}
	return i.current.Key
}

func (i *trtlStreamingIterator) Value() []byte {
	if i.current == nil {
		return nil
	}
	return i.current.Value
}

func (i *trtlStreamingIterator) Error() error {
	return i.err
}

func (i *trtlStreamingIterator) Release() {
	// If the cursor is nil it's likely that the connection erred and i.err is not nil.
	// Release() should be called safely before Error() to ensure any resources are
	// cleaned up (e.g. canceling the context).
	if i.cursor != nil {
		i.cursor.CloseSend()
	}
	i.cancel()
}

func (i *vaspIterator) VASP() (*pb.VASP, error) {
	vasp := new(pb.VASP)
	if err := proto.Unmarshal(i.Value(), vasp); err != nil {
		sentry.Error(nil).Err(err).Str("type", wire.NamespaceVASPs).Str("key", string(i.Key())).Msg("corrupted data encountered")
		return nil, err
	}
	return vasp, nil
}

func (i *vaspIterator) All() (vasps []*pb.VASP, err error) {
	vasps = make([]*pb.VASP, 0)
	defer i.Release()

	for i.Next() {
		vasp := new(pb.VASP)
		if err = proto.Unmarshal(i.Value(), vasp); err != nil {
			return nil, err
		}
		vasps = append(vasps, vasp)
	}

	if err = i.Error(); err != nil {
		return nil, err
	}
	return vasps, nil
}

func (i *vaspIterator) Id() string {
	return string(i.Key())
}

func (i *vaspIterator) SeekId(vaspID string) bool {
	return i.Seek([]byte(vaspID))
}

func (i *certIterator) Cert() (*models.Certificate, error) {
	c := new(models.Certificate)
	if err := proto.Unmarshal(i.Value(), c); err != nil {
		sentry.Error(nil).Err(err).Str("type", wire.NamespaceCerts).Str("key", string(i.Key())).Msg("corrupted data encountered")
		return nil, err
	}
	return c, nil
}

func (i *certIterator) All() (certs []*models.Certificate, err error) {
	certs = make([]*models.Certificate, 0)
	defer i.Release()
	for i.Next() {
		c := new(models.Certificate)
		if err = proto.Unmarshal(i.Value(), c); err != nil {
			return nil, err
		}
		certs = append(certs, c)
	}

	if err = i.Error(); err != nil {
		return nil, err
	}

	return certs, nil
}

func (i *certReqIterator) CertReq() (*models.CertificateRequest, error) {
	r := new(models.CertificateRequest)
	if err := proto.Unmarshal(i.Value(), r); err != nil {
		sentry.Error(nil).Err(err).Str("type", wire.NamespaceCertReqs).Str("key", string(i.Key())).Msg("corrupted data encountered")
		return nil, err
	}
	return r, nil
}

func (i *certReqIterator) All() (reqs []*models.CertificateRequest, err error) {
	reqs = make([]*models.CertificateRequest, 0)
	defer i.Release()
	for i.Next() {
		r := new(models.CertificateRequest)
		if err = proto.Unmarshal(i.Value(), r); err != nil {
			return nil, err
		}
		reqs = append(reqs, r)
	}

	if err = i.Error(); err != nil {
		return nil, err
	}

	return reqs, nil
}

func (i *organizationIterator) ID() string {
	orgID, err := uuid.FromBytes(i.Key())
	if err != nil {
		panic(err)
	}
	return orgID.String()
}

func (i *organizationIterator) Organization() (o *bff.Organization, err error) {
	o = new(bff.Organization)
	if err = proto.Unmarshal(i.Value(), o); err != nil {
		return nil, err
	}
	return o, nil
}

func (i *emailIterator) Email() (c *models.Email, err error) {
	c = &models.Email{}
	if err = proto.Unmarshal(i.Value(), c); err != nil {
		return nil, err
	}
	return c, nil
}
