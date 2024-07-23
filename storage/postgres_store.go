package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/arhsxro/platform-go-challenge/config"
	"github.com/arhsxro/platform-go-challenge/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sqlx.DB
}

func isValidAssetType(paramType string) bool {
	for _, validType := range models.ValidAssetTypes {
		if models.AssetType(paramType) == validType {
			return true
		}
	}
	return false
}

func NewPostgresStore(cfg *config.Config) (*PostgresStore, error) {
	db, err := NewPostgresDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &PostgresStore{db: db}, nil
}

// Retrieves a user's favorite assets from the database
func (store *PostgresStore) GetUserFavorites(ctx context.Context, userID, filterType string, page, pageSize int) ([]models.Asset, error) {

	var assets []models.Asset
	var query string
	var err error
	offset := (page - 1) * pageSize
	if filterType != "" {
		if !isValidAssetType(filterType) {
			log.Println("Invalid asset type")
			return nil, errors.New("invalid asset type")
		}
		query = "SELECT asset_id, type, description, data FROM assets WHERE user_id = $1 and type = $2 LIMIT $3 OFFSET $4"
		err = store.db.SelectContext(ctx, &assets, query, userID, filterType, pageSize, offset)
	} else {
		query = "SELECT asset_id, type, description, data FROM assets WHERE user_id = $1 LIMIT $2 OFFSET $3"
		err = store.db.SelectContext(ctx, &assets, query, userID, pageSize, offset)
	}

	return assets, err
}

// Adds a new favorite asset for a user in the database
func (store *PostgresStore) AddFavorite(ctx context.Context, userID string, asset models.Asset) error {
	tx, err := store.db.Begin() //begin transaction
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO assets (user_id, asset_id, type, description, data)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (asset_id) DO NOTHING`

	_, err = tx.ExecContext(ctx, query, userID, asset.ID, asset.Type, asset.Description, string(asset.Data))
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Removes an asset from a user in the database
func (store *PostgresStore) RemoveFavorite(ctx context.Context, userID, assetID string) error {
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "DELETE FROM assets WHERE user_id = $1 AND asset_id = $2"
	_, err = store.db.ExecContext(ctx, query, userID, assetID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Updates an asset's description from a user in the database
func (store *PostgresStore) UpdateDescription(ctx context.Context, userID, assetID, newDescription string) error {
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := "UPDATE assets SET description = $1 WHERE user_id = $2 AND asset_id = $3"
	_, err = store.db.ExecContext(ctx, query, newDescription, userID, assetID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
