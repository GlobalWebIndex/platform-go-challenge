package opinionsrv

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"x-gwi/app/storage"
	srvpb "x-gwi/proto/serv/opinion/v1"
	"x-gwi/service/core/opinion"
)

type Service struct {
	srvpb.UnimplementedOpinionServiceServer
	opinion *opinion.CoreOpinion
	storage *storage.ServiceStorage
}

func RegisterGRPC(grpcServer *grpc.Server, storage *storage.ServiceStorage) (*Service, error) {
	var err error

	s := &Service{ //nolint:exhaustruct
		storage: storage,
	}

	s.opinion, err = opinion.NewCore(s.storage)
	if err != nil {
		return nil, fmt.Errorf("user.NewCore: %w", err)
	}

	srvpb.RegisterOpinionServiceServer(grpcServer, s)

	return s, nil
}

func (s *Service) Create(ctx context.Context, in *srvpb.CreateRequest) (*srvpb.CreateResponse, error) {
	_, _ = (srvpb.UnimplementedOpinionServiceServer{}).Create(ctx, in)

	// return (srvpb.UnimplementedOpinionServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedOpinionServiceServer{}).Create(ctx, in) //nolint:wrapcheck
}

func (s *Service) Get(ctx context.Context, in *srvpb.GetRequest) (*srvpb.GetResponse, error) {
	_, _ = (srvpb.UnimplementedOpinionServiceServer{}).Get(ctx, in)

	// return (srvpb.UnimplementedOpinionServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedOpinionServiceServer{}).Get(ctx, in) //nolint:wrapcheck
}

func (s *Service) Gett(ctx context.Context, in *srvpb.GetRequest) (*srvpb.GetResponse, error) {
	_, _ = (srvpb.UnimplementedOpinionServiceServer{}).Gett(ctx, in)

	// return (srvpb.UnimplementedOpinionServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedOpinionServiceServer{}).Gett(ctx, in) //nolint:wrapcheck
}

func (s *Service) Update(ctx context.Context, in *srvpb.UpdateRequest) (*srvpb.UpdateResponse, error) {
	_, _ = (srvpb.UnimplementedOpinionServiceServer{}).Update(ctx, in)

	// return (srvpb.UnimplementedOpinionServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedOpinionServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *Service) Delete(ctx context.Context, in *srvpb.DeleteRequest) (*srvpb.DeleteResponse, error) {
	_, _ = (srvpb.UnimplementedOpinionServiceServer{}).Delete(ctx, in)

	// return (srvpb.UnimplementedOpinionServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (srvpb.UnimplementedOpinionServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *Service) List(in *srvpb.ListRequest, stream srvpb.OpinionService_ListServer) error {
	_ = (srvpb.UnimplementedOpinionServiceServer{}).List(in, stream)

	// return s.list(in, stream)
	return (srvpb.UnimplementedOpinionServiceServer{}).List(in, stream) //nolint:wrapcheck
}
