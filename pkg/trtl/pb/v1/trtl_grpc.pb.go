// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TrtlClient is the client API for Trtl service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TrtlClient interface {
	// Get is a unary request to retrieve a value for a key.
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
	// Put is a unary request to store a value for a key.
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutReply, error)
	// Delete is a unary request to remove a value and key.
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error)
	// Iter is a unary request that returns a completely materialized list of key value pairs.
	Iter(ctx context.Context, in *IterRequest, opts ...grpc.CallOption) (*IterReply, error)
	// Batch is a client-side streaming request to issue multiple commands, usually Put and Delete.
	Batch(ctx context.Context, opts ...grpc.CallOption) (Trtl_BatchClient, error)
	// Cursor is a server-side streaming request to iterate in a memory safe fashion.
	Cursor(ctx context.Context, in *CursorRequest, opts ...grpc.CallOption) (Trtl_CursorClient, error)
	// Sync is a bi-directional streaming mechanism to issue access requests synchronously.
	Sync(ctx context.Context, opts ...grpc.CallOption) (Trtl_SyncClient, error)
}

type trtlClient struct {
	cc grpc.ClientConnInterface
}

func NewTrtlClient(cc grpc.ClientConnInterface) TrtlClient {
	return &trtlClient{cc}
}

func (c *trtlClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := c.cc.Invoke(ctx, "/trtl.v1.Trtl/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trtlClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutReply, error) {
	out := new(PutReply)
	err := c.cc.Invoke(ctx, "/trtl.v1.Trtl/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trtlClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error) {
	out := new(DeleteReply)
	err := c.cc.Invoke(ctx, "/trtl.v1.Trtl/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trtlClient) Iter(ctx context.Context, in *IterRequest, opts ...grpc.CallOption) (*IterReply, error) {
	out := new(IterReply)
	err := c.cc.Invoke(ctx, "/trtl.v1.Trtl/Iter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trtlClient) Batch(ctx context.Context, opts ...grpc.CallOption) (Trtl_BatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Trtl_ServiceDesc.Streams[0], "/trtl.v1.Trtl/Batch", opts...)
	if err != nil {
		return nil, err
	}
	x := &trtlBatchClient{stream}
	return x, nil
}

type Trtl_BatchClient interface {
	Send(*BatchRequest) error
	CloseAndRecv() (*BatchReply, error)
	grpc.ClientStream
}

type trtlBatchClient struct {
	grpc.ClientStream
}

func (x *trtlBatchClient) Send(m *BatchRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *trtlBatchClient) CloseAndRecv() (*BatchReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(BatchReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *trtlClient) Cursor(ctx context.Context, in *CursorRequest, opts ...grpc.CallOption) (Trtl_CursorClient, error) {
	stream, err := c.cc.NewStream(ctx, &Trtl_ServiceDesc.Streams[1], "/trtl.v1.Trtl/Cursor", opts...)
	if err != nil {
		return nil, err
	}
	x := &trtlCursorClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Trtl_CursorClient interface {
	Recv() (*KVPair, error)
	grpc.ClientStream
}

type trtlCursorClient struct {
	grpc.ClientStream
}

func (x *trtlCursorClient) Recv() (*KVPair, error) {
	m := new(KVPair)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *trtlClient) Sync(ctx context.Context, opts ...grpc.CallOption) (Trtl_SyncClient, error) {
	stream, err := c.cc.NewStream(ctx, &Trtl_ServiceDesc.Streams[2], "/trtl.v1.Trtl/Sync", opts...)
	if err != nil {
		return nil, err
	}
	x := &trtlSyncClient{stream}
	return x, nil
}

type Trtl_SyncClient interface {
	Send(*SyncRequest) error
	Recv() (*SyncReply, error)
	grpc.ClientStream
}

type trtlSyncClient struct {
	grpc.ClientStream
}

func (x *trtlSyncClient) Send(m *SyncRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *trtlSyncClient) Recv() (*SyncReply, error) {
	m := new(SyncReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TrtlServer is the server API for Trtl service.
// All implementations must embed UnimplementedTrtlServer
// for forward compatibility
type TrtlServer interface {
	// Get is a unary request to retrieve a value for a key.
	Get(context.Context, *GetRequest) (*GetReply, error)
	// Put is a unary request to store a value for a key.
	Put(context.Context, *PutRequest) (*PutReply, error)
	// Delete is a unary request to remove a value and key.
	Delete(context.Context, *DeleteRequest) (*DeleteReply, error)
	// Iter is a unary request that returns a completely materialized list of key value pairs.
	Iter(context.Context, *IterRequest) (*IterReply, error)
	// Batch is a client-side streaming request to issue multiple commands, usually Put and Delete.
	Batch(Trtl_BatchServer) error
	// Cursor is a server-side streaming request to iterate in a memory safe fashion.
	Cursor(*CursorRequest, Trtl_CursorServer) error
	// Sync is a bi-directional streaming mechanism to issue access requests synchronously.
	Sync(Trtl_SyncServer) error
	mustEmbedUnimplementedTrtlServer()
}

// UnimplementedTrtlServer must be embedded to have forward compatible implementations.
type UnimplementedTrtlServer struct {
}

func (UnimplementedTrtlServer) Get(context.Context, *GetRequest) (*GetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedTrtlServer) Put(context.Context, *PutRequest) (*PutReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedTrtlServer) Delete(context.Context, *DeleteRequest) (*DeleteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedTrtlServer) Iter(context.Context, *IterRequest) (*IterReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Iter not implemented")
}
func (UnimplementedTrtlServer) Batch(Trtl_BatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Batch not implemented")
}
func (UnimplementedTrtlServer) Cursor(*CursorRequest, Trtl_CursorServer) error {
	return status.Errorf(codes.Unimplemented, "method Cursor not implemented")
}
func (UnimplementedTrtlServer) Sync(Trtl_SyncServer) error {
	return status.Errorf(codes.Unimplemented, "method Sync not implemented")
}
func (UnimplementedTrtlServer) mustEmbedUnimplementedTrtlServer() {}

// UnsafeTrtlServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TrtlServer will
// result in compilation errors.
type UnsafeTrtlServer interface {
	mustEmbedUnimplementedTrtlServer()
}

func RegisterTrtlServer(s grpc.ServiceRegistrar, srv TrtlServer) {
	s.RegisterService(&Trtl_ServiceDesc, srv)
}

func _Trtl_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrtlServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trtl.v1.Trtl/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrtlServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Trtl_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrtlServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trtl.v1.Trtl/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrtlServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Trtl_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrtlServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trtl.v1.Trtl/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrtlServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Trtl_Iter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrtlServer).Iter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trtl.v1.Trtl/Iter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrtlServer).Iter(ctx, req.(*IterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Trtl_Batch_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TrtlServer).Batch(&trtlBatchServer{stream})
}

type Trtl_BatchServer interface {
	SendAndClose(*BatchReply) error
	Recv() (*BatchRequest, error)
	grpc.ServerStream
}

type trtlBatchServer struct {
	grpc.ServerStream
}

func (x *trtlBatchServer) SendAndClose(m *BatchReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *trtlBatchServer) Recv() (*BatchRequest, error) {
	m := new(BatchRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Trtl_Cursor_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CursorRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TrtlServer).Cursor(m, &trtlCursorServer{stream})
}

type Trtl_CursorServer interface {
	Send(*KVPair) error
	grpc.ServerStream
}

type trtlCursorServer struct {
	grpc.ServerStream
}

func (x *trtlCursorServer) Send(m *KVPair) error {
	return x.ServerStream.SendMsg(m)
}

func _Trtl_Sync_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TrtlServer).Sync(&trtlSyncServer{stream})
}

type Trtl_SyncServer interface {
	Send(*SyncReply) error
	Recv() (*SyncRequest, error)
	grpc.ServerStream
}

type trtlSyncServer struct {
	grpc.ServerStream
}

func (x *trtlSyncServer) Send(m *SyncReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *trtlSyncServer) Recv() (*SyncRequest, error) {
	m := new(SyncRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Trtl_ServiceDesc is the grpc.ServiceDesc for Trtl service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Trtl_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "trtl.v1.Trtl",
	HandlerType: (*TrtlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Trtl_Get_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _Trtl_Put_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Trtl_Delete_Handler,
		},
		{
			MethodName: "Iter",
			Handler:    _Trtl_Iter_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Batch",
			Handler:       _Trtl_Batch_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Cursor",
			Handler:       _Trtl_Cursor_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Sync",
			Handler:       _Trtl_Sync_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "trtl/v1/trtl.proto",
}
