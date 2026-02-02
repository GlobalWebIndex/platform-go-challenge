package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/kongrillas/platform-go-challenge/internal/model"
	"github.com/kongrillas/platform-go-challenge/internal/store"
)

type FavoritesHandler struct {
	Store *store.MemoryStore
}

func NewFavoritesHandler(s *store.MemoryStore) *FavoritesHandler {
	return &FavoritesHandler{Store: s}
}

// Expected paths:
// /users/{userId}/favorites
// /users/{userId}/favorites/{assetId}
func parsePath(path string) (userID string, assetID string, ok bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 || parts[0] != "users" || parts[2] != "favorites" {
		return "", "", false
	}
	userID = parts[1]
	if userID == "" {
		return "", "", false
	}
	if len(parts) >= 4 {
		assetID = parts[3]
	}
	return userID, assetID, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *FavoritesHandler) ListFavorites(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := parsePath(r.URL.Path)
	if !ok {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	items, err := h.Store.List(userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *FavoritesHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := parsePath(r.URL.Path)
	if !ok {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	var a model.Asset
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := model.ValidateAsset(a); err != nil {
		http.Error(w, "invalid asset payload", http.StatusBadRequest)
		return
	}

	err := h.Store.Add(userID, a)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			http.Error(w, "already exists", http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, a)
}

func (h *FavoritesHandler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	userID, assetID, ok := parsePath(r.URL.Path)
	if !ok || assetID == "" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	err := h.Store.Remove(userID, assetID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type updateDescriptionRequest struct {
	Description string `json:"description"`
}

func (h *FavoritesHandler) UpdateDescription(w http.ResponseWriter, r *http.Request) {
	userID, assetID, ok := parsePath(r.URL.Path)
	if !ok || assetID == "" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	var req updateDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Description) == "" {
		http.Error(w, "description required", http.StatusBadRequest)
		return
	}

	updated, err := h.Store.UpdateDescription(userID, assetID, req.Description)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}
