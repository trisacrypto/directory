package trtl

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"

	"github.com/rotationalio/honu"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/iterator"
	"github.com/rotationalio/honu/object"
	"github.com/rotationalio/honu/options"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/trtl/internal"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// A TrtlService implements the RPCs for interacting with a Honu database.
type TrtlService struct {
	pb.UnimplementedTrtlServer
	parent *Server
	db     *honu.DB
}

func NewTrtlService(s *Server) (*TrtlService, error) {
	return &TrtlService{parent: s, db: s.db}, nil
}

const (
	defaultPageSize = 100
)

// Get is a unary request to retrieve a value for a key.
// If metadata is requested in the GetRequest, the request will use honu.Object() to
// retrieve the entire object, including the metadata. If no metadata is requested, the
// request will use honu.Get() to get only the value.
// If a namespace is provided, the namespace is passed to the internal honu Options,
// to look in that namespace only.
func (h *TrtlService) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	var err error
	var ns string

	if in.Namespace == "" {
		ns = "default"
	} else {
		ns = in.Namespace
	}

	if _, found := reservedNamespaces[in.Namespace]; found {
		log.Warn().Str("namespace", in.Namespace).Msg("cannot use reserved namespace")
		return nil, status.Error(codes.PermissionDenied, "cannot use reserved namespace")
	}
	if len(in.Key) == 0 {
		log.Warn().Msg("missing key in Trtl Get request")
		return nil, status.Error(codes.InvalidArgument, "key must be provided in Get request")
	}

	if in.Options != nil && in.Options.ReturnMeta {
		// Retrieve and return the metadata (uses honu.Object())
		log.Debug().Str("key", string(in.Key)).Bool("return_meta", in.Options.ReturnMeta).Msg("Trtl Get")

		// Check if we have a namespace
		// NOTE: empty string in.Namespace will use default namespace after honu v0.2.4
		var object *object.Object
		if object, err = h.db.Object(in.Key, options.WithNamespace(in.Namespace)); err != nil {
			// TODO: Check for the honu not found error instead.
			if err == engine.ErrNotFound {
				log.Debug().Err(err).Str("key", string(in.Key)).Msg("specified key not found")
				return nil, status.Error(codes.NotFound, err.Error())
			}
			log.Error().Err(err).Str("key", string(in.Key)).Msg("unable to retrieve object")
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Update prometheus metrics with Get
		pmGets.WithLabelValues(ns).Inc()

		return &pb.GetReply{
			Value: object.Data,
			Meta:  returnMeta(object),
		}, nil
	}

	// No metadata requested; just return the value for the given key (uses honu.Get())
	log.Debug().Str("key", string(in.Key)).Msg("Trtl Get")

	// But we do have to check if we have a namespace
	// NOTE: empty string in.Namespace will use default namespace after honu v0.2.4
	var value []byte
	if value, err = h.db.Get(in.Key, options.WithNamespace(in.Namespace)); err != nil {
		// TODO: Check for the honu not found error instead.
		if err == leveldb.ErrNotFound {
			log.Debug().Err(err).Str("key", string(in.Key)).Msg("specified key not found")
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Error().Err(err).Str("key", string(in.Key)).Msg("unable to retrieve value")
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Increment prometheus Get count
	pmGets.WithLabelValues(ns).Inc()

	return &pb.GetReply{
		Value: value,
	}, nil
}

// Put is a unary request to store a value for a key.
// If a namespace is provided, the namespace is passed to the internal honu Options,
// to put the value to that namespace.
func (h *TrtlService) Put(ctx context.Context, in *pb.PutRequest) (out *pb.PutReply, err error) {
	var ns string
	if in.Namespace == "" {
		ns = "default"
	} else {
		ns = in.Namespace
	}

	if _, found := reservedNamespaces[in.Namespace]; found {
		log.Warn().Str("namespace", in.Namespace).Msg("cannot use reserved namespace")
		return nil, status.Error(codes.PermissionDenied, "cannot use reserved namespace")
	}
	if len(in.Key) == 0 {
		log.Warn().Msg("missing key in Trtl Put request")
		return nil, status.Error(codes.InvalidArgument, "key must be provided in Put request")
	}
	if len(in.Value) == 0 {
		log.Warn().Msg("missing value in Trtl Put request")
		return nil, status.Error(codes.InvalidArgument, "value must be provided in Put request")
	}

	if in.Options != nil {
		log.Debug().Bytes("key", in.Key).Bool("return_meta", in.Options.ReturnMeta).Msg("Trtl Put")
	} else {
		log.Debug().Bytes("key", in.Key).Msg("Trtl Put")
	}

	// Check if we have a namespace
	// NOTE: empty string in.Namespace will use default namespace after honu v0.2.4
	var object *object.Object
	if object, err = h.db.Put(in.Key, in.Value, options.WithNamespace(in.Namespace)); err != nil {
		log.Error().Err(err).Str("key", string(in.Key)).Msg("unable to put object")
		return nil, status.Error(codes.Internal, err.Error())
	}

	out = &pb.PutReply{Success: true}

	if in.Options != nil && in.Options.ReturnMeta {
		out.Meta = returnMeta(object)
	}

	// Increment prometheus Put counter
	pmPuts.WithLabelValues(ns).Inc()

	// TODO: prometheus
	// If in.Options.ReturnMeta is true, we will get metadata from honu
	// Unfortunately, we currently we don't get tombstone information, but if we did,
	// we could use it to decrement our `pmTombstone` counter here

	return out, nil
}

// Delete is a unary request to delete a key.
// If a namespace is provided, the namespace is passed to the internal honu Options,
// to delete the key from a specific namespace. Note that this does not delete tombstones.
func (h *TrtlService) Delete(ctx context.Context, in *pb.DeleteRequest) (out *pb.DeleteReply, err error) {
	var ns string

	if in.Namespace == "" {
		ns = "default"
	} else {
		ns = in.Namespace
	}

	if _, found := reservedNamespaces[in.Namespace]; found {
		log.Warn().Str("namespace", in.Namespace).Msg("cannot use reserved namespace")
		return nil, status.Error(codes.PermissionDenied, "cannot use reserved namespace")
	}
	if len(in.Key) == 0 {
		log.Warn().Msg("missing key in Trtl Delete request")
		return nil, status.Error(codes.InvalidArgument, "key must be provided in Delete request")
	}

	if in.Options != nil {
		log.Debug().Bytes("key", in.Key).Bool("return_meta", in.Options.ReturnMeta).Msg("Trtl Delete")
	} else {
		log.Debug().Bytes("key", in.Key).Msg("Trtl Delete")
	}

	// Check if we have a namespace
	// NOTE: empty string in.Namespace will use default namespace after honu v0.2.4
	var object *object.Object
	if object, err = h.db.Delete(in.Key, options.WithNamespace(in.Namespace)); err != nil {
		log.Error().Err(err).Str("key", string(in.Key)).Msg("unable to delete object")
		return nil, status.Error(codes.Internal, err.Error())
	}

	out = &pb.DeleteReply{Success: true}

	if in.Options != nil && in.Options.ReturnMeta {
		out.Meta = returnMeta(object)
	}

	// Increment Prometheus Delete counter
	pmDels.WithLabelValues(ns).Inc()

	// TODO: Increment Prometheus Tombstone counter
	// Unfortunately we can't decrement yet! (see note in `Put`)
	// pmTombstones.WithLabelValues(ns).Inc()

	return out, nil
}

// Iter is a unary request to fetch a materialized collection of key/value pairs based
// on a shared prefix. If no prefix is specified an entire namespace may be returned.
// This RPC supports pagination to ensure that replies do not get too large. The default
// page size is 100 items, though this can be modified in the options. The next page
// token in the result will contain the next page to request, or will be empty if there
// are no more results to be supplied.
//
// Note that there are no snapshot guarantees during iteration, meaning that if the
// underlying database changes between requests, these changes could be reflected during
// iteration. For snapshot isolation in iteration, use the Cursor RPC.
//
// There are several options that modulate the Iter response:
//   - return_meta: each key/value pair will contain the object metadata
//   - iter_no_keys: each key/value pair will not have a key associated with it
//   - iter_no_values: each key/value pair will not have a value associated with it
//   - page_token: the page of results that the user wishes to fetch
//   - page_size: the number of results to be returned in the request
func (h *TrtlService) Iter(ctx context.Context, in *pb.IterRequest) (out *pb.IterReply, err error) {
	var ns string

	if in.Namespace == "" {
		ns = "default"
	} else {
		ns = in.Namespace
	}

	// Ensure the namespace is not reserved
	if _, found := reservedNamespaces[in.Namespace]; found {
		log.Warn().Str("namespace", in.Namespace).Msg("cannot use reserved namespace")
		return nil, status.Error(codes.PermissionDenied, "cannot used reserved namespace")
	}

	// Load the options from the request
	var opts *pb.Options
	if in.Options != nil {
		opts = in.Options
	} else {
		// Create empty options
		opts = &pb.Options{}
	}

	// Load defaults into options
	if opts.PageSize == 0 {
		opts.PageSize = defaultPageSize
	}

	// Test valid options
	if opts.IterNoKeys && opts.IterNoValues && !opts.ReturnMeta {
		log.Debug().
			Str("namespace", in.Namespace).
			Bool("iter_no_keys", opts.IterNoKeys).
			Bool("iter_no_values", opts.IterNoValues).
			Bool("return_meta", opts.ReturnMeta).
			Msg("iter request would return no data")
		return nil, status.Error(codes.InvalidArgument, "cannot specify no keys, values, and no return meta: no data would be returned")
	} else {
		log.Debug().Msg("Trtl Iter")
	}

	// If a page cursor is provided load it, otherwise create the cursor for iteration
	cursor := &internal.PageCursor{}
	if opts.PageToken != "" {
		if err = cursor.Load(opts.PageToken); err != nil {
			log.Warn().Err(err).Msg("invalid page token on iter request")
			return nil, status.Error(codes.InvalidArgument, "invalid page token")
		}

		// Validate the request has not changed
		if cursor.PageSize != opts.PageSize {
			log.Debug().Int32("cursor", cursor.PageSize).Int32("opts", opts.PageSize).Msg("invalid iter request: mismatched page size")
			return nil, status.Error(codes.InvalidArgument, "page size cannot change between requests")
		}

		// Note - prefix check happens on next key, but namespace check must match Honu iterator
		if !bytes.HasPrefix(cursor.NextKey, in.Prefix) {
			log.Debug().Msg("invalid iter request: mismatched prefix")
			return nil, status.Error(codes.InvalidArgument, "prefix cannot change between requests")
		}

	} else {
		// Create a new cursor
		cursor.PageSize = opts.PageSize
	}

	// Create response
	out = &pb.IterReply{
		Values: make([]*pb.KVPair, 0, cursor.PageSize),
	}

	// Create the honu iterator to begin collecting data with the specified prefix.
	// NOTE: empty string in.Namespace will use default namespace after honu v0.2.4
	var iter iterator.Iterator
	if iter, err = h.db.Iter(in.Prefix, options.WithNamespace(in.Namespace)); err != nil {
		log.Error().Err(err).Str("namespace", in.Namespace).Msg("could not create honu iterator")
		return nil, status.Errorf(codes.FailedPrecondition, "could not create iterator: %s", err)
	}
	defer iter.Release()

	// Perform namespace check to ensure that the page cursor matches the iterator namespace
	// NOTE: this section must come after the iterator is created, though it would be preferable
	// if it was in the section where the page cursor is created. This is because we want to
	// check the namespace that the iterator is operating on rather than the one specified by
	// the user in the first request (e.g. because a default namespace may be used).
	if opts.PageToken != "" && cursor.Namespace != iter.Namespace() {
		log.Debug().Msg("invalid iter request: mismatched namespace")
		return nil, status.Error(codes.InvalidArgument, "namespace cannot change between requests")
	}

	// If necessary seek to the next key specified by the cursor.
	if len(cursor.NextKey) > 0 {
		// If iter.Seek returns false (e.g. seek did not find the specified key) then
		// iter.Next() should also return false, so it isn't necessary to check the return.
		// NOTE: next key must be set to nil after it's used for seeking so that the last
		// page doesn't retain the old key and loop forever.
		iter.Seek(cursor.NextKey)
		cursor.NextKey = nil

		// Because we're going to be calling Next, we need to back up one key to ensure
		// that we start on the right key in the for loop.
		iter.Prev()
	}

	for iter.Next() {
		// Check if we're done iterating (e.g. at the end of the page with a next page)
		if len(out.Values) == int(cursor.PageSize) {
			// The current key is the next key for the next page, stop iteration and
			// prepare the page cursor to be returned.
			cursor.NextKey = iter.Key()
			cursor.Namespace = iter.Namespace()
			break
		}

		// Otherwise append the current key value pair to the page.
		// Fetch the metadata since it will need to be loaded for the response anyway.
		var object *object.Object
		if object, err = iter.Object(); err != nil {
			log.Error().Err(err).Str("key", base64.RawURLEncoding.EncodeToString(iter.Key())).Msg("could not fetch object metadata")
			return nil, status.Error(codes.FailedPrecondition, "database is in invalid state")
		}

		// Create the key value pair
		pair := &pb.KVPair{}
		if !opts.IterNoKeys {
			pair.Key = object.Key
			pair.Namespace = object.Namespace
		}

		if !opts.IterNoValues {
			pair.Value = object.Data
		}

		if opts.ReturnMeta {
			pair.Meta = returnMeta(object)
		}

		out.Values = append(out.Values, pair)
	}

	if err = iter.Error(); err != nil {
		log.Error().Err(err).Str("namespace", in.Namespace).Msg("could not iterate")
		return nil, status.Errorf(codes.FailedPrecondition, "iteration failure: %s", err)
	}

	// Check if there is a next page cursor
	if len(cursor.NextKey) != 0 {
		if out.NextPageToken, err = cursor.Dump(); err != nil {
			log.Error().Err(err).Str("namespace", in.Namespace).Msg("could not serialize next page token")
			return nil, status.Error(codes.FailedPrecondition, "could not serialize next page token")
		}
	}

	// Increment Prometheus Iter counter
	pmIters.WithLabelValues(ns).Inc()

	// Request complete
	log.Info().
		Str("namespace", in.Namespace).
		Int("count", len(out.Values)).
		Bool("has_next_page", out.NextPageToken != "").
		Msg("iter request complete")
	return out, nil
}

// Batch is a client-side streaming request to issue multiple commands, usually Put and Delete.
func (h *TrtlService) Batch(stream pb.Trtl_BatchServer) error {
	log.Debug().Msg("Trtl Batch")
	out := &pb.BatchReply{}
	for {
		// Read the next request from the stream.
		in, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(out)
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		out.Operations++
		if in.Request == nil {
			out.Failed++
			out.Errors = append(out.Errors, &pb.BatchReply_Error{
				Id:    in.Id,
				Error: "missing request field",
			})
			continue
		}

		// Process the request.
		switch in.Request.(type) {
		case *pb.BatchRequest_Put:
			var reply *pb.PutReply
			if reply, err = h.Put(stream.Context(), in.GetPut()); err != nil || !reply.Success {
				out.Failed++
				var errMsg string
				if err != nil {
					errMsg = err.Error()
				} else {
					errMsg = "unexpected Put error"
				}
				out.Errors = append(out.Errors, &pb.BatchReply_Error{
					Id:    in.Id,
					Error: errMsg,
				})
				continue
			}
		case *pb.BatchRequest_Delete:
			var reply *pb.DeleteReply
			if reply, err = h.Delete(stream.Context(), in.GetDelete()); err != nil || !reply.Success {
				out.Failed++
				var errMsg string
				if err != nil {
					errMsg = err.Error()
				} else {
					errMsg = "unexpected Delete error"
				}
				out.Errors = append(out.Errors, &pb.BatchReply_Error{
					Id:    in.Id,
					Error: errMsg,
				})
				continue
			}
		default:
			out.Failed++
			out.Errors = append(out.Errors, &pb.BatchReply_Error{
				Id:    in.Id,
				Error: "unknown request type",
			})
			continue
		}

		// Each case continues on failure so if we get here, the request was successful.
		out.Successful++
	}
}

// Cursor is a server-side streaming request to fetch a collection of key/value pairs
// based on a shared prefix. If no prefix is specified an entire namespace may be
// returned. The pairs are streamed one value at a time so that the client can control
// iteration and materialization without overloading either the server or the network.
//
// Note that there is a snapshot guarantee during iteration, meaning that if the
// underlying database changes via a concurrent request, the cursor stream will not be
// effected. Use Cursor instead of Iter if you require snapshot isolation reads.
//
// There are several options that modulate the Cursor stream:
//   - return_meta: each key/value pair will contain the object metadata
//   - iter_no_keys: each key/value pair will not have a key associated with it
//   - iter_no_values: each key/value pair will not have a value associated with it
//   - page_token: the page of results that the user wishes to fetch
//   - page_size: the number of results to be returned in the request
func (h *TrtlService) Cursor(in *pb.CursorRequest, stream pb.Trtl_CursorServer) (err error) {
	// Fetch the stream context
	ctx := stream.Context()

	// Ensure the namespace is not reserved
	if _, found := reservedNamespaces[in.Namespace]; found {
		log.Warn().Str("namespace", in.Namespace).Msg("cannot use reserved namespace")
		return status.Error(codes.PermissionDenied, "cannot used reserved namespace")
	}

	// Load the options from the request
	var opts *pb.Options
	if in.Options != nil {
		opts = in.Options
	} else {
		// Create empty options
		opts = &pb.Options{}
	}

	// Test valid options
	if opts.IterNoKeys && opts.IterNoValues && !opts.ReturnMeta {
		log.Debug().
			Str("namespace", in.Namespace).
			Bool("iter_no_keys", opts.IterNoKeys).
			Bool("iter_no_values", opts.IterNoValues).
			Bool("return_meta", opts.ReturnMeta).
			Msg("cursor request would return no data")
		return status.Error(codes.InvalidArgument, "cannot specify no keys, values, and no return meta: no data would be returned")
	} else {
		log.Debug().Msg("Trtl Cursor")
	}

	// NOTE: empty string in.Namespace will use default namespace after honu v0.2.4
	var iter iterator.Iterator
	if iter, err = h.db.Iter(in.Prefix, options.WithNamespace(in.Namespace)); err != nil {
		log.Error().Err(err).Str("namespace", in.Namespace).Msg("could not create honu iterator")
		return status.Errorf(codes.FailedPrecondition, "could not create iterator: %s", err)
	}
	defer iter.Release()

	// If a seek key is provided, seek to that key before iteration
	// NOTE: that because we'll be calling iter.Next to start the loop, we need set the
	// iterator to the key previous to the seek key.
	if len(in.SeekKey) > 0 {
		iter.Seek(in.SeekKey)
		iter.Prev()
	}

	var nMessages uint64
	for iter.Next() {
		// Check if the client has closed the stream
		select {
		case <-ctx.Done():
			if err = ctx.Err(); err != nil && err != io.EOF {
				log.Error().Err(err).Msg("cursor canceled by client with error")
				return status.Errorf(codes.Canceled, "cursor canceled by client: %s", err)
			}
			log.Info().
				Str("namespace", in.Namespace).
				Uint64("count", nMessages).
				Msg("cursor request canceled by client")
			return nil
		default:
		}

		// Fetch the metadata since it will need to be loaded for the response anyway.
		var object *object.Object
		if object, err = iter.Object(); err != nil {
			log.Error().Err(err).Str("key", base64.RawURLEncoding.EncodeToString(iter.Key())).Msg("could not fetch object metadata")
			return status.Error(codes.FailedPrecondition, "database is in invalid state")
		}

		// Create the key value pair to send in the cursor stream
		// NOTE: cannot call iter.Next() here or the iterator will advance
		msg := &pb.KVPair{}
		if !opts.IterNoKeys {
			msg.Key = object.Key
			msg.Namespace = object.Namespace
		}

		if !opts.IterNoValues {
			msg.Value = object.Data
		}

		if opts.ReturnMeta {
			msg.Meta = returnMeta(object)
		}

		// Send the message on the stream
		if err = stream.Send(msg); err != nil {
			log.Error().Err(err).Msg("could not send cursor reply during iteration")
			return status.Errorf(codes.Aborted, "send error occurred: %s", err)
		}

		// Count the number of messages successfully sent
		nMessages++

		// TODO: Prometheus
		// If in.Prefix is nil, nMessages will be all the objects in in.Namespace so we
		// could use this opportunity to update our Prometheus counter if we can find
		// out how to call something like update on the counter rather than increment
	}

	if err = iter.Error(); err != nil {
		log.Error().Err(err).Str("namespace", in.Namespace).Msg("could not iterate")
		return status.Errorf(codes.FailedPrecondition, "iteration failure: %s", err)
	}

	// Cursor stream complete
	log.Info().
		Str("namespace", in.Namespace).
		Uint64("count", nMessages).
		Msg("cursor request complete")
	return nil
}

func (h *TrtlService) Sync(stream pb.Trtl_SyncServer) (err error) {
	return status.Error(codes.Unimplemented, "not implemented")
}

// returnMeta is a helper function for returning the metadata on an object
func returnMeta(object *object.Object) *pb.Meta {
	meta := &pb.Meta{
		Key:       object.Key,
		Namespace: object.Namespace,
		Region:    object.Region,
		Owner:     object.Owner,
		Version: &pb.Version{
			Pid:     object.Version.Pid,
			Version: object.Version.Version,
			Region:  object.Version.Region,
		},
	}

	// If it is the first version, the parent will be nil.
	if object.Version.Parent != nil {
		meta.Parent = &pb.Version{
			Pid:     object.Version.Parent.Pid,
			Version: object.Version.Parent.Version,
			Region:  object.Version.Parent.Region,
		}
	}
	return meta
}
