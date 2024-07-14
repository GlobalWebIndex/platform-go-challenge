package server

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"
)

// MockDB is a mock implementation of the database.DB interface
type MockDB struct{}

func (m *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}

func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (m *MockDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, nil
}

func (m *MockDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return nil, nil
}

// Mock InitDB function
func mockInitDB() error {
	return nil
}

func TestNewApp(t *testing.T) {
	app := NewApp()

	if app == nil {
		t.Error("Expected NewApp to return a non-nil App")
	}

	if app.server == nil {
		t.Error("Expected App to have a non-nil server")
	}

	if app.server.Addr != ":8080" {
		t.Errorf("Expected server address to be :8080, got %s", app.server.Addr)
	}
}

func TestAppStart(t *testing.T) {
	app := NewApp()

	go func() {
		err := app.Start()
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Unexpected error from app.Start(): %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/ping")
	if err != nil {
		t.Fatalf("Failed to send request to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}

func TestAppShutdown(t *testing.T) {
	app := NewApp()

	go app.Start()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.Shutdown(ctx)
	if err != nil {
		t.Errorf("Unexpected error from app.Shutdown(): %v", err)
	}
}
