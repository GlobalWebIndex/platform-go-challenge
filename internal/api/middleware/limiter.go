package middleware

import (
	"encoding/json"
	"net/http"
	"platform-go-challenge/internal/api/types"
	"strings"
	"sync"
	"time"
)

type visitor struct {
	count    int
	lastSeen time.Time
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	limit    int           // requests per window
	window   time.Duration // time window
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
}

// RateLimitMiddleware creates a middleware that limits the number of requests from each IP address.
// It takes a requestsPerMinute parameter to specify the maximum number of allowed requests per IP
// within a one minute window. When the limit is exceeded, it returns a 429 Too Many Requests status
// with a JSON error response.
func RateLimitMiddleware(requestsPerMinute int) MiddlewareFunc {
	limiter := NewRateLimiter(requestsPerMinute, time.Minute)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/swagger/") {
				next.ServeHTTP(w, r)
				return
			}

			ip := r.RemoteAddr

			exceeded := limiter.isRateLimitExceeded(ip)
			if exceeded {
				response := types.NewErrorResponse(
					types.ErrCodeIPRateLimit,
					"Rate limit exceeded",
				)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(response)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (r *RateLimiter) isRateLimitExceeded(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	v, exists := r.visitors[ip]

	if !exists {
		// First visit
		r.visitors[ip] = &visitor{
			count:    1,
			lastSeen: now,
		}
		return false
	}

	// Reset if window has passed
	if now.Sub(v.lastSeen) > r.window {
		// Create new visitor instead of modifying existing one
		r.visitors[ip] = &visitor{
			count:    1,
			lastSeen: now,
		}
		return false
	}

	// Increment count and check limit
	v.count++
	v.lastSeen = now

	return v.count > r.limit
}
