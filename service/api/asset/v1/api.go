package apiv1asset

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"x-gwi/app/storage"
	assetpb "x-gwi/proto/core/asset/v1"
	pbsrv "x-gwi/proto/serv/asset/v1"
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

func (s *ServiceAPI) Create(ctx context.Context, in *pbsrv.CreateRequest) (*pbsrv.CreateResponse, error) {
	_, _ = (pbsrv.UnimplementedAssetServiceServer{}).Create(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process by core
	err := s.assetCore.Create(ctx, in.GetAsset())
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "assetCore.Create: %v", err)
	}

	// 3. adapt output
	out := &pbsrv.CreateResponse{
		Asset: in.GetAsset(),
	}

	return out, nil
}

func (s *ServiceAPI) Get(ctx context.Context, in *pbsrv.GetRequest) (*pbsrv.GetResponse, error) {
	_, _ = (pbsrv.UnimplementedAssetServiceServer{}).Get(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &pbsrv.GetResponse{
		Asset: &assetpb.AssetCore{ //nolint:exhaustruct
			Qid: in.GetQid(),
		},
	}

	// 3. process by core
	err := s.assetCore.Get(ctx, out.Asset)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "assetCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Gett(ctx context.Context, in *pbsrv.GetRequest) (*pbsrv.GetResponse, error) {
	_, _ = (pbsrv.UnimplementedAssetServiceServer{}).Gett(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &pbsrv.GetResponse{
		Asset: &assetpb.AssetCore{ //nolint:exhaustruct
			Qid: in.GetQid(),
		},
	}

	// 3. process by core
	err := s.assetCore.Get(ctx, out.Asset)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "assetCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Update(ctx context.Context, in *pbsrv.UpdateRequest) (*pbsrv.UpdateResponse, error) {
	_, _ = (pbsrv.UnimplementedAssetServiceServer{}).Update(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process by core
	err := s.assetCore.Update(ctx, in.Asset)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "assetCore.Update: %v", err)
	}

	// 3. adapt output
	out := &pbsrv.UpdateResponse{
		Asset: in.Asset,
	}

	return out, nil
}

func (s *ServiceAPI) Delete(ctx context.Context, in *pbsrv.DeleteRequest) (*pbsrv.DeleteResponse, error) {
	_, _ = (pbsrv.UnimplementedAssetServiceServer{}).Delete(ctx, in)

	// return (pbsrv.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (pbsrv.UnimplementedAssetServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}
