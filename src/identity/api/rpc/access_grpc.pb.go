// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package rpc

import (
	context "context"
	presenter "github.com/timoth-y/chainmetric-contracts/src/identity/api/presenter"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AccessServiceClient is the client API for AccessService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccessServiceClient interface {
	RequestFabricCredentials(ctx context.Context, in *presenter.FabricCredentialsRequest, opts ...grpc.CallOption) (*presenter.FabricCredentialsResponse, error)
	UpdatePassword(ctx context.Context, in *presenter.UpdatePasswordRequest, opts ...grpc.CallOption) (*presenter.StatusResponse, error)
}

type accessServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccessServiceClient(cc grpc.ClientConnInterface) AccessServiceClient {
	return &accessServiceClient{cc}
}

func (c *accessServiceClient) RequestFabricCredentials(ctx context.Context, in *presenter.FabricCredentialsRequest, opts ...grpc.CallOption) (*presenter.FabricCredentialsResponse, error) {
	out := new(presenter.FabricCredentialsResponse)
	err := c.cc.Invoke(ctx, "/chainmetric.identity.service.AccessService/requestFabricCredentials", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessServiceClient) UpdatePassword(ctx context.Context, in *presenter.UpdatePasswordRequest, opts ...grpc.CallOption) (*presenter.StatusResponse, error) {
	out := new(presenter.StatusResponse)
	err := c.cc.Invoke(ctx, "/chainmetric.identity.service.AccessService/updatePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccessServiceServer is the server API for AccessService service.
// All implementations must embed UnimplementedAccessServiceServer
// for forward compatibility
type AccessServiceServer interface {
	RequestFabricCredentials(context.Context, *presenter.FabricCredentialsRequest) (*presenter.FabricCredentialsResponse, error)
	UpdatePassword(context.Context, *presenter.UpdatePasswordRequest) (*presenter.StatusResponse, error)
	mustEmbedUnimplementedAccessServiceServer()
}

// UnimplementedAccessServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAccessServiceServer struct {
}

func (UnimplementedAccessServiceServer) RequestFabricCredentials(context.Context, *presenter.FabricCredentialsRequest) (*presenter.FabricCredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestFabricCredentials not implemented")
}
func (UnimplementedAccessServiceServer) UpdatePassword(context.Context, *presenter.UpdatePasswordRequest) (*presenter.StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePassword not implemented")
}
func (UnimplementedAccessServiceServer) mustEmbedUnimplementedAccessServiceServer() {}

// UnsafeAccessServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccessServiceServer will
// result in compilation errors.
type UnsafeAccessServiceServer interface {
	mustEmbedUnimplementedAccessServiceServer()
}

func RegisterAccessServiceServer(s grpc.ServiceRegistrar, srv AccessServiceServer) {
	s.RegisterService(&AccessService_ServiceDesc, srv)
}

func _AccessService_RequestFabricCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(presenter.FabricCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServiceServer).RequestFabricCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chainmetric.identity.service.AccessService/requestFabricCredentials",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServiceServer).RequestFabricCredentials(ctx, req.(*presenter.FabricCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessService_UpdatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(presenter.UpdatePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServiceServer).UpdatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chainmetric.identity.service.AccessService/updatePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServiceServer).UpdatePassword(ctx, req.(*presenter.UpdatePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccessService_ServiceDesc is the grpc.ServiceDesc for AccessService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccessService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chainmetric.identity.service.AccessService",
	HandlerType: (*AccessServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "requestFabricCredentials",
			Handler:    _AccessService_RequestFabricCredentials_Handler,
		},
		{
			MethodName: "updatePassword",
			Handler:    _AccessService_UpdatePassword_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "identity/api/rpc/access.proto",
}