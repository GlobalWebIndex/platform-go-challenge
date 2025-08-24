// Package router provides routing functionality
package router

import (
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var MongoClient *mongo.Client
var AppLogger *zap.Logger

func NewRouter(client *mongo.Client,logger *zap.Logger) http.Handler {

	// mongo client and logger to be used from controllers
	MongoClient = client
	AppLogger = logger

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	corsHandler := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept-Language"}),
	)(router)

	securedHandler := addSecurityHeaders(corsHandler)
	return securedHandler
}

// Security Headers middleware
func addSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'; object-src 'none'; script-src 'self'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		next.ServeHTTP(w, r)
	})
}

var routes = Routes{
	Route{
		"GetUserFavourites",
		strings.ToUpper("Get"),
		"/users/{userId}/favourites",
		getUserFavourites,
	},
	Route{
		"AddFavouriteAsset",
		strings.ToUpper("Post"),
		"/users/{userId}/favourites/{assetId}",
		addFavouriteAsset,
	},
	Route{
		"DelFavouriteAsset",
		strings.ToUpper("Delete"),
		"/users/{userId}/favourites/{assetId}",
		delFavouriteAsset,
	},
	Route{
		"PatchFavouriteAsset",
		strings.ToUpper("PATCH"),
		"/users/{userId}/favourites/{assetId}",
		patchFavouriteAsset,
	},
}






