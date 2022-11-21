package favouriteasset

// Asset contains the asset data.
type Asset struct {
	ID   string      `json:"id"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// FavouriteAsset contains the favourite asset data.
type FavouriteAsset struct {
	ID          *string `json:"id"`
	Description string  `json:"description"`
	Asset       *Asset  `json:"asset"`
}

// EditFavouriteAsset contains the favourite asset data for editing.
type EditFavouriteAsset struct {
	Description string `json:"description"`
}

// GetFavouriteAssetsRes contains the response of get favourite assets.
type GetFavouriteAssetsRes struct {
	FavouriteAssets *[]FavouriteAsset `json:"favourite_assets"`
}

// AddFavouriteAssetRes contains the response of a favourite asset edit.
type AddFavouriteAssetRes struct {
	ID string `json:"favourite_asset_id"`
}
