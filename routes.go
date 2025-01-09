package main

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger" // Swagger handler
	"platform-go-challenge/handlers"
)

// ConfigureRoutes sets up the routes for the application.
func ConfigureRoutes(userFavouritesHandler handlers.UserFavouritesHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users/{id}/favourites", userFavouritesHandler.GetUserFavourites).Methods("GET")
	r.HandleFunc("/users/{id}/favourites", userFavouritesHandler.AddAssetToFavourites).Methods("POST")
	r.HandleFunc("/users/{id}/favourites/{asset_id}", userFavouritesHandler.RemoveAssetFromFavourites).Methods("DELETE")
	r.HandleFunc("/users/{id}/favourites/{asset_id}", userFavouritesHandler.EditAssetDescription).Methods("PUT")

	// Add the Swagger endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
