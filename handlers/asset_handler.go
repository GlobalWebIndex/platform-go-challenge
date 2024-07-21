package handlers

import (
	"favourites/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssetHandler struct {
	service database.AssetService
}

func NewAssetHandler(service database.AssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

func (h *AssetHandler) GetAll(ctx *gin.Context) {
	assets, err := h.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Found Assets", "assets": assets})
}
