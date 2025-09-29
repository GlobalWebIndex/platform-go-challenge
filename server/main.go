package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"favourite_assets/server/authentication"
	"favourite_assets/server/controllers"
	"favourite_assets/server/repositories"
	"favourite_assets/server/routes"
	"favourite_assets/server/services"
)

func main() {
	// --- Initialize repositories ---
	userRepo := repositories.NewUserRepository()
	assetRepo := repositories.NewAssetRepository()
	favRepo := repositories.NewFavoriteRepository()

	// --- Initialize services ---
	userService := services.NewUserService(userRepo)
	assetService := services.NewAssetService(assetRepo)
	favService := services.NewFavouriteService(favRepo, userService, assetService)

	// --- Initialize Keycloak service ---
	keycloakService := services.NewKeycloakService()

	// --- Initialize controllers ---
	userController := controllers.NewUserController(userService)
	assetController := controllers.NewAssetController(assetService)
	favController := controllers.NewFavouriteController(favService)

	// --- Setup router ---
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- Register routes ---
	routes.RegisterRoutes(r, userController, assetController, favController, authentication.KeycloakAuth(keycloakService))

	// --- Start server ---
	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
