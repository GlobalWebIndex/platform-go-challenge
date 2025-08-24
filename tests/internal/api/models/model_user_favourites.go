package models

type Favourites struct {
	UserID string `json:"userID" bson:"userID"`
	AssetID string `json:"assetID" bson:"assetID"`
	AssetType string `json:"assetType" bson:"assetType"`
	Description string `json:"description" bson:"description"`
}

type AddFavouriteBody struct {
	AssetType string `json:"assetType"`
	Description string `json:"description"`
}

type PatchFavouriteBody struct {
	Desription string `json:"description"`
}
