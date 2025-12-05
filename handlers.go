package main

import (
	"net/http"
	"log"
	"time"
	"encoding/json"
	"github.com/google/uuid"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming request: method=%s path=%s", r.Method, r.URL.Path)

	userID, assetInfo, err := parseUserPath(r)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
		case http.MethodGet:
			handleAssetsGet(w, userID)
		case http.MethodPost:
			if assetInfo == "" {
				handleUserPost(w, userID)
			} else {
				assetType := assetInfo
				handleAssetPost(w, r, userID, assetType)
			}
		case http.MethodPatch:
			assetID := assetInfo
			handleAssetPatch(w, r, userID, assetID)
		case http.MethodDelete:
			assetID := assetInfo
			handleAssetDelete(w, userID, assetID)
		default:
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Request handle functions
func handleAssetsGet(w http.ResponseWriter, userID string)  {
	userFavourites := allUserFavourites[userID]

	if userFavourites == nil {
		log.Printf("Error: User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	returnedAssets := []Asset{}

	for _, asset := range userFavourites {
		returnedAssets = append(returnedAssets, asset)
	}

	log.Printf("Returning %d favourites for userID=%s", len(returnedAssets), userID)

	json.NewEncoder(w).Encode(returnedAssets)
}

func handleUserPost(w http.ResponseWriter, userID string) {
	if allUserFavourites[userID] == nil {
		allUserFavourites[userID] = map[string]Asset{}
		w.WriteHeader(http.StatusOK)
	} else {
		log.Printf("Error: User already exists")
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
}

func handleChartPost(w http.ResponseWriter, r *http.Request, userID string) {
	if allUserFavourites[userID] == nil {
		log.Printf("Error: User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var payload struct {
		Description string `json:"description"`
		Data Chart `json:"data"`
	}
	json.NewDecoder(r.Body).Decode(&payload)

	asset := Asset{ID: uuid.New().String(), Type: "chart", Description: payload.Description, Data: payload.Data, CreatedAt: time.Now(), UpdatedAt: time.Now(),}

	allUserFavourites[userID][asset.ID] = asset

	log.Printf("Added asset: id=%s type=%s description=%s userID=%s", asset.ID, asset.Type, payload.Description, userID)

	json.NewEncoder(w).Encode(asset)
}

func handleInsightPost(w http.ResponseWriter, r *http.Request, userID string) {
	if allUserFavourites[userID] == nil {
		log.Printf("Error: User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var payload struct {
		Description string `json:"description"`
		Data Insight `json:"data"`
	}
	json.NewDecoder(r.Body).Decode(&payload)

	asset := Asset{ID: uuid.New().String(), Type: "insight", Description: payload.Description, Data: payload.Data, CreatedAt: time.Now(), UpdatedAt: time.Now(),}

	allUserFavourites[userID][asset.ID] = asset

	log.Printf("Added asset: id=%s type=%s description=%s userID=%s", asset.ID, asset.Type, payload.Description, userID)

	json.NewEncoder(w).Encode(asset)
}

func handleAudiencePost(w http.ResponseWriter, r *http.Request, userID string) {
	if allUserFavourites[userID] == nil {
		log.Printf("Error: User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var payload struct {
		Description string `json:"description"`
		Data Audience `json:"data"`
	}
	json.NewDecoder(r.Body).Decode(&payload)

	asset := Asset{ID: uuid.New().String(), Type: "audience", Description: payload.Description, Data: payload.Data, CreatedAt: time.Now(), UpdatedAt: time.Now(),}

	allUserFavourites[userID][asset.ID] = asset

	log.Printf("Added asset: id=%s type=%s description=%s userID=%s", asset.ID, asset.Type, payload.Description, userID)

	json.NewEncoder(w).Encode(asset)
}

func handleAssetPost(w http.ResponseWriter, r *http.Request, userID string, assetType string) {
	switch assetType {
		case "chart":
			handleChartPost(w, r, userID)
		case "insight":
			handleInsightPost(w, r, userID)
		case "audience":
			handleAudiencePost(w, r, userID)
		default:
			http.Error(w, "invalid asset type", http.StatusBadRequest)
	}
}

func handleAssetDelete(w http.ResponseWriter, userID string, assetID string) {
	userFavourites := allUserFavourites[userID]
	
	if userFavourites == nil {
		log.Printf("Error: User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	_, exists := userFavourites[assetID]
	if !exists {
		log.Printf("Error: Asset not found")
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}

	delete(userFavourites, assetID)

	log.Printf("Deleted asset %s for user %s", assetID, userID)

	w.WriteHeader(http.StatusOK)
}

func handleAssetPatch(w http.ResponseWriter, r *http.Request, userID string, assetID string) {
	userFavourites := allUserFavourites[userID]

	if userFavourites == nil {
		log.Printf("Error: User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	asset, exists := userFavourites[assetID]
	if !exists {
		log.Printf("Error: Asset not found")
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}

	// patch payload
	var payload struct {
		Description *string `json:"description"`
		Data interface{} `json:"data"`
	}

	// deserialize JSON payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	// description patch
	if payload.Description != nil {
		asset.Description = *payload.Description
	}

	// asset data patch
	if payload.Data != nil {
		switch asset.Type {
		case "chart":
			// partial update of the chart payload (e.g., only Chart.Title or only Chart.Xaxis to be updated) is not supported
			var updated Chart
			b, _ := json.Marshal(payload.Data)
			err := json.Unmarshal(b, &updated)
			if err != nil {
				log.Printf("Error: %v", err)
				http.Error(w, "invalid chart data", http.StatusBadRequest)
				return
			}
			asset.Data = updated
		case "insight":
			var updated Insight
			b, _ := json.Marshal(payload.Data)
			err := json.Unmarshal(b, &updated)
			if err != nil {
				log.Printf("Error: %v", err)
				http.Error(w, "invalid insight data", http.StatusBadRequest)
				return
			}
			asset.Data = updated
		case "audience":
			// partial update of the audience payload (e.g., only Audience.Gender or only Audience.BirthCountry to be updated) is not supported
			var updated Audience
			b, _ := json.Marshal(payload.Data)
			err := json.Unmarshal(b, &updated)
			if err != nil {
				log.Printf("Error: %v", err)
				http.Error(w, "invalid audience data", http.StatusBadRequest)
				return
			}
			asset.Data = updated
		}
	}

	asset.UpdatedAt = time.Now()
	userFavourites[assetID] = asset

	log.Printf("Updated asset %s for user %s", assetID, userID)
	json.NewEncoder(w).Encode(asset)
}