package handlers

import (
	"net/http"
	"strings"
)

func UserRouter(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Expected paths:
		// POST /users - Create a new user
		// GET /users - List all users
		// GET /users/{userID} - Get user by ID

		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		if parts[0] != "users" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodPost:
			// POST /users - Create a new user
			if len(parts) == 1 {
				println("Create a new user")
				AddUser(db)(w, r)
				return
			}

		case http.MethodGet:
			// GET /users - List all users
			if len(parts) == 1 {
				println("List all users")
				GetUsers(db)(w, r)
				return
			}

			// case http.MethodDelete:
			// 	// DELETE /users/{userID} - Delete user
			// 	if len(parts) == 2 {
			// 		println("Delete user")
			// 		DeleteUser(db)(w, r)
			// 		return
			// 	}
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func FavouritesRouter(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Expected paths:
		// /users/{userID}/favourites
		// /users/{userID}/favourites/{assetID}

		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		// Minimum: users/{id}/favourites
		if len(parts) < 3 || parts[0] != "users" || parts[2] != "favourites" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {

		case http.MethodGet:
			// GET /users/{id}/favourites
			if len(parts) == 3 {
				println("Get lists of favorite assets for a user")
				GetUserFavourites(db)(w, r)
				return
			}

		case http.MethodPost:
			// POST /users/{id}/favourites
			if len(parts) == 3 {
				println("Add new favorite asset to a user")
				AddFavourite(db)(w, r)
				return
			}

		case http.MethodPatch:
			// PATCH /users/{id}/favourites/{assetID}
			if len(parts) == 4 {
				println("Update favorite assets for a user")
				UpdateFavourite(db)(w, r)
				return
			}

		case http.MethodDelete:
			// DELETE /users/{id}/favourites/{assetID}
			if len(parts) == 4 {
				println("Remove asset for a user")
				RemoveFavourite(db)(w, r)
				return
			}
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
