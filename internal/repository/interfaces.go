package repository

import "gwi-favorites-service/internal/domain"

// FavoritesRepository defines the interface for favorites storage operations
type FavoritesRepository interface {
	// Asset operations
	CreateAsset(asset domain.Asset) error
	GetAsset(assetID string) (domain.Asset, error)
	UpdateAsset(asset domain.Asset) error
	DeleteAsset(assetID string) error
	ListAssets(limit, offset int) ([]domain.Asset, error)

	// User operations
	CreateUser(user *domain.User) error
	GetUser(userID string) (*domain.User, error)

	// Favorites operations
	AddFavorite(userID string, asset domain.Asset) error
	RemoveFavorite(userID, assetID string) error
	GetUserFavorites(userID string, limit, offset int) ([]*domain.UserFavorite, error)
	IsFavorite(userID, assetID string) (bool, error)
	GetFavoriteCount(userID string) (int, error)
	UpdateFavoriteAsset(userID, assetID string, asset domain.Asset) error
}
