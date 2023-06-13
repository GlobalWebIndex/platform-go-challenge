package usersrv

import (
	"context"

	"google.golang.org/grpc"

	"x-gwi/app/storage"
	srvpb "x-gwi/proto/serv/user/v1"
)

type Service struct {
	srvpb.UnimplementedUserServiceServer
	storage *storage.Storage
}

// Register srvpb.AuthServiceServer interface
func RegisterGRPC(grpcServer *grpc.Server, storage *storage.Storage) {
	//nolint:exhaustruct
	s := &Service{
		storage: storage,
	}
	srvpb.RegisterUserServiceServer(grpcServer, s)
}

func (s *Service) Create(ctx context.Context, in *srvpb.CreateRequest) (*srvpb.CreateResponse, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Create(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Create(ctx, in) //nolint:wrapcheck
}

func (s *Service) Get(ctx context.Context, in *srvpb.GetRequest) (*srvpb.GetResponse, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Get(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Get(ctx, in) //nolint:wrapcheck
}

func (s *Service) Gett(ctx context.Context, in *srvpb.GetRequest) (*srvpb.GetResponse, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Gett(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Gett(ctx, in) //nolint:wrapcheck
}

func (s *Service) Update(ctx context.Context, in *srvpb.UpdateRequest) (*srvpb.UpdateResponse, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Update(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *Service) Delete(ctx context.Context, in *srvpb.DeleteRequest) (*srvpb.DeleteResponse, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Delete(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *Service) ListFavourites(in *srvpb.ListFavouritesRequest, stream srvpb.UserService_ListFavouritesServer) error {
	_ = (srvpb.UnimplementedUserServiceServer{}).ListFavourites(in, stream)

	// return s.list(in, stream)
	return (srvpb.UnimplementedUserServiceServer{}).ListFavourites(in, stream) //nolint:wrapcheck
}

func (s *Service) ListOpinions(in *srvpb.ListOpinionsRequest, stream srvpb.UserService_ListOpinionsServer) error {
	_ = (srvpb.UnimplementedUserServiceServer{}).ListOpinions(in, stream)

	// return s.list(in, stream)
	return (srvpb.UnimplementedUserServiceServer{}).ListOpinions(in, stream) //nolint:wrapcheck
}
