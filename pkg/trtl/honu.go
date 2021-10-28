package trtl

import (
	"context"

	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// A HonuService implements the RPCs for interacting with a Honu database.
type HonuService struct {
	pb.UnimplementedTrtlServer
}

func NewHonuService() *HonuService {
	return &HonuService{}
}

// Get is a unary request to retrieve a value for a key.
func (h *HonuService) Get(ctx context.Context, in *pb.GetRequest) (out *pb.GetReply, err error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
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
