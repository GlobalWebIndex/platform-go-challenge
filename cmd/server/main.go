package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"platform-go-challenge/internal/auth"
	"platform-go-challenge/internal/handler"
	"platform-go-challenge/internal/middleware"
	"platform-go-challenge/internal/storage"
)

func main() {
	// Initialize storage
	store := storage.NewMemoryStore()

	// Initialize JWT manager
	jwtConfig := auth.Config{
		SecretKey:  os.Getenv("JWT_SECRET"),
		ExpiryTime: 24 * time.Hour,
	}

	if jwtConfig.SecretKey == "" {
		jwtConfig.SecretKey = "your-secret-key-change-this-in-production"
		fmt.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable in production")
	}

	jwtManager := auth.NewJWTManager(jwtConfig)

	// Create metrics collector
	metricsCollector := middleware.NewMetrics()

	// Create router
	router := http.NewServeMux()

	// Public endpoints (no auth required)
	router.HandleFunc("GET /api/health", handler.HealthCheck(metricsCollector))
	router.HandleFunc("POST /api/login", handler.Login(jwtManager))

	// Protected endpoints (auth required)
	authMiddleware := middleware.AuthMiddleware(jwtManager)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /api/users/{id}/favorites", handler.GetFavorites(store))
	protectedMux.HandleFunc("POST /api/users/{id}/favorites", handler.AddFavorite(store))
	protectedMux.HandleFunc("DELETE /api/users/{id}/favorites/{assetId}", handler.RemoveFavorite(store))
	protectedMux.HandleFunc("PUT /api/users/{id}/favorites/{assetId}", handler.EditFavorite(store))

	router.Handle("/api/users/", authMiddleware(protectedMux))

	// Start server
	port := ":8080"
	
	loggingHandler := middleware.RequestLogger(metricsCollector)(router)
	
	server := &http.Server{
		Addr:    port,
		Handler: loggingHandler,
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

	fmt.Printf("Server started on %s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}