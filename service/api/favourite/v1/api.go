package apiv1favourite

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"x-gwi/app/storage"
	favouritepb "x-gwi/proto/core/favourite/v1"
	pbsrv "x-gwi/proto/serv/favourite/v1"
	"x-gwi/service/core/favourite"
)

type ServiceAPI struct {
	pbsrv.UnimplementedFavouriteServiceServer
	favouriteCore *favourite.CoreFavourite
	// storage   *storage.ServiceStorage
}

func RegisterGRPC(grpcServer *grpc.Server, storage *storage.ServiceStorage) (*ServiceAPI, error) {
	var err error

	s := &ServiceAPI{ //nolint:exhaustruct
		// storage: storage,
	}

	s.favouriteCore, err = favourite.NewCore(storage)
	if err != nil {
		return nil, fmt.Errorf("favourite.NewCore: %w", err)
	}

	pbsrv.RegisterFavouriteServiceServer(grpcServer, s)

	return s, nil
}

func (s *ServiceAPI) Create(ctx context.Context, in *pbsrv.CreateRequest) (*pbsrv.CreateResponse, error) {
	_, _ = (pbsrv.UnimplementedFavouriteServiceServer{}).Create(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process by core
	err := s.favouriteCore.Create(ctx, in.GetFavourite())
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "favouriteCore.Create: %v", err)
	}

	// 3. adapt output
	out := &pbsrv.CreateResponse{
		Favourite: in.GetFavourite(),
	}

	return out, nil
}

func (s *ServiceAPI) Get(ctx context.Context, in *pbsrv.GetRequest) (*pbsrv.GetResponse, error) {
	_, _ = (pbsrv.UnimplementedFavouriteServiceServer{}).Get(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &pbsrv.GetResponse{
		Favourite: &favouritepb.FavouriteCore{ //nolint:exhaustruct
			Qid: in.GetQid(),
		},
	}

	// 3. process by core
	err := s.favouriteCore.Get(ctx, out.Favourite)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "favouriteCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Gett(ctx context.Context, in *pbsrv.GetRequest) (*pbsrv.GetResponse, error) {
	_, _ = (pbsrv.UnimplementedFavouriteServiceServer{}).Gett(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &pbsrv.GetResponse{
		Favourite: &favouritepb.FavouriteCore{ //nolint:exhaustruct
			Qid: in.GetQid(),
		},
	}

	// 3. process by core
	err := s.favouriteCore.Get(ctx, out.Favourite)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "favouriteCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Update(ctx context.Context, in *pbsrv.UpdateRequest) (*pbsrv.UpdateResponse, error) {
	_, _ = (pbsrv.UnimplementedFavouriteServiceServer{}).Update(ctx, in)

	// return (pbsrv.UnimplementedFavouriteServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (pbsrv.UnimplementedFavouriteServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *ServiceAPI) Delete(ctx context.Context, in *pbsrv.DeleteRequest) (*pbsrv.DeleteResponse, error) {
	_, _ = (pbsrv.UnimplementedFavouriteServiceServer{}).Delete(ctx, in)

	// return (pbsrv.UnimplementedFavouriteServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (pbsrv.UnimplementedFavouriteServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *ServiceAPI) List(in *pbsrv.ListRequest, stream pbsrv.FavouriteService_ListServer) error {
	_ = (pbsrv.UnimplementedFavouriteServiceServer{}).List(in, stream)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process stream by core
	err := s.favouriteCore.ListV1(in.GetFavourite(), stream)
	if err != nil {
		return status.Errorf(codes.Unknown, "favouriteCore.ListV1: %v", err)
	}

	return nil
}
