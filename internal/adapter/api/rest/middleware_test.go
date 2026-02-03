package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Context().Value(requestIDKey)
		assert.NotNil(t, rid, "RequestID should be in context")
		assert.NotEmpty(t, rid.(string), "RequestID should not be empty")

		respRid := w.Header().Get("X-Request-ID")
		assert.Equal(t, rid.(string), respRid, "Header should match context")
	})

	handlerToTest := RequestID(nextHandler)

	t.Run("generates new id when missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handlerToTest.ServeHTTP(w, req)
	})

	t.Run("preserves existing id", func(t *testing.T) {
		existingID := "existing-id"
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Request-ID", existingID)
		w := httptest.NewRecorder()

		nextHandlerWithCheck := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Context().Value(requestIDKey).(string)
			assert.Equal(t, existingID, rid)
		})

		RequestID(nextHandlerWithCheck).ServeHTTP(w, req)
		assert.Equal(t, existingID, w.Header().Get("X-Request-ID"))
	})
}

func TestChain(t *testing.T) {
	var calls []string
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "mw1")
			next.ServeHTTP(w, r)
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "mw2")
			next.ServeHTTP(w, r)
		})
	}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "final")
	})

	chained := Chain(final, mw1, mw2)
	chained.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	assert.Equal(t, []string{"mw1", "mw2", "final"}, calls, "Middleware should be called in order")
}
