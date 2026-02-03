package observability

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Latency of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	dbConnectionPoolStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_connection_pool_stats",
			Help: "Database connection pool statistics (total, idle, active)",
		},
		[]string{"state"},
	)

	cacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)
	cacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)
)

func init() {
	// Register metrics
	prometheus.MustRegister(httpRequestLatency)
	prometheus.MustRegister(dbConnectionPoolStats)
	prometheus.MustRegister(cacheHits)
	prometheus.MustRegister(cacheMisses)
}

// Middleware records HTTP request latency.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriterSpy{ResponseWriter: w, code: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()
		httpRequestLatency.WithLabelValues(r.Method, r.Pattern, fmt.Sprint(ww.code)).Observe(duration)
	})
}

type responseWriterSpy struct {
	http.ResponseWriter
	code int
}

func (w *responseWriterSpy) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// StartDBStatsCollector starts a background goroutine to collect DB stats.
func StartDBStatsCollector(dbPool *pgxpool.Pool) {
	go func() {
		for {
			stats := dbPool.Stat()
			dbConnectionPoolStats.WithLabelValues("total").Set(float64(stats.TotalConns()))
			dbConnectionPoolStats.WithLabelValues("idle").Set(float64(stats.IdleConns()))
			dbConnectionPoolStats.WithLabelValues("acquired").Set(float64(stats.AcquiredConns()))
			time.Sleep(5 * time.Second)
		}
	}()
}
