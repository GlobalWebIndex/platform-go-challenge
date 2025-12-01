package routes

import (
    "net/http"
    "github.com/dintzeler/platform-go-challenge/controllers"
    "github.com/dintzeler/platform-go-challenge/middleware"
)

func SetupFavoritesRoutes() {
    http.HandleFunc("/favorites", middleware.JWTMiddleware(controllers.FavoritesHandler))
}