package assetsrv2

import (
	"context"

	"google.golang.org/grpc"

	"x-gwi/app/storage"
	pb "x-gwi/proto/core/asset/v1"
	pb_idx "x-gwi/proto/core/idx/v1"
	srvpb "x-gwi/proto/serv/asset/v2"
)

type Service struct {
	srvpb.UnimplementedAssetServiceServer
	storage *storage.Storage
}

// Register srvpb.AuthServiceServer interface
func RegisterGRPC(grpcServer *grpc.Server, storage *storage.Storage) {
	//nolint:exhaustruct
	s := &Service{
		storage: storage,
	}
	srvpb.RegisterAssetServiceServer(grpcServer, s)
}

func (s *Service) Create(ctx context.Context, in *pb.AssetInstance) (*pb.AssetInstance, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Create(ctx, in)

	// return s.create(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Create(ctx, in) //nolint:wrapcheck
}

func (s *Service) Get(ctx context.Context, in *pb_idx.IDX) (*pb.AssetInstance, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Get(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Get(ctx, in) //nolint:wrapcheck
}

func (s *Service) Gett(ctx context.Context, in *pb_idx.IDX) (*pb.AssetInstance, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Gett(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Gett(ctx, in) //nolint:wrapcheck
}

func (s *Service) Update(ctx context.Context, in *pb.AssetInstance) (*pb.AssetInstance, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Update(ctx, in)

	// return s.update(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *Service) Delete(ctx context.Context, in *pb_idx.IDX) (*pb.AssetInstance, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Delete(ctx, in)

	// return s.delete(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}
