package presentation

import (
	"log/slog"
	"net/http"

	"pgc/app"
)

// RegisterRoutes sets up the HTTP routes for the FavoriteHandler.
// It maps HTTP methods and paths to their corresponding handler functions.
func RegisterRoutes(mux *http.ServeMux, fs *app.FavoriteService) {
	favoriteHandler := NewFavoriteHandler(fs)

	mux.HandleFunc("GET /favorite/", favoriteHandler.GetUserFavorites)
	slog.Info("Registered GET /favorite/")

	mux.HandleFunc("POST /favorite/{aid}", favoriteHandler.AddAssetToUserFavorites)
	slog.Info("Registered POST /favorite/{aid}")

	mux.HandleFunc("PATCH /favorite/{id}", favoriteHandler.UpdateUserFavoriteAsset)
	slog.Info("Registered PATCH /favorite/{id}")

	mux.HandleFunc("DELETE /favorite/{id}", favoriteHandler.RemoveAssetFromUserFavorites)
	slog.Info("Registered DELETE /favorite/{id}")
}
