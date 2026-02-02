package auth

import "net/http"

// RequireUserMatch ensures token sub == userId in URL path (prevents IDOR).
func RequireUserMatch(userIDFromPath func(*http.Request) (string, bool)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sub, ok := SubjectFromContext(r.Context())
			if !ok || sub == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			userID, ok := userIDFromPath(r)
			if !ok || userID == "" {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			if sub != userID {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
