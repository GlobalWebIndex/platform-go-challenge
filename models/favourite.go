package models

import "time"

type Favourite struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	AssetID   string    `json:"asset_id"`
	CreatedAt time.Time `json:"created_at"`
}
