package app

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type validatable interface {
	Validate() error
}

func GrpcInterceptor() grpc.ServerOption {
	grpcServerOptions := grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if v, ok := req.(validatable); ok {
			err := v.Validate()
			if err != nil {
				return nil, err
			}
		}
		resp, err = handler(ctx, req)
		return handler(ctx, req)
	})
	return grpcServerOptions
}

func HttpInterceptor() runtime.ServeMuxOption {
	httpServerOptions := runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		return nil
	})
	return httpServerOptions
}

func (m *MicroserviceServer) getUserIdFromToken(ctx context.Context) (string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	token := md.Get("Authorization")
	if token == nil {
		return "", status.Errorf(codes.PermissionDenied, "user isn't authorized")
	}

	userID, err := m.tokenManager.ValidateFirebase(token[0])
	if err != nil {
		return "", err
	}
	return *userID, nil
}
