package dto

// AddFavouriteRequest represents the request body for adding a favourite.
type AddFavouriteRequest struct {
	AssetID string `json:"assetID" binding:"required"`
}

// UpdateFavouriteDescriptionRequest represents the request body for updating a favourite description.
type UpdateFavouriteDescriptionRequest struct {
	Description string `json:"description" binding:"required"`
}

func (req *AddFavouriteRequest) GetAssetID() string {
	return req.AssetID
}

func (req *UpdateFavouriteDescriptionRequest) GetDescription() string {
	return req.Description
}
