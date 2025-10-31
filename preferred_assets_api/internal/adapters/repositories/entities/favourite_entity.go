package entities

import (
	"time"
)

type FavouriteEntity struct {
	UserId    string
	AssetId   string
	CreatedAt time.Time
}
