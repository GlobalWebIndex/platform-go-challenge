package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"platform-go-challenge/internal/handler"
	"platform-go-challenge/internal/middleware"
	"platform-go-challenge/internal/storage"
)

// Main is the entry point for the server.
// It initializes the in-memory storage, creates a router and registers the handlers.
// Then it starts the server and sets up a goroutine to listen for shutdown signals.
func main() {
	// Initialize storage
	store := storage.NewMemoryStore()

	// Create router
	router := http.NewServeMux()

	// Register handlers
	router.HandleFunc("GET /api/health", handler.HealthCheck())
	router.HandleFunc("GET /api/users/{id}/favorites", handler.GetFavorites(store))
	router.HandleFunc("POST /api/users/{id}/favorites", handler.AddFavorite(store))
	router.HandleFunc("DELETE /api/users/{id}/favorites/{assetId}", handler.RemoveFavorite(store))
	router.HandleFunc("PUT /api/users/{id}/favorites/{assetId}", handler.EditFavorite(store))
	
	// Start server
	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: middleware.LoggingMiddleware(router),
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nShutting down server...")
		if err := server.Close(); err != nil {
			log.Printf("Error closing server: %v", err)
		}
	}()

	log.Printf("Server started on %s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}