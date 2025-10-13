package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter is a simple in-memory rate limiter per IP.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]int
	resetAt  map[string]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new limiter allowing `limit` requests per `window`.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]int),
		resetAt:  make(map[string]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow checks whether the IP can make a request.
func (rl *RateLimiter) Allow(ip string) (bool, time.Time, int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Reset if the window expired
	if t, exists := rl.resetAt[ip]; !exists || now.After(t) {
		rl.visitors[ip] = 0
		rl.resetAt[ip] = now.Add(rl.window)
	}

	rl.visitors[ip]++
	count := rl.visitors[ip]
	reset := rl.resetAt[ip]

	// Return false if over limit
	if count > rl.limit {
		return false, reset, count
	}

	return true, reset, count
}

// RateLimitMiddleware applies rate limiting.
// - /api/login: strict (5 req/min per IP)
// - /api/health: no limit
// - everything else: relaxed (50 req/min per IP)
func RateLimitMiddleware(loginLimiter, defaultLimiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var limiter *RateLimiter

			// Choose limiter based on path
			if strings.HasPrefix(r.URL.Path, "/api/login") {
				limiter = loginLimiter
			} else if strings.HasPrefix(r.URL.Path, "/api/health") {
				// no limit for health checks
				next.ServeHTTP(w, r)
				return
			} else {
				limiter = defaultLimiter
			}

			ip := r.RemoteAddr
			allowed, reset, _ := limiter.Allow(ip)

			if !allowed {
				retryAfter := int(time.Until(reset).Seconds())
				w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
