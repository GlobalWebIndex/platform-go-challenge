// Package favourite contains the Crag model.
package favourite

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AssetType string

const (
	AssetChart    AssetType = "chart"
	AssetInsight  AssetType = "insight"
	AssetAudience AssetType = "audience"
)

func (t AssetType) IsValid() bool {
	switch t {
	case AssetChart, AssetInsight, AssetAudience:
		return true
	default:
		return false
	}
}

type Favorite struct {
	ID          uuid.UUID       `json:"id,omitempty"`
	Type        AssetType       `json:"type"`
	Description string          `json:"description"`
	Data        json.RawMessage `json:"data"`
	CreatedAt   time.Time       `json:"createdAt,omitempty"`
	UpdatedAt   time.Time       `json:"updatedAt,omitempty"`
}
