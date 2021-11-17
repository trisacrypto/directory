package trtl

import (
	"bytes"
	"context"
	"io"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/object"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// A HonuService implements the RPCs for interacting with a Honu database.
type HonuService struct {
	pb.UnimplementedTrtlServer
	parent *Server
	db     *honu.DB
}

func NewHonuService(s *Server) (*HonuService, error) {
	return &HonuService{parent: s, db: s.db}, nil
}

// Get is a unary request to retrieve a value for a key.
func (h *HonuService) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	var err error

	if _, found := reservedNamespaces[in.Namespace]; found {
		log.Warn().Msg("cannot use reserved namespace")
		return nil, status.Error(codes.PermissionDenied, "cannot use reserved namespace")
	}
	if len(in.Key) == 0 {
		log.Warn().Msg("missing key in Trtl Get request")
		return nil, status.Error(codes.InvalidArgument, "key must be provided in Get request")
	}

	var key []byte
	if len(in.Namespace) > 0 {
		key = prepend(in.Namespace, in.Key)
	} else {
		key = prepend("default", in.Key)
	}

	if in.Options != nil {
		log.Debug().Str("key", string(key)).Bool("return_meta", in.Options.ReturnMeta).Msg("Trtl Get")
	} else {
		log.Debug().Str("key", string(key)).Msg("Trtl Get")
	}
	if in.Options != nil && in.Options.ReturnMeta {
		// Retrieve and return the metadata.
		var object *object.Object
		if object, err = h.db.Object(key); err != nil {
			// TODO: Check for the honu not found error instead.
			if err == leveldb.ErrNotFound {
				log.Debug().Err(err).Str("key", string(key)).Msg("specified key not found")
				return nil, status.Error(codes.NotFound, err.Error())
			}
			log.Error().Err(err).Str("key", string(key)).Msg("unable to retrieve object")
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &pb.GetReply{
			Value: object.Data,
			Meta: &pb.Meta{
				Key:       object.Key,
				Namespace: object.Namespace,
				Region:    object.Region,
				Owner:     object.Owner,
				Version: &pb.Version{
					Pid:     object.Version.Pid,
					Version: object.Version.Version,
					Region:  object.Version.Region,
				},
				Parent: &pb.Version{
					Pid:     object.Version.Parent.Pid,
					Version: object.Version.Parent.Version,
					Region:  object.Version.Parent.Region,
				},
			},
		}, nil
	}

	// Just return the value for the given key.
	var value []byte
	if value, err = h.db.Get(key); err != nil {
		// TODO: Check for the honu not found error instead.
		if err == leveldb.ErrNotFound {
			log.Debug().Err(err).Str("key", string(key)).Msg("specified key not found")
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Error().Err(err).Str("key", string(key)).Msg("unable to retrieve value")
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.GetReply{
		Value: value,
	}, nil
}

// Put is a unary request to store a value for a key.
func (h *HonuService) Put(ctx context.Context, in *pb.PutRequest) (out *pb.PutReply, err error) {
	if _, found := reservedNamespaces[in.Namespace]; found {
		log.Warn().Msg("cannot use reserved namespace")
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

	var key []byte
	if len(in.Namespace) > 0 {
		key = prepend(in.Namespace, in.Key)
	} else {
		key = prepend("default", in.Key)
	}

	if in.Options != nil {
		log.Debug().Str("key", string(key)).Bool("return_meta", in.Options.ReturnMeta).Msg("Trtl Put")
	} else {
		log.Debug().Str("key", string(key)).Msg("Trtl Put")
	}

	var success bool
	if err := h.db.Put(key, in.Value); err != nil {
		success = false
	} else {
		success = true
	}

	out = &pb.PutReply{Success: success}

	// if Options include a request for metadata, we need to do a Get
	if in.Options != nil && in.Options.ReturnMeta {
		opts := &pb.Options{ReturnMeta: true}
		req := &pb.GetRequest{Key: in.Key, Namespace: in.Namespace, Options: opts}
		get, err := h.Get(ctx, req)
		if err != nil {
			// We failed to get the metadata, but still return the PutReply
			return out, err
		}

		out.Meta = &pb.Meta{
			Key:       get.Meta.Key,
			Namespace: get.Meta.Namespace,
			Region:    get.Meta.Region,
			Owner:     get.Meta.Owner,
			Version: &pb.Version{
				Pid:     get.Meta.Version.Pid,
				Version: get.Meta.Version.Version,
				Region:  get.Meta.Version.Region,
			},
			Parent: &pb.Version{
				Pid:     get.Meta.Parent.Pid,
				Version: get.Meta.Parent.Version,
				Region:  get.Meta.Parent.Region,
			},
		}
	}
	return out, nil
}

func (h *HonuService) Delete(ctx context.Context, in *pb.DeleteRequest) (out *pb.DeleteReply, err error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (h *HonuService) Iter(ctx context.Context, in *pb.IterRequest) (out *pb.IterReply, err error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// Batch is a client-side streaming request to issue multiple commands, usually Put and Delete.
func (h *HonuService) Batch(stream pb.Trtl_BatchServer) error {
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

func (h *HonuService) Cursor(in *pb.CursorRequest, stream pb.Trtl_CursorServer) (err error) {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (h *HonuService) Sync(stream pb.Trtl_SyncServer) (err error) {
	return status.Error(codes.Unimplemented, "not implemented")
}

// prepend the namespace to the key
func prepend(namespace string, key []byte) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(namespace),
			key,
		}, []byte("::"),
	)
}
