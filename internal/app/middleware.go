package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"ownify_api/internal/utils"
	desc "ownify_api/pkg"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

var limiter = rate.NewLimiter(10, 1)

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

	return runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		meta := map[string]string{}
		if req.Header.Get("Authorization") != "" {
			meta["Authorization"] = req.Header.Get("Authorization")
		}
		if req.Header.Get("x-api-key") != "" {
			meta["x-api-key"] = req.Header.Get("x-api-key")
		}
		if len(meta) == 0 {
			return nil
		}
		return metadata.New(meta)
	})
}

func (m *MicroserviceServer) getUserIdFromToken(ctx context.Context) (string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	token := md.Get("Authorization")

	if token == nil {
		return "", status.Errorf(codes.PermissionDenied, "user isn't authorized")
	}
	removeBearerFromToken := strings.Split(token[0], " ")
	return removeBearerFromToken[1], nil
}
func (m *MicroserviceServer) getUserInfoFromApiKey(ctx context.Context) (string, string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	apiKey := md.Get("x-api-key")
	if apiKey == nil {
		return "", "", status.Errorf(codes.PermissionDenied, "user isn't authorized")
	}
	contents := strings.Split(apiKey[0], "-")
	if len(contents) != 4 {
		return "", "", status.Errorf(codes.PermissionDenied, "invalid api key")
	}
	email, err := utils.Decrypt(contents[1], contents[3])
	if err != nil {
		return "", "", status.Errorf(codes.PermissionDenied, "invalid api key")
	}
	userId, err := utils.Decrypt(contents[1], contents[3])
	if err != nil {
		return "", "", status.Errorf(codes.PermissionDenied, "invalid api key")
	}
	return email, userId, nil
}

func (m *MicroserviceServer) TokenInterceptor(ctx context.Context) (*string, error) {
	// check rate limitation
	if !limiter.Allow() {
		return nil, status.Errorf(codes.ResourceExhausted, "Too many requests")
	}
	//validate token.
	token, err := m.getUserIdFromToken(ctx)
	if err != nil {
		_, userId, err := m.getUserInfoFromApiKey(ctx)
		if err != nil {
			log.Println("user isn't authorized")
			return nil, err
		}
		return &userId, err
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
