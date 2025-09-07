package routes

import (
	"platform-go-challenge/internal/handlers"
	"platform-go-challenge/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(s *service.FavouriteService) *gin.Engine {
	router := gin.Default()

	favouritesHandler := handlers.NewFavouriteHandler(s)
	healthHandler := handlers.NewHealthHandler()

	router.GET("/health", healthHandler)
	favorites := router.Group("/users/:userID/favourites")
	{
		favorites.GET("", favouritesHandler.GetFavourites)
		favorites.POST("", favouritesHandler.AddFavourite)
		favorites.DELETE("/:assetID", favouritesHandler.DeleteFavourite)
		favorites.PATCH("/:assetID", favouritesHandler.UpdateFavouriteDescription)
	}

	return router
}
