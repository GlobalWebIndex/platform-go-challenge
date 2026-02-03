package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"

	"go-favorites-app/internal/core/domain/favorites"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository implements ports.FavoriteRepository using PostgreSQL.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new postgres repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// entity represents the database structure for an asset.
type entity struct {
	ID        string          `db:"id"`
	Type      string          `db:"type"`
	AssetData json.RawMessage `db:"asset_data"`
}

// Save persists a generic Asset.
func (r *Repository) Save(ctx context.Context, asset favorites.Asset) error {
	data, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset: %w", err)
	}

	query := `
		INSERT INTO favorites (id, type, asset_data, user_id)
		VALUES ($1, $2, $3, $4)
	`
	_, err = r.db.Exec(ctx, query, asset.GetID(), string(asset.GetType()), data, asset.GetUserID())
	if err != nil {
		return fmt.Errorf("failed to insert asset: %w", err)
	}
	return nil
}

// FindByID retrieves an asset by its ID.
func (r *Repository) FindByID(ctx context.Context, id string) (favorites.Asset, error) {
	query := `SELECT type, asset_data FROM favorites WHERE id = $1`

	var typeStr string
	var data []byte

	err := r.db.QueryRow(ctx, query, id).Scan(&typeStr, &data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("asset not found")
		}
		return nil, fmt.Errorf("failed to fetch asset: %w", err)
	}

	return unmarshalAsset(typeStr, data)
}

// FindAll returns an iterator of Assets to stream results.
func (r *Repository) FindAll(ctx context.Context, limit, offset int) (iter.Seq2[favorites.Asset, error], error) {
	query := `
		SELECT type, asset_data 
		FROM favorites 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}

	// Return an iterator compatible with Go 1.23 ranges
	return func(yield func(favorites.Asset, error) bool) {
		defer rows.Close()

		for rows.Next() {
			var typeStr string
			var data []byte

			if err := rows.Scan(&typeStr, &data); err != nil {
				yield(nil, fmt.Errorf("failed to scan row: %w", err))
				return
			}

			asset, err := unmarshalAsset(typeStr, data)
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(asset, nil) {
				return
			}
		}

		if err := rows.Err(); err != nil {
			yield(nil, fmt.Errorf("rows iteration error: %w", err))
		}
	}, nil
}

// Delete removes an asset by ID.
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM favorites WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("asset not found")
	}
	return nil
}

// UpdateDescription updates just the description of an asset.
func (r *Repository) UpdateDescription(ctx context.Context, id, description string) (favorites.Asset, error) {
	query := `
		UPDATE favorites
		SET asset_data = jsonb_set(asset_data, '{description}', to_jsonb($1::text)),
		    updated_at = NOW()
		WHERE id = $2
		RETURNING type, asset_data
	`
	var typeStr string
	var data []byte

	err := r.db.QueryRow(ctx, query, description, id).Scan(&typeStr, &data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("asset not found")
		}
		return nil, fmt.Errorf("failed to update description: %w", err)
	}

	return unmarshalAsset(typeStr, data)
}

// FindByUser returns an iterator of Assets for a specific user.
func (r *Repository) FindByUser(ctx context.Context, userID string, limit, offset int) (iter.Seq2[favorites.Asset, error], error) {
	query := `
		SELECT id, type, asset_data 
		FROM favorites 
		WHERE user_id = $1
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query favorites: %w", err)
	}

	return func(yield func(favorites.Asset, error) bool) {
		defer rows.Close()
		for rows.Next() {
			var ent entity
			if err := rows.Scan(&ent.ID, &ent.Type, &ent.AssetData); err != nil {
				yield(nil, fmt.Errorf("scan error: %w", err))
				return
			}
			asset, err := unmarshalAsset(ent.Type, ent.AssetData)
			if err != nil {
				yield(nil, fmt.Errorf("unmarshal error: %w", err))
				return
			}
			if !yield(asset, nil) {
				return
			}
		}
		if err := rows.Err(); err != nil {
			yield(nil, err)
		}
	}, nil
}

// unmarshalAsset is a helper to deserialize JSON into the correct concrete type.
func unmarshalAsset(t string, data []byte) (favorites.Asset, error) {
	var asset favorites.Asset

	switch favorites.AssetType(t) {
	case favorites.AssetTypeChart:
		var c favorites.Chart
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		asset = c
	case favorites.AssetTypeInsight:
		var i favorites.Insight
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}
		asset = i
	case favorites.AssetTypeAudience:
		var a favorites.Audience
		if err := json.Unmarshal(data, &a); err != nil {
			return nil, err
		}
		asset = a
	default:
		return nil, fmt.Errorf("unknown asset type: %s", t)
	}
	return asset, nil
}
