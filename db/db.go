package db

import (
	"errors"
	"sync"

	"platform-go-challenge/models"
)

type DB struct {
	mu        sync.RWMutex
	favorites map[string][]models.Asset
}

func NewDB() *DB {
	return &DB{
		favorites: make(map[string][]models.Asset),
	}
}

func (db *DB) GetFavorites(userId string) ([]models.Asset, error) {
	if userId == "" {
		return nil, errors.New("userId is cannot be empty")
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	assets, ok := db.favorites[userId]
	if !ok {
		return nil, errors.New("user not found")
	}
	return assets, nil
}

func (db *DB) AddFavorite(userId string, asset models.Asset) error {
	//Adding if user exists, otherwise it creates an asset with not userID (Orphaned asset?)
	if userId == "" {
		return errors.New("userId cannot be empty")
	}

	if asset.ID == "" {
		return errors.New("asset ID cannot be empty")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	db.favorites[userId] = append(db.favorites[userId], asset)
	return nil
}

func (db *DB) RemoveFavorite(userId string, assetId string) error {
	if userId == "" {
		return errors.New("userId cannot be empty")
	}

	if assetId == "" {
		return errors.New("asset ID cannot be empty")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	assets := db.favorites[userId]
	for i, asset := range assets {
		if asset.ID == assetId {
			db.favorites[userId] = append(assets[:i], assets[i+1:]...)
			return nil
		}
	}
	return errors.New("asset not found")
}

func (db *DB) UpdateFavorites(userId string, asset models.Asset) error {
	if userId == "" {
		return errors.New("userId cannot be empty")
	}
	if asset.ID == "" {
		return errors.New("asset ID cannot be empty")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	assets := db.favorites[userId]
	for i, a := range assets {
		if a.ID == asset.ID {
			db.favorites[userId][i] = asset
			return nil
		}
	}
	return errors.New("asset not found")
}
