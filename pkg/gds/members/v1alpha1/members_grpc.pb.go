// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: gds/members/v1alpha1/members.proto

package members

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

// TRISAMembersClient is the client API for TRISAMembers service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TRISAMembersClient interface {
	// List all verified VASP members in the Directory Service.
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListReply, error)
	// Get a short summary of the verified VASP members in the Directory Service.
	Summary(ctx context.Context, in *SummaryRequest, opts ...grpc.CallOption) (*SummaryReply, error)
	// Get details for a VASP member in the Directory Service.
	Details(ctx context.Context, in *DetailsRequest, opts ...grpc.CallOption) (*MemberDetails, error)
}

type tRISAMembersClient struct {
	cc grpc.ClientConnInterface
}

func NewTRISAMembersClient(cc grpc.ClientConnInterface) TRISAMembersClient {
	return &tRISAMembersClient{cc}
}

func (c *tRISAMembersClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListReply, error) {
	out := new(ListReply)
	err := c.cc.Invoke(ctx, "/gds.members.v1alpha1.TRISAMembers/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tRISAMembersClient) Summary(ctx context.Context, in *SummaryRequest, opts ...grpc.CallOption) (*SummaryReply, error) {
	out := new(SummaryReply)
	err := c.cc.Invoke(ctx, "/gds.members.v1alpha1.TRISAMembers/Summary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tRISAMembersClient) Details(ctx context.Context, in *DetailsRequest, opts ...grpc.CallOption) (*MemberDetails, error) {
	out := new(MemberDetails)
	err := c.cc.Invoke(ctx, "/gds.members.v1alpha1.TRISAMembers/Details", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TRISAMembersServer is the server API for TRISAMembers service.
// All implementations must embed UnimplementedTRISAMembersServer
// for forward compatibility
type TRISAMembersServer interface {
	// List all verified VASP members in the Directory Service.
	List(context.Context, *ListRequest) (*ListReply, error)
	// Get a short summary of the verified VASP members in the Directory Service.
	Summary(context.Context, *SummaryRequest) (*SummaryReply, error)
	// Get details for a VASP member in the Directory Service.
	Details(context.Context, *DetailsRequest) (*MemberDetails, error)
	mustEmbedUnimplementedTRISAMembersServer()
}

// UnimplementedTRISAMembersServer must be embedded to have forward compatible implementations.
type UnimplementedTRISAMembersServer struct {
}

func (UnimplementedTRISAMembersServer) List(context.Context, *ListRequest) (*ListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedTRISAMembersServer) Summary(context.Context, *SummaryRequest) (*SummaryReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Summary not implemented")
}
func (UnimplementedTRISAMembersServer) Details(context.Context, *DetailsRequest) (*MemberDetails, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Details not implemented")
}
func (UnimplementedTRISAMembersServer) mustEmbedUnimplementedTRISAMembersServer() {}

// UnsafeTRISAMembersServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TRISAMembersServer will
// result in compilation errors.
type UnsafeTRISAMembersServer interface {
	mustEmbedUnimplementedTRISAMembersServer()
}

func RegisterTRISAMembersServer(s grpc.ServiceRegistrar, srv TRISAMembersServer) {
	s.RegisterService(&TRISAMembers_ServiceDesc, srv)
}

func _TRISAMembers_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TRISAMembersServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gds.members.v1alpha1.TRISAMembers/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TRISAMembersServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TRISAMembers_Summary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SummaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TRISAMembersServer).Summary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gds.members.v1alpha1.TRISAMembers/Summary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TRISAMembersServer).Summary(ctx, req.(*SummaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TRISAMembers_Details_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetailsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TRISAMembersServer).Details(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gds.members.v1alpha1.TRISAMembers/Details",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TRISAMembersServer).Details(ctx, req.(*DetailsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TRISAMembers_ServiceDesc is the grpc.ServiceDesc for TRISAMembers service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TRISAMembers_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gds.members.v1alpha1.TRISAMembers",
	HandlerType: (*TRISAMembersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _TRISAMembers_List_Handler,
		},
		{
			MethodName: "Summary",
			Handler:    _TRISAMembers_Summary_Handler,
		},
		{
			MethodName: "Details",
			Handler:    _TRISAMembers_Details_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gds/members/v1alpha1/members.proto",
}
