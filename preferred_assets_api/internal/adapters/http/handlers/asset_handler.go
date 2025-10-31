package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/http/middleware"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/mapping"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/ports"
	"github.com/go-chi/chi/v5"
)

var _ ports.AssetHandler = (*AssetHandler)(nil)

type AssetHandler struct {
	service ports.AssetService
}

func NewAssetHandler(s ports.AssetService) *AssetHandler {
	return &AssetHandler{service: s}
}

// Create creates a new asset
// @Summary Create a new asset
// @Description Creates a new asset of specified type (audience, chart, or insight)
// @Tags Assets
// @Accept json
// @Produce json
// @Param request body dto.AssetRequest true "Asset creation request"
// @Success 201 {object} dto.AssetCreationResponse
// @Failure 400 {string} string "Invalid input data"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /assets [post]
func (h *AssetHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, ok := middleware.GetValidatedBody[dto.AssetRequest](r)
	if !ok {
		http.Error(w, "missing validated body", http.StatusBadRequest)
		return
	}

	asset, err := mapping.AssetReqToDomain(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAsset, err := h.service.CreateAsset(asset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := mapping.AssetDomainToCreationResponse(createdAsset)

	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("JSON marshaling error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Assets JSON response: %s", string(jsonBytes))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

// Delete removes an asset by ID
// @Summary Delete an asset
// @Description Permanently removes an asset from the system
// @Tags Assets
// @Accept json
// @Produce json
// @Param assetId path string true "Asset ID"
// @Success 204 "Asset deleted successfully"
// @Failure 400 {string} string "Invalid asset ID"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /assets/{assetId} [delete]
func (h *AssetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	assetID := chi.URLParam(r, "assetId")
	if assetID == "" {
		http.Error(w, "missing asset id", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteAsset(assetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
