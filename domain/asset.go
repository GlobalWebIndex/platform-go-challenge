package domain

type AssetType string

type Asset struct {
	AssetId     uint      `json:"id" gorm:"primaryKey"`
	Type        AssetType `json:"type"`
	TypeId      uint      `json:"type_id"`
	Description string    `json:"description"`
}

const (
	AudienceAsset AssetType = "audience"
	ChartAsset    AssetType = "chart"
	InisightAsset AssetType = "insight"
)
