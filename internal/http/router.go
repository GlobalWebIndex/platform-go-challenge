package http

import (
	"net/http"

	"github.com/SpBalis/platform-go-challenge/internal/service"
	"github.com/go-chi/chi/v5"
)

func NewRouter(favs *service.FavouritesService, users *service.UsersService) http.Handler {
	r := chi.NewRouter()
	h := &Handler{Favs: favs, Users: users}

	r.Route("/v1", func(r chi.Router) {

		//Handle users
		r.Post("/users", h.CreateUser)
		r.Get("/users/{id}", h.GetUser)
		r.Get("/users", h.ListUsers)

		//Handle favourites
		r.Get("/users/{id}/favourites", h.List)
		r.Post("/users/{id}/favourites", h.Add)
		r.Delete("/users/{id}/favourites/{assetId}", h.Remove)
		r.Patch("/users/{id}/favourites/{assetId}", h.EditDesc)
		r.Delete("/users/{id}/favourites", h.ClearAll)
	})

	return r
}
