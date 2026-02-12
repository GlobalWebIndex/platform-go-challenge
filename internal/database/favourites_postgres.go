package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/giannis84/platform-go-challenge/internal/models"
	"github.com/lib/pq"
)

var (
	ErrNotFound      = errors.New("favourite not found")
	ErrAlreadyExists = errors.New("favourite already exists")
)

// DB is the package-level database connection.
var DB *sql.DB

func GetUserFavouritesFromDB(userID string) ([]*models.FavouriteAsset, error) {
	const query = `
		SELECT id, user_id, asset_type, description, data, created_at, updated_at
		FROM favourites
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("querying user favourites: %w", err)
	}
	defer rows.Close()

	var favourites []*models.FavouriteAsset
	for rows.Next() {
		fav, err := scanFavourite(rows)
		if err != nil {
			return nil, err
		}
		favourites = append(favourites, fav)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating user favourites: %w", err)
	}

	if favourites == nil {
		favourites = []*models.FavouriteAsset{}
	}
	return favourites, nil
}

func GetFavouriteFromDB(userID, assetID string) (*models.FavouriteAsset, error) {
	const query = `
		SELECT id, user_id, asset_type, description, data, created_at, updated_at
		FROM favourites
		WHERE user_id = $1 AND id = $2`

	row := DB.QueryRow(query, userID, assetID)

	var fav models.FavouriteAsset
	var rawData []byte

	err := row.Scan(
		&fav.ID, &fav.UserID, &fav.AssetType,
		&fav.Description, &rawData,
		&fav.CreatedAt, &fav.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scanning favourite: %w", err)
	}

	asset, err := unmarshalAssetData(fav.AssetType, rawData)
	if err != nil {
		return nil, err
	}
	fav.Data = asset

	return &fav, nil
}

func AddFavouriteInDB(ctx context.Context, favourite *models.FavouriteAsset) error {
	dataJSON, err := json.Marshal(favourite.Data)
	if err != nil {
		return fmt.Errorf("marshalling asset data: %w", err)
	}

	const query = `
		INSERT INTO favourites (id, user_id, asset_type, description, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = DB.ExecContext(ctx, query,
		favourite.ID, favourite.UserID, string(favourite.AssetType),
		favourite.Description, dataJSON,
		favourite.CreatedAt, favourite.UpdatedAt,
	)
	if err != nil {
		// Check for unique-violation (PG error code 23505)
		if isUniqueViolation(err) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("inserting favourite: %w", err)
	}
	return nil
}

func UpdateFavouriteInDB(favourite *models.FavouriteAsset) error {
	dataJSON, err := json.Marshal(favourite.Data)
	if err != nil {
		return fmt.Errorf("marshalling asset data: %w", err)
	}

	const query = `
		UPDATE favourites
		SET description = $1, data = $2, updated_at = $3
		WHERE user_id = $4 AND id = $5`

	result, err := DB.Exec(query,
		favourite.Description, dataJSON, favourite.UpdatedAt,
		favourite.UserID, favourite.ID,
	)
	if err != nil {
		return fmt.Errorf("updating favourite: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func DeleteFavouriteFromDB(userID, assetID string) error {
	const query = `DELETE FROM favourites WHERE user_id = $1 AND id = $2`

	result, err := DB.Exec(query, userID, assetID)
	if err != nil {
		return fmt.Errorf("deleting favourite: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// scanFavourite scans a single row from the favourites table into a FavouriteAsset.
func scanFavourite(rows *sql.Rows) (*models.FavouriteAsset, error) {
	var fav models.FavouriteAsset
	var rawData []byte

	err := rows.Scan(
		&fav.ID, &fav.UserID, &fav.AssetType,
		&fav.Description, &rawData,
		&fav.CreatedAt, &fav.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning favourite row: %w", err)
	}

	asset, err := unmarshalAssetData(fav.AssetType, rawData)
	if err != nil {
		return nil, err
	}
	fav.Data = asset

	return &fav, nil
}

// unmarshalAssetData deserialises JSONB data into the correct Asset implementation
// based on the asset_type column.
func unmarshalAssetData(assetType models.AssetType, data []byte) (models.Asset, error) {
	if data == nil {
		return nil, nil
	}

	switch assetType {
	case models.AssetTypeChart:
		var chart models.Chart
		if err := json.Unmarshal(data, &chart); err != nil {
			return nil, fmt.Errorf("unmarshalling chart data: %w", err)
		}
		return &chart, nil
	case models.AssetTypeInsight:
		var insight models.Insight
		if err := json.Unmarshal(data, &insight); err != nil {
			return nil, fmt.Errorf("unmarshalling insight data: %w", err)
		}
		return &insight, nil
	case models.AssetTypeAudience:
		var audience models.Audience
		if err := json.Unmarshal(data, &audience); err != nil {
			return nil, fmt.Errorf("unmarshalling audience data: %w", err)
		}
		return &audience, nil
	default:
		return nil, fmt.Errorf("unknown asset type: %s", assetType)
	}
}

// isUniqueViolation checks if a PostgreSQL error is a unique constraint violation (23505).
func isUniqueViolation(err error) bool {
	var pge *pq.Error
	if errors.As(err, &pge) {
		return pge.Code == "23505"
	}
	return false
}
