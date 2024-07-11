package services

import (
	"fmt"
	"gwi-test/domain"
	"gwi-test/ports"
)

type AssetsHandler struct {
	Db ports.Db
}

func NewAssetHandler(db ports.Db) *AssetsHandler {
	return &AssetsHandler{
		Db: db,
	}
}

func (a *AssetsHandler) GetAssetsById(userId uint) ([]domain.Asset, error) {
	assets, err := a.Db.GetAssetsById(userId)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (a *AssetsHandler) GetAssetsObjects(userId uint) (*domain.AssetsResponse, error) {
	resp := &domain.AssetsResponse{}
	assets, err := a.Db.GetAssetsById(userId)
	if err != nil {
		return nil, err
	}

	for _, asset := range assets {
		fmt.Println("[LOG] asset: ", asset)
		switch asset.Type {
		case domain.AudienceAsset:
			audience, err := a.Db.GetAudience(asset.TypeId)
			if err != nil {
				return nil, err
			}
			resp.Audiences = append(resp.Audiences, *audience)
		case domain.ChartAsset:
			chart, err := a.Db.GetChart(asset.TypeId)
			if err != nil {
				return nil, err
			}
			resp.Charts = append(resp.Charts, *chart)
		case domain.InisightAsset:
			inisight, err := a.Db.GetInsight(asset.TypeId)
			if err != nil {
				return nil, err
			}
			resp.Insights = append(resp.Insights, *inisight)
		}
	}

	return resp, nil
}

func (a *AssetsHandler) DeleteAsset(assetId uint) error {
	return a.Db.DeleteAsset(assetId)
}

func (a *AssetsHandler) AddAsset(asset *domain.Asset) error {
	return a.Db.AddAsset(asset)
}
func (a *AssetsHandler) UpdateAsset(asset *domain.Asset) error {
	return a.Db.UpdateAsset(asset)
}

// I should create a different file for each service but for now I will just add the functions here

func (a *AssetsHandler) GetChart(id uint) (*domain.Chart, error) {
	return a.Db.GetChart(id)
}

func (a *AssetsHandler) GetAudience(id uint) (*domain.Audience, error) {
	return a.Db.GetAudience(id)
}

func (a *AssetsHandler) GetInisight(id uint) (*domain.Insight, error) {
	return a.Db.GetInsight(id)
}
