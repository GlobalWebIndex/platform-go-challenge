package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	desc "ownify_api/pkg"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
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

// func HttpInterceptor() runtime.ServeMuxOption {
// 	httpServerOptions := runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
// 		return nil
// 	})
// 	return httpServerOptions
// }

func HttpInterceptor() runtime.ServeMuxOption {
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "http://localhost:3000"
	return runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		return metadata.New(headers)
	})
}

func (m *MicroserviceServer) getUserIdFromToken(ctx context.Context) (string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	token := md.Get("Authorization")
	if token == nil {
		return "", status.Errorf(codes.PermissionDenied, "user isn't authorized")
	}
	return token[0], nil
}

func (m *MicroserviceServer) TokenInterceptor(ctx context.Context) (*string, error) {
	// validate token.
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println(md.Get("test"))
	token, err := m.getUserIdFromToken(ctx)
	if err != nil {
		log.Println("user isn't authorized")
		return nil, err
	}
	uid, err := m.tokenManager.ValidateFirebase(token)
	if err != nil {
		log.Println("user isn't authorized")
		return nil, err
	}

	return uid, nil
}

func BuildRes[T any](data T, message string, isSuccess bool) (*desc.NetWorkResponse, error) {
	rawData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	s := &structpb.Struct{}
	_ = protojson.Unmarshal(rawData, s)
	return &desc.NetWorkResponse{
		Msg:     message,
		Success: isSuccess,
		Data:    s,
	}, nil
}
