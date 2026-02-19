package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"

	"github.com/gwi/favourites-api/internal/cache"
	"github.com/gwi/favourites-api/internal/service"
	"github.com/gwi/favourites-api/internal/store"
)

// executeRequest executes an HTTP request against the test application
func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

// checkResponseCode asserts the response code that matches expected
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}

// newContractTestApplication creates a test app with real in-memory stores + cache.
// This supports stable API contract tests
func newContractTestApplication(t *testing.T) *application {
	logger := zaptest.NewLogger(t).Sugar()

	storage := store.NewStorage()        // real in-memory stores
	memCache := cache.NewInMemoryCache() // real cache
	ttl := 5 * time.Minute

	favouritesService := service.NewFavouritesService(storage, memCache, ttl)

	return &application{
		config: config{
			addr: ":8080",
			env:  "test",
		},
		storage:           storage,
		cache:             memCache,
		favouritesService: favouritesService,
		logger:            logger,
	}
}

func decodeJSON(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var m map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &m); err != nil {
		t.Fatalf("failed to decode json: %v\nbody=%s", err, rr.Body.String())
	}
	return m
}

// For endpoints that return { "data": { ... } }
func envelopeObject(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	m := decodeJSON(t, rr)

	data, ok := m["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected {data:{...}} envelope; got: %s", rr.Body.String())
	}

	return data
}

// For paginated endpoints that return:
// { "data": [...], "meta": {...} }
func envelopePage(t *testing.T, rr *httptest.ResponseRecorder) ([]any, map[string]any) {
	t.Helper()

	m := decodeJSON(t, rr)

	items, ok := m["data"].([]any)
	if !ok {
		t.Fatalf("expected data array; got: %s", rr.Body.String())
	}

	meta, ok := m["meta"].(map[string]any)
	if !ok {
		t.Fatalf("expected meta object; got: %s", rr.Body.String())
	}

	return items, meta
}
