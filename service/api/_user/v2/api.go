package apiv2user

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"x-gwi/app/storage"
	sharepb "x-gwi/proto/core/_share/v1"
	userpb "x-gwi/proto/core/_user/v1"
	pbsrv "x-gwi/proto/serv/_user/v2"
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

	// user "x-gwi/service/core/_user"
	s.userCore, err = user.NewCore(storage)
	if err != nil {
		return nil, fmt.Errorf("user.NewCore: %w", err)
	}

	pbsrv.RegisterUserServiceServer(grpcServer, s)

	return s, nil
}

func (s *ServiceAPI) Create(ctx context.Context, in *userpb.UserCore) (*userpb.UserCore, error) {
	_, _ = (pbsrv.UnimplementedUserServiceServer{}).Create(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process by core
	err := s.userCore.Create(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "userCore.Create: %v", err)
	}

	return in, nil
}

func (s *ServiceAPI) Get(ctx context.Context, in *sharepb.ShareQID) (*userpb.UserCore, error) {
	_, _ = (pbsrv.UnimplementedUserServiceServer{}).Get(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &userpb.UserCore{ //nolint:exhaustruct
		Qid: in,
	}

	// 3. process by core
	err := s.userCore.Get(ctx, out)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "userCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Gett(ctx context.Context, in *sharepb.ShareQID) (*userpb.UserCore, error) {
	_, _ = (pbsrv.UnimplementedUserServiceServer{}).Gett(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &userpb.UserCore{ //nolint:exhaustruct
		Qid: in,
	}

	// 3. process by core
	err := s.userCore.Get(ctx, out)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "userCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Update(ctx context.Context, in *userpb.UserCore) (*userpb.UserCore, error) {
	_, _ = (pbsrv.UnimplementedUserServiceServer{}).Update(ctx, in)

	// return s.update(ctx, in)
	return (pbsrv.UnimplementedUserServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *ServiceAPI) Delete(ctx context.Context, in *sharepb.ShareQID) (*userpb.UserCore, error) {
	_, _ = (pbsrv.UnimplementedUserServiceServer{}).Delete(ctx, in)

	// return s.delete(ctx, in)
	return (pbsrv.UnimplementedUserServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *ServiceAPI) ListFavourites(in *sharepb.ShareQID, stream pbsrv.UserService_ListFavouritesServer) error {
	_ = (pbsrv.UnimplementedUserServiceServer{}).ListFavourites(in, stream)

	// return s.list(in, stream)
	return (pbsrv.UnimplementedUserServiceServer{}).ListFavourites(in, stream) //nolint:wrapcheck
}

func (s *ServiceAPI) ListOpinions(in *sharepb.ShareQID, stream pbsrv.UserService_ListOpinionsServer) error {
	_ = (pbsrv.UnimplementedUserServiceServer{}).ListOpinions(in, stream)

	// return s.list(in, stream)
	return (pbsrv.UnimplementedUserServiceServer{}).ListOpinions(in, stream) //nolint:wrapcheck
}
