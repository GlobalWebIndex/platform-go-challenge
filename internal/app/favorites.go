package app

import (
	"context"
	"fmt"

	"gwi_api/internal/domain"
	desc "gwi_api/pkg"
)

func (m *MicroserviceServer) AddFavorites(ctx context.Context, req *desc.AddFavoriteRequest) (*desc.NetWorkResponse, error) {
	userId, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}
	err = m.favoriteService.AddFavorites(uint64(*userId), req.AssetId)
	if err != nil {
		return nil, err
	}
	return &desc.NetWorkResponse{
		Msg:     "successfully added",
		Success: true,
	}, nil
}

func (m *MicroserviceServer) GetFavorites(ctx context.Context, req *desc.GetFavoritesRequest) (*desc.NetWorkResponse, error) {
	userId, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}
	assets, err := m.favoriteService.GetFavorites(uint64(*userId), domain.PaginationParams{
		PageNumber: int(req.Pagination.Page),
		PageSize:   int(req.Pagination.PageSize),
	})
	if err != nil {
		fmt.Println("[Err]:", err)
		assets = []domain.Asset{}
	}
	return BuildRes[[]domain.Asset](assets, "successfully login", true)
}
