package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"pgc/app"
	"pgc/internal/config"
	"pgc/persistence/postgres"
	"pgc/presentation"
)

// Run initializes the GWI application, sets up the database, services, and HTTP server,
// and starts the server with graceful shutdown capabilities.
func Run(ctx context.Context, cfg *config.Config) error {
	slog.Info("Starting GWI application...")

	db, err := setupDatabase(cfg.Database)
	if err != nil {
		return err
	}
	defer closeDatabase(db)

	fs, err := setupServices(db)
	if err != nil {
		return err
	}

	httpServer := setupHTTPServer(ctx, cfg.Server, fs)

	return runServer(ctx, httpServer, cfg.Server.ShutdownTimeout)
}

// setupDatabase initializes the database connection using the provided configuration.
func setupDatabase(cfg config.DatabaseConfig) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithNetwork(cfg.Network),
		pgdriver.WithAddr(cfg.Address),
		pgdriver.WithTLSConfig(nil),
		pgdriver.WithUser(cfg.User),
		pgdriver.WithPassword(cfg.Password),
		pgdriver.WithDatabase(cfg.Name),
		pgdriver.WithApplicationName(cfg.AppName),
	))

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", "error", err)
		return nil, err
	}

	return db, nil
}

// closeDatabase closes the database connection and logs any errors.
func closeDatabase(db *bun.DB) {
	slog.Info("closing database connection")
	if err := db.Close(); err != nil {
		slog.Error("failed to close database connection", "error", err)
	}
}

// setupServices initializes the application services with the database repository.
func setupServices(db *bun.DB) (*app.FavoriteService, error) {
	fr, err := postgres.NewFavoriteRepository(db)
	if err != nil {
		slog.Error("failed to initialize favorite repository", "error", err)
		return nil, err
	}

	return app.NewFavoriteService(fr), nil
}

// setupHTTPServer initializes the HTTP server with routes and context.
func setupHTTPServer(ctx context.Context, cfg config.ServerConfig, fs *app.FavoriteService) *http.Server {
	mux := http.NewServeMux()
	presentation.RegisterRoutes(mux, fs)

	return &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
}

// runServer starts the HTTP server and handles graceful shutdown.
func runServer(ctx context.Context, httpServer *http.Server, shutdownTimeout time.Duration) error {
	// Create an error channel for server startup errors
	serverErr := make(chan error, 1)

	// Start the HTTP server in a separate goroutine
	go func() {
		slog.Info("starting HTTP server", slog.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		slog.Info("HTTP server shut down gracefully")
	}()

	// Wait for either shutdown signal OR server error
	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received - shutting down service")
	case err := <-serverErr:
		slog.Error("HTTP server failed to start", "error", err.Error())
		return err
	}

	// Create a timeout context for shutting down the server gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt to gracefully shut down the HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("server shutdown timed out", "error", err.Error())
		} else {
			slog.Error("failed to gracefully shutdown HTTP server", "error", err.Error())
		}
		return err
	}

	slog.Info("server shutdown complete")
	return nil
}
