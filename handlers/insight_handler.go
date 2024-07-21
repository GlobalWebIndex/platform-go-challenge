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

type InsightHandler struct {
	service database.InsightService
}

func NewInsightHandler(service database.InsightService) *InsightHandler {
	return &InsightHandler{service: service}
}

func (h *InsightHandler) GetAll(ctx *gin.Context) {
	insights, err := h.service.GetAll(nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if insights == nil {
		// will return "null" if empty, with this "trick" we return "[]" json.
		insights = make([]models.Insight, 0)
	}

	ctx.JSON(http.StatusOK, insights)
}

func (h *InsightHandler) Get(ctx *gin.Context) {
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

func (h *InsightHandler) Add(ctx *gin.Context) {

	byteValue, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	var result *models.Insight
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

func (h *InsightHandler) AddAll(ctx *gin.Context) {

	byteValue, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		ctx.Abort()
		return
	}
	var result []*models.Insight
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
