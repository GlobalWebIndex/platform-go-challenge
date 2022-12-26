package app

import (
	"context"

	desc "ownify_api/pkg"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MicroserviceServer) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	userID, err := m.getUserIdFromToken(ctx)
	if err != nil {
		return nil, err
	}

	err = m.userService.DeleteUser(req.GetId(), userID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
