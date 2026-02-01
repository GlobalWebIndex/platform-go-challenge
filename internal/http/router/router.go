package router

import (
	"net/http"
)

type Handler interface {
	ListFavorites(w http.ResponseWriter, r *http.Request)
	AddFavorite(w http.ResponseWriter, r *http.Request)
	RemoveFavorite(w http.ResponseWriter, r *http.Request)
	UpdateDescription(w http.ResponseWriter, r *http.Request)
}

func New(h Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.ListFavorites(w, r)
		case http.MethodPost:
			h.AddFavorite(w, r)
		case http.MethodDelete:
			h.RemoveFavorite(w, r)
		case http.MethodPatch:
			h.UpdateDescription(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
