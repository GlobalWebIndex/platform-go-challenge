package http

import (
	"net/http"

	"golang.org/x/time/rate"
)

// LimitRequestSize returns a middleware that limits the request body size to maxBytes.
// This prevents DoS attacks from large request bodies.
func LimitRequestSize(next http.Handler, maxBytes int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		}
		next.ServeHTTP(w, r)
	})
}

// RateLimitMiddleware returns a middleware that applies a global rate limit to all requests.
// Uses a token bucket algorithm: rps requests per second with burst capacity.
func RateLimitMiddleware(next http.Handler, rps, burst int) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			WriteError(w, http.StatusTooManyRequests, ErrCodeRateLimit, "rate limit exceeded, please try again later")
			return
		}
		next.ServeHTTP(w, r)
	})
}
