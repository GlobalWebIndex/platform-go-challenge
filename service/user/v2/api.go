package usersrv2

import (
	"context"

	"google.golang.org/grpc"

	"x-gwi/app/storage"
	pb_idx "x-gwi/proto/core/idx/v1"
	pb "x-gwi/proto/core/user/v1"
	srvpb "x-gwi/proto/serv/user/v2"
)

type Service struct {
	srvpb.UnimplementedUserServiceServer
	storage *storage.Storage
}

// Register srvpb.UserServiceServer interface
func RegisterGRPC(grpcServer *grpc.Server, storage *storage.Storage) {
	//nolint:exhaustruct
	s := &Service{
		storage: storage,
	}
	srvpb.RegisterUserServiceServer(grpcServer, s)
}

func (s *Service) Create(ctx context.Context, in *pb.UserInstance) (*pb.UserInstance, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Create(ctx, in)

	// return s.create(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Create(ctx, in) //nolint:wrapcheck
}

func (s *Service) Get(ctx context.Context, in *pb_idx.IDX) (*pb.UserInstance, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Get(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Get(ctx, in) //nolint:wrapcheck
}

func (s *Service) Gett(ctx context.Context, in *pb_idx.IDX) (*pb.UserInstance, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Gett(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Gett(ctx, in) //nolint:wrapcheck
}

func (s *Service) Update(ctx context.Context, in *pb.UserInstance) (*pb.UserInstance, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Update(ctx, in)

	// return s.update(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *Service) Delete(ctx context.Context, in *pb_idx.IDX) (*pb.UserInstance, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Delete(ctx, in)

	// return s.delete(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *Service) ListFavourites(in *pb_idx.IDX, stream srvpb.UserService_ListFavouritesServer) error {
	_ = (srvpb.UnimplementedUserServiceServer{}).ListFavourites(in, stream)

	// return s.list(in, stream)
	return (srvpb.UnimplementedUserServiceServer{}).ListFavourites(in, stream) //nolint:wrapcheck
}

func (s *Service) ListOpinions(in *pb_idx.IDX, stream srvpb.UserService_ListOpinionsServer) error {
	_ = (srvpb.UnimplementedUserServiceServer{}).ListOpinions(in, stream)

	// return s.list(in, stream)
	return (srvpb.UnimplementedUserServiceServer{}).ListOpinions(in, stream) //nolint:wrapcheck
}
