package server

import (
	"github.com/loukaspe/platform-go-challenge/internal/core/services"
	"github.com/loukaspe/platform-go-challenge/internal/handlers"
	userFavourite "github.com/loukaspe/platform-go-challenge/internal/handlers/userFavourites"
	"github.com/loukaspe/platform-go-challenge/internal/repositories"
	"github.com/loukaspe/platform-go-challenge/pkg/auth"
	"github.com/loukaspe/platform-go-challenge/pkg/helpers"
	"net/http"
	"os"
)

func (s *Server) initializeRoutes() {
	// health check
	healthCheckHandler := handlers.NewHealthCheckHandler(s.DB)
	s.router.HandleFunc("/health-check", healthCheckHandler.HealthCheckController).Methods("GET")

	// auth
	jwtMechanism := auth.NewAuthMechanism(
		os.Getenv("JWT_SECRET_KEY"),
		os.Getenv("JWT_SIGNING_METHOD"),
	)
	jwtService := services.NewJwtService(jwtMechanism)
	jwtMiddleware := handlers.NewAuthenticationMw(jwtMechanism)
	jwtHandler := handlers.NewJwtClaimsHandler(jwtService, s.logger)

	s.router.HandleFunc("/token", jwtHandler.JwtTokenController).Methods(http.MethodPost)

	protected := s.router.PathPrefix("/").Subrouter()
	protected.Use(jwtMiddleware.AuthenticationMW)

	// user favourites
	userRepository := repositories.NewUserRepository(s.DB)
	userFavouriteService := services.NewUserFavouriteService(s.logger, userRepository)
	assetTypeValidator := helpers.NewAssetTypeValidator()

	createUserFavouriteHandler := userFavourite.NewAddUserFavouriteHandler(userFavouriteService, assetTypeValidator, s.logger)
	getUserFavouriteHandler := userFavourite.NewGetUserFavouriteHandler(userFavouriteService, s.logger)
	updateUserFavouriteHandler := userFavourite.NewUpdateUserFavouriteHandler(userFavouriteService, assetTypeValidator, s.logger)
	deleteUserFavouriteHandler := userFavourite.NewDeleteUserFavouriteHandler(userFavouriteService, assetTypeValidator, s.logger)

	protected.HandleFunc("/users/{user_id:[0-9]+}/favourites", createUserFavouriteHandler.AddUserFavouriteAssetController).Methods("POST")
	protected.HandleFunc("/users/{user_id:[0-9]+}/favourites", getUserFavouriteHandler.GetUserFavouriteController).Methods("GET")
	protected.HandleFunc("/users/{user_id:[0-9]+}/favourites/{asset_id:[0-9]+}", updateUserFavouriteHandler.UpdateUserFavouriteController).Methods("PATCH")
	protected.HandleFunc("/users/{user_id:[0-9]+}/favourites", deleteUserFavouriteHandler.DeleteUserFavouriteController).Methods("DELETE")
}
