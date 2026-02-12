package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/giannis84/platform-go-challenge/internal"
	"github.com/giannis84/platform-go-challenge/internal/config"
	"github.com/giannis84/platform-go-challenge/internal/database"
	"github.com/giannis84/platform-go-challenge/internal/logging"
	"github.com/giannis84/platform-go-challenge/internal/routes"
)

func main() {
	// Initialize shared dependencies
	logger := logging.NewLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", slog.String(logging.ErrorKey, err.Error()))
		os.Exit(1)
	}
	logger.Info("configuration loaded",
		slog.String("api_addr", cfg.APIAddr()),
		slog.String("health_addr", cfg.HealthAddr()),
	)

	// Connect to PostgreSQL and initialise schema
	db, err := database.Connect(cfg.PostgresConnString())
	if err != nil {
		logger.Error("failed to initialise database", slog.String(logging.ErrorKey, err.Error()))
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database ready")

	// Create health check and favourites http services
	healthService := &internal.Service{
		Addr:         cfg.HealthAddr(),
		Logger:       logger,
		DB:           db,
		Routes:       routes.RegisterHealthRoutes(cfg.RateLimitConfig()),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	healthService.Init()

	apiService := &internal.Service{
		Addr:         cfg.APIAddr(),
		Logger:       logger,
		DB:           db,
		Routes:       routes.RegisterFavouritesRoutes(cfg.AuthConfig(), cfg.RateLimitConfig()),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	apiService.Init()

	// Start http service threads
	go func() {
		if err := healthService.ListenAndServeWrapper("health check api"); err != nil && err != http.ErrServerClosed {
			logger.Error("health check service failed", slog.String(logging.ErrorKey, err.Error()))
			os.Exit(1)
		}
	}()
	go func() {
		if err := apiService.ListenAndServeWrapper("favourites api"); err != nil && err != http.ErrServerClosed {
			logger.Error("favourites service failed", slog.String(logging.ErrorKey, err.Error()))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	receivedSignal := <-quit

	// Shutdown http service threads gracefully
	logger.Info("shutting down service", slog.Any("OS signal received", os.Signal.String(receivedSignal)))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := apiService.HTTPServer.Shutdown(ctx); err != nil {
		logger.Error("API service shutdown error", slog.String(logging.ErrorKey, err.Error()))
	}
	if err := healthService.HTTPServer.Shutdown(ctx); err != nil {
		logger.Error("health service shutdown error", slog.String(logging.ErrorKey, err.Error()))
	}
	logger.Info("exiting...")
}