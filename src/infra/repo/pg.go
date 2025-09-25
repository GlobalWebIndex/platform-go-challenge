package repo

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PG struct{ DB *sql.DB }

func NewPG(dsn string) (*PG, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Simple, sensible pool defaults
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &PG{DB: db}, nil
}

func (p *PG) Close() error { return p.DB.Close() }
