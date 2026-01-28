package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"platform-go-challenge/models"
)

func CreateAsset(
	ctx context.Context,
	db *sql.DB,
	asset models.Asset,
	description *string,
) (string, error) {
	query := `
	INSERT INTO assets (type, title, description, data)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	var assetID string
	err := db.QueryRowContext(
		ctx,
		query,
		asset.Type,
		asset.Title,
		description,
		asset.Data,
	).Scan(&assetID)
	if err != nil {
		return "", err
	}

	return assetID, nil
}

func GetAssetByID(
	ctx context.Context,
	db *sql.DB,
	assetID string,
) (*models.Asset, error) {
	query := `
	SELECT id, type, title, data, created_at
	FROM assets
	WHERE id = $1;
	`

	var asset models.Asset
	err := db.QueryRowContext(ctx, query, assetID).Scan(
		&asset.ID,
		&asset.Type,
		&asset.Title,
		&asset.Data,
		&asset.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &asset, nil
}

func ListAssets(
	ctx context.Context,
	db *sql.DB,
) ([]models.Asset, error) {
	query := `
	SELECT id, type, title, data, created_at
	FROM assets
	ORDER BY created_at DESC;
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var a models.Asset
		if err := rows.Scan(&a.ID, &a.Type, &a.Title, &a.Data, &a.CreatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}

func unmarshalAssetData(
	assetType models.AssetType,
	raw json.RawMessage,
) (any, error) {

	switch assetType {
	case models.AssetChart:
		var c struct {
			XAxisTitle string    `json:"x_axis_title"`
			YAxisTitle string    `json:"y_axis_title"`
			DataPoints []float64 `json:"data_points"`
		}
		return c, json.Unmarshal(raw, &c)

	case models.AssetInsight:
		var i struct {
			Text string `json:"text"`
		}
		return i, json.Unmarshal(raw, &i)

	case models.AssetAudience:
		var a struct {
			Gender                string   `json:"gender"`
			BirthCountry          string   `json:"birth_country"`
			AgeGroups             []string `json:"age_groups"`
			HoursOnSocialMediaMin int      `json:"hours_on_social_media_min"`
			PurchasesLastMonthMin int      `json:"purchases_last_month_min"`
		}
		return a, json.Unmarshal(raw, &a)

	default:
		return nil, fmt.Errorf("unknown asset type: %s", assetType)
	}
}
