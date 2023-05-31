// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: trtl/peers/v1/peers.proto

package peers

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

// PeerManagementClient is the client API for PeerManagement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PeerManagementClient interface {
	GetPeers(ctx context.Context, in *PeersFilter, opts ...grpc.CallOption) (*PeersList, error)
	AddPeers(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*PeersStatus, error)
	RmPeers(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*PeersStatus, error)
}

type peerManagementClient struct {
	cc grpc.ClientConnInterface
}

func NewPeerManagementClient(cc grpc.ClientConnInterface) PeerManagementClient {
	return &peerManagementClient{cc}
}

func (c *peerManagementClient) GetPeers(ctx context.Context, in *PeersFilter, opts ...grpc.CallOption) (*PeersList, error) {
	out := new(PeersList)
	err := c.cc.Invoke(ctx, "/trtl.peers.v1.PeerManagement/GetPeers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *peerManagementClient) AddPeers(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*PeersStatus, error) {
	out := new(PeersStatus)
	err := c.cc.Invoke(ctx, "/trtl.peers.v1.PeerManagement/AddPeers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *peerManagementClient) RmPeers(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*PeersStatus, error) {
	out := new(PeersStatus)
	err := c.cc.Invoke(ctx, "/trtl.peers.v1.PeerManagement/RmPeers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PeerManagementServer is the server API for PeerManagement service.
// All implementations must embed UnimplementedPeerManagementServer
// for forward compatibility
type PeerManagementServer interface {
	GetPeers(context.Context, *PeersFilter) (*PeersList, error)
	AddPeers(context.Context, *Peer) (*PeersStatus, error)
	RmPeers(context.Context, *Peer) (*PeersStatus, error)
	mustEmbedUnimplementedPeerManagementServer()
}

// UnimplementedPeerManagementServer must be embedded to have forward compatible implementations.
type UnimplementedPeerManagementServer struct {
}

func (UnimplementedPeerManagementServer) GetPeers(context.Context, *PeersFilter) (*PeersList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPeers not implemented")
}
func (UnimplementedPeerManagementServer) AddPeers(context.Context, *Peer) (*PeersStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPeers not implemented")
}
func (UnimplementedPeerManagementServer) RmPeers(context.Context, *Peer) (*PeersStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RmPeers not implemented")
}
func (UnimplementedPeerManagementServer) mustEmbedUnimplementedPeerManagementServer() {}

// UnsafePeerManagementServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PeerManagementServer will
// result in compilation errors.
type UnsafePeerManagementServer interface {
	mustEmbedUnimplementedPeerManagementServer()
}

func RegisterPeerManagementServer(s grpc.ServiceRegistrar, srv PeerManagementServer) {
	s.RegisterService(&PeerManagement_ServiceDesc, srv)
}

func _PeerManagement_GetPeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PeersFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerManagementServer).GetPeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trtl.peers.v1.PeerManagement/GetPeers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerManagementServer).GetPeers(ctx, req.(*PeersFilter))
	}
	return interceptor(ctx, in, info, handler)
}

func _PeerManagement_AddPeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Peer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerManagementServer).AddPeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trtl.peers.v1.PeerManagement/AddPeers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerManagementServer).AddPeers(ctx, req.(*Peer))
	}
	return interceptor(ctx, in, info, handler)
}

func _PeerManagement_RmPeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Peer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerManagementServer).RmPeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trtl.peers.v1.PeerManagement/RmPeers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerManagementServer).RmPeers(ctx, req.(*Peer))
	}
	return interceptor(ctx, in, info, handler)
}

// PeerManagement_ServiceDesc is the grpc.ServiceDesc for PeerManagement service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PeerManagement_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "trtl.peers.v1.PeerManagement",
	HandlerType: (*PeerManagementServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPeers",
			Handler:    _PeerManagement_GetPeers_Handler,
		},
		{
			MethodName: "AddPeers",
			Handler:    _PeerManagement_AddPeers_Handler,
		},
		{
			MethodName: "RmPeers",
			Handler:    _PeerManagement_RmPeers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "trtl/peers/v1/peers.proto",
}
