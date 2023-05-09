package userFavourite

import "github.com/loukaspe/platform-go-challenge/internal/core/domain"

// request for adding asset to user's favourite <chart|audience|insights>
//
// swagger:parameters addUserFavouriteAsset
type AddUserFavouriteAssetRequest struct {
	// in:body
	// Required: true
	AssetType string `json:"assetType"`
	// in:body
	// Required: true
	AssetId uint `json:"assetId"`
}

// Response when we add asset to user's favorite
// swagger:model AddUserFavouriteAssetResponse
type AddUserFavouriteAssetResponse struct {
	// possible error message
	//
	// Required: false
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// Response when we get user's favorite assets
// swagger:model GetUserFavouriteAssetsResponse
type GetUserFavouriteAssetsResponse struct {
	// user's favourite assets
	//
	// Required: false
	FavouritesAssets []domain.Asset `json:"favouritesAssets,omitempty"`
	// possible error message
	//
	// Required: false
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// request for deleting asset from user's favourite <chart|audience|insights>
//
// swagger:parameters deleteUserFavouriteAsset
type DeleteUserFavouriteAssetRequest struct {
	// in:body
	// Required: true
	AssetType string `json:"assetType"`
	// in:body
	// Required: true
	AssetId uint `json:"assetId"`
}

// Response when we delete user's favorite asset
// swagger:model DeleteUserFavouriteAssetResponse
type DeleteUserFavouriteAssetResponse struct {
	// possible error message
	//
	// Required: false
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// request for updating asset's description to user's favourite <chart|audience|insight>
//
// swagger:parameters updateUserFavouriteAssetDescription
type UpdateUserFavouriteAssetDescriptionRequest struct {
	// in:body
	// Required: true
	AssetType string `json:"assetType"`
	// in:body
	// Required: true
	Description string `json:"description"`
}

// Response when we get user's favorite assets
// swagger:model UpdateUserFavouriteAssetDescriptionResponse
type UpdateUserFavouriteAssetDescriptionResponse struct {
	// possible error message
	//
	// Required: false
	ErrorMessage string `json:"errorMessage,omitempty"`
}
