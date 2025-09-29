package routes

import (
	"favourite_assets/server/controllers"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	r chi.Router,
	userController *controllers.UserController,
	assetController *controllers.AssetController,
	favController *controllers.FavouriteController,
	authMiddleware func(next http.Handler) http.Handler,
) {
	r.Use(authMiddleware)

	// Users
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userController.CreateUserHandler)
		r.Get("/", userController.ListUsersHandler)
		r.Get("/by-id", userController.GetUserHandler)
		r.Put("/", userController.UpdateUserHandler)
		r.Delete("/", userController.DeleteUserHandler)
	})

	// Assets
	r.Route("/assets", func(r chi.Router) {
		r.Post("/", assetController.CreateAssetHandler)
		r.Get("/", assetController.ListAssetsHandler)
		r.Get("/by-id", assetController.GetAssetHandler)
		r.Put("/", assetController.UpdateAssetHandler)
		r.Delete("/", assetController.DeleteAssetHandler)
	})

	// Favourites
	r.Route("/favorites", func(r chi.Router) {
		r.Post("/", favController.AddFavouriteHandler)
		r.Delete("/", favController.RemoveFavouriteHandler)
		r.Get("/", favController.ListFavouritesHandler)
		r.Get("/by-id", favController.GetFavouriteHandler)
	})
}