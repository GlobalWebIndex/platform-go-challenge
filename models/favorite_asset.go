package models

type FavouriteAsset struct {
	AssetID     string    `json:"id"`
	Type        AssetType `json:"type"`
	Title       *string   `json:"title"`
	Data        any       `json:"data"`
	Description *string   `json:"description"`
}
