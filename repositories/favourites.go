package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"platform-go-challenge/models"
)

func GetUserFavourites(
	ctx context.Context,
	db *sql.DB,
	userID string,
) ([]models.FavouriteAsset, error) {

	query := `
	SELECT
		a.id,
		a.type,
		a.title,
		a.data,
		f.description
	FROM favourites f
	JOIN assets a ON a.id = f.asset_id
	WHERE f.user_id = $1
	ORDER BY f.created_at DESC;
	`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.FavouriteAsset

	for rows.Next() {
		var (
			id          string
			assetType   models.AssetType
			title       *string
			rawData     json.RawMessage
			description *string
		)

		if err := rows.Scan(
			&id,
			&assetType,
			&title,
			&rawData,
			&description,
		); err != nil {
			return nil, err
		}

		data, err := unmarshalAssetData(assetType, rawData)
		if err != nil {
			return nil, err
		}

		result = append(result, models.FavouriteAsset{
			AssetID:     id,
			Type:        assetType,
			Title:       title,
			Data:        data,
			Description: description,
		})
	}

	return result, nil
}

func AddFavourite(
	ctx context.Context,
	db *sql.DB,
	userID, assetID string,
	description *string,
) error {

	query := `
	INSERT INTO favourites (user_id, asset_id, description)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id, asset_id)
	DO UPDATE SET description = EXCLUDED.description;
	`

	_, err := db.ExecContext(
		ctx,
		query,
		userID,
		assetID,
		description,
	)

	return err
}

func RemoveFavourite(
	ctx context.Context,
	db *sql.DB,
	userID, assetID string,
) error {

	query := `
	DELETE FROM favourites
	WHERE user_id = $1 AND asset_id = $2;
	`

	res, err := db.ExecContext(ctx, query, userID, assetID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("favourite not found")
	}

	return nil
}

func UpdateFavouriteDescription(
	ctx context.Context,
	db *sql.DB,
	userID, assetID string,
	description *string,
) error {

	query := `
	UPDATE favourites
	SET description = $3
	WHERE user_id = $1 AND asset_id = $2;
	`

	res, err := db.ExecContext(
		ctx,
		query,
		userID,
		assetID,
		description,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("favourite not found")
	}

	return nil
}
