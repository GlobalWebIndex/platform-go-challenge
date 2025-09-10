package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/SokratisChaimanas/platform-go-challenge/docs"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/adapters/ent" // ent adapters
	chihttp "github.com/SokratisChaimanas/platform-go-challenge/internal/adapters/http/chi"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/app"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/platform/config"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/platform/db"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/shared/logger"
)

func main() {
	cfg := config.LoadFromEnv()

	// Logger
	level := parseLevel(cfg.LogLevel)
	log := logger.NewFileLogger(level, cfg.LogPath, 10) // 10MB rotate
	log.Info("starting service",
		"env", cfg.AppEnv,
		"http_addr", cfg.HTTPAddr,
		"db_host", cfg.DBHost,
		"db_name", cfg.DBName,
	)

	// Open DB (Ent) with a startup timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	entClient, err := db.NewClient(ctx, cfg)
	if err != nil {
		log.Error("failed to initialize database client", "err", err)
		os.Exit(1)
	}
	defer func() {
		if cerr := entClient.Close(); cerr != nil {
			log.Error("failed to close ent client", "err", cerr)
		}
	}()
	log.Info("database client initialized")

	// Wire adapters (repos)
	userRepo := entadapter.NewUserRepo(entClient)
	assetRepo := entadapter.NewAssetRepo(entClient)
	favRepo := entadapter.NewFavouriteRepo(entClient)

	// Wire services (use cases)
	userSvc := app.NewUserService(userRepo)
	assetSvc := app.NewAssetService(assetRepo)
	favSvc := app.NewFavouritesService(userRepo, assetRepo, favRepo)

	// Build HTTP router
	router := chihttp.NewRouter(userSvc, assetSvc, favSvc)

	// HTTP server
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Run server in background
	go func() {
		log.Info("http server listening", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server error", "err", err)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Info("shutdown signal received", "signal", sig.String())

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Stop accepting new connections, wait for in-flight to complete or timeout.
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown error", "err", err)
	} else {
		log.Info("server shutdown complete")
	}
}

// parseLevel converts a string like "debug", "info", "warn", "error" to slog.Level.
func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "err", "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
