package service

type AssetType string

const (
	AssetChart    AssetType = "chart"
	AssetInsight  AssetType = "insight"
	AssetAudience AssetType = "audience"
)

type Asset struct {
	ID          int64     `json:"id"`
	Type        AssetType `json:"type"`
	Description string    `json:"description,omitempty"`
	Data        any       `json:"data"`
}

type Favourite struct {
	UserID            int64  `json:"user_id"`
	Asset             Asset  `json:"asset"`
	CustomDescription string `json:"custom_description,omitempty"`
}
