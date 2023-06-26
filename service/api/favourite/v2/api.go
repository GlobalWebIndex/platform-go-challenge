package apiv2favourite

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"x-gwi/app/storage"
	sharepb "x-gwi/proto/core/_share/v1"
	favouritepb "x-gwi/proto/core/favourite/v1"
	pbsrv "x-gwi/proto/serv/favourite/v2"
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

func (s *ServiceAPI) Create(ctx context.Context, in *favouritepb.FavouriteCore) (*favouritepb.FavouriteCore, error) {
	_, _ = (pbsrv.UnimplementedFavouriteServiceServer{}).Create(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process by core
	err := s.favouriteCore.Create(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "favouriteCore.Create: %v", err)
	}

	return in, nil
}

func (s *ServiceAPI) Get(ctx context.Context, in *sharepb.ShareQID) (*favouritepb.FavouriteCore, error) {
	_, _ = (pbsrv.UnimplementedFavouriteServiceServer{}).Get(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &favouritepb.FavouriteCore{ //nolint:exhaustruct
		Qid: in,
	}

	// 3. process by core
	err := s.favouriteCore.Get(ctx, out)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "favouriteCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Gett(ctx context.Context, in *sharepb.ShareQID) (*favouritepb.FavouriteCore, error) {
	_, _ = (pbsrv.UnimplementedFavouriteServiceServer{}).Gett(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &favouritepb.FavouriteCore{ //nolint:exhaustruct
		Qid: in,
	}

	// 3. process by core
	err := s.favouriteCore.Get(ctx, out)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "favouriteCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Update(ctx context.Context, in *favouritepb.FavouriteCore) (*favouritepb.FavouriteCore, error) {
	_, _ = (pbsrv.UnimplementedFavouriteServiceServer{}).Update(ctx, in)

	// return s.update(ctx, in)
	return (pbsrv.UnimplementedFavouriteServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *ServiceAPI) Delete(ctx context.Context, in *sharepb.ShareQID) (*favouritepb.FavouriteCore, error) {
	_, _ = (pbsrv.UnimplementedFavouriteServiceServer{}).Delete(ctx, in)

	// return s.delete(ctx, in)
	return (pbsrv.UnimplementedFavouriteServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *ServiceAPI) List(in *favouritepb.FavouriteCore, stream pbsrv.FavouriteService_ListServer) error {
	_ = (pbsrv.UnimplementedFavouriteServiceServer{}).List(in, stream)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process stream by core
	err := s.favouriteCore.ListV2(in, stream)
	if err != nil {
		return status.Errorf(codes.Unknown, "favouriteCore.ListV2: %v", err)
	}

	return nil
}
