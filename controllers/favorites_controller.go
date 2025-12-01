package controllers

import (
    "encoding/json"
    "net/http"
    "github.com/dintzeler/platform-go-challenge/services"
	"github.com/dintzeler/platform-go-challenge/models"
	"github.com/dintzeler/platform-go-challenge/validators"
	"github.com/dintzeler/platform-go-challenge/customerrors"
	"github.com/dintzeler/platform-go-challenge/middleware"
)

type AddFavoriteResponse struct {
	UserID int `json:"user_id"`
	AssetID int `json:"asset_id"`
	AssetType models.AssetType `json:"asset_type"`
}

type UpdateFavoriteResponse struct {
	UserID int `json:"user_id"`
	AssetID int `json:"asset_id"`
	AssetType models.AssetType `json:"asset_type"`
	UpdatedDescription string `json:"updated_description"`
}


func FavoritesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetFavorites(w, r)
	case http.MethodPost:
		handleAddFavorite(w, r)
	case http.MethodDelete:
		handleDeleteFavorite(w, r)
	case http.MethodPatch:
		handleUpdateFavorite(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAddFavorite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&customerrors.ValidationError{
			Message:   "Invalid User-ID",
			ErrorCode: "INVALID_USER_ID",
		})
		return
	}
	
	addFavoriteRequest, err := validators.ValidateAddFavoriteRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	status, err := services.AddFavorite(userID, *addFavoriteRequest.AssetID, *addFavoriteRequest.AssetType)
	if err != nil {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(AddFavoriteResponse{
		UserID: userID,
		AssetID: *addFavoriteRequest.AssetID,
		AssetType: *addFavoriteRequest.AssetType,
	})
}

func handleGetFavorites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&customerrors.ValidationError{
			Message:   "Invalid User-ID",
			ErrorCode: "INVALID_USER_ID",
		})
		return
	}

    favorites, err := services.GetFavorites(userID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(err)
        return
    }

    json.NewEncoder(w).Encode(favorites)
}

func handleDeleteFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&customerrors.ValidationError{
			Message:   "Invalid User-ID",
			ErrorCode: "INVALID_USER_ID",
		})
		return
	}

	deleteFavoriteRequest, err := validators.ValidateDeleteFavorite(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	status, err := services.DeleteFavorite(userID, *deleteFavoriteRequest.AssetID, *deleteFavoriteRequest.AssetType)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(status)
}

func handleUpdateFavorite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&customerrors.ValidationError{
			Message:   "Invalid User-ID",
			ErrorCode: "INVALID_USER_ID",
		})
		return
	}

	updateFavorite, err := validators.ValidateUpdateFavoriteRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	status, err := services.UpdateFavorite(userID, *updateFavorite.AssetType, *updateFavorite.AssetID, updateFavorite.Description)
	if err != nil {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(UpdateFavoriteResponse{
		UserID: userID,
		AssetID: *updateFavorite.AssetID,
		AssetType: *updateFavorite.AssetType,
		UpdatedDescription: updateFavorite.Description,
	})
}