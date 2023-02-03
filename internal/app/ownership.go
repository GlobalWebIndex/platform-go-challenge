package app

import (
	"context"

	"ownify_api/internal/utils"
	desc "ownify_api/pkg"
)

func (m *MicroserviceServer) GetOwnership(ctx context.Context, req *desc.GetOwnershipRequest) (*desc.NetWorkResponse, error) {
	err := utils.IsPubKey(req.WalletAddress)
	if err != nil {
		return nil, err
	}
	b, u, err := m.ownershipService.GetOwnerShip(req.WalletAddress)
	if err != nil {
		return nil, err
	}
	if b.Valid() {
		return BuildRes(b, "business", true)
	}
	return BuildRes(u, "user", true)
}
