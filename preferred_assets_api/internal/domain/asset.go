package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type AssetBase struct {
	ID          string
	Type        AssetType
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Asset interface {
	GetID() string
	GetType() AssetType
	GetTitle() string
	GetDescription() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetID(id string)
	SetType(typ AssetType)
	SetTitle(title string)
	SetDescription(desc string)
	SetCreatedAt(t time.Time)
	SetUpdatedAt(t time.Time)
	Validate() error
}

// AssetType Parser
func ParseAssetType(s string) (AssetType, error) {
	switch strings.ToLower(s) {
	case "audience":
		return AssetTypeAudience, nil
	case "chart":
		return AssetTypeChart, nil
	case "insight":
		return AssetTypeInsight, nil
	default:
		return -99, fmt.Errorf("invalid asset type: %s", s)
	}
}

// String returns the string representation of the AssetType
func (at AssetType) String() string {
	switch at {
	case AssetTypeAudience:
		return "audience"
	case AssetTypeChart:
		return "chart"
	case AssetTypeInsight:
		return "insight"
	default:
		return "unknown"
	}
}

// Common Getter methods
func (a AssetBase) GetID() string           { return a.ID }
func (a AssetBase) GetType() AssetType      { return a.Type }
func (a AssetBase) GetTitle() string        { return a.Title }
func (a AssetBase) GetDescription() string  { return a.Description }
func (a AssetBase) GetCreatedAt() time.Time { return a.CreatedAt }
func (a AssetBase) GetUpdatedAt() time.Time { return a.UpdatedAt }

// Common Setter methods
func (a *AssetBase) SetID(id string)            { a.ID = id }
func (a *AssetBase) SetType(typ AssetType)      { a.Type = typ }
func (a *AssetBase) SetTitle(title string)      { a.Title = title }
func (a *AssetBase) SetDescription(desc string) { a.Description = desc }
func (a *AssetBase) SetCreatedAt(t time.Time)   { a.CreatedAt = t }
func (a *AssetBase) SetUpdatedAt(t time.Time)   { a.UpdatedAt = t }

// Domain validation
func (a AssetBase) Validate() error {
	if a.ID == "" {
		return errors.New("asset id is required")
	}

	if strings.TrimSpace(a.Title) == "" {
		return errors.New("asset title is required")
	}
	return nil
}
