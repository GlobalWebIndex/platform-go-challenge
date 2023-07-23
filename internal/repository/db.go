package repository

import (
	"fmt"
	"gwi_api/internal/domain"
	"gwi_api/internal/dto"

	_ "github.com/go-sql-driver/mysql"
)

type MemoryDB struct {
	nextUserID  uint64
	nextAssetID uint64
	users       map[uint64]domain.User
	userEmails  map[string]uint64
	assets      map[uint64]domain.Asset
}

// user related methods

func (d *MemoryDB) IncreaseUserID() {
	d.nextUserID += 1
}

func (d *MemoryDB) RegisterUser(user dto.UserDto) (*uint64, error) {
	_, found := DB.userEmails[user.Email]
	if found {
		return nil, fmt.Errorf("already exist user")
	}
	id := d.nextUserID
	DB.userEmails[user.Email] = id
	DB.users[id] = domain.User{
		ID:        id,
		Password:  user.Password,
		Email:     user.Email,
		Favorites: make(map[uint64]bool),
	}
	d.IncreaseUserID()
	return &id, nil
}

func (d *MemoryDB) DeleteUser(userId uint64) error {
	user, found := DB.users[userId]
	if !found {
		return fmt.Errorf("[Err] user don't exist")
	}
	delete(DB.userEmails, user.Email)
	delete(DB.users, userId)
	return nil
}

// asset related method
func (d *MemoryDB) IncreaseAssetID() {
	d.nextAssetID += 1
}

func (d *MemoryDB) AddFavorites(userId uint64, assetID uint64) error {
	user, found := d.users[userId]
	if !found {
		return fmt.Errorf("[Err] did not find user %d", userId)
	}
	if _, found = user.Favorites[assetID]; found {
		return fmt.Errorf("[Err] already added this asset %d", assetID)
	}
	d.users[userId].Favorites[assetID] = true
	return nil
}

func (d *MemoryDB) GetAllFavorites(userId uint64, pagination domain.PaginationParams) ([]domain.Asset, error) {
	user, found := d.users[userId]

	if !found {
		return nil, fmt.Errorf("[Err] did not find user %d", userId)
	}

	// Define the start and end index based on the page and pageSize
	startIndex := (pagination.PageNumber - 1) * pagination.PageSize
	if startIndex < 0 {
		startIndex = 0
	}
	endIndex := pagination.PageNumber * pagination.PageSize

	var favoriteAssets []domain.Asset
	index := 0
	for assetID := range user.Favorites {
		if index >= startIndex && index < endIndex {
			asset, found := d.assets[assetID]
			if found {
				favoriteAssets = append(favoriteAssets, asset)
			} else {
				fmt.Printf("[Warn] could not find asset %d", assetID)
			}
		}
		index++
	}

	return favoriteAssets, nil
}

// asset
// AddAsset implements AssetService.
func (d *MemoryDB) AddAsset(asset dto.AssetDto) (*uint64, error) {
	id := d.nextAssetID
	d.assets[d.nextAssetID] = domain.Asset{
		ID:          id,
		CreatedBy:   asset.CreatedBy,
		Type:        asset.Type,
		Description: asset.Description,
		Data:        asset.Data,
	}
	d.IncreaseAssetID()
	return &id, nil
}

// DeleteAsset implements AssetService.
func (d *MemoryDB) DeleteAsset(assetId uint64) error {
	_, found := d.assets[assetId]
	if !found {
		return fmt.Errorf("[Err] doesn't exist asset %d", assetId)
	}
	delete(d.assets, assetId)
	return nil
}

// GetAsset implements AssetService.
// GetAsset retrieves all assets with pagination.
func (d *MemoryDB) GetAssets(params domain.PaginationParams) ([]domain.Asset, error) {
	// Define the start and end index based on the page and pageSize
	startIndex := (params.PageNumber - 1) * params.PageSize
	if startIndex < 0 {
		startIndex = 0
	}
	endIndex := params.PageNumber * params.PageSize

	// To preserve the order of assets as they were inserted, we would need to keep a separate slice of asset IDs
	// or a slice of assets themselves. For this example, we'll iterate through the map, but keep in mind it does not guarantee order.
	// The code might be a bit different if you have that slice.

	var assets []domain.Asset
	index := 0
	for _, asset := range d.assets {
		if index >= startIndex && index < endIndex {
			assets = append(assets, asset)
		}
		index++
	}

	if len(assets) == 0 {
		return nil, fmt.Errorf("[Err] No assets found")
	}

	return assets, nil
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		nextUserID:  0,
		nextAssetID: 0,
		users:       make(map[uint64]domain.User),
		userEmails:  make(map[string]uint64),
		assets:      make(map[uint64]domain.Asset),
	}
}

type DBHandler interface {
	NewUserQuery() UserQuery
	NewFavoritesQuery() FavoritesQuery
}

type dbHandler struct {
	db *MemoryDB
}

// NewFavoritesQuery implements DBHandler.
func (*dbHandler) NewFavoritesQuery() FavoritesQuery {
	return &favoritesQuery{}
}

var DB *MemoryDB

func NewDBHandler(db *MemoryDB) DBHandler {
	DB = db
	return &dbHandler{db}
}

func NewDB() (*MemoryDB, error) {
	DB := NewMemoryDB()
	return DB, nil
}

func (d *dbHandler) NewUserQuery() UserQuery {
	return &userQuery{}
}
