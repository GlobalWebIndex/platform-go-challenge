package chihttp

import (
	"net/http"
	"time"

	_ "github.com/SokratisChaimanas/platform-go-challenge/docs"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/adapters/http/chi/handlers"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// NewRouter builds the chi router and mounts all routes.
func NewRouter(
	userService *app.UserService,
	assetService *app.AssetService,
	favService *app.FavouritesService,
) http.Handler {
	router := chi.NewRouter()

	// Basic middlewares.
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	// Health.
	healthHandler := handlers.NewHealthHandler()
	router.Method(http.MethodGet, "/api/healthz", healthHandler)

	// Users.
	userHandler := handlers.NewUserHandler(userService)
	router.Route("/api/users", func(r chi.Router) {
		r.Get("/{user_id}", userHandler.Get)
		// Favourites
		favouritesHandler := handlers.NewFavouritesHandler(favService)
		r.Get("/{user_id}/favourites", favouritesHandler.ListByUser)
		r.Post("/{user_id}/favourites", favouritesHandler.Add)
		r.Delete("/{user_id}/favourites/{asset_id}", favouritesHandler.Remove)
	})

	// Assets.
	assetHandler := handlers.NewAssetHandler(assetService)
	router.Patch("/api/assets/{asset_id}/description", assetHandler.EditDescription)

	// 404 fallback.
	router.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		handlers.WriteJsonError(w, "not found", http.StatusNotFound)
	})

	router.Get("/docs/*", httpSwagger.WrapHandler)

	return router
}
