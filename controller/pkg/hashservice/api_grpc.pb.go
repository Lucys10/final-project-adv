// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package hashservice

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

// HashClient is the controller API for Hash service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HashClient interface {
	CalculateHash(ctx context.Context, in *StrList, opts ...grpc.CallOption) (*HashList, error)
}

type hashClient struct {
	cc grpc.ClientConnInterface
}

func NewHashClient(cc grpc.ClientConnInterface) HashClient {
	return &hashClient{cc}
}

func (c *hashClient) CalculateHash(ctx context.Context, in *StrList, opts ...grpc.CallOption) (*HashList, error) {
	out := new(HashList)
	err := c.cc.Invoke(ctx, "/hashserver.Hash/CalculateHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HashServer is the httpserver API for Hash service.
// All implementations must embed UnimplementedHashServer
// for forward compatibility
type HashServer interface {
	CalculateHash(context.Context, *StrList) (*HashList, error)
	mustEmbedUnimplementedHashServer()
}

// UnimplementedHashServer must be embedded to have forward compatible implementations.
type UnimplementedHashServer struct {
}

func (UnimplementedHashServer) CalculateHash(context.Context, *StrList) (*HashList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculateHash not implemented")
}
func (UnimplementedHashServer) mustEmbedUnimplementedHashServer() {}

// UnsafeHashServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HashServer will
// result in compilation errors.
type UnsafeHashServer interface {
	mustEmbedUnimplementedHashServer()
}

func RegisterHashServer(s grpc.ServiceRegistrar, srv HashServer) {
	s.RegisterService(&Hash_ServiceDesc, srv)
}

func _Hash_CalculateHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StrList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashServer).CalculateHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hashserver.Hash/CalculateHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashServer).CalculateHash(ctx, req.(*StrList))
	}
	return interceptor(ctx, in, info, handler)
}

// Hash_ServiceDesc is the grpc.ServiceDesc for Hash service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Hash_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hashserver.Hash",
	HandlerType: (*HashServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CalculateHash",
			Handler:    _Hash_CalculateHash_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/api.proto",
}
