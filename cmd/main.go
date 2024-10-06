package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	docs "svenvermeulen/platform-go-challenge/api/docs"
	"svenvermeulen/platform-go-challenge/internal/handler"
	"svenvermeulen/platform-go-challenge/internal/repository/audience"
	"svenvermeulen/platform-go-challenge/internal/repository/chart"
	"svenvermeulen/platform-go-challenge/internal/repository/favourite"
	"svenvermeulen/platform-go-challenge/internal/repository/insight"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Favourites API
func main() {

	// set up router to map http routes to handler functions
	router := SetupRouter()
	router.Run("0.0.0.0:8086")
}

func SetupRouter() *gin.Engine {
	// set up repositories
	favouriteRepo := favourite.NewRepository()
	audienceRepo := audience.NewRepository()
	chartRepo := chart.NewRepository()
	insightRepo := insight.NewRepository()

	// set up http handlers
	favouritesHandler := handler.NewFavouritesHandler(favouriteRepo, audienceRepo, chartRepo, insightRepo)

	log.Info("Setting up gin router")
	router := gin.Default()

	// set up documentation endpoint
	docs.SwaggerInfo.Host = "0.0.0.0:8086"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// setup routes for API endpoints
	router.GET("/favourites", favouritesHandler.GetFavourites)
	router.POST("/favourites", favouritesHandler.CreateUserFavourite)
	router.DELETE("/favourites/:favouriteid", favouritesHandler.DeleteUserFavourite)
	router.PATCH("/favourites/:favouriteid", favouritesHandler.UpdateUserFavourite)

	return router
}
