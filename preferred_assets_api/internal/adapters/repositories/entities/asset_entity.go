package entities

import (
	"errors"
	"time"
)

type AssetBaseEntity struct {
	ID          string    `db:"id"`
	Type        AssetType `db:"type"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type AssetEntity interface {
	GetID() string
	GetType() AssetType
	GetTitle() string
	GetDescription() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	Validate() error
}

func (a AssetBaseEntity) GetID() string           { return a.ID }
func (a AssetBaseEntity) GetType() AssetType      { return AssetType(a.Type) }
func (a AssetBaseEntity) GetTitle() string        { return a.Title }
func (a AssetBaseEntity) GetDescription() string  { return a.Description }
func (a AssetBaseEntity) GetCreatedAt() time.Time { return a.CreatedAt }
func (a AssetBaseEntity) GetUpdatedAt() time.Time { return a.UpdatedAt }

// Validate Data Consistency Validation
func (a AssetBaseEntity) Validate() error {
	if a.ID == "" {
		return errors.New("asset ID is required")
	}
	if a.Title == "" {
		return errors.New("asset title is required")
	}
	if a.Type == 0 {
		return errors.New("asset type is required")
	}
	return nil
}
