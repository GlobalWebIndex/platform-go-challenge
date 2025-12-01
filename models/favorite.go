package models

type Favorite struct {
	UserID int `json:"user_id"`
	AssetID int `json:"asset_id"`
	AssetType AssetType `json:"asset_type"`
}