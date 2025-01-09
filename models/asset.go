package models

// AssetType represents the type of asset
// @Description AssetType represents the type of asset
type AssetType string

const (
	ChartAsset    AssetType = "chart"
	InsightAsset  AssetType = "insight"
	AudienceAsset AssetType = "audience"
)

// Asset represents an asset with a type, description, and data
// @Description Asset represents an asset with a type, description, and data
type Asset struct {
	ID          string      `json:"id" example:"1"`
	Type        AssetType   `json:"type" example:"chart"`
	Description string      `json:"description" example:"Sample Chart"`
	Data        interface{} `json:"data"`
}
