// @title GWI Favourite Assets API
// @version 1.0
// @description This is a simple API for managing user favourite assets.
// @contact.name F. Daniil
// @contact.email daniilfanis@gmail.com
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"platform-go-challenge/internal/repository"
	"platform-go-challenge/internal/seed"
	"platform-go-challenge/internal/service"
	"platform-go-challenge/routes"
	"syscall"
	"time"

	"platform-go-challenge/docs" // This import is crucial!

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	defaultServerAddr = "0.0.0.0:8080"
	shutdownTimeout   = 30 * time.Second
)

func main() {
	// Set Gin to release mode for production
	// This disables debug output.
	gin.SetMode(gin.ReleaseMode)

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = defaultServerAddr
	}

	logger := InitLogger()
	defer func() {
		if err := logger.Sync(); err != nil && err.Error() != "sync" {
			log.Printf("failed to sync logger: %v", err)
		}
	}()

	// Initialize the repository implementation.
	favouriteRepo := repository.NewInMemoryRepository()
	favouriteService := service.NewFavouriteService(favouriteRepo, logger)

	seed.Seed(context.Background(), favouriteRepo, logger)

	r := routes.NewRouter(favouriteService)

	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r, // Use the Gin router as the handler
		// Best practice timeouts
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	// Graceful shutdown logic remains the same.
	go func() {
		logger.Info("server starting", zap.String("addr", serverAddr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("could not start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exiting gracefully")
}

func InitLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return logger
}
