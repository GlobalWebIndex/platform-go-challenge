package domain

import "time"

type Favourite struct {
	From      string    `json:"_from,omitempty"` // Reference to the user document
	To        string    `json:"_to,omitempty"`   // Reference to the asset document
	CreatedAt time.Time `json:"created_at"`
}

type FavouriteService interface {
	AddToFavourites(userID, assetID string) (Favourite, error)
	GetFavourites(userID string) ([]Asset, error)
	EditFavourites(userID, assetID, description string) (Asset, error)
	DeleteFavourite(userID, assetID string) (Favourite, error)
}

type FavouriteRepository interface {
	AddToFavourites(userID, assetID string) (Favourite, error)
	GetFavourites(userID string) ([]Asset, error)
	EditFavourites(userID, assetID, description string) (Asset, error)
	DeleteFavourite(userID, assetID string) (Favourite, error)
}
