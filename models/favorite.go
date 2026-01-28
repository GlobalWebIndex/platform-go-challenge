package models

import "time"

type Favourite struct {
	UserID      string    `db:"user_id"`
	AssetID     string    `db:"asset_id"`
	Description *string   `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}
