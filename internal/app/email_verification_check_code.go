package app

import (
	"context"

	desc "ownify_api/pkg"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MicroserviceServer) EmailVerificationCheckCode(ctx context.Context, req *desc.EmailVerificationCheckCodeRequest) (*emptypb.Empty, error) {
	userID, err := m.getUserIdFromToken(ctx)
	if err != nil {
		return nil, err
	}

	err = m.emailService.CheckCode(userID, req.GetCode())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
