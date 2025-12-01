package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"platform-go-challenge/internal/common"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/storage"
)

type UserAdminHandler struct {
	cfg       *config.Config
	users     *storage.UserStore
	favorites *storage.FavoritesStore
}

func NewUserAdminHandler(cfg *config.Config, users *storage.UserStore, favorites *storage.FavoritesStore) *UserAdminHandler {
	return &UserAdminHandler{cfg: cfg, users: users, favorites: favorites}
}

func (uah *UserAdminHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	// AdminAPIKeyMiddleware already validated admin access, no need to check role
	// Fetch the list of users (usernames only)
	users := uah.users.GetAllUsernames()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		// Create an error response for failure to encode the response
		errorResponse := common.ErrorResponse{
			Error: "failed to encode response: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
	}
}

func (uah *UserAdminHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	// AdminAPIKeyMiddleware already validated admin access, no need to check role
	// Extract user ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		errorResponse := common.ErrorResponse{
			Error: "missing user ID in path",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		errorResponse := common.ErrorResponse{
			Error: "invalid user ID",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Delete user from store
	if err := uah.users.DeleteUser(userID); err != nil {
		errorResponse := common.ErrorResponse{
			Error: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Remove user favorites
	uah.favorites.RemoveUserFavorites(userID)

	// Success: no content
	w.WriteHeader(http.StatusNoContent)
}
