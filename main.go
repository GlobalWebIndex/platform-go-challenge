package main

import (
	"log"
	"net/http"

	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/handlers"
	"platform-go-challenge/internal/storage"
)

func main() {
	// Load config
	cfg := config.NewConfig()

	// Initialize storage
	users := storage.NewUserStore()
	favorites := storage.NewFavoritesStore()

	// Initialize router from handlers package
	r := handlers.SetupRouter(cfg, favorites, users)

	// Start server
	log.Printf("Server starting on %s...\n", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
