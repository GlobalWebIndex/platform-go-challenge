package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds Prometheus collectors for HTTP observability.
type Metrics struct {
	reg        *prometheus.Registry
	reqTotal   *prometheus.CounterVec
	reqLatency *prometheus.HistogramVec
}

// New creates a Metrics instance with HTTP request metrics.
func New() *Metrics {
	reg := prometheus.NewRegistry()

	reqTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests served.",
		},
		[]string{"method", "endpoint", "status"},
	)

	reqLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)

	reg.MustRegister(reqTotal, reqLatency)

	return &Metrics{
		reg:        reg,
		reqTotal:   reqTotal,
		reqLatency: reqLatency,
	}
}

// Handler returns an HTTP handler that exposes metrics in Prometheus format.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})
}

// observeRequest records request metrics. Endpoint should be the route pattern
// when available, falling back to the raw URL path.
func (m *Metrics) observeRequest(method, endpoint string, status int, duration time.Duration) {
	labels := prometheus.Labels{
		"method":   method,
		"endpoint": endpoint,
		"status":   strconv.Itoa(status),
	}
	m.reqTotal.With(labels).Inc()
	m.reqLatency.With(labels).Observe(duration.Seconds())
}
