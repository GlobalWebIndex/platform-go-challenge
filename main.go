package main

import (
	"go-challenge/handlers"
	"go-challenge/repository"
	"log"
	"net/http"
)

func main() {
	// Init repository.
	favRepo := repository.NewInMemoryFavouriteRepository()

	// Init handler.
	favHandler := handlers.NewFavouriteHandler(favRepo)

	// Setup routes
	mux := http.NewServeMux()

	// FAVOURITE ENDPOINTS
	mux.HandleFunc("GET /api/users/{userID}/favourites", favHandler.GetUserFavourites)
	mux.HandleFunc("POST /api/users/{userID}/favourites", favHandler.AddFavourite)
	mux.HandleFunc("DELETE /api/users/{userID}/favourites/{assetID}", favHandler.RemoveFavourite)

	mux.HandleFunc("PATCH /api/assets/{assetID}/description", favHandler.UpdateAssetDescription)

	// Start server
	log.Println("Server starting on :8080")
	log.Println("")
	log.Println("  Endpoints:")
	log.Println("  GET    /api/users/{userID}/favourites               - Get user's favourites")
	log.Println("  POST  /api/users/{userID}/favourites                - Add asset to favourites")
	log.Println("  DELETE /api/users/{userID}/favourites/{assetID}     - Remove favourite")
	log.Println("  PATCH  /api/assets/{assetID}/description            - Edit asset description")
	log.Println("")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
