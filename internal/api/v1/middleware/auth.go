package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/auth"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/shared/httpx"
)

// ctxKeyUserID is a context key for storing the authenticated user ID.
type ctxKeyUserID struct{}

// UserIDFromContext retrieves the authenticated user ID from the context,
// if present. Handlers can use this if they need the ID without re-reading
// headers or path parameters.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(ctxKeyUserID{}).(string)
	return userID, ok
}

// NewAuthUserMiddleware creates middleware that authenticates requests using
// the provided Authenticator. Returns 401 if unauthenticated, 403 if the
// authenticated user doesn't match the {userId} path parameter, and stores
// the user ID in the request context.
// This keeps auth concerns at the edge and allows swapping the concrete
// Authenticator implementation (header-based, JWT, OAuth2, etc.).
func NewAuthUserMiddleware(authenticator auth.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := authenticator.Authenticate(r)
			if err != nil {
				if err == auth.ErrUnauthenticated {
					httpx.WriteProblem(
						w,
						http.StatusUnauthorized,
						"unauthorized",
						"missing authentication",
						"authentication is required for this endpoint",
						r.URL.Path,
					)
					return
				}

				// Any other authentication error is treated as internal.
				httpx.WriteProblem(
					w,
					http.StatusInternalServerError,
					"internal_error",
					"authentication error",
					"",
					r.URL.Path,
				)
				return
			}

			pathUserID := chi.URLParam(r, "userId")
			if pathUserID != "" && pathUserID != userID {
				httpx.WriteProblem(
					w,
					http.StatusForbidden,
					"forbidden",
					"user mismatch",
					"authenticated user does not match requested userId",
					r.URL.Path,
				)
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyUserID{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
