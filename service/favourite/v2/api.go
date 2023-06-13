package favouritesrv2

import (
	"context"

	"google.golang.org/grpc"

	"x-gwi/app/storage"
	pb "x-gwi/proto/core/favourite/v1"
	pb_idx "x-gwi/proto/core/idx/v1"
	srvpb "x-gwi/proto/serv/favourite/v2"
)

type Service struct {
	srvpb.UnimplementedFavouriteServiceServer
	storage *storage.Storage
}

// Register srvpb.FavouriteServiceServer interface
func RegisterGRPC(grpcServer *grpc.Server, storage *storage.Storage) {
	//nolint:exhaustruct
	s := &Service{
		storage: storage,
	}
	srvpb.RegisterFavouriteServiceServer(grpcServer, s)
}

func (s *Service) Create(ctx context.Context, in *pb.FavouriteAsset) (*pb.FavouriteAsset, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Create(ctx, in)

	// return s.create(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Create(ctx, in) //nolint:wrapcheck
}

func (s *Service) Get(ctx context.Context, in *pb_idx.IDX) (*pb.FavouriteAsset, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Get(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Get(ctx, in) //nolint:wrapcheck
}

func (s *Service) Gett(ctx context.Context, in *pb_idx.IDX) (*pb.FavouriteAsset, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Gett(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Gett(ctx, in) //nolint:wrapcheck
}

func (s *Service) Update(ctx context.Context, in *pb.FavouriteAsset) (*pb.FavouriteAsset, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Update(ctx, in)

	// return s.update(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *Service) Delete(ctx context.Context, in *pb_idx.IDX) (*pb.FavouriteAsset, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Delete(ctx, in)

	// return s.delete(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *Service) List(in *pb_idx.IDX, stream srvpb.FavouriteService_ListServer) error {
	_ = (srvpb.UnimplementedFavouriteServiceServer{}).List(in, stream)

	// return s.list(in, stream)
	return (srvpb.UnimplementedFavouriteServiceServer{}).List(in, stream) //nolint:wrapcheck
}
