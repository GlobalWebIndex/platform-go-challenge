package webserver

import (
	"challenge/middleware"
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterRoutes(srv *WebServer) { // router *mux.Router
	if srv.Configuration.Metrics {
		srv.Server.Path("/metrics").Handler(promhttp.Handler())
	}
	// Register the endpoints
	srv.Server.HandleFunc("/signup", CreateUserHandler).Methods("POST")
	srv.Server.HandleFunc("/signin", SigninExistingUser).Methods("POST")
	srv.Server.HandleFunc("/signout", SignoutExistingUser).Methods("GET")
	srv.Server.HandleFunc("/version", GetBuildVersion).Methods("GET")
	// Create a subrouter for the `/users` paths to authorize user access to their lists
	usersRouter := srv.Server.PathPrefix("/users").Subrouter()
	usersRouter.Use(middleware.VerifyToken)
	usersRouter.HandleFunc("/{user_id}/create", CreateUserHandler).Methods("POST")
	usersRouter.HandleFunc("/{user_id}/delete", DeleteUserHandler).Methods("DELETE")
	usersRouter.HandleFunc("/{user_id}/favourites", GetFavouritesHandler).Methods("GET")
	usersRouter.HandleFunc("/{user_id}/favourites", AddFavouriteHandler).Methods("POST")
	usersRouter.HandleFunc("/{user_id}/favourites/{asset_id}", EditAssetDescriptionHandler).Methods("PUT")
	usersRouter.HandleFunc("/{user_id}/favourites/{asset_id}", RemoveFavouriteHandler).Methods("DELETE")
}

func AttachProfiler(srv *WebServer) {
	srv.Server.NewRoute().PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)
}
