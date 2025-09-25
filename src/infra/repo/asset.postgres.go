package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"platform-go-challenge/src/core"
)

// Compile-time assertion: *AssetRepo implements core.AssetRepository
var _ core.AssetRepository = (*AssetRepo)(nil)

type AssetRepo struct{ pg *PG }

func NewAssetRepo(pg *PG) *AssetRepo { return &AssetRepo{pg: pg} }

func (r *AssetRepo) Upsert(a core.AnyAsset) error {
	base := a.GetAssetBase()
	payload, typ := encodePayload(a)

	_, err := r.pg.DB.Exec(`
INSERT INTO assets (id, type, title, description, payload, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7)
ON CONFLICT(id) DO UPDATE SET
  title=EXCLUDED.title,
  description=EXCLUDED.description,
  payload=EXCLUDED.payload,
  updated_at=EXCLUDED.updated_at
`, string(base.ID), string(typ), base.Title, base.Description, payload, base.CreatedAt, base.UpdatedAt)
	return err
}

func (r *AssetRepo) GetByID(id core.AssetID) (core.AnyAsset, error) {
	row := r.pg.DB.QueryRow(`SELECT id,type,title,description,payload,created_at,updated_at FROM assets WHERE id=$1`, string(id))
	return scanAsset(row)
}

func (r *AssetRepo) UpdateDescription(id core.AssetID, desc string, now time.Time) error {
	_, err := r.pg.DB.Exec(`UPDATE assets SET description=$1, updated_at=$2 WHERE id=$3`, desc, now.UTC(), string(id))
	return err
}

func (r *AssetRepo) ListByIDsPreserveOrder(ids []core.AssetID) ([]core.AnyAsset, error) {
	// Minimal/clear implementation (N lookups). Fine for small pages.
	out := make([]core.AnyAsset, 0, len(ids))
	for _, id := range ids {
		a, err := r.GetByID(id)
		if err == nil {
			out = append(out, a)
		}
	}
	return out, nil
}

func (r *AssetRepo) Exists(id core.AssetID) (bool, error) {
	row := r.pg.DB.QueryRow(`SELECT 1 FROM assets WHERE id=$1`, string(id))
	var x int
	if err := row.Scan(&x); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ----- helpers (asset payload <-> JSON) -----

func encodePayload(a core.AnyAsset) ([]byte, core.AssetType) {
	switch x := a.(type) {
	case *core.Chart:
		b, _ := json.Marshal(map[string]any{
			"xTitle": x.XAxisTitle,
			"yTitle": x.YAxisTitle,
			"series": x.Series,
		})
		return b, core.AssetTypeChart
	case *core.Insight:
		b, _ := json.Marshal(map[string]any{"text": x.Text})
		return b, core.AssetTypeInsight
	case *core.Audience:
		b, _ := json.Marshal(map[string]any{
			"gender":                x.Gender,
			"birthCountry":          x.BirthCountry,
			"ageMin":                x.AgeMin,
			"ageMax":                x.AgeMax,
			"minHoursSocialPerDay":  x.MinHoursSocialPerDay,
			"minPurchasesLastMonth": x.MinPurchasesLastMonth,
		})
		return b, core.AssetTypeAudience
	default:
		return []byte("{}"), core.AssetType("unknown")
	}
}

func scanAsset(scanner interface{ Scan(dest ...any) error }) (core.AnyAsset, error) {
	var id, typ, title, desc string
	var payload []byte
	var created, updated time.Time
	if err := scanner.Scan(&id, &typ, &title, &desc, &payload, &created, &updated); err != nil {
		return nil, err
	}
	switch core.AssetType(typ) {
	case core.AssetTypeChart:
		var p struct {
			XTitle string             `json:"xTitle"`
			YTitle string             `json:"yTitle"`
			Series []core.ChartSeries `json:"series"`
		}
		if err := json.Unmarshal(payload, &p); err != nil {
			return nil, err
		}
		return &core.Chart{
			AssetBase:  core.AssetBase{ID: core.AssetID(id), Type: core.AssetTypeChart, Title: title, Description: desc, CreatedAt: created, UpdatedAt: updated},
			XAxisTitle: p.XTitle,
			YAxisTitle: p.YTitle,
			Series:     p.Series,
		}, nil
	case core.AssetTypeInsight:
		var p struct{ Text string `json:"text"` }
		if err := json.Unmarshal(payload, &p); err != nil {
			return nil, err
		}
		return &core.Insight{
			AssetBase: core.AssetBase{ID: core.AssetID(id), Type: core.AssetTypeInsight, Title: title, Description: desc, CreatedAt: created, UpdatedAt: updated},
			Text:      p.Text,
		}, nil
	case core.AssetTypeAudience:
		var p struct {
			Gender                *core.Gender `json:"gender"`
			BirthCountry          *string      `json:"birthCountry"`
			AgeMin                *int         `json:"ageMin"`
			AgeMax                *int         `json:"ageMax"`
			MinHoursSocialPerDay  *int         `json:"minHoursSocialPerDay"`
			MinPurchasesLastMonth *int         `json:"minPurchasesLastMonth"`
		}
		if err := json.Unmarshal(payload, &p); err != nil {
			return nil, err
		}
		return &core.Audience{
			AssetBase:             core.AssetBase{ID: core.AssetID(id), Type: core.AssetTypeAudience, Title: title, Description: desc, CreatedAt: created, UpdatedAt: updated},
			Gender:                p.Gender,
			BirthCountry:          p.BirthCountry,
			AgeMin:                p.AgeMin,
			AgeMax:                p.AgeMax,
			MinHoursSocialPerDay:  p.MinHoursSocialPerDay,
			MinPurchasesLastMonth: p.MinPurchasesLastMonth,
		}, nil
	default:
		return nil, fmt.Errorf("unknown asset type %q", typ)
	}
}