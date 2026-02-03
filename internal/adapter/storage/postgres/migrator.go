package postgres

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations executes the embedded SQL migration files.
// For a real production app, use golang-migrate or goose.
// For this challenge, we'll just execute the up script ensuring idempotency (IF NOT EXISTS).
func RunMigrations(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger) error {
	logger.Info("running database migrations")

	// 1. Create Favorites Table
	content, err := migrationsFS.ReadFile("migrations/000001_create_favorites_table.up.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file 1: %w", err)
	}
	if _, err := db.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("failed to execute migration 1: %w", err)
	}

	// 2. Add Users Table
	contentUsers, err := migrationsFS.ReadFile("migrations/000002_add_users_table.up.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file 2: %w", err)
	}
	if _, err := db.Exec(ctx, string(contentUsers)); err != nil {
		return fmt.Errorf("failed to execute migration 2: %w", err)
	}

	logger.Info("migrations completed successfully")
	return nil
}
