package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"platform-go-challenge/internal/domain"
	"platform-go-challenge/internal/repository/testdata"
	"platform-go-challenge/internal/service/types"
	"strings"

	"github.com/mattn/go-sqlite3"
)

type sqliteRepository struct {
	db *sql.DB
}

var dbPath = "app.db"

func New() (domain.Repository, error) {
	// Remove existing database if it already exists
	os.Remove(dbPath)

	// Create new database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Ensure the connection is available
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables
	if err := initSchema(db); err != nil {
		return nil, err
	}

	// Seed db
	err = testdata.SeedDb(db)
	if err != nil {
		return nil, err
	}

	return &sqliteRepository{db: db}, nil
}

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS assets (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			description TEXT,
			data JSON NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS favorites (
			user_id TEXT NOT NULL,
			asset_id TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			PRIMARY KEY (user_id, asset_id),
			FOREIGN KEY (asset_id) REFERENCES assets(id)
		);
		
		CREATE INDEX IF NOT EXISTS idx_favorites_user 
		ON favorites(user_id);
	`)
	return err
}

func (r *sqliteRepository) GetAssets(ctx context.Context, ids []any) ([]types.Asset, error) {
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}

	args := strings.Join(placeholders, ",")

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, type, description, data, created_at, updated_at 
		FROM assets WHERE id in (`+args+`)`, ids...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	defer rows.Close()

	var assets []types.Asset

	for rows.Next() {
		var data []byte
		var base types.BaseAsset

		rows.Scan(&base.ID, &base.Type, &base.Description, &data, &base.CreatedAt, &base.UpdatedAt)

		switch base.Type {
		case types.AssetTypeChart:
			var chart types.Chart
			if err := json.Unmarshal(data, &chart); err != nil {
				return nil, err
			}

			chart.BaseAsset = base

			assets = append(assets, &chart)

		case types.AssetTypeInsight:
			var insight types.Insight
			if err := json.Unmarshal(data, &insight); err != nil {
				return nil, err
			}

			insight.BaseAsset = base

			assets = append(assets, &insight)

		case types.AssetTypeAudience:
			var audience types.Audience
			if err := json.Unmarshal(data, &audience); err != nil {
				return nil, err
			}
			audience.BaseAsset = base

			assets = append(assets, &audience)
		}

	}

	return assets, nil
}

func (r *sqliteRepository) UpdateAsset(ctx context.Context, asset types.Asset) error {
	data, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE assets 
		SET description = ?, data = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		asset.GetDescription(), data, asset.GetID())
	return err
}

func (r *sqliteRepository) AddFavorite(ctx context.Context, fav types.Favorite) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO favorites (user_id, asset_id, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)`,
		fav.UserID, fav.AssetID)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return domain.ErrAlreadyExists

		}
	}

	return err
}

func (r *sqliteRepository) RemoveFavorite(ctx context.Context, userID, assetID string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM favorites 
		WHERE user_id = ? AND asset_id = ?`,
		userID, assetID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *sqliteRepository) ListFavorites(ctx context.Context, userID string, page, limit int) ([]types.Favorite, error) {
	offset := (page - 1) * limit
	rows, err := r.db.QueryContext(ctx, `
		SELECT asset_id, created_at 
		FROM favorites
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []types.Favorite
	for rows.Next() {
		var fav types.Favorite
		fav.UserID = userID
		if err := rows.Scan(&fav.AssetID, &fav.CreatedAt); err != nil {
			return nil, err
		}
		favorites = append(favorites, fav)
	}
	return favorites, rows.Err()
}
