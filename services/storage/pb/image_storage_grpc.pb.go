// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.2
// source: image_storage.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ImageStorage_CreateImage_FullMethodName = "/gateway.ImageStorage/CreateImage"
	ImageStorage_GetImage_FullMethodName    = "/gateway.ImageStorage/GetImage"
	ImageStorage_GetImages_FullMethodName   = "/gateway.ImageStorage/GetImages"
)

// ImageStorageClient is the client API for ImageStorage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ImageStorageClient interface {
	CreateImage(ctx context.Context, in *CreateImageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetImage(ctx context.Context, in *FindImageRequest, opts ...grpc.CallOption) (*FindImageResponse, error)
	GetImages(ctx context.Context, in *FindImagesRequest, opts ...grpc.CallOption) (*FindImagesResponse, error)
}

type imageStorageClient struct {
	cc grpc.ClientConnInterface
}

func NewImageStorageClient(cc grpc.ClientConnInterface) ImageStorageClient {
	return &imageStorageClient{cc}
}

func (c *imageStorageClient) CreateImage(ctx context.Context, in *CreateImageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ImageStorage_CreateImage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageStorageClient) GetImage(ctx context.Context, in *FindImageRequest, opts ...grpc.CallOption) (*FindImageResponse, error) {
	out := new(FindImageResponse)
	err := c.cc.Invoke(ctx, ImageStorage_GetImage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageStorageClient) GetImages(ctx context.Context, in *FindImagesRequest, opts ...grpc.CallOption) (*FindImagesResponse, error) {
	out := new(FindImagesResponse)
	err := c.cc.Invoke(ctx, ImageStorage_GetImages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImageStorageServer is the server API for ImageStorage service.
// All implementations must embed UnimplementedImageStorageServer
// for forward compatibility
type ImageStorageServer interface {
	CreateImage(context.Context, *CreateImageRequest) (*emptypb.Empty, error)
	GetImage(context.Context, *FindImageRequest) (*FindImageResponse, error)
	GetImages(context.Context, *FindImagesRequest) (*FindImagesResponse, error)
	mustEmbedUnimplementedImageStorageServer()
}

// UnimplementedImageStorageServer must be embedded to have forward compatible implementations.
type UnimplementedImageStorageServer struct {
}

func (UnimplementedImageStorageServer) CreateImage(context.Context, *CreateImageRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateImage not implemented")
}
func (UnimplementedImageStorageServer) GetImage(context.Context, *FindImageRequest) (*FindImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetImage not implemented")
}
func (UnimplementedImageStorageServer) GetImages(context.Context, *FindImagesRequest) (*FindImagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetImages not implemented")
}
func (UnimplementedImageStorageServer) mustEmbedUnimplementedImageStorageServer() {}

// UnsafeImageStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ImageStorageServer will
// result in compilation errors.
type UnsafeImageStorageServer interface {
	mustEmbedUnimplementedImageStorageServer()
}

func RegisterImageStorageServer(s grpc.ServiceRegistrar, srv ImageStorageServer) {
	s.RegisterService(&ImageStorage_ServiceDesc, srv)
}

func _ImageStorage_CreateImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageStorageServer).CreateImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ImageStorage_CreateImage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageStorageServer).CreateImage(ctx, req.(*CreateImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageStorage_GetImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageStorageServer).GetImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ImageStorage_GetImage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageStorageServer).GetImage(ctx, req.(*FindImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageStorage_GetImages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindImagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageStorageServer).GetImages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ImageStorage_GetImages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageStorageServer).GetImages(ctx, req.(*FindImagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ImageStorage_ServiceDesc is the grpc.ServiceDesc for ImageStorage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ImageStorage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gateway.ImageStorage",
	HandlerType: (*ImageStorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateImage",
			Handler:    _ImageStorage_CreateImage_Handler,
		},
		{
			MethodName: "GetImage",
			Handler:    _ImageStorage_GetImage_Handler,
		},
		{
			MethodName: "GetImages",
			Handler:    _ImageStorage_GetImages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "image_storage.proto",
}
