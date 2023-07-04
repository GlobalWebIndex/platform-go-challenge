package apiv2opinion

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"x-gwi/app/storage"
	sharepb "x-gwi/proto/core/_share/v1"
	opinionpb "x-gwi/proto/core/opinion/v1"
	pbsrv "x-gwi/proto/serv/opinion/v2"
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

func (s *ServiceAPI) Create(ctx context.Context, in *opinionpb.OpinionCore) (*opinionpb.OpinionCore, error) {
	_, _ = (pbsrv.UnimplementedOpinionServiceServer{}).Create(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process by core
	err := s.opinionCore.Create(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "opinionCore.Create: %v", err)
	}

	return in, nil
}

func (s *ServiceAPI) Get(ctx context.Context, in *sharepb.ShareQID) (*opinionpb.OpinionCore, error) {
	_, _ = (pbsrv.UnimplementedOpinionServiceServer{}).Get(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &opinionpb.OpinionCore{ //nolint:exhaustruct
		Qid: in,
	}

	// 3. process by core
	err := s.opinionCore.Get(ctx, out)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "opinionCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Gett(ctx context.Context, in *sharepb.ShareQID) (*opinionpb.OpinionCore, error) {
	_, _ = (pbsrv.UnimplementedOpinionServiceServer{}).Gett(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. adapt output to fill
	out := &opinionpb.OpinionCore{ //nolint:exhaustruct
		Qid: in,
	}

	// 3. process by core
	err := s.opinionCore.Get(ctx, out)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "opinionCore.Get: %v", err)
	}

	return out, nil
}

func (s *ServiceAPI) Update(ctx context.Context, in *opinionpb.OpinionCore) (*opinionpb.OpinionCore, error) {
	_, _ = (pbsrv.UnimplementedOpinionServiceServer{}).Update(ctx, in)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process by core
	err := s.opinionCore.Update(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "opinionCore.Update: %v", err)
	}

	return in, nil
}

func (s *ServiceAPI) Delete(ctx context.Context, in *sharepb.ShareQID) (*opinionpb.OpinionCore, error) {
	_, _ = (pbsrv.UnimplementedOpinionServiceServer{}).Delete(ctx, in)

	// return s.delete(ctx, in)
	return (pbsrv.UnimplementedOpinionServiceServer{}).Delete(ctx, in) //nolint:wrapcheck
}

func (s *ServiceAPI) List(in *opinionpb.OpinionCore, stream pbsrv.OpinionService_ListServer) error {
	_ = (pbsrv.UnimplementedOpinionServiceServer{}).List(in, stream)
	// 1. validate input
	if err := in.ValidateAll(); err != nil {
		return status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// 2. process stream by core
	err := s.opinionCore.ListV2(in, stream)
	if err != nil {
		return status.Errorf(codes.Unknown, "opinionCore.ListV2: %v", err)
	}

	return nil
}
