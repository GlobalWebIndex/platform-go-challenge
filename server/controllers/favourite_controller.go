package controllers

import (
	"net/http"

	"favourite_assets/server/errors"
	"favourite_assets/server/authentication"
	"favourite_assets/server/services"

	"github.com/google/uuid"
)

type FavouriteController struct {
	FavouriteService *services.FavouriteService
}

func NewFavouriteController(favService *services.FavouriteService) *FavouriteController {
	return &FavouriteController{
		FavouriteService: favService,
	}
}

// (all-roles)
func (c *FavouriteController) AddFavouriteHandler(w http.ResponseWriter, r *http.Request) {
	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	userIDStr := r.URL.Query().Get("userId")
	assetIDStr := r.URL.Query().Get("assetId")
	if userIDStr == "" || assetIDStr == "" {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	fav, err := c.FavouriteService.AddFavourite(userID, assetID)
	if err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	errors.WriteJSON(w, http.StatusCreated, fav)
}

// (all-roles)
func (c *FavouriteController) RemoveFavouriteHandler(w http.ResponseWriter, r *http.Request) {

	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	favID, err := uuid.Parse(r.URL.Query().Get("favouriteId"))
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	if err := c.FavouriteService.RemoveFavourite(favID); err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// (all-roles)
func (c *FavouriteController) ListFavouritesHandler(w http.ResponseWriter, r *http.Request) {

	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	userID, err := uuid.Parse(r.URL.Query().Get("userId"))
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	favourites, err := c.FavouriteService.ListFavouritesByUser(userID)
	if err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	errors.WriteJSON(w, http.StatusOK, favourites)
}

// (all-roles)
func (c *FavouriteController) GetFavouriteHandler(w http.ResponseWriter, r *http.Request) {
	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	favIDStr := r.URL.Query().Get("favouriteId")
	if favIDStr == "" {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	favID, err := uuid.Parse(favIDStr)
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	fav, err := c.FavouriteService.GetFavourite(favID)
	if err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	errors.WriteJSON(w, http.StatusOK, fav)
}
