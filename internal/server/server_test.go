package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KostasDasios/platform-go-challenge/internal/config"
	"github.com/KostasDasios/platform-go-challenge/internal/models"
)

// newTestServer constructs a server configured for testing.
// It disables logging and rate limiting, uses short timeouts, and avoids .env dependencies.
// This allows tests to run fast, deterministically, and without external side effects.
func newTestServer() *Server {
	cfg := &config.Config{
		Port:            "0",
		AppEnv:          "test",
		LogEnabled:      false,       // silence middleware logs during test runs
		RateLimitMillis: 0,           // disable rate limiting for tests
		MaxBodyBytes:    1 << 20,     // 1 MB max body size
		ReadTimeout:     2 * time.Second,
		WriteTimeout:    2 * time.Second,
		IdleTimeout:     2 * time.Second,
		LogLevel:        "info",
	}
	return NewServer(cfg)
}

// TestHealthz validates that the /healthz endpoint responds correctly and fast.
// This test ensures that the service is wired up and responds to liveness probes as expected.
func TestHealthz(t *testing.T) {
	s := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	// We use s.handler instead of s.mux to include middleware in the test flow.
	s.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status from /healthz: got=%d, want=%d", rr.Code, http.StatusOK)
	}
	if rr.Body.Len() == 0 {
		t.Fatalf("/healthz returned empty response body")
	}
}

// TestFavouritesCRUD_HTTP verifies the full HTTP flow for CRUD operations on favourites.
// It exercises POST → GET → PATCH → DELETE, ensuring that routing, validation,
// and service integration are functioning correctly.
func TestFavouritesCRUD_HTTP(t *testing.T) {
	s := newTestServer()
	user := "kostas"

	// --- CREATE ---
	body := []byte(`{"asset":{"type":"insight","text":"hello","description":"d"}}`)
	req := httptest.NewRequest(http.MethodPost, "/users/"+user+"/favourites", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	s.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("POST /favourites: unexpected status %d body=%s", rr.Code, rr.Body.String())
	}

	var created models.Favourite
	if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to unmarshal create response: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("create response missing ID")
	}

	// --- LIST ---
	req = httptest.NewRequest(http.MethodGet, "/users/"+user+"/favourites", nil)
	rr = httptest.NewRecorder()
	s.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("GET /favourites: unexpected status %d", rr.Code)
	}

	// --- UPDATE (PATCH) ---
	patch := []byte(`{"description":"updated"}`)
	req = httptest.NewRequest(http.MethodPatch, "/users/"+user+"/favourites/"+created.ID, bytes.NewReader(patch))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	s.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("PATCH /favourites: unexpected status %d body=%s", rr.Code, rr.Body.String())
	}

	// --- DELETE ---
	req = httptest.NewRequest(http.MethodDelete, "/users/"+user+"/favourites/"+created.ID, nil)
	rr = httptest.NewRecorder()
	s.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("DELETE /favourites: unexpected status %d", rr.Code)
	}
}

// TestFavourites_ListPagination_EmptyDefaults verifies that default limit/offset are applied
// and that an empty list returns total=0 with a proper shape.
func TestFavourites_ListPagination_EmptyDefaults(t *testing.T) {
    s := newTestServer()

    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/users/kostas/favourites", nil)
    s.handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("GET default status=%d body=%s", rr.Code, rr.Body.String())
    }

    var resp struct {
        Favourites []models.Favourite `json:"favourites"`
        Total      int                `json:"total"`
        Limit      int                `json:"limit"`
        Offset     int                `json:"offset"`
    }
    if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
        t.Fatalf("unmarshal: %v", err)
    }
    if resp.Total != 0 || len(resp.Favourites) != 0 {
        t.Fatalf("expected empty list, got total=%d len=%d", resp.Total, len(resp.Favourites))
    }
    // defaultLimit currently 100 in server; assert that default is >0 and equals 100
    if resp.Limit != 100 || resp.Offset != 0 {
        t.Fatalf("defaults mismatch: limit=%d offset=%d", resp.Limit, resp.Offset)
    }
}

// TestFavourites_ListPagination_WithData creates multiple items and asserts limit/offset slicing.
func TestFavourites_ListPagination_WithData(t *testing.T) {
    s := newTestServer()
    user := "kostas"

    // create 5 favourites
    for i := 0; i < 5; i++ {
        body := []byte(`{"asset":{"type":"insight","text":"x","description":"d"}}`)
        req := httptest.NewRequest(http.MethodPost, "/users/"+user+"/favourites", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        rr := httptest.NewRecorder()
        s.handler.ServeHTTP(rr, req)
        if rr.Code != http.StatusCreated {
            t.Fatalf("POST status=%d body=%s", rr.Code, rr.Body.String())
        }
    }

    // fetch with limit=2 offset=2
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/users/"+user+"/favourites?limit=2&offset=2", nil)
    s.handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Fatalf("GET status=%d body=%s", rr.Code, rr.Body.String())
    }

    var resp struct {
        Favourites []models.Favourite `json:"favourites"`
        Total      int                `json:"total"`
        Limit      int                `json:"limit"`
        Offset     int                `json:"offset"`
    }
    if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
        t.Fatalf("unmarshal: %v", err)
    }
    if resp.Total != 5 {
        t.Fatalf("total mismatch: %d", resp.Total)
    }
    if resp.Limit != 2 || resp.Offset != 2 {
        t.Fatalf("limit/offset mismatch: %d/%d", resp.Limit, resp.Offset)
    }
    if len(resp.Favourites) != 2 {
        t.Fatalf("expected 2 items, got %d", len(resp.Favourites))
    }
}

