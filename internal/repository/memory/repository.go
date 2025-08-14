package memory

import (
	"sync"
	"time"

	"gwi-favorites-service/internal/domain"
	"gwi-favorites-service/internal/repository"
)

// Repository implements FavoritesRepository using in-memory storage
type Repository struct {
	mu        sync.RWMutex
	assets    map[string]domain.Asset
	users     map[string]*domain.User
	favorites map[string]map[string]*domain.UserFavorite // userID -> assetID -> UserFavorite
}

// NewRepository creates a new in-memory repository
func NewRepository() *Repository {
	return &Repository{
		assets:    make(map[string]domain.Asset),
		users:     make(map[string]*domain.User),
		favorites: make(map[string]map[string]*domain.UserFavorite),
	}
}

// Asset operations
func (r *Repository) CreateAsset(asset domain.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.assets[asset.GetID()]; exists {
		return domain.ErrAssetAlreadyExists
	}

	r.assets[asset.GetID()] = asset
	return nil
}

func (r *Repository) GetAsset(assetID string) (domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	asset, exists := r.assets[assetID]
	if !exists {
		return nil, domain.ErrAssetNotFound
	}

	return asset, nil
}

func (r *Repository) UpdateAsset(asset domain.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.assets[asset.GetID()]; !exists {
		return domain.ErrAssetNotFound
	}

	asset.SetUpdatedAt(time.Now())
	r.assets[asset.GetID()] = asset

	// Update in all user favorites
	for userID := range r.favorites {
		if favorite, exists := r.favorites[userID][asset.GetID()]; exists {
			favorite.Asset = asset
			favorite.UpdatedAt = time.Now()
		}
	}

	return nil
}

func (r *Repository) DeleteAsset(assetID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.assets[assetID]; !exists {
		return domain.ErrAssetNotFound
	}

	delete(r.assets, assetID)

	// Remove from all user favorites
	for userID := range r.favorites {
		delete(r.favorites[userID], assetID)
	}

	return nil
}

func (r *Repository) ListAssets(limit, offset int) ([]domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var assets []domain.Asset
	count := 0

	for _, asset := range r.assets {
		if count < offset {
			count++
			continue
		}
		if len(assets) >= limit {
			break
		}
		assets = append(assets, asset)
		count++
	}

	return assets, nil
}

// User operations
func (r *Repository) CreateUser(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user
	if r.favorites[user.ID] == nil {
		r.favorites[user.ID] = make(map[string]*domain.UserFavorite)
	}
	return nil
}

func (r *Repository) GetUser(userID string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[userID]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

// Favorites operations
func (r *Repository) AddFavorite(userID string, asset domain.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Ensure user exists
	if _, exists := r.users[userID]; !exists {
		return domain.ErrUserNotFound
	}

	// Ensure asset exists
	if _, exists := r.assets[asset.GetID()]; !exists {
		return domain.ErrAssetNotFound
	}

	// Initialize user favorites if needed
	if r.favorites[userID] == nil {
		r.favorites[userID] = make(map[string]*domain.UserFavorite)
	}

	// Check if already a favorite
	if _, exists := r.favorites[userID][asset.GetID()]; exists {
		return domain.ErrFavoriteAlreadyExists
	}

	// Add to favorites
	favorite := domain.NewUserFavorite(userID, asset)
	r.favorites[userID][asset.GetID()] = favorite

	return nil
}

func (r *Repository) RemoveFavorite(userID, assetID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user exists
	if _, exists := r.users[userID]; !exists {
		return domain.ErrUserNotFound
	}

	// Check if favorite exists
	if _, exists := r.favorites[userID][assetID]; !exists {
		return domain.ErrFavoriteNotFound
	}

	delete(r.favorites[userID], assetID)
	return nil
}

func (r *Repository) GetUserFavorites(userID string, limit, offset int) ([]*domain.UserFavorite, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if user exists
	if _, exists := r.users[userID]; !exists {
		return nil, domain.ErrUserNotFound
	}

	userFavorites := r.favorites[userID]
	if userFavorites == nil {
		return []*domain.UserFavorite{}, nil
	}

	var favorites []*domain.UserFavorite
	count := 0

	for _, favorite := range userFavorites {
		if count < offset {
			count++
			continue
		}
		if len(favorites) >= limit {
			break
		}
		favorites = append(favorites, favorite)
		count++
	}

	return favorites, nil
}

func (r *Repository) IsFavorite(userID, assetID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if userFavorites := r.favorites[userID]; userFavorites != nil {
		_, exists := userFavorites[assetID]
		return exists, nil
	}

	return false, nil
}

func (r *Repository) GetFavoriteCount(userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if userFavorites := r.favorites[userID]; userFavorites != nil {
		return len(userFavorites), nil
	}

	return 0, nil
}

func (r *Repository) UpdateFavoriteAsset(userID, assetID string, asset domain.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user exists
	if _, exists := r.users[userID]; !exists {
		return domain.ErrUserNotFound
	}

	// Check if favorite exists
	favorite, exists := r.favorites[userID][assetID]
	if !exists {
		return domain.ErrFavoriteNotFound
	}

	// Update the asset in the favorite
	favorite.Asset = asset
	favorite.UpdatedAt = time.Now()

	return nil
}

// Ensure Repository implements the interface
var _ repository.FavoritesRepository = (*Repository)(nil)
