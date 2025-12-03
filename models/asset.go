package models

import "time"

type AssetType string

// constants of type AssetType
const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

type Asset struct {
	ID          string    `json:"id"`
	Type        AssetType `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Only one of these will be populated based on Type.
	Chart    *Chart    `json:"chart,omitempty"`
	Insight  *Insight  `json:"insight,omitempty"`
	Audience *Audience `json:"audience,omitempty"`
}
