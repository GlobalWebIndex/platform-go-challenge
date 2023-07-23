package app

import (
	"context"

	"gwi_api/internal/domain"
	"gwi_api/internal/dto"
	desc "gwi_api/pkg"
)

func (m *MicroserviceServer) AddAsset(ctx context.Context, req *desc.AddAssetRequest) (*desc.NetWorkResponse, error) {
	userId, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}
	assetId, err := m.assetService.AddAsset(dto.AssetDto{
		CreatedBy:   uint64(*userId),
		Type:        domain.AssetType(req.Type),
		Description: req.Description,
		Data:        req.Data,
	})

	if err != nil {
		return nil, err
	}

	type AddAssetResponse struct {
		AssetId uint64
	}
	data := AddAssetResponse{
		AssetId: *assetId,
	}
	return BuildRes[AddAssetResponse](data, "successfully added!", true)
}

func (m *MicroserviceServer) GetAssets(ctx context.Context, req *desc.GetAssetRequest) (*desc.NetWorkResponse, error) {
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}

	assets, err := m.assetService.GetAsset(domain.PaginationParams{
		PageNumber: int(req.Page), PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	return BuildRes[[]domain.Asset](assets, "assets", true)
}

func (m *MicroserviceServer) DeleteAsset(ctx context.Context, req *desc.DeleteAssetRequest) (*desc.NetWorkResponse, error) {
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}

	err = m.assetService.DeleteAsset(req.Id)
	if err != nil {
		return nil, err
	}
	return &desc.NetWorkResponse{
		Msg:     "successfully deleted",
		Success: true,
	}, nil
}
