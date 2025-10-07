package asset_route

import (
	"gwi/constants"
	"gwi/middleware"
	asset_service "gwi/services/asset"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes sets up the routes for asset management with JWT authentication.
func RegisterRoutes(router *mux.Router) {
	router.Handle("/"+constants.ENDPOINT_ASSETS, middleware.JWTAuth(http.HandlerFunc(asset_service.GetAssets))).Methods("GET")
	router.Handle("/"+constants.ENDPOINT_ASSETS, middleware.JWTAuth(http.HandlerFunc(asset_service.AddAsset))).Methods("POST")
	router.Handle("/"+constants.ENDPOINT_ASSETS+"/{id}", middleware.JWTAuth(http.HandlerFunc(asset_service.EditAsset))).Methods("PUT")
}
