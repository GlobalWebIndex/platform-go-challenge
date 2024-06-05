package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func resetRateLimiter() {
	limiter = rate.NewLimiter(1, 3) // Reset limiter before each test
}

func TestRateLimit(t *testing.T) {
	tests := []struct {
		name           string
		requests       int
		expectedStatus int
		delay          time.Duration
	}{
		{"Below Rate Limit", 1, http.StatusOK, 0},
		{"Above Rate Limit", 4, http.StatusTooManyRequests, 0},
		{"Within Rate Limit with Delay", 4, http.StatusOK, time.Second / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetRateLimiter() // Reset limiter before each test
			router := mux.NewRouter()
			router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			router.Use(RateLimit)

			var lastStatus int
			for i := 0; i < tt.requests; i++ {
				rr := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatal(err)
				}
				router.ServeHTTP(rr, req)
				lastStatus = rr.Code

				if tt.delay > 0 {
					time.Sleep(tt.delay)
				}
			}

			if lastStatus != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", lastStatus, tt.expectedStatus)
			}
		})
	}
}
