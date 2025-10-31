package handlers

import (
	"net/http"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/http/middleware"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/mapping"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/ports"
	"github.com/go-chi/chi/v5"
)

var _ ports.FavouriteHandler = (*FavouriteHandler)(nil)

type FavouriteHandler struct {
	service ports.FavouriteService
}

func NewFavouriteHandler(s ports.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{service: s}
}

// Create adds a favourite for a user
// @Summary Add a favourite
// @Description Adds an asset to a user's favourites list
// @Tags Favourites
// @Accept json
// @Produce json
// @Param request body dto.FavouriteRequest true "Favourite creation request"
// @Success 201 "Favourite added successfully"
// @Failure 400 {string} string "Invalid input data"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /favourites [post]
func (f *FavouriteHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, ok := middleware.GetValidatedBody[dto.FavouriteRequest](r)
	if !ok {
		http.Error(w, "missing validated body", http.StatusBadRequest)
		return
	}

	favourite := mapping.FavouriteReqToDomain(req)

	err := f.service.CreateFavourite(favourite)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Delete removes a favourite for a user
// @Summary Remove a favourite
// @Description Removes an asset from a user's favourites list
// @Tags Favourites
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param assetId path string true "Asset ID"
// @Success 200 "Favourite removed successfully"
// @Failure 400 {string} string "Invalid user ID or asset ID"
// @Failure 404 {string} string "User or favourite not found"
// @Failure 405 {string} string "Method not allowed"
// @Security BearerAuth
// @Router /favourites/{userId}/{assetId} [delete]
func (f *FavouriteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	usrId := chi.URLParam(r, "userId")
	if usrId == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	assetId := chi.URLParam(r, "assetId")
	if assetId == "" {
		http.Error(w, "missing asset id", http.StatusBadRequest)
		return
	}

	err := f.service.DeleteFavourite(usrId, assetId)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
}
