package model

import "encoding/json"

type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

type Asset struct {
	ID          string          `json:"id"`
	Type        AssetType       `json:"type"`
	Description string          `json:"description,omitempty"`
	Payload     json.RawMessage `json:"payload"` // type-specific
}
