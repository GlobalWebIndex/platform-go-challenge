package routes

import (
	"github.com/artemis13/platform-go-challenge/handlers"

	"github.com/labstack/echo/v4"

	myMiddleware "github.com/artemis13/platform-go-challenge/middleware" // Import my custom middleware
)

func RegisterRoutes(e *echo.Echo) {

	// Public routes (no authentication required)
	e.GET("/public", handlers.PublicHandler)

	// Group routes that require authentication
	authGroup := e.Group("/users")
	authGroup.Use(myMiddleware.AuthMiddleware)

	// Protected routes (require authentication)
	authGroup.GET("/:id/favorites", handlers.GetUserFavorites)
	authGroup.POST("/:id/favorites", handlers.AddUserFavorite)
	authGroup.PUT("/:id/favorites/:asset_id", handlers.EditUserFavorite)
	authGroup.DELETE("/:id/favorites/:asset_id", handlers.RemoveUserFavorite)
}
