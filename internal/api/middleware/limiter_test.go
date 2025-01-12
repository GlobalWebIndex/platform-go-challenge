package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLimiterMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		requests int
		failAt   int // 0 means should never fail
	}{
		{
			name:     "rate exceeded",
			limit:    50,
			requests: 51,
			failAt:   51,
		},
		{
			name:     "not exceeded",
			limit:    50,
			requests: 49,
			failAt:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RateLimitMiddleware(tt.limit)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			for i := 1; i <= tt.requests; i++ {
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))

				if tt.failAt == 0 {
					if rr.Code != http.StatusOK {
						t.Errorf("request %d: got status %d, want %d", i, rr.Code, http.StatusOK)
					}
				} else {
					wantCode := http.StatusOK
					if i >= tt.failAt {
						wantCode = http.StatusTooManyRequests
					}
					if rr.Code != wantCode {
						t.Errorf("request %d: got status %d, want %d", i, rr.Code, wantCode)
					}
				}
			}
		})
	}
}
