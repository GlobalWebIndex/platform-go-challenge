package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"platform-go-challenge/db"
	"platform-go-challenge/handlers"
	"syscall"
	"time"

	_ "platform-go-challenge/docs" // Import swagger docs

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/joho/godotenv"
)

// @title Platform Go Challenge API
// @version 1.0
// @description REST API for managing users, assets, and favorites with JWT authentication
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the 'Bearer ' prefix, e.g. 'Bearer abc123'

func initDatabase() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	database, err := db.New(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return database, nil
}

func initServer(database *sql.DB) *http.Server {
	mux := http.NewServeMux()

	// Swagger documentation
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Health
	mux.HandleFunc("/health", handlers.HealthCheck)

	// Public routes
	mux.HandleFunc("/login", handlers.Login(database))
	mux.HandleFunc("/register", handlers.Register(database))

	// Protected routes
	mux.HandleFunc("/users", handlers.AuthMiddleware(handlers.UserRouter(database)))
	mux.Handle("/users/", handlers.AuthMiddleware(handlers.FavouritesRouter(database)))
	mux.HandleFunc("/assets", handlers.AuthMiddleware(handlers.AssetsRouter(database)))
	mux.Handle("/assets/", handlers.AuthMiddleware(handlers.AssetsRouter(database)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

func main() {
	// Load .env (no-op in prod)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found (this is fine in prod)")
	}

	// ---- database ----
	database, err := initDatabase()
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}
	defer database.Close()

	// ---- server ----
	server := initServer(database)

	// ---- graceful shutdown ----
	go func() {
		log.Println("🚀 Server running on port: " + os.Getenv("PORT"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("✅ Server exited cleanly")
}
