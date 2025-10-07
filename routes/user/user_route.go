package user_route

import (
	"gwi/constants"
	"gwi/middleware"
	user_service "gwi/services/user"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes sets up the routes for user management with JWT authentication.
func RegisterRoutes(router *mux.Router) {
	router.Handle("/"+constants.ENDPOINT_USERS, middleware.JWTAuth(http.HandlerFunc(user_service.GetUsers))).Methods("GET")
	router.Handle("/"+constants.ENDPOINT_USERS, middleware.JWTAuth(http.HandlerFunc(user_service.AddUser))).Methods("POST")
	router.Handle("/"+constants.ENDPOINT_USERS+"/{id}", middleware.JWTAuth(http.HandlerFunc(user_service.DeleteUser))).Methods("DELETE")
}
