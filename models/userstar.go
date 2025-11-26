package models

import (
	"database/sql/driver"
	"fmt"
)

// AssetType represents the type of asset that can be starred
type AssetType string

// Asset type constants
const (
	AssetTypeAudience AssetType = "Audience"
	AssetTypeChart    AssetType = "Chart"
	AssetTypeInsight  AssetType = "Insight"
)

// IsValid checks if the AssetType is one of the valid types
func (at AssetType) IsValid() bool {
	switch at {
	case AssetTypeAudience, AssetTypeChart, AssetTypeInsight:
		return true
	}
	return false
}

// String returns the string representation of AssetType
func (at AssetType) String() string {
	return string(at)
}

// Value implements the driver.Valuer interface for database serialization
func (at AssetType) Value() (driver.Value, error) {
	if !at.IsValid() {
		return nil, fmt.Errorf("invalid asset type: %s", at)
	}
	return string(at), nil
}

// Scan implements the sql.Scanner interface for database deserialization
func (at *AssetType) Scan(value any) error {
	if value == nil {
		return fmt.Errorf("asset type cannot be null")
	}

	str, ok := value.(string)
	if !ok {
		// Handle []byte as well (some drivers return bytes)
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("asset type must be a string, got %T", value)
		}
		str = string(bytes)
	}

	*at = AssetType(str)
	if !at.IsValid() {
		return fmt.Errorf("invalid asset type: %s", str)
	}
	return nil
}

type UserStar struct {
	ID      uint      `json:"id" gorm:"primaryKey"`
	UserID  uint      `json:"userid"`
	Type    AssetType `json:"type"`
	AssetID uint      `json:"assetid"`
}
