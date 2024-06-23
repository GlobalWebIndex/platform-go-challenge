package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexandros-georgantas/platform-go-challenge/internal/serializers"
	"github.com/alexandros-georgantas/platform-go-challenge/internal/services"
	"github.com/gin-gonic/gin"
)

type AssetController interface {
	GetAsset(c *gin.Context)
	GetAssets(c *gin.Context)
	GetCharts(c *gin.Context)
	GetAudiences(c *gin.Context)
	GetInsights(c *gin.Context)
	UpdateDescription(c *gin.Context)
}

type assetController struct {
	assetService services.AssetService
}

func NewAssetController(assetService services.AssetService) (AssetController, error) {
	return &assetController{assetService: assetService}, nil
}

func (ac *assetController) GetAssets(c *gin.Context) {
	q := c.Request.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("limit"))
	t := q.Get("type")

	assets, err := ac.assetService.GetAssets(page, pageSize, t)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"assets": assets,
	})
}

func (ac *assetController) GetCharts(c *gin.Context) {
	q := c.Request.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("limit"))

	charts, err := ac.assetService.GetCharts(page, pageSize)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"charts": charts,
	})
}
func (ac *assetController) GetAudiences(c *gin.Context) {
	q := c.Request.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("limit"))

	audiences, err := ac.assetService.GetAudiences(page, pageSize)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"audiences": audiences,
	})
}

func (ac *assetController) GetInsights(c *gin.Context) {
	q := c.Request.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("limit"))

	insights, err := ac.assetService.GetInsights(page, pageSize)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"insights": insights,
	})
}

func (ac *assetController) UpdateDescription(c *gin.Context) {
	var uD serializers.UpdateDescription
	aId, _ := strconv.Atoi(c.Param("assetId"))
	bErr := c.BindJSON(&uD)
	uA, aErr := ac.assetService.UpdateDescription(uint(aId), uD.Description)

	if uD.Description == "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("invalid request").Error(),
		})
		return
	}

	if aErr != nil || bErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("something went wrong while updating description of asset with id %v", aId).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"asset": uA,
	})

}

func (ac *assetController) GetAsset(c *gin.Context) {
	aid, _ := strconv.Atoi(c.Param("assetId"))

	a, aErr := ac.assetService.GetAsset(uint(aid))

	if aErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("something went wrong while updating description of asset with id %v", aid).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"asset": a,
	})

}
