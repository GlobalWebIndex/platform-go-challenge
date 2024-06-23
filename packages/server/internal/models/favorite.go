package models

import "time"

type Favorite struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint `gorm:"primaryKey;uniqueIndex:user_asset,priority:1"`
	AssetID   uint `gorm:"primaryKey;uniqueIndex:user_asset,priority:2"`
	Asset     Asset
}

type FavoritesResponse struct {
	ID      uint
	UserID  uint
	AssetID uint
	Asset   AssetResponse
}
