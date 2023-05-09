package domain

type AssetType string

const (
	ChartType    AssetType = "chart"
	AudienceType AssetType = "audience"
	InsightType  AssetType = "insight"
)

type Asset interface {
	GetDescription() string
	GetType() string
}
