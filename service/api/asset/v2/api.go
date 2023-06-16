package assetsrv2

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"x-gwi/app/storage"
	pb_idx "x-gwi/proto/core/_store/v1"
	pb "x-gwi/proto/core/asset/v1"
	srvpb "x-gwi/proto/serv/asset/v2"
	"x-gwi/service/core/asset"
)

type Service struct {
	srvpb.UnimplementedAssetServiceServer
	asset   *asset.CoreAsset
	storage *storage.ServiceStorage
}

func RegisterGRPC(grpcServer *grpc.Server, storage *storage.ServiceStorage) (*Service, error) {
	var err error

	s := &Service{ //nolint:exhaustruct
		storage: storage,
	}

	s.asset, err = asset.NewCore(s.storage)
	if err != nil {
		return nil, fmt.Errorf("asset.NewCore: %w", err)
	}

	srvpb.RegisterAssetServiceServer(grpcServer, s)

	return s, nil
}

func (s *Service) Create(ctx context.Context, in *pb.AssetInstance) (*pb.AssetInstance, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Create(ctx, in)

	// return s.create(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Create(ctx, in) //nolint:wrapcheck
}

func (s *Service) Get(ctx context.Context, in *pb_idx.StoreIDX) (*pb.AssetInstance, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Get(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Get(ctx, in) //nolint:wrapcheck
}

func (s *Service) Gett(ctx context.Context, in *pb_idx.StoreIDX) (*pb.AssetInstance, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Gett(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Gett(ctx, in) //nolint:wrapcheck
}

func (s *Service) Update(ctx context.Context, in *pb.AssetInstance) (*pb.AssetInstance, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Update(ctx, in)

	// return s.update(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *Service) Delete(ctx context.Context, in *pb_idx.StoreIDX) (*pb.AssetInstance, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Delete(ctx, in)

	// return s.delete(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}
