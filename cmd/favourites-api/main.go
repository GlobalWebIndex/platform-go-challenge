package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/flowchartsman/swaggerui"
	amem "github.com/gwi/platform-go-challenge/assets/adapters/memory"
	"github.com/gwi/platform-go-challenge/config"
	"github.com/gwi/platform-go-challenge/favourites/adapters/assetcatalog"
	fmem "github.com/gwi/platform-go-challenge/favourites/adapters/memory"
	"github.com/gwi/platform-go-challenge/favourites/application"
	httppkg "github.com/gwi/platform-go-challenge/http"
)

func initLogger() {
	levelStr := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	if levelStr == "" {
		levelStr = "info"
	}
	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
}

func main() {
	initLogger()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = filepath.Join("config", env+".yaml")
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("config load failed", "err", err)
		os.Exit(1)
	}
	slog.Info("config loaded", "path", configPath, "env", env)

	signer, verifier, err := httppkg.LoadJWTAuth()
	if err != nil {
		slog.Error("JWT auth init failed", "err", err)
		os.Exit(1)
	}
	internalEventSecret := os.Getenv("INTERNAL_EVENT_SECRET")

	repo := fmem.NewRepository()
	catalog := amem.NewInMemoryCatalog()
	assetCatalogClient := assetcatalog.NewClient(catalog)
	svc := application.NewFavouritesService(repo, assetCatalogClient)
	openAPIServer := &httppkg.OpenAPIServer{
		UC:           svc,
		EventsUC:     svc,
		EventsSecret: internalEventSecret,
		Signer:       signer,
		Verifier:     verifier,
		TokenExpiry:  cfg.TokenExpiry(),
	}

	specPath := os.Getenv("OPENAPI_SPEC")
	if specPath == "" {
		specPath = "swagger.yaml"
	}
	spec, err := os.ReadFile(specPath)
	if err != nil {
		if os.IsNotExist(err) {
			alt := filepath.Join(filepath.Dir(configPath), "..", "swagger.yaml")
			spec, err = os.ReadFile(alt)
		}
		if err != nil {
			slog.Warn("could not load OpenAPI spec, Swagger UI disabled", "path", specPath, "err", err)
			spec = nil
		}
	}

	mux := http.NewServeMux()
	if len(spec) > 0 {
		mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
		mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
		})
	}
	handler := httppkg.HandlerFromMuxWithBaseURL(openAPIServer, mux, "/v1")
	handler = httppkg.LimitRequestSize(handler, cfg.Server.MaxRequestBodyBytes)
	handler = httppkg.RateLimitMiddleware(handler, cfg.Server.RateLimitPerSecond, cfg.Server.RateLimitBurst)

	srv := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout(),
		WriteTimeout: cfg.WriteTimeout(),
		IdleTimeout:  cfg.IdleTimeout(),
	}

	go func() {
		slog.Info("server listening", "addr", cfg.Addr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout())
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "err", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}