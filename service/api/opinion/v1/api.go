package apiv1opinion

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"x-gwi/app/storage"
	opinionpb "x-gwi/proto/core/opinion/v1"
	pbsrv "x-gwi/proto/serv/opinion/v1"
	"x-gwi/service/core/opinion"
)

type ServiceAPI struct {
	pbsrv.UnimplementedOpinionServiceServer
	opinionCore *opinion.CoreOpinion
	// storage *storage.ServiceStorage
}

func RegisterGRPC(grpcServer *grpc.Server, storage *storage.ServiceStorage) (*ServiceAPI, error) {
	var err error

	s := &ServiceAPI{ //nolint:exhaustruct
		// storage: storage,
	}

	s.opinionCore, err = opinion.NewCore(storage)
	if err != nil {
		return nil, fmt.Errorf("user.NewCore: %w", err)
	}

	pbsrv.RegisterOpinionServiceServer(grpcServer, s)

	return s, nil
}

func (s *ServiceAPI) Create(ctx context.Context, in *pbsrv.CreateRequest) (*pbsrv.CreateResponse, error) {
	_, _ = (pbsrv.UnimplementedOpinionServiceServer{}).Create(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process by core
	err := s.opinionCore.Create(ctx, in.GetOpinion())
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "opinionCore.Create: %v", err)
	}

	// 3. adapt output
	out := &pbsrv.CreateResponse{
		Opinion: in.GetOpinion(),
	}

	return out, nil
}

func (s *ServiceAPI) Get(ctx context.Context, in *pbsrv.GetRequest) (*pbsrv.GetResponse, error) {
	_, _ = (pbsrv.UnimplementedOpinionServiceServer{}).Get(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &pbsrv.GetResponse{
		Opinion: &opinionpb.OpinionCore{ //nolint:exhaustruct
			Qid: in.GetQid(),
		},
	}

	// 3. process by core
	err := s.opinionCore.Get(ctx, out.Opinion)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "opinionCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Gett(ctx context.Context, in *pbsrv.GetRequest) (*pbsrv.GetResponse, error) {
	_, _ = (pbsrv.UnimplementedOpinionServiceServer{}).Gett(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &pbsrv.GetResponse{
		Opinion: &opinionpb.OpinionCore{ //nolint:exhaustruct
			Qid: in.GetQid(),
		},
	}

	// 3. process by core
	err := s.opinionCore.Get(ctx, out.Opinion)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "opinionCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Update(ctx context.Context, in *pbsrv.UpdateRequest) (*pbsrv.UpdateResponse, error) {
	_, _ = (pbsrv.UnimplementedOpinionServiceServer{}).Update(ctx, in)

	// return (pbsrv.UnimplementedOpinionServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (pbsrv.UnimplementedOpinionServiceServer{}).Update(ctx, in) //nolint:wrapcheck
}

func (s *ServiceAPI) Delete(ctx context.Context, in *pbsrv.DeleteRequest) (*pbsrv.DeleteResponse, error) {
	_, _ = (pbsrv.UnimplementedOpinionServiceServer{}).Delete(ctx, in)

	// return (pbsrv.UnimplementedOpinionServiceServer{}).Create(ctx, in)
	// return s.createJWT(ctx, in)
	return (pbsrv.UnimplementedOpinionServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *ServiceAPI) List(in *pbsrv.ListRequest, stream pbsrv.OpinionService_ListServer) error {
	_ = (pbsrv.UnimplementedOpinionServiceServer{}).List(in, stream)

	// return s.list(in, stream)
	return (pbsrv.UnimplementedOpinionServiceServer{}).List(in, stream) //nolint:wrapcheck
}
