package usersrv

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"x-gwi/app/logs"
	"x-gwi/app/storage"
	user_pb "x-gwi/proto/core/user/v1"
	srvpb "x-gwi/proto/serv/user/v1"
	"x-gwi/service/core/user"
)

type Service struct {
	srvpb.UnimplementedUserServiceServer
	user    *user.CoreUser
	storage *storage.ServiceStorage
}

func RegisterGRPC(grpcServer *grpc.Server, storage *storage.ServiceStorage) (*Service, error) {
	var err error

	s := &Service{ //nolint:exhaustruct
		storage: storage,
	}

	s.user, err = user.NewCore(s.storage)
	if err != nil {
		return nil, fmt.Errorf("user.NewCore: %w", err)
	}

	srvpb.RegisterUserServiceServer(grpcServer, s)

	return s, nil
}

func (s *Service) Create(ctx context.Context, in *srvpb.CreateRequest) (*srvpb.CreateResponse, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Create(ctx, in)

	log := logs.LogC2.With().Interface("in", in).Logger()
	t := time.Now()

	// validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	if in.GetUser() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "in == nil")
	}

	if err := in.GetUser().ValidateAll(); err != nil {
		log.Error().Err(err).Dur("dur", time.Since(t)).Send()

		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	wrk := &user.WorkUser{
		In:  in.GetUser(),
		Out: new(user_pb.UserInstance),
	}

	err := s.user.Create(ctx, wrk)
	if err != nil {
		return nil, fmt.Errorf("user.Create: %w", err)
	}

	out := &srvpb.CreateResponse{
		User: wrk.Out,
	}

	log.Debug().Interface("out", out).Dur("dur", time.Since(t)).Send()

	return out, nil
}

func (s *Service) Get(ctx context.Context, in *srvpb.GetRequest) (*srvpb.GetResponse, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Get(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Get(ctx, in) //nolint:wrapcheck
}

func (s *Service) Gett(ctx context.Context, in *srvpb.GetRequest) (*srvpb.GetResponse, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Gett(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Gett(ctx, in) //nolint:wrapcheck
}

func (s *Service) Update(ctx context.Context, in *srvpb.UpdateRequest) (*srvpb.UpdateResponse, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Update(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *Service) Delete(ctx context.Context, in *srvpb.DeleteRequest) (*srvpb.DeleteResponse, error) {
	_, _ = (srvpb.UnimplementedUserServiceServer{}).Delete(ctx, in)

	// return (srvpb.UnimplementedAuthServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedUserServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *Service) ListFavourites(in *srvpb.ListFavouritesRequest, stream srvpb.UserService_ListFavouritesServer) error {
	_ = (srvpb.UnimplementedUserServiceServer{}).ListFavourites(in, stream)

	// return s.list(in, stream)
	return (srvpb.UnimplementedUserServiceServer{}).ListFavourites(in, stream) //nolint:wrapcheck
}

func (s *Service) ListOpinions(in *srvpb.ListOpinionsRequest, stream srvpb.UserService_ListOpinionsServer) error {
	_ = (srvpb.UnimplementedUserServiceServer{}).ListOpinions(in, stream)

	// return s.list(in, stream)
	return (srvpb.UnimplementedUserServiceServer{}).ListOpinions(in, stream) //nolint:wrapcheck
}
