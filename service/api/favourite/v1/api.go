package favouritesrv

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"x-gwi/app/storage"
	srvpb "x-gwi/proto/serv/favourite/v1"
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

func (s *Service) Create(ctx context.Context, in *srvpb.CreateRequest) (*srvpb.CreateResponse, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Create(ctx, in)

	// return (srvpb.UnimplementedFavouriteServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Create(ctx, in) //nolint:wrapcheck
}

func (s *Service) Get(ctx context.Context, in *srvpb.GetRequest) (*srvpb.GetResponse, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Get(ctx, in)

	// return (srvpb.UnimplementedFavouriteServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Get(ctx, in) //nolint:wrapcheck
}

func (s *Service) Gett(ctx context.Context, in *srvpb.GetRequest) (*srvpb.GetResponse, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Gett(ctx, in)

	// return (srvpb.UnimplementedFavouriteServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Gett(ctx, in) //nolint:wrapcheck
}

func (s *Service) Update(ctx context.Context, in *srvpb.UpdateRequest) (*srvpb.UpdateResponse, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Update(ctx, in)

	// return (srvpb.UnimplementedFavouriteServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *Service) Delete(ctx context.Context, in *srvpb.DeleteRequest) (*srvpb.DeleteResponse, error) {
	_, _ = (srvpb.UnimplementedFavouriteServiceServer{}).Delete(ctx, in)

	// return (srvpb.UnimplementedFavouriteServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedFavouriteServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *Service) List(in *srvpb.ListRequest, stream srvpb.FavouriteService_ListServer) error {
	_ = (srvpb.UnimplementedFavouriteServiceServer{}).List(in, stream)

	// return s.list(in, stream)
	return (srvpb.UnimplementedFavouriteServiceServer{}).List(in, stream) //nolint:wrapcheck
}
