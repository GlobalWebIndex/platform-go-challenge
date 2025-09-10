package db

import (
	"context"
	"database/sql" // stdlib
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/SokratisChaimanas/platform-go-challenge/ent"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/platform/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewClient(ctx context.Context, cfg *config.Config) (*ent.Client, error) {
	dsn := buildDSN(cfg)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("running schema migration: %w", err)
	}

	if err := seedDevOnce(ctx, client, cfg.AppEnv); err != nil {
		return nil, fmt.Errorf("seeding dev data: %w", err)
	}

	return client, nil
}

// buildDSN constructs the Postgres connection string (DSN) from config.
// Example: postgres://user:pass@host:5432/dbname?sslmode=disable
func buildDSN(cfg *config.Config) string {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
	)

	return dsn
}
