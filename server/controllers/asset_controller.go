package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"favourite_assets/server/errors"
	"favourite_assets/server/authentication"
	"favourite_assets/server/models"
	"favourite_assets/server/services"
)

type AssetController struct {
	AssetService *services.AssetService
}

func NewAssetController(assetService *services.AssetService) *AssetController {
	return &AssetController{AssetService: assetService}
}

// (admin- only)
func (c *AssetController) CreateAssetHandler(w http.ResponseWriter, r *http.Request) {

	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	if err := authentication.RequireRole(r.Context(), "admin"); err != nil {
		errors.WriteError(w, errors.ErrForbidden)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	var asset models.Asset
	assetTypeStr, ok := req["type"].(string)
	if !ok {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}
	assetType := models.AssetType(strings.ToLower(assetTypeStr))

	switch assetType {
	case models.AssetChart:
		asset = &models.Chart{
			BaseAsset: models.BaseAsset{Description: req["description"].(string)},
			Title:     req["title"].(string),
			XAxis:     req["xAxis"].(string),
			YAxis:     req["yAxis"].(string),
		}
	case models.AssetInsight:
		asset = &models.Insight{
			BaseAsset: models.BaseAsset{Description: req["description"].(string)},
			Text:      req["text"].(string),
		}
	case models.AssetAudience:
		asset = &models.Audience{
			BaseAsset:          models.BaseAsset{Description: req["description"].(string)},
			Gender:             req["gender"].(string),
			BirthCountry:       req["birthCountry"].(string),
			AgeGroup:           req["ageGroup"].(string),
			HoursOnSocial:      int(req["hoursSocialDaily"].(float64)),
			PurchasesLastMonth: int(req["purchasesLastMonth"].(float64)),
		}
	default:
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	created, err := c.AssetService.CreateAsset(asset)
	if err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	errors.WriteJSON(w, http.StatusCreated, created)
}

// (all-roles)
func (c *AssetController) GetAssetHandler(w http.ResponseWriter, r *http.Request) {
	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	idStr := r.URL.Query().Get("assetId")
	assetID, err := uuid.Parse(idStr)
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	asset, err := c.AssetService.GetAsset(assetID)
	if err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	errors.WriteJSON(w, http.StatusOK, asset)
}

// (admin-only)
func (c *AssetController) UpdateAssetHandler(w http.ResponseWriter, r *http.Request) {
	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	if err := authentication.RequireRole(r.Context(), "admin"); err != nil {
		errors.WriteError(w, errors.ErrForbidden)
		return
	}

	idStr := r.URL.Query().Get("assetId")
	assetID, err := uuid.Parse(idStr)
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}
	updated, err := c.AssetService.UpdateAsset(assetID, req)
	if err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	errors.WriteJSON(w, http.StatusOK, updated)
}

// (admin- only)
func (c *AssetController) DeleteAssetHandler(w http.ResponseWriter, r *http.Request) {

	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	if err := authentication.RequireRole(r.Context(), "admin"); err != nil {
		errors.WriteError(w, errors.ErrForbidden)
		return
	}

	idStr := r.URL.Query().Get("assetId")
	assetID, err := uuid.Parse(idStr)
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	if err := c.AssetService.DeleteAsset(assetID); err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// (all-roles)
func (c *AssetController) ListAssetsHandler(w http.ResponseWriter, r *http.Request) {
	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	queryType := r.URL.Query().Get("type")

	var assets []models.Asset
	if queryType != "" {
		assets = c.AssetService.ListAssetsByType(models.AssetType(queryType))
	} else {
		assets = c.AssetService.ListAssets()
	}

	errors.WriteJSON(w, http.StatusOK, assets)
}
