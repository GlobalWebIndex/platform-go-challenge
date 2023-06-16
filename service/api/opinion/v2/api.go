package opinionsrv2

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"x-gwi/app/storage"
	pb_idx "x-gwi/proto/core/_store/v1"
	pb "x-gwi/proto/core/opinion/v1"
	srvpb "x-gwi/proto/serv/opinion/v2"
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

func (s *Service) Create(ctx context.Context, in *pb.OpinionAsset) (*pb.OpinionAsset, error) {
	_, _ = (srvpb.UnimplementedOpinionServiceServer{}).Create(ctx, in)

	// return s.create(ctx, in)
	return (srvpb.UnimplementedOpinionServiceServer{}).Create(ctx, in) //nolint:wrapcheck
}

func (s *Service) Get(ctx context.Context, in *pb_idx.StoreIDX) (*pb.OpinionAsset, error) {
	_, _ = (srvpb.UnimplementedOpinionServiceServer{}).Get(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedOpinionServiceServer{}).Get(ctx, in) //nolint:wrapcheck
}

func (s *Service) Gett(ctx context.Context, in *pb_idx.StoreIDX) (*pb.OpinionAsset, error) {
	_, _ = (srvpb.UnimplementedOpinionServiceServer{}).Gett(ctx, in)

	// return s.get(ctx, in)
	return (srvpb.UnimplementedOpinionServiceServer{}).Gett(ctx, in) //nolint:wrapcheck
}

func (s *Service) Update(ctx context.Context, in *pb.OpinionAsset) (*pb.OpinionAsset, error) {
	_, _ = (srvpb.UnimplementedOpinionServiceServer{}).Update(ctx, in)

	// return s.update(ctx, in)
	return (srvpb.UnimplementedOpinionServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *Service) Delete(ctx context.Context, in *pb_idx.StoreIDX) (*pb.OpinionAsset, error) {
	_, _ = (srvpb.UnimplementedOpinionServiceServer{}).Delete(ctx, in)

	// return s.delete(ctx, in)
	return (srvpb.UnimplementedOpinionServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *Service) List(in *pb_idx.StoreIDX, stream srvpb.OpinionService_ListServer) error {
	_ = (srvpb.UnimplementedOpinionServiceServer{}).List(in, stream)

	// return s.list(in, stream)
	return (srvpb.UnimplementedOpinionServiceServer{}).List(in, stream) //nolint:wrapcheck
}
