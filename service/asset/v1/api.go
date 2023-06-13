package assetsrv

import (
	"context"

	"google.golang.org/grpc"

	"x-gwi/app/storage"
	srvpb "x-gwi/proto/serv/asset/v1"
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

func (s *Service) Create(ctx context.Context, in *srvpb.CreateRequest) (*srvpb.CreateResponse, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Create(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Create(ctx, in) //nolint:wrapcheck
}

func (s *Service) Get(ctx context.Context, in *srvpb.GetRequest) (*srvpb.GetResponse, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Get(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Get(ctx, in) //nolint:wrapcheck
}

func (s *Service) Gett(ctx context.Context, in *srvpb.GetRequest) (*srvpb.GetResponse, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Gett(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Gett(ctx, in) //nolint:wrapcheck
}

func (s *Service) Update(ctx context.Context, in *srvpb.UpdateRequest) (*srvpb.UpdateResponse, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Update(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *Service) Delete(ctx context.Context, in *srvpb.DeleteRequest) (*srvpb.DeleteResponse, error) {
	_, _ = (srvpb.UnimplementedAssetServiceServer{}).Delete(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedAssetServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}
