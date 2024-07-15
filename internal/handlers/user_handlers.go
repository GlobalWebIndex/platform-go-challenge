package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/ceciivanov/platform-go-challenge/internal/service"
	"github.com/ceciivanov/platform-go-challenge/internal/utils"
	"github.com/gorilla/mux"
)

// Handler struct
type UserHandler struct {
	UserService *service.UserService
}

// NewUserHandler initializes and returns a new Handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// RegisterRoutes registers the routes (endpoints) for the user handler
func (handler *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/{id}/favorites", handler.GetUserFavorites).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}/favorites", handler.AddUserFavorite).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}/favorites/{assetID}", handler.DeleteUserFavorite).Methods(http.MethodDelete)
	r.HandleFunc("/users/{id}/favorites/{assetID}", handler.EditUserFavorite).Methods(http.MethodPut)
}

// GetUserFavorites returns a map of user's favorite assets
func (h *UserHandler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	favorites, err := h.UserService.GetUserFavorites(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favorites)
}

// AddUserFavorite adds an asset to the user's favorites
func (h *UserHandler) AddUserFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	if r.Body == nil {
		http.Error(w, "no request body", http.StatusBadRequest)
		return
	}

	newAssetData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	newAsset, err := utils.DecodeAsset(newAssetData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.UserService.AddUserFavorite(userID, newAsset)
	if err != nil {
		switch err.Error() {
		case "user not found":
			http.Error(w, "user not found", http.StatusNotFound)
			return
		case "asset already exists":
			http.Error(w, "asset already exists", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAsset)
}

// DeleteUserFavorite deletes an asset from the user's favorites
func (h *UserHandler) DeleteUserFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])
	assetID, _ := strconv.Atoi(vars["assetID"])

	err := h.UserService.DeleteUserFavorite(userID, assetID)
	if err != nil {
		switch err.Error() {
		case "user not found":
			http.Error(w, "user not found", http.StatusNotFound)
			return
		case "asset not found":
			http.Error(w, "asset not found", http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// EditUserFavorite edits an asset in the user's favorites
func (h *UserHandler) EditUserFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])
	assetID, _ := strconv.Atoi(vars["assetID"])

	updatedAssetData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	updatedAsset, err := utils.DecodeAsset(updatedAssetData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.UserService.EditUserFavorite(userID, assetID, updatedAsset)
	if err != nil {
		switch err.Error() {
		case "user not found":
			http.Error(w, "user not found", http.StatusNotFound)
			return
		case "asset not found":
			http.Error(w, "asset not found", http.StatusNotFound)
			return
		case "edited asset ID does not match existing asset ID":
			http.Error(w, "edited asset ID does not match existing asset ID", http.StatusBadRequest)
			return
		case "edited asset type does not match existing asset type":
			http.Error(w, "edited asset type does not match existing asset type", http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedAsset)
}
