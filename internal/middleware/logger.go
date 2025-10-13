package middleware

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Metrics tracks application metrics
type Metrics struct {
	totalRequests    int64
	requestsByStatus map[int]int64
	mutex            sync.RWMutex
	startTime        time.Time
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		requestsByStatus: make(map[int]int64),
		startTime:        time.Now(),
	}
}

// RecordRequest increments the total request count
func (m *Metrics) RecordRequest() {
	atomic.AddInt64(&m.totalRequests, 1)
}

// RecordStatusCode records the status code of a request
func (m *Metrics) RecordStatusCode(statusCode int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.requestsByStatus[statusCode]++
}

// GetHealthStatus returns a health status object
func (m *Metrics) GetHealthStatus() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	statusCodes := make(map[int]int64)
	for k, v := range m.requestsByStatus {
		statusCodes[k] = v
	}

	return map[string]interface{}{
		"status":           "healthy",
		"total_requests":   atomic.LoadInt64(&m.totalRequests),
		"uptime_seconds":   int64(time.Since(m.startTime).Seconds()),
		"status_codes":     statusCodes,
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestLogger is middleware that logs HTTP requests
func RequestLogger(metricsCollector *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Record the request
			metricsCollector.RecordRequest()

			// Create a custom response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Call the next handler
			next.ServeHTTP(wrapped, r)

			// Record the status code
			metricsCollector.RecordStatusCode(wrapped.statusCode)

			// Log the request
			duration := time.Since(start)
			fmt.Printf("%s %s %s - Status: %d - Duration: %v\n",
				r.Method,
				r.RequestURI,
				r.RemoteAddr,
				wrapped.statusCode,
				duration,
			)
			os.Stdout.Sync()
		})
	}
}