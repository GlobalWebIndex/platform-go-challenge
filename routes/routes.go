package routes

import (
	"github.com/gin-gonic/gin"
	"platform-go-challenge/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Public endpoint to get a token
	r.POST("/login", handlers.GenerateToken)

	// Protected routes with JWT middleware
	auth := r.Group("/")
	auth.Use(handlers.JWTAuthMiddleware())
	auth.GET("/favorites/:user_id", handlers.GetFavoritesHandler)
	auth.POST("/favorites", handlers.AddFavoriteHandler)
	auth.DELETE("/favorites/:user_id/:asset_id", handlers.RemoveFavoriteHandler)
	auth.PUT("/favorites/:user_id/:asset_id", handlers.EditFavoriteHandler)

	return r
}
