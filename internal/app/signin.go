package app

import (
	"context"
	"fmt"

	desc "ownify_api/pkg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MicroserviceServer) Login(ctx context.Context, req *desc.SignInRequest) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println(md.Get("test"))

	token, err := m.authService.SignIn(req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	err = grpc.SendHeader(ctx, metadata.New(map[string]string{
		"Token": *token,
	}))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
