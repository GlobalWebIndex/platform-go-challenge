package apiv2asset

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"x-gwi/app/storage"
	sharepb "x-gwi/proto/core/_share/v1"
	assetpb "x-gwi/proto/core/asset/v1"
	pbsrv "x-gwi/proto/serv/asset/v2"
	"x-gwi/service/core/asset"
)

type ServiceAPI struct {
	pbsrv.UnimplementedAssetServiceServer
	assetCore *asset.CoreAsset
	// storage *storage.ServiceStorage
}

func RegisterGRPC(grpcServer *grpc.Server, storage *storage.ServiceStorage) (*ServiceAPI, error) {
	var err error

	s := &ServiceAPI{ //nolint:exhaustruct
		// storage: storage,
	}

	s.assetCore, err = asset.NewCore(storage)
	if err != nil {
		return nil, fmt.Errorf("asset.NewCore: %w", err)
	}

	pbsrv.RegisterAssetServiceServer(grpcServer, s)

	return s, nil
}

func (s *ServiceAPI) Create(ctx context.Context, in *assetpb.AssetCore) (*assetpb.AssetCore, error) {
	_, _ = (pbsrv.UnimplementedAssetServiceServer{}).Create(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process by core
	err := s.assetCore.Create(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "assetCore.Create: %v", err)
	}

	return in, nil
}

func (s *ServiceAPI) Get(ctx context.Context, in *sharepb.ShareQID) (*assetpb.AssetCore, error) {
	_, _ = (pbsrv.UnimplementedAssetServiceServer{}).Get(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &assetpb.AssetCore{ //nolint:exhaustruct
		Qid: in,
	}

	// 3. process by core
	err := s.assetCore.Get(ctx, out)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "assetCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Gett(ctx context.Context, in *sharepb.ShareQID) (*assetpb.AssetCore, error) {
	_, _ = (pbsrv.UnimplementedAssetServiceServer{}).Gett(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &assetpb.AssetCore{ //nolint:exhaustruct
		Qid: in,
	}

	// 3. process by core
	err := s.assetCore.Get(ctx, out)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "assetCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Update(ctx context.Context, in *assetpb.AssetCore) (*assetpb.AssetCore, error) {
	_, _ = (pbsrv.UnimplementedAssetServiceServer{}).Update(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process by core
	err := s.assetCore.Update(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "assetCore.Update: %v", err)
	}

	return in, nil
}

func (s *ServiceAPI) Delete(ctx context.Context, in *sharepb.ShareQID) (*assetpb.AssetCore, error) {
	_, _ = (pbsrv.UnimplementedAssetServiceServer{}).Delete(ctx, in)

	// return s.delete(ctx, in)
	return (pbsrv.UnimplementedAssetServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}
