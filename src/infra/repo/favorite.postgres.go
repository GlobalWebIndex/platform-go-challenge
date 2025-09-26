package repo

import (
	"encoding/json"
	"fmt"
	"time"

	"platform-go-challenge/src/core"
)

var _ core.FavoriteRepository = (*FavoriteRepo)(nil)

type FavoriteRepo struct{ pg *PG }
func NewFavoriteRepo(pg *PG) *FavoriteRepo { return &FavoriteRepo{pg: pg} }

func (r *FavoriteRepo) Add(userID core.UserID, assetID core.AssetID, now time.Time) error {
	_, err := r.pg.DB.Exec(`
INSERT INTO favorites (user_id, asset_id, created_at, override, override_updated_at)
VALUES ($1,$2,$3, COALESCE('{}'::jsonb, '{}'::jsonb), NULL)
ON CONFLICT(user_id,asset_id) DO NOTHING`,
		string(userID), string(assetID), now.UTC())
	return err
}

func (r *FavoriteRepo) Remove(userID core.UserID, assetID core.AssetID) error {
	_, err := r.pg.DB.Exec(`DELETE FROM favorites WHERE user_id=$1 AND asset_id=$2`,
		string(userID), string(assetID))
	return err
}

// NEW: merge in override JSON (only on existing favorite)
func (r *FavoriteRepo) SetOverride(userID core.UserID, assetID core.AssetID, override map[string]any, now time.Time) error {
	if override == nil || len(override) == 0 {
		// No-op
		return nil
	}
	b, err := json.Marshal(override)
	if err != nil { return err }

	res, err := r.pg.DB.Exec(`
UPDATE favorites
   SET override = COALESCE(override, '{}'::jsonb) || $3::jsonb,
       override_updated_at = $4
 WHERE user_id = $1 AND asset_id = $2`,
		string(userID), string(assetID), string(b), now.UTC())
	if err != nil { return err }

	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("favorite not found for user_id=%s asset_id=%s", userID, assetID)
	}
	return nil
}

func (r *FavoriteRepo) ListByCursor(userID core.UserID, limit int, after *core.FavoriteEdge, typ *core.AssetType) ([]core.FavoriteEdge, error) {
	args := []any{string(userID)}
	q := `
SELECT f.asset_id, f.created_at, COALESCE(f.override, '{}'::jsonb) AS override
FROM favorites f
`
	if typ != nil {
		q += `JOIN assets a ON a.id = f.asset_id AND a.type = $2
WHERE f.user_id = $1 `
		args = append(args, string(*typ))
	} else {
		q += `WHERE f.user_id = $1 `
	}

	if after != nil {
		args = append(args, after.CreatedAt, string(after.AssetID))
		q += fmt.Sprintf(`AND (f.created_at < $%d OR (f.created_at = $%d AND f.asset_id > $%d)) `,
			len(args)-1, len(args)-1, len(args))
	}

	q += fmt.Sprintf(`ORDER BY f.created_at DESC, f.asset_id ASC LIMIT $%d`, len(args)+1)
	args = append(args, limit)

	rows, err := r.pg.DB.Query(q, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var out []core.FavoriteEdge
	for rows.Next() {
		var id string
		var at time.Time
		var ovRaw []byte
		if err := rows.Scan(&id, &at, &ovRaw); err != nil { return nil, err }
		var ov map[string]any
		if len(ovRaw) > 0 {
			_ = json.Unmarshal(ovRaw, &ov)
		}
		out = append(out, core.FavoriteEdge{
			AssetID:   core.AssetID(id),
			CreatedAt: at,
			Override:  ov,
		})
	}
	return out, rows.Err()
}
