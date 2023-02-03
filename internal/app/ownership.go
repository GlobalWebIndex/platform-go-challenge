package app

import (
	"context"

	desc "ownify_api/pkg"
)

func (m *MicroserviceServer) GetOwnership(ctx context.Context, req *desc.GetOwnershipRequest) (*desc.NetWorkResponse, error) {

	b, u, err := m.ownershipService.GetOwnerShip(req.WalletAddress)
	if err != nil {
		return nil, err
	}
	if b.Valid() {
		return BuildRes(b, "business", true)
	}
	return BuildRes(u, "user", true)
}
