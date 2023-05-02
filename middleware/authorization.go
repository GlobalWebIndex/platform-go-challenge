package middleware

import (
	"challenge/models"
	"challenge/storage"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

// JwtSecret is the provided secret for users authorization
var JwtSecret []byte

// VerifyToken handles the user authorization for the protected resources
func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Get the JWT string from the cookie
		tknStr := c.Value
		// Initialize a new instance of `Claims`
		claims := &models.Claims{}
		// Parse the JWT string and store the result in `claims`.
		// Return an error if the token is invalid or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Extract the user ID from the token claims
		userID := claims.ID
		_, adminCheck, err := storage.Db.CheckAdmin(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		path := r.URL.RawPath
		if strings.HasSuffix(path, "/create") || strings.HasSuffix(path, "/delete") {
			if !adminCheck {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		// Extract the user ID from the URL path
		userIDPathParam := mux.Vars(r)["user_id"]
		if userIDPathParam == "" || userIDPathParam == " " {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Check if user is an administrator
		if !adminCheck {
			// Compare the user ID in the token with the user ID in the URL path
			if userID != userIDPathParam {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
