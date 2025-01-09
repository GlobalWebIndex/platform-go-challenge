package handlers

import (
	"net/http"
)

// UserFavouritesHandler defines the methods for handling user favourites.
type UserFavouritesHandler interface {

	// GetUserFavourites godoc
	// @Summary Get user's favourites
	// @Description Get the list of user's favourite assets
	// @Tags users
	// @Accept  json
	// @Produce  json
	// @Param id path string true "User ID"
	// @Success 200 {array} models.Asset
	// @Failure 404 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /users/{id}/favourites [get]
	GetUserFavourites(w http.ResponseWriter, r *http.Request)

	// AddAssetToFavourites godoc
	// @Summary Add asset to user's favourites
	// @Description Add an asset to the user's list of favourite assets
	// @Tags users
	// @Accept  json
	// @Produce  json
	// @Param id path string true "User ID"
	// @Param asset body models.Asset true "Asset"
	// @Success 200 {array} models.Asset
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /users/{id}/favourites [post]
	AddAssetToFavourites(w http.ResponseWriter, r *http.Request)

	// RemoveAssetFromFavourites godoc
	// @Summary Remove asset from user's favourites
	// @Description Remove an asset from the user's list of favourite assets
	// @Tags users
	// @Accept  json
	// @Produce  json
	// @Param id path string true "User ID"
	// @Param asset_id path string true "Asset ID"
	// @Success 200 {array} models.Asset
	// @Failure 404 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /users/{id}/favourites/{asset_id} [delete]
	RemoveAssetFromFavourites(w http.ResponseWriter, r *http.Request)

	// EditAssetDescription godoc
	// @Summary Edit asset description
	// @Description Edit the description of an asset in the user's list of favourite assets
	// @Tags users
	// @Accept  json
	// @Produce  json
	// @Param id path string true "User ID"
	// @Param asset_id path string true "Asset ID"
	// @Param asset body models.Asset true "Asset"
	// @Success 200 {array} models.Asset
	// @Failure 404 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /users/{id}/favourites/{asset_id} [put]
	EditAssetDescription(w http.ResponseWriter, r *http.Request)
}
