package database

import (
	"context"
	"database/sql"
	"fmt"
	"gwi-platform/utils"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type DB interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// InitDB initializes the database connection
func InitDB() error {
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPassword,
		Net:                  "tcp",
		Addr:                 dbHost,
		DBName:               dbName,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	utils.ErrorLogger.Println("Connected to the database successfully")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			utils.ErrorLogger.Printf("Error closing database connection: %v", err)
		} else {
			utils.InfoLogger.Println("Database connection closed")
		}
	}
}

// GetDB returns the database connection
func GetDB() DB {
	if db == nil {
		utils.ErrorLogger.Println("Database connection is not initialized")
	}
	return db
}
