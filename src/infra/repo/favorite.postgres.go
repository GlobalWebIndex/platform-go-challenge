package repo

import (
	"fmt"
	"time"

	"platform-go-challenge/src/core"
)

// Compile-time assertion: *FavoriteRepo implements core.FavoriteRepository
var _ core.FavoriteRepository = (*FavoriteRepo)(nil)

type FavoriteRepo struct{ pg *PG }

func NewFavoriteRepo(pg *PG) *FavoriteRepo { return &FavoriteRepo{pg: pg} }

func (r *FavoriteRepo) Add(userID core.UserID, assetID core.AssetID, now time.Time) error {
	_, err := r.pg.DB.Exec(`
INSERT INTO favorites (user_id, asset_id, created_at)
VALUES ($1,$2,$3)
ON CONFLICT(user_id,asset_id) DO UPDATE SET created_at=EXCLUDED.created_at`,
		string(userID), string(assetID), now.UTC())
	return err
}

func (r *FavoriteRepo) Remove(userID core.UserID, assetID core.AssetID) error {
	_, err := r.pg.DB.Exec(`DELETE FROM favorites WHERE user_id=$1 AND asset_id=$2`,
		string(userID), string(assetID))
	return err
}

func (r *FavoriteRepo) ListByCursor(userID core.UserID, limit int, after *core.FavoriteEdge, typ *core.AssetType) ([]core.FavoriteEdge, error) {
	args := []any{string(userID)}
	q := `
SELECT f.asset_id, f.created_at
FROM favorites f
`
	if typ != nil {
		q += `JOIN assets a ON a.id = f.asset_id AND a.type = $2
WHERE f.user_id = $1 `
		args = append(args, string(*typ))
	} else {
		q += `WHERE f.user_id = $1 `
	}

	// Keyset for ORDER BY created_at DESC, asset_id ASC
	if after != nil {
		args = append(args, after.CreatedAt, string(after.AssetID))
		q += fmt.Sprintf(`AND (f.created_at < $%d OR (f.created_at = $%d AND f.asset_id > $%d)) `,
			len(args)-1, len(args)-1, len(args))
	}

	q += fmt.Sprintf(`ORDER BY f.created_at DESC, f.asset_id ASC LIMIT $%d`, len(args)+1)
	args = append(args, limit)

	rows, err := r.pg.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []core.FavoriteEdge
	for rows.Next() {
		var id string
		var at time.Time
		if err := rows.Scan(&id, &at); err != nil {
			return nil, err
		}
		out = append(out, core.FavoriteEdge{AssetID: core.AssetID(id), CreatedAt: at})
	}
	return out, rows.Err()
}