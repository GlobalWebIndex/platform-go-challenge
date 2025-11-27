// Package app contains the application layer of the application. It is the layer that contains the business logic of the application. It is the layer that interacts with the domain layer and the interface adapters. It is the layer that contains the use cases of the application
package app

import (
	"github.com/akazantzidis/gwi-ass/internal/app/auth/command"
	"github.com/akazantzidis/gwi-ass/internal/app/favourite/commands"
	queries2 "github.com/akazantzidis/gwi-ass/internal/app/user/queries"
	"github.com/akazantzidis/gwi-ass/internal/domain/token"
	"github.com/akazantzidis/gwi-ass/internal/domain/user"

	"github.com/akazantzidis/gwi-ass/internal/app/favourite/queries"
	"github.com/akazantzidis/gwi-ass/internal/app/notification"
	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/akazantzidis/gwi-ass/internal/pkg/time"
	"github.com/akazantzidis/gwi-ass/internal/pkg/uuid"
)

// Queries Contains all available query handlers of this app
type Queries struct {
	GetAllFavoritesHandler queries.GetAllFavoritesRequestHandler
	GetFavoriteHandler     queries.GetFavoriteRequestHandler

	GetUserHandler queries2.GetUserHandler
}

// Commands Contains all available command handlers of this app
type Commands struct {
	CreateFavoriteHandler        commands.CreateFavoriteRequestHandler
	UpdateFavoriteHandler        commands.UpdateFavoriteRequestHandler
	UpdatePartialFavoriteHandler commands.UpdatePartialFavoriteRequestHandler

	DeleteFavoriteHandler commands.DeleteFavoriteRequestHandler

	LoginUserHandler        command.LoginHandler
	RefreshTokenUserHandler command.RefreshHandler
	LogoutUserHandler       command.LogoutHandler
}

// FavoriteServices Contains the grouped queries and command of the app layer
type FavoriteServices struct {
	Queries  Queries
	Commands Commands
}

type AuthServices struct {
	Queries  Queries
	Commands Commands
}

type UserServices struct {
	Queries  Queries
	Commands Commands
}

// Services contains all exposed services of the application layer
type Services struct {
	FavoriteServices FavoriteServices
	AuthServices     AuthServices
	UserServices     UserServices
}

// NewServices Bootstraps Application Layer dependencies
func NewServices(favoriteRepo favourite.Repository, ns notification.Service, userRepo user.Repository, refreshTokenRepo token.RefreshRepository, _ uuid.Provider, _ time.Provider) Services {
	return Services{
		FavoriteServices: FavoriteServices{
			Queries: Queries{
				GetAllFavoritesHandler: queries.NewGetAllFavoritesRequestHandler(favoriteRepo),
				GetFavoriteHandler:     queries.NewGetFavoriteRequestHandler(favoriteRepo),
			},
			Commands: Commands{
				CreateFavoriteHandler:        commands.NewAddFavoriteRequestHandler(favoriteRepo, ns),
				UpdateFavoriteHandler:        commands.NewUpdateFavoriteRequestHandler(favoriteRepo),
				UpdatePartialFavoriteHandler: commands.NewUpdatePartialFavoriteRequestHandler(favoriteRepo),

				DeleteFavoriteHandler: commands.NewDeleteFavoriteRequestHandler(favoriteRepo),
			},
		},
		AuthServices: AuthServices{
			Queries: Queries{},
			Commands: Commands{
				LoginUserHandler:        command.NewLoginHandler(userRepo, refreshTokenRepo),
				LogoutUserHandler:       command.NewLogoutHandler(refreshTokenRepo),
				RefreshTokenUserHandler: command.NewRefreshHandler(refreshTokenRepo),
			},
		},
		UserServices: UserServices{
			Queries: Queries{
				GetUserHandler: queries2.NewGetUserHandler(userRepo),
			},
			Commands: Commands{},
		},
	}
}
