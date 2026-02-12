package logging

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
)

func TestFromContext(t *testing.T) {
	customBuf := &bytes.Buffer{}
	customLogger := slog.New(slog.NewTextHandler(customBuf, nil))

	tests := []struct {
		name        string
		ctx         context.Context
		wantDefault bool
		logMessage  string // if not wantDefault, verify this message appears in customBuf
	}{
		{
			name:        "with logger in context",
			ctx:         NewContextWithLogger(context.Background(), customLogger),
			wantDefault: false,
			logMessage:  "custom logger test",
		},
		{
			name:        "without logger in context",
			ctx:         context.Background(),
			wantDefault: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear buffer for each test
			customBuf.Reset()

			got := FromContext(tt.ctx)

			if got == nil {
				t.Fatal("expected non-nil logger")
			}

			if tt.wantDefault {
				// Should return default logger - verify it's functional
				got.Info("fallback test")
			} else {
				// Should return custom logger - verify by checking buffer
				got.Info(tt.logMessage)
				if !strings.Contains(customBuf.String(), tt.logMessage) {
					t.Errorf("expected custom logger to write %q to buffer, got: %s", tt.logMessage, customBuf.String())
				}
			}
		})
	}
}

func TestRequestLogger(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		logMessage     string
		checkRequestID bool
	}{
		{
			name:           "attaches logger to context with request_id",
			method:         "GET",
			path:           "/test",
			logMessage:     "handler executed",
			checkRequestID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := slog.New(slog.NewJSONHandler(buf, nil))

			var capturedLogger *slog.Logger
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedLogger = FromContext(r.Context())
				capturedLogger.Info(tt.logMessage)
				w.WriteHeader(http.StatusOK)
			})

			// Chain: RequestID -> RequestLogger -> handler
			handler := middleware.RequestID(RequestLogger(logger)(testHandler))

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			// Verify handler was called
			if rec.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", rec.Code)
			}

			// Verify logger was captured
			if capturedLogger == nil {
				t.Error("expected logger to be attached to context")
			}

			output := buf.String()

			// Verify log message
			if !strings.Contains(output, tt.logMessage) {
				t.Errorf("expected log output to contain %q, got: %s", tt.logMessage, output)
			}

			// Verify request_id if required
			if tt.checkRequestID {
				if !strings.Contains(output, `"request_id":`) {
					t.Errorf("expected log output to contain request_id field, got: %s", output)
				}
				if strings.Contains(output, `"request_id":""`) {
					t.Errorf("expected request_id to have a value, got empty: %s", output)
				}
			}
		})
	}
}


