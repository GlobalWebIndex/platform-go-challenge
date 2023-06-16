// proto/serv/user/v2/user_srv.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.2
// source: proto/serv/user/v2/user_srv.proto

package user_srvpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	v11 "x-gwi/proto/core/_store/v1"
	v12 "x-gwi/proto/core/favourite/v1"
	v13 "x-gwi/proto/core/opinion/v1"
	v1 "x-gwi/proto/core/user/v1"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	UserService_Create_FullMethodName         = "/proto.serv.user.v2.UserService/Create"
	UserService_Get_FullMethodName            = "/proto.serv.user.v2.UserService/Get"
	UserService_Gett_FullMethodName           = "/proto.serv.user.v2.UserService/Gett"
	UserService_Update_FullMethodName         = "/proto.serv.user.v2.UserService/Update"
	UserService_Delete_FullMethodName         = "/proto.serv.user.v2.UserService/Delete"
	UserService_ListFavourites_FullMethodName = "/proto.serv.user.v2.UserService/ListFavourites"
	UserService_ListOpinions_FullMethodName   = "/proto.serv.user.v2.UserService/ListOpinions"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	// Create
	Create(ctx context.Context, in *v1.UserInstance, opts ...grpc.CallOption) (*v1.UserInstance, error)
	// Get
	Get(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.UserInstance, error)
	// Gett
	Gett(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.UserInstance, error)
	// Update
	Update(ctx context.Context, in *v1.UserInstance, opts ...grpc.CallOption) (*v1.UserInstance, error)
	// Delete
	Delete(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.UserInstance, error)
	// ListFavourites - stream favourites of a user
	ListFavourites(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (UserService_ListFavouritesClient, error)
	// ListOpinions - stream opinions of a user
	ListOpinions(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (UserService_ListOpinionsClient, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) Create(ctx context.Context, in *v1.UserInstance, opts ...grpc.CallOption) (*v1.UserInstance, error) {
	out := new(v1.UserInstance)
	err := c.cc.Invoke(ctx, UserService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Get(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.UserInstance, error) {
	out := new(v1.UserInstance)
	err := c.cc.Invoke(ctx, UserService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Gett(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.UserInstance, error) {
	out := new(v1.UserInstance)
	err := c.cc.Invoke(ctx, UserService_Gett_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Update(ctx context.Context, in *v1.UserInstance, opts ...grpc.CallOption) (*v1.UserInstance, error) {
	out := new(v1.UserInstance)
	err := c.cc.Invoke(ctx, UserService_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Delete(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (*v1.UserInstance, error) {
	out := new(v1.UserInstance)
	err := c.cc.Invoke(ctx, UserService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ListFavourites(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (UserService_ListFavouritesClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserService_ServiceDesc.Streams[0], UserService_ListFavourites_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &userServiceListFavouritesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UserService_ListFavouritesClient interface {
	Recv() (*v12.FavouriteAsset, error)
	grpc.ClientStream
}

type userServiceListFavouritesClient struct {
	grpc.ClientStream
}

func (x *userServiceListFavouritesClient) Recv() (*v12.FavouriteAsset, error) {
	m := new(v12.FavouriteAsset)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userServiceClient) ListOpinions(ctx context.Context, in *v11.StoreIDX, opts ...grpc.CallOption) (UserService_ListOpinionsClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserService_ServiceDesc.Streams[1], UserService_ListOpinions_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &userServiceListOpinionsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UserService_ListOpinionsClient interface {
	Recv() (*v13.OpinionAsset, error)
	grpc.ClientStream
}

type userServiceListOpinionsClient struct {
	grpc.ClientStream
}

func (x *userServiceListOpinionsClient) Recv() (*v13.OpinionAsset, error) {
	m := new(v13.OpinionAsset)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	// Create
	Create(context.Context, *v1.UserInstance) (*v1.UserInstance, error)
	// Get
	Get(context.Context, *v11.StoreIDX) (*v1.UserInstance, error)
	// Gett
	Gett(context.Context, *v11.StoreIDX) (*v1.UserInstance, error)
	// Update
	Update(context.Context, *v1.UserInstance) (*v1.UserInstance, error)
	// Delete
	Delete(context.Context, *v11.StoreIDX) (*v1.UserInstance, error)
	// ListFavourites - stream favourites of a user
	ListFavourites(*v11.StoreIDX, UserService_ListFavouritesServer) error
	// ListOpinions - stream opinions of a user
	ListOpinions(*v11.StoreIDX, UserService_ListOpinionsServer) error
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) Create(context.Context, *v1.UserInstance) (*v1.UserInstance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedUserServiceServer) Get(context.Context, *v11.StoreIDX) (*v1.UserInstance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedUserServiceServer) Gett(context.Context, *v11.StoreIDX) (*v1.UserInstance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Gett not implemented")
}
func (UnimplementedUserServiceServer) Update(context.Context, *v1.UserInstance) (*v1.UserInstance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedUserServiceServer) Delete(context.Context, *v11.StoreIDX) (*v1.UserInstance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedUserServiceServer) ListFavourites(*v11.StoreIDX, UserService_ListFavouritesServer) error {
	return status.Errorf(codes.Unimplemented, "method ListFavourites not implemented")
}
func (UnimplementedUserServiceServer) ListOpinions(*v11.StoreIDX, UserService_ListOpinionsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListOpinions not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.UserInstance)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Create(ctx, req.(*v1.UserInstance))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.StoreIDX)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Get(ctx, req.(*v11.StoreIDX))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Gett_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.StoreIDX)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Gett(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Gett_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Gett(ctx, req.(*v11.StoreIDX))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.UserInstance)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Update(ctx, req.(*v1.UserInstance))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.StoreIDX)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Delete(ctx, req.(*v11.StoreIDX))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ListFavourites_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(v11.StoreIDX)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServiceServer).ListFavourites(m, &userServiceListFavouritesServer{stream})
}

type UserService_ListFavouritesServer interface {
	Send(*v12.FavouriteAsset) error
	grpc.ServerStream
}

type userServiceListFavouritesServer struct {
	grpc.ServerStream
}

func (x *userServiceListFavouritesServer) Send(m *v12.FavouriteAsset) error {
	return x.ServerStream.SendMsg(m)
}

func _UserService_ListOpinions_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(v11.StoreIDX)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServiceServer).ListOpinions(m, &userServiceListOpinionsServer{stream})
}

type UserService_ListOpinionsServer interface {
	Send(*v13.OpinionAsset) error
	grpc.ServerStream
}

type userServiceListOpinionsServer struct {
	grpc.ServerStream
}

func (x *userServiceListOpinionsServer) Send(m *v13.OpinionAsset) error {
	return x.ServerStream.SendMsg(m)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.serv.user.v2.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _UserService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _UserService_Get_Handler,
		},
		{
			MethodName: "Gett",
			Handler:    _UserService_Gett_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _UserService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _UserService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListFavourites",
			Handler:       _UserService_ListFavourites_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListOpinions",
			Handler:       _UserService_ListOpinions_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/serv/user/v2/user_srv.proto",
}
