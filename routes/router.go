package routes

import (
	"platform-go-challenge/internal/handlers"
	"platform-go-challenge/internal/service"

	"platform-go-challenge/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Serve the Swagger UI
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
