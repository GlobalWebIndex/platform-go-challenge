package ports

import "gwi-test/domain"

// This file is for the "ports" of the app
// Basically the ports are the interfaces that the app will use to interact with the outside world

type Server interface {
	NewServer() *Server
	Initialize()
}

type AssetHandler interface {
	GetAssetsById(userId uint) ([]domain.Asset, error)
	GetAssetsObjects(userId uint) (*domain.AssetsResponse, error)
	DeleteAsset(assetId uint) error
	AddAsset(asset *domain.Asset) error
	UpdateAsset(asset *domain.Asset) error
}

type Db interface {
	AddAsset(asset *domain.Asset) error
	DeleteAsset(assetId uint) error
	UpdateAsset(asset *domain.Asset) error
	GetAssetsById(userId uint) ([]domain.Asset, error)

	GetAudience(id uint) (*domain.Audience, error)
	GetInsight(id uint) (*domain.Insight, error)
	GetChart(id uint) (*domain.Chart, error)
}
