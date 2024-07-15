package models

// AssetType is a custom type for asset types
type AssetType string

// Define constants for Asset Types
const (
	ChartType    AssetType = "Chart"
	InsightType  AssetType = "Insight"
	AudienceType AssetType = "Audience"
)

// Asset interface to be implemented by all asset types
type Asset interface {
	GetID() int
	GetType() AssetType
	GetDescription() string
}
