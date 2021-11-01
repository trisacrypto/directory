package trtl

import (
	"context"

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

	if len(in.Key) == 0 {
		log.Error().Msg("missing key in Trtl Get request")
		return nil, status.Error(codes.InvalidArgument, "key must be provided in Get request")
	}
	if in.Options == nil {
		log.Error().Msg("missing options in Trtl Get request")
		return nil, status.Error(codes.InvalidArgument, "options must be provided in Get request")
	}
	log.Debug().Str("key", string(in.Key)).Bool("return_meta", in.Options.ReturnMeta).Msg("Trtl Get")
	if in.Options.ReturnMeta {
		// Retrieve and return the metadata.
		var object *object.Object
		if object, err = h.db.Object(in.Key); err != nil {
			// TODO: Check for the honu not found error instead.
			if err == leveldb.ErrNotFound {
				log.Debug().Err(err).Str("key", string(in.Key)).Msg("specified key not found")
				return nil, status.Error(codes.NotFound, err.Error())
			}
			log.Error().Err(err).Str("key", string(in.Key)).Msg("unable to retrieve object")
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
	if value, err = h.db.Get(in.Key); err != nil {
		// TODO: Check for the honu not found error instead.
		if err == leveldb.ErrNotFound {
			log.Debug().Err(err).Str("key", string(in.Key)).Msg("specified key not found")
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Error().Err(err).Str("key", string(in.Key)).Msg("unable to retrieve value")
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.GetReply{
		Value: value,
	}, nil
}

// Put is a unary request to store a value for a key.
func (h *HonuService) Put(ctx context.Context, in *pb.PutRequest) (out *pb.PutReply, err error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (h *HonuService) Delete(ctx context.Context, in *pb.DeleteRequest) (out *pb.DeleteReply, err error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (h *HonuService) Iter(ctx context.Context, in *pb.IterRequest) (out *pb.IterReply, err error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (h *HonuService) Batch(stream pb.Trtl_BatchServer) (err error) {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (h *HonuService) Cursor(in *pb.CursorRequest, stream pb.Trtl_CursorServer) (err error) {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (h *HonuService) Sync(stream pb.Trtl_SyncServer) (err error) {
	return status.Error(codes.Unimplemented, "not implemented")
}
