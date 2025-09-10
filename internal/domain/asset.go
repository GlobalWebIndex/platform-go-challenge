package domain

import (
	"strings"

	"github.com/google/uuid"
)

// AssetType defines the allowed categories of assets.
type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

// Asset is the domain representation of an asset.
type Asset struct {
	ID          uuid.UUID
	Type        AssetType
	Description string
	Payload     map[string]any
}

// EditDescription updates the asset description.
// It enforces that the description is not empty.
func (a *Asset) EditDescription(newDesc string) error {
	newDesc = strings.TrimSpace(newDesc)
	if newDesc == "" {
		return ErrEmptyDescription
	}
	a.Description = newDesc
	return nil
}
