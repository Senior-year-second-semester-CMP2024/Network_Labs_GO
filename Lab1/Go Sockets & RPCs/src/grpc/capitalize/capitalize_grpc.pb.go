// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: capitalize/capitalize.proto

package capitalize

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

// TextServiceClient is the client API for TextService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TextServiceClient interface {
	Capitalize(ctx context.Context, in *TextRequest, opts ...grpc.CallOption) (*TextResponse, error)
}

type textServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTextServiceClient(cc grpc.ClientConnInterface) TextServiceClient {
	return &textServiceClient{cc}
}

func (c *textServiceClient) Capitalize(ctx context.Context, in *TextRequest, opts ...grpc.CallOption) (*TextResponse, error) {
	out := new(TextResponse)
	err := c.cc.Invoke(ctx, "/capitalize.TextService/Capitalize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TextServiceServer is the server API for TextService service.
// All implementations must embed UnimplementedTextServiceServer
// for forward compatibility
type TextServiceServer interface {
	Capitalize(context.Context, *TextRequest) (*TextResponse, error)
	mustEmbedUnimplementedTextServiceServer()
}

// UnimplementedTextServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTextServiceServer struct {
}

func (UnimplementedTextServiceServer) Capitalize(context.Context, *TextRequest) (*TextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Capitalize not implemented")
}
func (UnimplementedTextServiceServer) mustEmbedUnimplementedTextServiceServer() {}

// UnsafeTextServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TextServiceServer will
// result in compilation errors.
type UnsafeTextServiceServer interface {
	mustEmbedUnimplementedTextServiceServer()
}

func RegisterTextServiceServer(s grpc.ServiceRegistrar, srv TextServiceServer) {
	s.RegisterService(&TextService_ServiceDesc, srv)
}

func _TextService_Capitalize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TextServiceServer).Capitalize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/capitalize.TextService/Capitalize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TextServiceServer).Capitalize(ctx, req.(*TextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TextService_ServiceDesc is the grpc.ServiceDesc for TextService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TextService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "capitalize.TextService",
	HandlerType: (*TextServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Capitalize",
			Handler:    _TextService_Capitalize_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "capitalize/capitalize.proto",
}
