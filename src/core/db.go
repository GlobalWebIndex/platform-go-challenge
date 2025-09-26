package core

import "time"

type AssetRepository interface {
	Upsert(a AnyAsset) error
	GetByID(id AssetID) (AnyAsset, error)
	UpdateDescription(id AssetID, desc string, now time.Time) error
	ListByIDsPreserveOrder(ids []AssetID) ([]AnyAsset, error)
	Exists(id AssetID) (bool, error)
}

type FavoriteRepository interface {
	Add(userID UserID, assetID AssetID, now time.Time) error
	Remove(userID UserID, assetID AssetID) error
	ListByCursor(userID UserID, limit int, after *FavoriteEdge, typ *AssetType) ([]FavoriteEdge, error)
	SetOverride(userID UserID, assetID AssetID, override map[string]any, now time.Time) error
}