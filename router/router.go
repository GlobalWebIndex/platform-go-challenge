package router

import (
	"user-favorites-api/handlers"
	"user-favorites-api/middleware"
	"user-favorites-api/store"

	"github.com/gorilla/mux"
)

func NewRouter(store *store.Store) *mux.Router {
    h := &handlers.Handler{Store: store}
    r := mux.NewRouter()

    r.HandleFunc("/register", h.RegisterUser).Methods("POST")
    r.HandleFunc("/login", h.Login).Methods("POST")

    api := r.PathPrefix("/api").Subrouter()
    api.Use(middleware.JWTAuth)
    api.HandleFunc("/favorites/{userID}", h.GetFavorites).Methods("GET")
    api.HandleFunc("/favorites/{userID}", h.AddFavorite).Methods("POST")
    api.HandleFunc("/favorites/{userID}/{assetID}", h.RemoveFavorite).Methods("DELETE")
    api.HandleFunc("/favorites/{userID}/{assetID}", h.EditFavorite).Methods("PUT")
  

    return r
}
