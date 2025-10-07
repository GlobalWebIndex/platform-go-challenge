package routes

import (
	asset_route "gwi/routes/asset"
	favorite_route "gwi/routes/favorite"
	user_route "gwi/routes/user"

	"github.com/gorilla/mux"
)

// RegisterRoutes initializes and returns a new router with all the application routes registered.
// It sets up the main router and registers user, asset, and favorite routes.
func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	user_route.RegisterRoutes(router)
	asset_route.RegisterRoutes(router)
	favorite_route.RegisterRoutes(router)

	return router
}
