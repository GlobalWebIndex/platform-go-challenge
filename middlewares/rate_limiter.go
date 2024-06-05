package middlewares

import (
	"net/http"

	"golang.org/x/time/rate"
)

// 1 request per second with a burst of 3 requests
var limiter = rate.NewLimiter(1, 3)

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
