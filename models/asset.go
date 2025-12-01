package models

type Asset interface {
	GetID() int
	GetType() AssetType
	GetDescription() string
	SetDescription(desc string)
}

type AssetType string

const (
	AssetTypeChart AssetType = "chart"
	AssetTypeInsight AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)


