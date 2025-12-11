package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLoggingMiddleware_LogsRequestAndCallsNext(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	lg := NewWithWriter(&buf)

	var nextCalled bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusTeapot) // Use a non-200 status
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/users/user-1/favourites", nil)
	rr := httptest.NewRecorder()

	mw := RequestLoggingMiddleware(lg)
	mw(next).ServeHTTP(rr, req)

	if !nextCalled {
		t.Fatalf("expected next handler to be called")
	}
	if rr.Code != http.StatusTeapot {
		t.Fatalf("expected status 418 from next handler, got %d", rr.Code)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "http_request") {
		t.Fatalf("expected log output to contain 'http_request', got %q", logOutput)
	}
	// Loose checks: just ensure status code and path are present somewhere.
	if !strings.Contains(logOutput, "418") {
		t.Fatalf("expected log output to contain status code 418, got %q", logOutput)
	}
	if !strings.Contains(logOutput, "/v1/users/user-1/favourites") {
		t.Fatalf("expected log output to contain path, got %q", logOutput)
	}
}
