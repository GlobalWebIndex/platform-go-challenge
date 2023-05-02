package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// responseData defines a struct for the request's status and response size
type responseData struct {
	status int
	size   int
}

// loggingResponse composes the original responseWriter
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write writes the response using the original responseWriter
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader writes the request status code using the original responseWriter.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging defines a middleware that logs a request along with additional
// fields.
func WithLogging(h http.Handler) http.Handler {
	loggingFn := func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		logRespWriter := loggingResponseWriter{
			ResponseWriter: rw,
			responseData:   responseData,
		}
		// Inject the response writer
		h.ServeHTTP(&logRespWriter, req)
		duration := time.Since(start)
		logrus.WithFields(logrus.Fields{
			"uri":      req.RequestURI,
			"method":   req.Method,
			"status":   responseData.status,
			"duration": duration,
			"size":     responseData.size,
		}).Info("HTTP request")
	}
	return http.HandlerFunc(loggingFn)
}
