package favorites

import (
	"errors"
	"fmt"
)

// ErrValidation is the sentinel error for validation failures.
var ErrValidation = errors.New("validation failed")

// AssetType defines the supported asset types.
type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

// Asset is a sealed interface for domain assets.
type Asset interface {
	Validate() error
	GetID() string
	GetUserID() string
	GetType() AssetType
	isAsset() // Sealed interface method
}

// BaseAsset contains common fields for all assets.
type BaseAsset struct {
	ID          string    `json:"id,omitzero"`
	UserID      string    `json:"user_id,omitzero"`
	Type        AssetType `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitzero"`
}

func (b BaseAsset) GetID() string {
	return b.ID
}

func (b BaseAsset) GetUserID() string {
	return b.UserID
}

func (b BaseAsset) GetType() AssetType {
	return b.Type
}

// isAsset implements the sealed interface marker for all embedding types.
func (b BaseAsset) isAsset() {}

// ValidateCommon validates the fields common to all assets and ensures the type matches.
func (b BaseAsset) ValidateCommon(expectedType AssetType) error {
	if b.ID == "" {
		return fmt.Errorf("%w: id is required", ErrValidation)
	}
	if b.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	if b.Type != expectedType {
		return fmt.Errorf("%w: invalid asset type: expected %s, got %s", ErrValidation, expectedType, b.Type)
	}
	return nil
}
