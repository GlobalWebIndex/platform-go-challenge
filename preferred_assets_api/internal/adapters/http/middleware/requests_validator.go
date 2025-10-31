package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type contextKey string

const validatedBodyKey contextKey = "validatedBody"

// Validator
var validate = validator.New()

// BodyGetter Interface (non-generic method)
// BodyGetter defines a contract to retrieve validated DTOs from a request
type BodyGetter interface {
	GetValidatedBody(r *http.Request) (any, bool)
}

// DefaultBodyGetter uses request context to retrieve validated bodies
type DefaultBodyGetter struct{}

// GetValidatedBody retrieves the validated DTO from context
func (d DefaultBodyGetter) GetValidatedBody(r *http.Request) (any, bool) {
	val := r.Context().Value(validatedBodyKey)
	if val == nil {
		return nil, false
	}
	return val, true
}

// Global BodyGetter instance (can be replaced in tests)
var Body BodyGetter = DefaultBodyGetter{}

// Middleware
// ValidateBody returns a middleware that validates a JSON request body into type T
func ValidateBody[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body T

			// Decode JSON request
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}

			// Validate struct fields using tags
			if err := validate.Struct(body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Add the validated body to context
			ctx := context.WithValue(r.Context(), validatedBodyKey, body)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetValidatedBody retrieves the validated DTO using the current BodyGetter
func GetValidatedBody[T any](r *http.Request) (T, bool) {
	val, ok := Body.GetValidatedBody(r)
	if !ok {
		var zero T
		return zero, false
	}
	tVal, ok := val.(T)
	if !ok {
		var zero T
		return zero, false
	}
	return tVal, true
}
