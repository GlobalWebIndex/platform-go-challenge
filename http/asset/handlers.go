package asset

import (
	"encoding/json"
	"fmt"
	"net/http"

	"platform-go-challenge/internal/app/assets"
)

type Handler struct {
	service *assets.Service
}

func NewAssetHandler(assetService *assets.Service) Handler {
	return Handler{
		service: assetService,
	}
}

// GetDashboard
func (h *Handler) ListAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	ctx := r.Context()
	assets, err := h.service.List(ctx)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(assets)
}
