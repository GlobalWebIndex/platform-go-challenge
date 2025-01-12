package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestContextMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.Handler
		timeOut        time.Duration
		withCancel     bool
		expectedStatus int
	}{
		{
			name: "context time out",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(1 * time.Millisecond)
				select {
				case <-r.Context().Done():
					w.WriteHeader(http.StatusGatewayTimeout)
				default:
					w.WriteHeader(http.StatusOK)
				}

			}),
			withCancel:     false,
			timeOut:        1 * time.Nanosecond,
			expectedStatus: http.StatusGatewayTimeout,
		},
		{
			name: "cancel context",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(1 * time.Millisecond)
				select {
				case <-r.Context().Done():
					w.WriteHeader(http.StatusGatewayTimeout)
				default:
					w.WriteHeader(http.StatusOK)
				}
			}),
			withCancel:     true,
			timeOut:        50 * time.Second,
			expectedStatus: http.StatusGatewayTimeout,
		},
		{
			name: "normal execution",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(1 * time.Second)
				select {
				case <-r.Context().Done():
					w.WriteHeader(http.StatusGatewayTimeout)
				default:
					w.WriteHeader(http.StatusOK)
				}
			}),
			withCancel:     false,
			timeOut:        3 * time.Second,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			rr := httptest.NewRecorder()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.withCancel {
				cancel()
			}

			contextTimeOut = tt.timeOut

			handler := ContextMiddleware(ctx)(tt.handler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
