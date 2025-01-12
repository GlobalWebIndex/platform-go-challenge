package middleware

import (
	"net/http"
	"slices"
)

// MiddlewareFunc represents a standard middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// Chain chains multiple middleware functions together
// We reverse the middleware slice because the execution order needs to match
// the order in which middlewares are provided to the function.
func Chain(h http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	// Apply middlewares in reverse order
	slices.Reverse(middlewares)
	for _, m := range middlewares {
		h = m(h)
	}

	return h
}
