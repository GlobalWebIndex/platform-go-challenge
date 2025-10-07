package repo

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/SpBalis/platform-go-challenge/internal/service"
)

type FavouritesRepo struct{ DB *sql.DB }

func NewFavouritesRepo(db *sql.DB) *FavouritesRepo { return &FavouritesRepo{DB: db} }

func (r *FavouritesRepo) ListByUser(ctx context.Context, userID int64, limit, offset int) ([]service.Favourite, int, error) {
	const qc = `SELECT COUNT(*) FROM favourites WHERE user_id=$1`
	var total int
	if err := r.DB.QueryRowContext(ctx, qc, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `
      SELECT f.user_id, f.custom_description,
             a.id, a.type, a.description, a.data
      FROM favourites f
      JOIN assets a ON a.id = f.asset_id
      WHERE f.user_id = $1
      ORDER BY f.created_at DESC
      LIMIT $2 OFFSET $3`
	rows, err := r.DB.QueryContext(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]service.Favourite, 0)
	for rows.Next() {
		var fav service.Favourite
		var t string
		var raw []byte
		if err := rows.Scan(&fav.UserID, &fav.CustomDescription, &fav.Asset.ID, &t, &fav.Asset.Description, &raw); err != nil {
			return nil, 0, err
		}
		fav.Asset.Type = service.AssetType(t)
		var v any
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, 0, err
		}
		fav.Asset.Data = v
		out = append(out, fav)
	}
	return out, total, rows.Err()
}

func (r *FavouritesRepo) Add(ctx context.Context, userID int64, asset service.Asset, customDesc *string) (service.Favourite, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return service.Favourite{}, err
	}
	defer tx.Rollback()

	b, err := json.Marshal(asset.Data)
	if err != nil {
		return service.Favourite{}, err
	}
	const qa = `INSERT INTO assets (type, description, data) VALUES ($1,$2,$3) RETURNING id`
	if err := tx.QueryRowContext(ctx, qa, string(asset.Type), asset.Description, b).Scan(&asset.ID); err != nil {
		return service.Favourite{}, err
	}

	const qf = `INSERT INTO favourites (user_id, asset_id, custom_description) VALUES ($1,$2,$3)`
	if _, err := tx.ExecContext(ctx, qf, userID, asset.ID, customDesc); err != nil {
		return service.Favourite{}, err
	}

	if err := tx.Commit(); err != nil {
		return service.Favourite{}, err
	}
	desc := ""
	if customDesc != nil {
		desc = *customDesc
	}
	return service.Favourite{UserID: userID, Asset: asset, CustomDescription: desc}, nil
}

func (r *FavouritesRepo) Remove(ctx context.Context, userID, assetID int64) error {
	const q = `DELETE FROM favourites WHERE user_id=$1 AND asset_id=$2`
	_, err := r.DB.ExecContext(ctx, q, userID, assetID)
	return err
}

func (r *FavouritesRepo) UpdateDescription(ctx context.Context, userID, assetID int64, desc string) error {
	const q = `UPDATE favourites SET custom_description=$1 WHERE user_id=$2 AND asset_id=$3`
	_, err := r.DB.ExecContext(ctx, q, desc, userID, assetID)
	return err
}

func (r *FavouritesRepo) ClearAll(ctx context.Context, userID int64) error {
	const q = `DELETE FROM favourites WHERE user_id=$1`
	_, err := r.DB.ExecContext(ctx, q, userID)
	return err
}
