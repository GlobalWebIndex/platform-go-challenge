package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SokratisChaimanas/platform-go-challenge/internal/app"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// AssetHandler serves asset endpoints.
type AssetHandler struct {
	svc *app.AssetService
}

func NewAssetHandler(assetService *app.AssetService) *AssetHandler {
	return &AssetHandler{svc: assetService}
}

// EditDescription godoc
// @Summary      Edit asset description
// @Description  Updates the description of an asset.
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        asset_id  path   string                       true  "Asset ID (UUID)"
// @Param        payload   body   handlers.AssetEditRequest     true  "New description payload"
// @Success      200       {object} handlers.AssetResponse
// @Failure      400       {object} handlers.ErrorResponse
// @Failure      404       {object} handlers.ErrorResponse
// @Failure      500       {object} handlers.ErrorResponse
// @Router       /assets/{asset_id}/description [patch]
func (handler *AssetHandler) EditDescription(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(req, "asset_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteJsonError(writer, "invalid asset_id", http.StatusBadRequest)
		return
	}

	var body struct {
		Description string `json:"description"`
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		WriteJsonError(writer, "invalid json", http.StatusBadRequest)
		return
	}

	a, err := handler.svc.EditDescription(req.Context(), id, body.Description)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAssetNotFound):
			WriteJsonError(writer, "asset not found", http.StatusNotFound)
			return
		case errors.Is(err, domain.ErrEmptyDescription):
			WriteJsonError(writer, "description cannot be empty", http.StatusBadRequest)
			return
		default:
			WriteJsonError(writer, "internal error", http.StatusInternalServerError)
			return
		}
	}

	_ = json.NewEncoder(writer).Encode(AssetResponse{
		ID:          a.ID,
		Type:        a.Type,
		Description: a.Description,
		Payload:     a.Payload,
	})
}
