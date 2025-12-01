package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"platform-go-challenge/internal/auth"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/storage"
)

func SetupRouter(cfg *config.Config, favorites *storage.FavoritesStore, users *storage.UserStore) *mux.Router {
	jwtGen := &auth.RealJWTGenerator{}
	favoritesHandler := NewAssetHandler(cfg, favorites)
	userHandler := NewUserHandler(cfg, users, jwtGen)
	userAdminHandler := NewUserAdminHandler(cfg, users, favorites)

	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/login", userHandler.login).Methods("POST")
	r.HandleFunc("/register", userHandler.register).Methods("POST")

	// Routes that require JWT authentication
	api := r.PathPrefix("/api").Subrouter()
	api.Use(func(next http.Handler) http.Handler { return auth.TokenMiddleware(cfg, next) }) // enforce JWT

	api.HandleFunc("/favorites/list", favoritesHandler.list).Methods("GET")
	api.HandleFunc("/favorites/add", favoritesHandler.add).Methods("POST")
	api.HandleFunc("/favorites/remove", favoritesHandler.remove).Methods("DELETE")
	api.HandleFunc("/favorites/update", favoritesHandler.update).Methods("PATCH")

	// Admin-only route using API key middleware
	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(func(next http.Handler) http.Handler { return auth.AdminAPIKeyMiddleware(cfg, next) })
	admin.HandleFunc("/users/list", userAdminHandler.listUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", userAdminHandler.deleteUser).Methods("DELETE")

	return r
}
