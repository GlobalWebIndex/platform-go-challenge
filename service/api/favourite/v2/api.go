package favouritesrv2

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"x-gwi/app/storage"
	pb_idx "x-gwi/proto/core/_store/v1"
	pb "x-gwi/proto/core/favourite/v1"
	srvpb "x-gwi/proto/serv/favourite/v2"
	"x-gwi/service/core/favourite"
)

type Service struct {
	srvpb.UnimplementedFavouriteServiceServer
	favourite *favourite.CoreFavourite
	storage   *storage.ServiceStorage
}

func RegisterGRPC(grpcServer *grpc.Server, storage *storage.ServiceStorage) (*Service, error) {
	var err error

	s := &Service{ //nolint:exhaustruct
		storage: storage,
	}

	s.favourite, err = favourite.NewCore(s.storage)
	if err != nil {
		return nil, fmt.Errorf("favourite.NewCore: %w", err)
	}

	srvpb.RegisterFavouriteServiceServer(grpcServer, s)

	return s, nil
}

func (s *Service) Create(ctx context.Context, in *pb.FavouriteAsset) (*pb.FavouriteAsset, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Create(ctx, in)

	// return s.create(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Create(ctx, in) //nolint:wrapcheck
}

func (s *Service) Get(ctx context.Context, in *pb_idx.StoreIDX) (*pb.FavouriteAsset, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Get(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Get(ctx, in) //nolint:wrapcheck
}

func (s *Service) Gett(ctx context.Context, in *pb_idx.StoreIDX) (*pb.FavouriteAsset, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Gett(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Gett(ctx, in) //nolint:wrapcheck
}

func (s *Service) Update(ctx context.Context, in *pb.FavouriteAsset) (*pb.FavouriteAsset, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Update(ctx, in)

	// return s.update(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *Service) Delete(ctx context.Context, in *pb_idx.StoreIDX) (*pb.FavouriteAsset, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Delete(ctx, in)

	// return s.delete(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *Service) List(in *pb_idx.StoreIDX, stream srvpb.FavouriteService_ListServer) error {
	_ = (srvpb.UnimplementedFavouriteServiceServer{}).List(in, stream)

	// return s.list(in, stream)
	return (srvpb.UnimplementedFavouriteServiceServer{}).List(in, stream) //nolint:wrapcheck
}
