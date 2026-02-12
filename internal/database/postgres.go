package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const schema = `
	CREATE TABLE IF NOT EXISTS favourites (
		id         TEXT        NOT NULL,
		user_id    TEXT        NOT NULL,
		asset_type TEXT        NOT NULL,
		description TEXT,
		data       JSONB,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		PRIMARY KEY (user_id, id)
	);
`

// Connect opens a PostgreSQL connection pool, verifies connectivity,
// initialises the schema, and returns the ready-to-use *sql.DB.
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// Connection pool defaults, normally these values could be made configurable in production.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	if _, err := db.ExecContext(ctx, schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("initializing schema: %w", err)
	}

	return db, nil
}

// PingDB checks database connectivity. Intended for health check endpoints.
func PingDB(ctx context.Context) error {
	return DB.PingContext(ctx)
}
