package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

// parseUserID extracts and parses the user_id query parameter as uuid.UUID
func parseUserID(r *http.Request, w http.ResponseWriter) (uuid.UUID, bool) {
	userIDStr := r.URL.Query().Get("user_id")
	log.Printf("Parsed user_id: %s of type %T\n", userIDStr, userIDStr)
	userID, err := uuid.Parse(userIDStr)
	log.Printf("Parsed userID: %s of type %T\n", userID, userID)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return uuid.UUID{}, false
	}
	return userID, true
}

// parseAssetID extracts and parses the asset_id query parameter as uuid.UUID
func parseAssetID(r *http.Request, w http.ResponseWriter) (uuid.UUID, bool) {
	assetIDStr := r.URL.Query().Get("asset_id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		http.Error(w, "Invalid asset_id", http.StatusBadRequest)
		return uuid.UUID{}, false
	}
	return assetID, true
}

func setupRoutes() {
	http.HandleFunc("/token", TokenHandler)
	http.HandleFunc("/favourites", AuthMiddleware(handleFavourites))
	http.HandleFunc("/favourites/add", AuthMiddleware(handleAddFavourite))
	http.HandleFunc("/favourites/remove", AuthMiddleware(handleRemoveFavourite))
	http.HandleFunc("/favourites/edit", AuthMiddleware(handleEditFavourite))
	http.HandleFunc("/favourites/delete", AuthMiddleware(handleDeleteFavourite))

}

// List all the assets of the user with Favorite == true
func handleFavourites(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		log.Printf("handleFavourites: invalid user_id from token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user := store.GetUser(userID)
	if user == nil {
		log.Printf("handleFavourites: user not found %s", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Only return assets with Favorite == true
	favs := make([]Asset, 0)
	for _, asset := range user.Favourites {
		if asset.IsFavorite() {
			favs = append(favs, asset)
		}
	}
	log.Printf("handleFavourites: returning %d assets for user %s", len(favs), userID)
	// Pagination: limit and offset query params
	limit := 0
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}
	start := offset
	if start > len(favs) {
		start = len(favs)
	}
	end := len(favs)
	if limit > 0 && start+limit < end {
		end = start + limit
	}
	paged := favs[start:end]
	json.NewEncoder(w).Encode(paged)
}

// Create an asset
func handleAddFavourite(w http.ResponseWriter, r *http.Request) {

	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		log.Printf("handleAddFavourite: invalid user_id from token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user := store.GetUser(userID)
	if user == nil {
		log.Printf("handleAddFavourite: user not found %s", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	var req struct {
		Type     string          `json:"type"`
		Asset    json.RawMessage `json:"asset"`
		Favorite bool            `json:"favorite"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("handleAddFavourite: invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var asset Asset
	switch req.Type {
	case ChartType:
		log.Printf("Adding Chart asset: %s\n", req.Asset)
		var c Chart
		if err := json.Unmarshal(req.Asset, &c); err != nil {
			log.Printf("handleAddFavourite: invalid chart asset: %v", err)
			http.Error(w, "Invalid chart asset", http.StatusBadRequest)
			return
		}
		if c.ID == uuid.Nil {
			c.ID = uuid.New()
		}
		c.Favorite = req.Favorite
		asset = &c
	case InsightType:
		log.Printf("Adding Insight asset: %s\n", req.Asset)
		var i Insight
		if err := json.Unmarshal(req.Asset, &i); err != nil {
			log.Printf("handleAddFavourite: invalid insight asset: %v", err)
			http.Error(w, "Invalid insight asset", http.StatusBadRequest)
			return
		}
		if i.ID == uuid.Nil {
			i.ID = uuid.New()
		}
		i.Favorite = req.Favorite
		asset = &i
	case AudienceType:
		log.Printf("Adding Audience asset: %s\n", req.Asset)
		var a Audience
		if err := json.Unmarshal(req.Asset, &a); err != nil {
			log.Printf("handleAddFavourite: invalid audience asset: %v", err)
			http.Error(w, "Invalid audience asset", http.StatusBadRequest)
			return
		}
		if a.ID == uuid.Nil {
			a.ID = uuid.New()
		}
		a.Favorite = req.Favorite
		asset = &a
	default:
		log.Printf("handleAddFavourite: unknown asset type %s", req.Type)
		http.Error(w, "Unknown asset type", http.StatusBadRequest)
		return
	}
	log.Printf("handleAddFavourite: asset added for user %s, type %s, id %s", userID, req.Type, asset.GetID())
	user.Favourites = append(user.Favourites, asset)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(asset)
}

// Edit the isFavorite field of an asset
func handleRemoveFavourite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Printf("handleRemoveFavourite: method not allowed %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		log.Printf("handleRemoveFavourite: invalid user_id from token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	assetID, ok := parseAssetID(r, w)
	if !ok {
		log.Printf("handleRemoveFavourite: invalid asset_id")
		return
	}
	user := store.GetUser(userID)
	if user == nil {
		log.Printf("handleRemoveFavourite: user not found %s", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	var req struct {
		Favorite bool `json:"favorite"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("handleRemoveFavourite: invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	for _, fav := range user.Favourites {
		if fav.GetID() == assetID {
			log.Printf("handleRemoveFavourite: updating favorite for asset %s to %v", assetID, req.Favorite)
			fav.SetFavorite(req.Favorite)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(fav)
			return
		}
	}
	log.Printf("handleRemoveFavourite: asset not found %s", assetID)
	http.Error(w, "Asset not found in favourites", http.StatusNotFound)
}

// Edits the description of an asset in general
func handleEditFavourite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Printf("handleEditFavourite: method not allowed %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		log.Printf("handleEditFavourite: invalid user_id from token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	assetID, ok := parseAssetID(r, w)
	if !ok {
		log.Printf("handleEditFavourite: invalid asset_id")
		return
	}
	user := store.GetUser(userID)
	if user == nil {
		log.Printf("handleEditFavourite: user not found %s", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	var req struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("handleEditFavourite: invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	for _, fav := range user.Favourites {
		if fav.GetID() == assetID {
			log.Printf("handleEditFavourite: updating description for asset %s to '%s'", assetID, req.Description)
			fav.SetDescription(req.Description)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(fav)
			return
		}
	}
	log.Printf("handleEditFavourite: asset not found %s", assetID)
	http.Error(w, "Asset not found in favourites", http.StatusNotFound)
}

// Deletes an asset in general
func handleDeleteFavourite(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		log.Printf("handleDeleteFavourite: invalid user_id from token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	assetID, ok := parseAssetID(r, w)
	if !ok {
		log.Printf("handleDeleteFavourite: invalid asset_id")
		return
	}
	user := store.GetUser(userID)
	if user == nil {
		log.Printf("handleDeleteFavourite: user not found %s", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	found := false
	newFavs := make([]Asset, 0, len(user.Favourites))
	for _, fav := range user.Favourites {
		if fav.GetID() == assetID {
			log.Printf("handleDeleteFavourite: deleting asset %s for user %s", assetID, userID)
			found = true
			continue
		}
		newFavs = append(newFavs, fav)
	}
	if !found {
		log.Printf("handleDeleteFavourite: asset not found %s", assetID)
		http.Error(w, "Asset not found in favourites", http.StatusNotFound)
		return
	}
	user.Favourites = newFavs
	log.Printf("handleDeleteFavourite: asset deleted, %d assets remain for user %s", len(user.Favourites), userID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.Favourites)
}
