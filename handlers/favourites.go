package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"platform-go-challenge/repositories"
)

// ensureUserExists returns an error if the user does not exist.
func ensureUserExists(ctx context.Context, db *sql.DB, userID string) error {
	user, err := repositories.GetUserByID(ctx, db, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return sql.ErrNoRows
	}
	return nil
}

func GetUserFavourites(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		userID := parts[2]

		if err := ensureUserExists(r.Context(), db.(*sql.DB), userID); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to verify user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		favs, err := repositories.GetUserFavourites(
			r.Context(),
			db.(*sql.DB),
			userID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(favs)
	}
}

func AddFavourite(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		userID := parts[2]

		if err := ensureUserExists(r.Context(), db.(*sql.DB), userID); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to verify user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var input struct {
			AssetID     string  `json:"asset_id"`
			Description *string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		err := repositories.AddFavourite(
			r.Context(),
			db.(*sql.DB),
			userID,
			input.AssetID,
			input.Description,
		)
		if err != nil {
			println("Error adding favourite:", err.Error())
			http.Error(w, "Failed to add favourite: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"user_id":  userID,
			"asset_id": input.AssetID,
			"message":  "Favourite added successfully",
		})
	}
}

func UpdateFavourite(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		userID := parts[2]
		assetID := parts[4]

		if err := ensureUserExists(r.Context(), db.(*sql.DB), userID); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to verify user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var input struct {
			Description *string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		err := repositories.UpdateFavouriteDescription(
			r.Context(),
			db.(*sql.DB),
			userID,
			assetID,
			input.Description,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func RemoveFavourite(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		userID := parts[2]
		assetID := parts[4]

		if err := ensureUserExists(r.Context(), db.(*sql.DB), userID); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to verify user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err := repositories.RemoveFavourite(
			r.Context(),
			db.(*sql.DB),
			userID,
			assetID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

type DB interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
