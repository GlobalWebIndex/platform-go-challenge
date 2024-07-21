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

type AudienceHandler struct {
	service database.AudienceService
}

func NewAudienceHandler(service database.AudienceService) *AudienceHandler {
	return &AudienceHandler{service: service}
}

func (h *AudienceHandler) GetAll(ctx *gin.Context) {
	audiences, err := h.service.GetAll(nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	if audiences == nil {
		// will return "null" if empty, with this "trick" we return "[]" json.
		audiences = make([]models.Audience, 0)
	}

	ctx.JSON(http.StatusOK, audiences)
}

func (h *AudienceHandler) Get(ctx *gin.Context) {
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

func (h *AudienceHandler) Add(ctx *gin.Context) {

	byteValue, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	var result *models.Audience
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
		ctx.Abort()
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *AudienceHandler) AddAll(ctx *gin.Context) {

	byteValue, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	var result []*models.Audience
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
