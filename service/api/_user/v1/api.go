package apiv1user

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"x-gwi/app/storage"
	userpb "x-gwi/proto/core/_user/v1"
	pbsrv "x-gwi/proto/serv/_user/v1"
	user "x-gwi/service/core/_user"
)

type ServiceAPI struct {
	pbsrv.UnimplementedUserServiceServer
	userCore *user.CoreUser // user "x-gwi/service/core/_user"
	// storage  *storage.ServiceStorage
}

func RegisterGRPC(grpcServer *grpc.Server, storage *storage.ServiceStorage) (*ServiceAPI, error) {
	var err error

	s := &ServiceAPI{ //nolint:exhaustruct
		// storage: storage,
	}

	s.userCore, err = user.NewCore(storage)
	if err != nil {
		return nil, fmt.Errorf("user.NewCore: %w", err)
	}

	pbsrv.RegisterUserServiceServer(grpcServer, s)

	return s, nil
}

func (s *ServiceAPI) Create(ctx context.Context, in *pbsrv.CreateRequest) (*pbsrv.CreateResponse, error) {
	_, _ = (pbsrv.UnimplementedUserServiceServer{}).Create(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process by core
	err := s.userCore.Create(ctx, in.GetUser())
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "userCore.Create: %v", err)
	}

	// 3. adapt output
	out := &pbsrv.CreateResponse{
		User: in.GetUser(),
	}

	return out, nil
}

func (s *ServiceAPI) Get(ctx context.Context, in *pbsrv.GetRequest) (*pbsrv.GetResponse, error) {
	_, _ = (pbsrv.UnimplementedUserServiceServer{}).Get(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &pbsrv.GetResponse{ //nolint:exhaustruct
		User: &userpb.UserCore{
			Qid: in.GetQid(),
		},
	}

	// 3. process by core
	err := s.userCore.Get(ctx, out.User)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "userCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Gett(ctx context.Context, in *pbsrv.GetRequest) (*pbsrv.GetResponse, error) {
	_, _ = (pbsrv.UnimplementedUserServiceServer{}).Gett(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &pbsrv.GetResponse{ //nolint:exhaustruct
		User: &userpb.UserCore{
			Qid: in.GetQid(),
		},
	}

	// 3. process by core
	err := s.userCore.Get(ctx, out.User)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "userCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Update(ctx context.Context, in *pbsrv.UpdateRequest) (*pbsrv.UpdateResponse, error) {
	_, _ = (pbsrv.UnimplementedUserServiceServer{}).Update(ctx, in)

	// return (pbsrv.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (pbsrv.UnimplementedUserServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *ServiceAPI) Delete(ctx context.Context, in *pbsrv.DeleteRequest) (*pbsrv.DeleteResponse, error) {
	_, _ = (pbsrv.UnimplementedUserServiceServer{}).Delete(ctx, in)

	// return (pbsrv.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (pbsrv.UnimplementedUserServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *ServiceAPI) List(in *pbsrv.ListRequest, stream pbsrv.UserService_ListServer) error {
	_ = (pbsrv.UnimplementedUserServiceServer{}).List(in, stream)

	// return s.list(in, stream)
	return (pbsrv.UnimplementedUserServiceServer{}).List(in, stream) //nolint:wrapcheck
}
