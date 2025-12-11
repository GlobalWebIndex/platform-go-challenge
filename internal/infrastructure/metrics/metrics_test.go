package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestRequestMetricsMiddleware_RecordsMetrics(t *testing.T) {
	t.Parallel()

	m := New()

	var nextCalled bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusCreated) // 201
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/users/user-1/favourites", nil)
	rr := httptest.NewRecorder()

	mw := RequestMetricsMiddleware(m)
	mw(next).ServeHTTP(rr, req)

	if !nextCalled {
		t.Fatalf("expected next handler to be called")
	}
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201 from next handler, got %d", rr.Code)
	}

	metricFamilies, err := m.reg.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}

	found := false
	for _, mf := range metricFamilies {
		if mf.GetName() == "http_requests_total" {
			found = true
			if len(mf.GetMetric()) == 0 {
				t.Fatalf("expected http_requests_total to have at least one sample")
			}
			break
		}
	}
	if !found {
		t.Fatalf("expected to find http_requests_total in registry")
	}
}

func TestMetricsHandler_ExposesMetrics(t *testing.T) {
	t.Parallel()

	m := New()

	// Prime the registry with a sample so /metrics is not empty.
	m.observeRequest("GET", "/health", http.StatusOK, 10*time.Millisecond)

	handler := m.Handler()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	ct := rr.Header().Get("Content-Type")
	if ct == "" {
		t.Fatalf("expected Content-Type to be set, got empty")
	}
	if !strings.Contains(ct, "text/plain") && !strings.Contains(ct, "openmetrics") {
		t.Fatalf("expected Content-Type to indicate text or openmetrics, got %q", ct)
	}

	body := rr.Body.String()
	if body == "" {
		t.Fatalf("expected non-empty response body")
	}
}

func TestMetrics_RegistersCollectors(t *testing.T) {
	t.Parallel()

	m := New()

	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_counter_extra",
		Help: "extra counter for testing registry behaviour",
	})

	if err := m.reg.Register(counter); err != nil {
		t.Fatalf("expected to register additional metric successfully, got %v", err)
	}
}
