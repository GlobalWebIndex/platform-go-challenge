package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexandros-georgantas/platform-go-challenge/internal/serializers"
	"github.com/alexandros-georgantas/platform-go-challenge/internal/services"
	"github.com/gin-gonic/gin"
)

type FavoriteController interface {
	AddToFavorites(c *gin.Context)
	DeleteFavorite(c *gin.Context)
	GetFavorite(c *gin.Context)
	GetFavorites(c *gin.Context)
}

type favoriteController struct {
	favoriteService services.FavoriteService
}

func NewFavoriteController(favoriteService services.FavoriteService) (FavoriteController, error) {
	return &favoriteController{favoriteService: favoriteService}, nil
}

func (fc *favoriteController) AddToFavorites(c *gin.Context) {
	var aF serializers.AddToFavorites
	pUId, _ := strconv.Atoi(c.Param("userId"))
	cUId := c.MustGet("userId").(int)
	bErr := c.BindJSON(&aF)
	if pUId != cUId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token action"})
		return
	}

	favorite, aErr := fc.favoriteService.AddToFavorites(uint(pUId), uint(aF.ID))

	if aErr != nil || bErr != nil {
		em := aErr.Error()

		if em == "ERROR: duplicate key value violates unique constraint \"user_asset\" (SQLSTATE 23505)" {

			c.JSON(http.StatusConflict, gin.H{
				"error": em,
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": em,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"favorite": favorite,
	})
}

func (fc *favoriteController) DeleteFavorite(c *gin.Context) {
	pUId, _ := strconv.Atoi(c.Param("userId"))
	pFId, _ := strconv.Atoi(c.Param("favoriteId"))
	cUId := c.MustGet("userId").(int)

	if pUId != cUId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token action"})
		return
	}

	dId, aErr := fc.favoriteService.RemoveFromFavorites(uint(pFId))

	if aErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("something went wrong while removing asset with id %v from favorites", pFId).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deletedId": dId,
	})
}

func (fc *favoriteController) GetFavorite(c *gin.Context) {
	pUId, _ := strconv.Atoi(c.Param("userId"))
	pFId, _ := strconv.Atoi(c.Param("favoriteId"))

	cUId := c.MustGet("userId").(int)

	if pUId != cUId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token action"})
		return
	}

	favorite, aErr := fc.favoriteService.GetFavorite(uint(pUId), uint(pFId))

	if aErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("something went wrong while fetching favorite  with id %v ", pFId).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"favorite": favorite,
	})
}
func (fc *favoriteController) GetFavorites(c *gin.Context) {
	pUId, _ := strconv.Atoi(c.Param("userId"))

	cUId := c.MustGet("userId").(int)

	if pUId != cUId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token action"})
		return
	}

	favorites, aErr := fc.favoriteService.GetFavorites(uint(pUId))

	if aErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("something went wrong while fetching favorites for user with id %v ", pUId).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"favorites": favorites,
	})
}
