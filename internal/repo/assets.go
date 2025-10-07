package repo

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/SpBalis/platform-go-challenge/internal/service"
)

type AssetsRepo struct{ DB *sql.DB }

func NewAssetsRepo(db *sql.DB) *AssetsRepo { return &AssetsRepo{DB: db} }

func (r *AssetsRepo) Create(ctx context.Context, a service.Asset) (service.Asset, error) {
	b, err := json.Marshal(a.Data)
	if err != nil {
		return service.Asset{}, err
	}
	const q = `INSERT INTO assets (type, description, data) VALUES ($1,$2,$3) RETURNING id`
	if err := r.DB.QueryRowContext(ctx, q, string(a.Type), a.Description, b).Scan(&a.ID); err != nil {
		return service.Asset{}, err
	}
	return a, nil
}

func (r *AssetsRepo) GetByID(ctx context.Context, id int64) (service.Asset, error) {
	const q = `SELECT id, type, description, data FROM assets WHERE id=$1`
	var a service.Asset
	var t string
	var raw []byte
	if err := r.DB.QueryRowContext(ctx, q, id).Scan(&a.ID, &t, &a.Description, &raw); err != nil {
		return service.Asset{}, err
	}
	a.Type = service.AssetType(t)
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return service.Asset{}, err
	}
	a.Data = v
	return a, nil
}
