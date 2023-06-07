package app

import (
	"context"
	"fmt"

	"ownify_api/internal/constants"
	"ownify_api/internal/utils"
	desc "ownify_api/pkg"
)

func (m *MicroserviceServer) GetOwnership(ctx context.Context, req *desc.GetOwnershipRequest) (*desc.NetWorkResponse, error) {
	err := utils.IsPubKey(req.WalletAddress)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}
	b, u, err := m.ownershipService.GetOwnerShip(req.WalletAddress)
	if err != nil {
		return nil, err
	}
	if b.Email != "" {
		return BuildRes(b, "business", true)
	}
	if u.ValidMainInfo() == nil {
		return BuildRes(u, "user", true)
	}
	return nil, fmt.Errorf("[ERR] Did not find ownership: %s", req.WalletAddress)

}
