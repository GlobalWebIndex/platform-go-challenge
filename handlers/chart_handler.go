package handlers

import (
	"encoding/json"
	"favourites/database"
	"favourites/models"
	"favourites/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type ChartHandler struct {
	service database.ChartService
}

func NewChartHandler(service database.ChartService) *ChartHandler {
	return &ChartHandler{service: service}
}

func (h *ChartHandler) GetAll(ctx *gin.Context) {
	charts, err := h.service.GetAll(nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if charts == nil {
		// will return "null" if empty, with this "trick" we return "[]" json.
		charts = make([]models.Chart, 0)
	}

	ctx.JSON(http.StatusOK, charts)
}

func (h *ChartHandler) Get(ctx *gin.Context) {
	id, _ := ctx.Params.Get("id")

	m, err := h.service.GetByID(nil, id)
	if err != nil {
		if err.Error() == utils.ErrorNotFound {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, m)
}

func (h *ChartHandler) Add(ctx *gin.Context) {

	byteValue, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	var result *models.Chart
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	err = h.service.Create(ctx, result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *ChartHandler) AddAll(ctx *gin.Context) {

	byteValue, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	var result []*models.Chart
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	err = h.service.CreateAll(ctx, result)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.Status(http.StatusCreated)
}
