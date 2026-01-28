package models

import (
	"encoding/json"
	"time"
)

type AssetType string

const (
	AssetChart    AssetType = "chart"
	AssetInsight  AssetType = "insight"
	AssetAudience AssetType = "audience"
)

type Asset struct {
	ID        string          `db:"id"`
	Type      AssetType       `db:"type"`
	Title     *string         `db:"title"`
	Data      json.RawMessage `db:"data"` // JSONB
	CreatedAt time.Time       `db:"created_at"`
}
