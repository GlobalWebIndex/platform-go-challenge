package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

// BasicConnection interface to database connection object.
type BasicConnection interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// BasicConnectionWithTransactions interface.
type BasicConnectionWithTransactions interface {
	BasicConnection
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Scannable interface (e.g. *sql.Rows).
type Scannable interface {
	Scan(dest ...interface{}) error
}

// Client is a wrap that enhances the standard sql client some transactional utilities.
type Client struct {
	*sql.DB
}

// New returns a new sql client.
func New(connection string) (*Client, error) {
	config, err := mysql.ParseDSN(connection)
	if err != nil {
		return nil, err
	}

	connector, err := mysql.NewConnector(config)
	if err != nil {
		return nil, err
	}

	client := sql.OpenDB(connector)
	client.SetMaxOpenConns(100)
	client.SetMaxIdleConns(5)
	client.SetConnMaxIdleTime(40 * time.Second)
	client.SetConnMaxLifetime(4 * time.Hour)

	return &Client{
		DB: client,
	}, nil
}

// Query wrapper that also counts the operations.
func (c *Client) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.QueryContext(context.Background(), query, args...)
}

// QueryContext wrapper that also counts the operations.
func (c *Client) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	r, err := c.DB.QueryContext(ctx, query, args...)

	return r, err
}

// QueryRow wrapper that also counts the operations.
func (c *Client) QueryRow(query string, args ...interface{}) *sql.Row {
	return c.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext wrapper that also counts the operations.
func (c *Client) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	r := c.DB.QueryRowContext(ctx, query, args...)

	return r
}

// Exec wrapper that also counts the operations.
func (c *Client) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.ExecContext(context.Background(), query, args...)
}

// ExecContext wrapper that also counts the operations.
func (c *Client) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	r, err := c.DB.ExecContext(ctx, query, args...)

	return r, err
}
