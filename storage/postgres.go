package storage

import (
	"fmt"
	"time"

	"github.com/arhsxro/platform-go-challenge/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Creates a new db connection.
func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		cfg.DBUsername, cfg.DBName, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	//These values can be set accordingly to the usecase in order to optimize resource usage
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(15 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (store *PostgresStore) Close() error {
	if store.db != nil {
		return store.db.Close()
	}
	return nil
}
