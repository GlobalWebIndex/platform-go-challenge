package ports

import (
	"net/http"
)

type UserHandler interface {
	// Create handles HTTP POST /users requests
	Create(w http.ResponseWriter, r *http.Request)

	// Get handles HTTP GET /users/{id} requests
	Get(w http.ResponseWriter, r *http.Request)

	// List handles HTTP GET /users requests
	List(w http.ResponseWriter, r *http.Request)

	// Update handles HTTP PUT /users/{id} requests
	Update(w http.ResponseWriter, r *http.Request)

	// Delete handles HTTP DELETE /users/{id} requests
	Delete(w http.ResponseWriter, r *http.Request)

	// Get handles HTTP GET /users/{id}/favourites requests
	GetFavourites(w http.ResponseWriter, r *http.Request)
}

type FavouriteHandler interface {
	// Create handles HTTP POST /favourites requests
	Create(w http.ResponseWriter, r *http.Request)

	// Delete handles HTTP DELETE /favourites/{id} requests
	Delete(w http.ResponseWriter, r *http.Request)
}

type AssetHandler interface {
	// Create handles HTTP POST /assets requests
	Create(w http.ResponseWriter, r *http.Request)

	// Delete handles HTTP DELETE /assets/{id} requests
	Delete(w http.ResponseWriter, r *http.Request)
}
