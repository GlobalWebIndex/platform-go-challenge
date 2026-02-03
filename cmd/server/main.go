package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go-favorites-app/internal/adapter/api/rest"
	"go-favorites-app/internal/adapter/cache/redis"
	repo "go-favorites-app/internal/adapter/storage/postgres"
	"go-favorites-app/internal/config"
	"go-favorites-app/internal/core/domain/favorites"
	"go-favorites-app/internal/core/service"
	"go-favorites-app/internal/observability"
)

// -- MAIN --

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, relying on environment variables")
	}

	// Config
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Init Tracing
	tpShutdown, err := observability.InitTracerProvider(ctx, "favorites-service", cfg.OtelExporterEndpoint)
	if err != nil {
		logger.Error("failed to init tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := tpShutdown(ctx); err != nil {
			logger.Error("failed to shutdown tracer", "error", err)
		}
	}()

	// Init DB
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Run Migrations (Apply on Startup)
	if err := repo.RunMigrations(ctx, dbPool, logger); err != nil {
		logger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Metrics: DB Stats Poller
	observability.StartDBStatsCollector(dbPool)

	// Init Cache
	redisAdapter := redis.NewAdapter(cfg.RedisAddr)
	// Wrap with metrics
	cacheSvc := observability.NewInstrumentedCache(redisAdapter)

	// Init Service
	// Mock Enricher for now
	enricher := &NoOpEnricher{}

	// Repository Init
	favRepo := repo.NewRepository(dbPool)
	userRepo := repo.NewUserRepository(dbPool)

	// Service Init
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	favSvc := service.NewService(favRepo, cacheSvc, enricher, logger)

	// Init Handlers
	favHandler := rest.NewHandler(favSvc, logger)
	authHandler := rest.NewAuthHandler(authSvc)

	// Init Router
	router := rest.NewRouter(favHandler, authHandler, cfg.JWTSecret, rest.RequestID, rest.Logger(logger), observability.Middleware)

	// Add /metrics endpoint
	// Note: Usually /metrics is on a separate admin port or protected, adding to main mux for simplicity
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", router)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	// Graceful Shutdown
	go func() {
		logger.Info("Starting server", "addr", srv.Addr, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}

// Simple NoOp Enricher
type NoOpEnricher struct{}

func (e *NoOpEnricher) Enrich(ctx context.Context, asset favorites.Asset) error { return nil }
