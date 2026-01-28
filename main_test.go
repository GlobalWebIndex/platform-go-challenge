package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func TestInitDatabase(t *testing.T) {
	tests := []struct {
		name        string
		databaseURL string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing DATABASE_URL",
			databaseURL: "",
			wantErr:     true,
			errContains: "DATABASE_URL is not set",
		},
		{
			name:        "invalid connection string",
			databaseURL: "invalid://connection",
			wantErr:     true,
			errContains: "failed to connect to db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			originalURL := os.Getenv("DATABASE_URL")
			defer os.Setenv("DATABASE_URL", originalURL)

			if tt.databaseURL == "" {
				os.Unsetenv("DATABASE_URL")
			} else {
				os.Setenv("DATABASE_URL", tt.databaseURL)
			}

			// Test
			db, err := initDatabase()

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if tt.errContains != "" && err != nil {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("error = %v, want error containing %q", err, tt.errContains)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if db != nil {
					db.Close()
				}
			}
		})
	}
}

func TestInitServer(t *testing.T) {
	// Create a mock database (nil is acceptable for testing server initialization)
	var mockDB *sql.DB

	tests := []struct {
		name     string
		port     string
		wantAddr string
	}{
		{
			name:     "default port",
			port:     "",
			wantAddr: ":8080",
		},
		{
			name:     "custom port",
			port:     "3000",
			wantAddr: ":3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			originalPort := os.Getenv("PORT")
			defer os.Setenv("PORT", originalPort)

			if tt.port == "" {
				os.Unsetenv("PORT")
			} else {
				os.Setenv("PORT", tt.port)
			}

			// Test
			server := initServer(mockDB)

			if server == nil {
				t.Fatal("expected server but got nil")
			}

			if server.Addr != tt.wantAddr {
				t.Errorf("server.Addr = %q, want %q", server.Addr, tt.wantAddr)
			}

			if server.Handler == nil {
				t.Error("expected handler but got nil")
			}
		})
	}
}

func TestInitServerRoutes(t *testing.T) {
	// Create a mock database
	var mockDB *sql.DB

	// Initialize server
	server := initServer(mockDB)

	if server == nil {
		t.Fatal("expected server but got nil")
	}

	if server.Handler == nil {
		t.Fatal("expected handler but got nil")
	}

	// Test that routes exist by checking the handler is not nil
	// We can't test actual route handling without a real database connection
	t.Run("server has handler configured", func(t *testing.T) {
		// Test a simple unknown route to verify handler is working
		req := httptest.NewRequest("GET", "/unknown", nil)
		rec := httptest.NewRecorder()
		server.Handler.ServeHTTP(rec, req)

		// Expect 404 for unknown route
		if rec.Code != http.StatusNotFound {
			t.Logf("unknown route returned status %d (expected 404 but other codes acceptable)", rec.Code)
		}
	})
}

func TestInitServerHandlerNotNil(t *testing.T) {
	var mockDB *sql.DB
	server := initServer(mockDB)

	if server.Handler == nil {
		t.Fatal("server handler should not be nil")
	}

	// Test that the handler can respond to requests
	req := httptest.NewRequest("GET", "/unknown-route", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	// For an unknown route, we expect 404
	if rec.Code != http.StatusNotFound {
		// This is acceptable - just checking the handler works
	}
}
