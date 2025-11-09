package middleware

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SecurityHeaders injects common HTTP headers that harden the API surface
// against basic attacks and content sniffing.
// It disables MIME sniffing, clickjacking, and referrer leaks.
// CORS is intentionally disabled by default.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		// To allow cross-origin requests, uncomment the following line:
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

// MaxBody limits the maximum size of a request body (in bytes).
// It prevents large or malicious payloads from exhausting memory.
func MaxBody(n int64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, n)
		next.ServeHTTP(w, r)
	})
}

// RequestID attaches a unique request identifier to every HTTP response.
// It helps correlate logs across distributed systems or concurrent requests.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strconv.FormatInt(time.Now().UnixNano(), 36) + "-" + strconv.Itoa(rand.Intn(999999))
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

// statusRecorder wraps ResponseWriter to record status code and written bytes
// for structured logging and observability.
type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	if sr.status == 0 {
		sr.status = http.StatusOK
	}
	n, err := sr.ResponseWriter.Write(b)
	sr.bytes += n
	return n, err
}

// Logger provides basic structured access logging with latency metrics.
// It can be disabled by setting ENABLE_HTTP_LOG=false in the environment.
// Example log line:
//   method=GET path=/users/kostas/favourites status=200 bytes=512 dur=3.1ms ua="curl/7.77" req_id=abc123
func Logger(next http.Handler) http.Handler {
	enabled := os.Getenv("ENABLE_HTTP_LOG")
	if strings.ToLower(enabled) == "false" {
		return next // logs disabled
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(sr, r)

		reqID := w.Header().Get("X-Request-ID")
		ua := r.UserAgent()
		duration := time.Since(start)

		log.Printf(
			`method=%s path=%s status=%d bytes=%d dur=%s ua="%s" req_id=%s`,
			r.Method, r.URL.Path, sr.status, sr.bytes, duration, ua, reqID,
		)
	})
}

// RateLimiter implements a simple per-IP or per-user token bucket
// with a minimum interval between requests. It is meant as a lightweight
// protection against abuse or accidental floods, not a full quota system.
type RateLimiter struct {
	mu   sync.Mutex
	last map[string]time.Time
	rate time.Duration // minimum duration between allowed requests
}

// NewRateLimiter constructs a new limiter enforcing one request every minInterval.
func NewRateLimiter(minInterval time.Duration) *RateLimiter {
	return &RateLimiter{last: make(map[string]time.Time), rate: minInterval}
}

// Middleware wraps the handler and enforces the rate limit policy.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr
		if u := parseUserFromPath(r.URL.Path); u != "" {
			key = "user:" + u
		}
		now := time.Now()

		rl.mu.Lock()
		if t, ok := rl.last[key]; ok && now.Sub(t) < rl.rate {
			rl.mu.Unlock()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		rl.last[key] = now
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// parseUserFromPath extracts the userID from URLs of the form /users/{userID}/...
func parseUserFromPath(p string) string {
	parts := strings.Split(strings.Trim(p, "/"), "/")
	if len(parts) >= 2 && parts[0] == "users" {
		return parts[1]
	}
	return ""
}

// APIKeyAuth enforces a simple shared-secret authentication via the X-API-Key header.
// If requiredKey is empty, the middleware is a no-op (auth disabled).
// This is intentionally lightweight for the challenge scope, and can be replaced by JWT or OAuth later.
func APIKeyAuth(requiredKey string, next http.Handler) http.Handler {
	if strings.TrimSpace(requiredKey) == "" {
		return next // auth off
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != requiredKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

