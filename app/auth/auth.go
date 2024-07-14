package auth

import (
	"gwi-platform/database"
	"gwi-platform/handlers"
	"net/http"
	"strconv"
)

func UserStatusAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		db := database.GetDB()

		status, err := handlers.GetUserStatus(ctx, db, userID)
		if err != nil {
			http.Error(w, "Couldnt find user", http.StatusInternalServerError)
			return
		}

		if status != "A" {
			http.Error(w, "Unauthorized: User is not active", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
