package main

import (
	"assetsApp/internal/config"
	"assetsApp/internal/handlers"
	assetServices "assetsApp/internal/services/asset"
	favouriteServices "assetsApp/internal/services/favourite"
	"assetsApp/internal/storage"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Println("Starting the application...")
	cfg := config.LoadConfig()
	// -------------------- STORAGE --------------------
	// Uncomment one depending on which store you want

	// Memory store
	// store := storage.NewMemoryStore()

	// Postgres store
	db, err := pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	dbStore := storage.NewPostgresStore(db)
	redisClient := storage.NewRedisClient(cfg.RedisAddr)
	store := storage.NewCachedStore(dbStore, redisClient)

	// -------------------- SERVICES --------------------
	assetService := assetServices.NewAssetService(store)
	favouriteService := favouriteServices.NewFavouriteService(store)

	// -------------------- HANDLERS --------------------
	assetHandler := handlers.NewAssetHandler(assetService)
	favouriteHandler := handlers.NewFavouriteHandler(favouriteService)

	// -------------------- ROUTER --------------------
	r := mux.NewRouter()

	// Asset routes
	r.HandleFunc("/users/{userId}/assets", assetHandler.GetAssets).Methods("GET")
	r.HandleFunc("/users/{userId}/assets", assetHandler.AddAsset).Methods("POST")
	r.HandleFunc("/users/{userId}/assets/{assetId}", assetHandler.EditAsset).Methods("PUT")
	r.HandleFunc("/users/{userId}/assets/{assetId}", assetHandler.RemoveAsset).Methods("DELETE")

	// Favourite routes
	r.HandleFunc("/users/{userId}/favourites", favouriteHandler.GetFavourites).Methods("GET")
	r.HandleFunc("/users/{userId}/favourites/{assetId}", favouriteHandler.AddFavourite).Methods("POST")
	r.HandleFunc("/users/{userId}/favourites/{assetId}", favouriteHandler.RemoveFavourite).Methods("DELETE")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// -------------------- START SERVER --------------------
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
