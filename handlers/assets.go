package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"platform-go-challenge/models"
	"platform-go-challenge/repositories"
)

func AssetsRouter(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Expected paths:
		// POST /assets - Create a new asset
		// GET /assets/{assetID} - Get asset by ID

		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		if parts[0] != "assets" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodPost:
			// POST /assets - Create a new asset
			if len(parts) == 1 {
				println("Create a new asset")
				CreateAsset(db)(w, r)
				return
			}

		case http.MethodGet:
			// GET /assets - List all assets
			if len(parts) == 1 {
				println("List all assets")
				GetAllAssets(db)(w, r)
				return
			}
			// GET /assets/{assetID} - Get asset by ID
			if len(parts) == 2 {
				println("Get asset by ID")
				GetAsset(db)(w, r)
				return
			}
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateAsset(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Type        models.AssetType `json:"type"`
			Title       string           `json:"title"`
			Description *string          `json:"description"`
			Data        json.RawMessage  `json:"data"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate input
		if input.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		if input.Type == "" {
			http.Error(w, "Type is required (chart, insight, or audience)", http.StatusBadRequest)
			return
		}

		asset := models.Asset{
			Type:  input.Type,
			Title: &input.Title,
			Data:  input.Data,
		}

		assetID, err := repositories.CreateAsset(r.Context(), db.(*sql.DB), asset, input.Description)
		if err != nil {
			println("Error creating asset:", err.Error())
			http.Error(w, "Failed to create asset: "+err.Error(), http.StatusInternalServerError)
			return
		}

		asset.ID = assetID
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":          assetID,
			"type":        asset.Type,
			"title":       asset.Title,
			"description": input.Description,
			"data":        asset.Data,
		})
	}
}

func GetAsset(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		assetID := parts[1]

		asset, err := repositories.GetAssetByID(r.Context(), db.(*sql.DB), assetID)
		if err != nil {
			println("Error getting asset:", err.Error())
			http.Error(w, "Failed to get asset: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if asset == nil {
			http.Error(w, "Asset not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(asset)
	}
}

func GetAllAssets(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assets, err := repositories.ListAssets(r.Context(), db.(*sql.DB))
		if err != nil {
			http.Error(w, "Failed to fetch assets: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(assets)
	}
}
