package handlers

import (
	"context"
	"encoding/json"
	"gwi-platform/assets"
	"gwi-platform/database"
	"gwi-platform/models"
	"gwi-platform/utils"
	"net/http"
	"strconv"
	"time"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	asset, ok := assets.GetAsset(req.Type)
	if !ok {
		http.Error(w, "Invalid asset type", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(req.Payload, asset); err != nil {
		http.Error(w, "Invalid payload for asset type", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	db := database.GetDB()
	if err := asset.Add(ctx, db); err != nil {
		utils.ErrorLogger.Printf("Failed to add asset: %v", err)
		http.Error(w, "Failed to add asset", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(asset)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	assetType := r.URL.Query().Get("type")
	idStr := r.URL.Query().Get("id")

	asset, ok := assets.GetAsset(assetType)
	if !ok {
		http.Error(w, "Invalid asset type", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	db := database.GetDB()
	if err := asset.Delete(ctx, db, models.ID(id)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete asset: %v", err)
		http.Error(w, "Failed to delete asset", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("DELETED"))
}

func ModifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	asset, ok := assets.GetAsset(req.Type)
	if !ok {
		http.Error(w, "Invalid asset type", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(req.Payload, asset); err != nil {
		http.Error(w, "Invalid payload for asset type", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	db := database.GetDB()
	if err := asset.Modify(ctx, db); err != nil {
		utils.ErrorLogger.Printf("Failed to modify asset: %v", err)
		http.Error(w, "Failed to modify asset", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(asset)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	assetType := r.URL.Query().Get("type")
	idStr := r.URL.Query().Get("id")

	asset, ok := assets.GetAsset(assetType)
	if !ok {
		http.Error(w, "Invalid asset type", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	db := database.GetDB()
	result, err := asset.Get(ctx, db, models.ID(id))
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get asset: %v", err)
		http.Error(w, "Failed to get asset", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
