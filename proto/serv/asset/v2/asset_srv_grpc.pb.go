// proto/serv/asset/v2/asset_srv.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.2
// source: proto/serv/asset/v2/asset_srv.proto

package asset_srvpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	v11 "x-gwi/proto/core/_store/v1"
	v1 "x-gwi/proto/core/asset/v1"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	AssetService_Create_FullMethodName = "/proto.serv.asset.v2.AssetService/Create"
	AssetService_Get_FullMethodName    = "/proto.serv.asset.v2.AssetService/Get"
	AssetService_Gett_FullMethodName   = "/proto.serv.asset.v2.AssetService/Gett"
	AssetService_Update_FullMethodName = "/proto.serv.asset.v2.AssetService/Update"
	AssetService_Delete_FullMethodName = "/proto.serv.asset.v2.AssetService/Delete"
)

// AssetServiceClient is the client API for AssetService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AssetServiceClient interface {
	// Create
	Create(ctx context.Context, in *v1.AssetInstance, opts ...grpc.CallOption) (*v1.AssetInstance, error)
	// Get
	Get(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.AssetInstance, error)
	// Gett
	Gett(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.AssetInstance, error)
	// Update
	Update(ctx context.Context, in *v1.AssetInstance, opts ...grpc.CallOption) (*v1.AssetInstance, error)
	// Delete
	Delete(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.AssetInstance, error)
}

type assetServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAssetServiceClient(cc grpc.ClientConnInterface) AssetServiceClient {
	return &assetServiceClient{cc}
}

func (c *assetServiceClient) Create(ctx context.Context, in *v1.AssetInstance, opts ...grpc.CallOption) (*v1.AssetInstance, error) {
	out := new(v1.AssetInstance)
	err := c.cc.Invoke(ctx, AssetService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assetServiceClient) Get(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.AssetInstance, error) {
	out := new(v1.AssetInstance)
	err := c.cc.Invoke(ctx, AssetService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assetServiceClient) Gett(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.AssetInstance, error) {
	out := new(v1.AssetInstance)
	err := c.cc.Invoke(ctx, AssetService_Gett_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assetServiceClient) Update(ctx context.Context, in *v1.AssetInstance, opts ...grpc.CallOption) (*v1.AssetInstance, error) {
	out := new(v1.AssetInstance)
	err := c.cc.Invoke(ctx, AssetService_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assetServiceClient) Delete(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.AssetInstance, error) {
	out := new(v1.AssetInstance)
	err := c.cc.Invoke(ctx, AssetService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AssetServiceServer is the server API for AssetService service.
// All implementations must embed UnimplementedAssetServiceServer
// for forward compatibility
type AssetServiceServer interface {
	// Create
	Create(context.Context, *v1.AssetInstance) (*v1.AssetInstance, error)
	// Get
	Get(context.Context, *v11.StoreIDX) (*v1.AssetInstance, error)
	// Gett
	Gett(context.Context, *v11.StoreIDX) (*v1.AssetInstance, error)
	// Update
	Update(context.Context, *v1.AssetInstance) (*v1.AssetInstance, error)
	// Delete
	Delete(context.Context, *v11.StoreIDX) (*v1.AssetInstance, error)
	mustEmbedUnimplementedAssetServiceServer()
}

// UnimplementedAssetServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAssetServiceServer struct {
}

func (UnimplementedAssetServiceServer) Create(context.Context, *v1.AssetInstance) (*v1.AssetInstance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedAssetServiceServer) Get(context.Context, *v11.StoreIDX) (*v1.AssetInstance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedAssetServiceServer) Gett(context.Context, *v11.StoreIDX) (*v1.AssetInstance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Gett not implemented")
}
func (UnimplementedAssetServiceServer) Update(context.Context, *v1.AssetInstance) (*v1.AssetInstance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedAssetServiceServer) Delete(context.Context, *v11.StoreIDX) (*v1.AssetInstance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedAssetServiceServer) mustEmbedUnimplementedAssetServiceServer() {}

// UnsafeAssetServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AssetServiceServer will
// result in compilation errors.
type UnsafeAssetServiceServer interface {
	mustEmbedUnimplementedAssetServiceServer()
}

func RegisterAssetServiceServer(s grpc.ServiceRegistrar, srv AssetServiceServer) {
	s.RegisterService(&AssetService_ServiceDesc, srv)
}

func _AssetService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.AssetInstance)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssetServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AssetService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssetServiceServer).Create(ctx, req.(*v1.AssetInstance))
	}
	return interceptor(ctx, in, info, handler)
}

func _AssetService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.StoreIDX)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssetServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AssetService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssetServiceServer).Get(ctx, req.(*v11.StoreIDX))
	}
	return interceptor(ctx, in, info, handler)
}

func _AssetService_Gett_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.StoreIDX)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssetServiceServer).Gett(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AssetService_Gett_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssetServiceServer).Gett(ctx, req.(*v11.StoreIDX))
	}
	return interceptor(ctx, in, info, handler)
}

func _AssetService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.AssetInstance)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssetServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AssetService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssetServiceServer).Update(ctx, req.(*v1.AssetInstance))
	}
	return interceptor(ctx, in, info, handler)
}

func _AssetService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.StoreIDX)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssetServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AssetService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssetServiceServer).Delete(ctx, req.(*v11.StoreIDX))
	}
	return interceptor(ctx, in, info, handler)
}

// AssetService_ServiceDesc is the grpc.ServiceDesc for AssetService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AssetService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.serv.asset.v2.AssetService",
	HandlerType: (*AssetServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _AssetService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _AssetService_Get_Handler,
		},
		{
			MethodName: "Gett",
			Handler:    _AssetService_Gett_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _AssetService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _AssetService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/serv/asset/v2/asset_srv.proto",
}
