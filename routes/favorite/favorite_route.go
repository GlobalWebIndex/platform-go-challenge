package favorite_route

import (
	"gwi/constants"
	"gwi/middleware"
	favorite_service "gwi/services/favorite"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes sets up the routes for favorite management with JWT authentication.
func RegisterRoutes(router *mux.Router) {
	router.Handle("/"+constants.ENDPOINT_FAVORITES, middleware.JWTAuth(http.HandlerFunc(favorite_service.GetFavorites))).Methods("GET")
	router.Handle("/"+constants.ENDPOINT_FAVORITES, middleware.JWTAuth(http.HandlerFunc(favorite_service.AddFavorite))).Methods("POST")
	router.Handle("/"+constants.ENDPOINT_FAVORITES+"/{id}", middleware.JWTAuth(http.HandlerFunc(favorite_service.DeleteFavorite))).Methods("DELETE")
}
