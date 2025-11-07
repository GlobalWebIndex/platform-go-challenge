package storage

import (
	"assetsApp/internal/models"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

// ----------------- Asset Methods -----------------

func (p *PostgresStore) Add(userID uuid.UUID, asset models.Asset) {
	ctx := context.Background()

	// Ensure user exists
	_, err := p.pool.Exec(ctx,
		"INSERT INTO users (id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		userID, "User "+userID.String(),
	)
	if err != nil {
		log.Println("Failed to ensure user exists:", err)
		return
	}

	switch a := asset.(type) {
	case *models.Chart:
		tx, err := p.pool.Begin(ctx)
		if err != nil {
			log.Println("Failed to start transaction:", err)
			return
		}
		defer tx.Rollback(ctx)

		_, err = tx.Exec(ctx,
			"INSERT INTO assets (asset_id, title, description, asset_type, user_id) VALUES ($1, $2, $3, $4, $5)",
			a.ID, a.Title, a.Description, "chart", userID,
		)
		if err != nil {
			log.Println("Failed to insert into assets (chart):", err)
			return
		}

		_, err = tx.Exec(ctx,
			"INSERT INTO charts (id, title, description, x_axis_title, y_axis_title) VALUES ($1,$2,$3,$4,$5)",
			a.ID, a.Title, a.Description, a.XAxisTitle, a.YAxisTitle,
		)
		if err != nil {
			log.Println("Failed to insert chart:", err)
			return
		}

		for _, d := range a.Data {
			_, err = tx.Exec(ctx,
				"INSERT INTO chart_data (chart_id, datapoint_code, value) VALUES ($1,$2,$3)",
				a.ID, d.DatapointCode, d.Value,
			)
			if err != nil {
				log.Println("Failed to insert chart data:", err)
				return
			}
		}

		if err = tx.Commit(ctx); err != nil {
			log.Println("Failed to commit chart transaction:", err)
		}

	case *models.Insight:
		tx, err := p.pool.Begin(ctx)
		if err != nil {
			log.Println("Failed to start transaction:", err)
			return
		}
		defer tx.Rollback(ctx)

		_, err = tx.Exec(ctx,
			"INSERT INTO assets (asset_id, title, description, asset_type, user_id) VALUES ($1, $2, $3, $4, $5)",
			a.ID, "Insight", a.Description, "insight", userID,
		)
		if err != nil {
			log.Println("Failed to insert into assets:", err)
			return
		}

		_, err = tx.Exec(ctx,
			"INSERT INTO insights (id, description) VALUES ($1,$2)",
			a.ID, a.Description,
		)
		if err != nil {
			log.Println("Failed to insert insight:", err)
			return
		}

		if err = tx.Commit(ctx); err != nil {
			log.Println("Failed to commit insight transaction:", err)
		}

	case *models.Audience:
		tx, err := p.pool.Begin(ctx)
		if err != nil {
			log.Println("Failed to start transaction:", err)
			return
		}
		defer tx.Rollback(ctx)

		_, err = tx.Exec(ctx,
			"INSERT INTO assets (asset_id, title, description, asset_type, user_id) VALUES ($1, $2, $3, $4, $5)",
			a.ID, "Audience", a.Description, "audience", userID,
		)
		if err != nil {
			log.Println("Failed to insert into assets:", err)
			return
		}

		_, err = tx.Exec(ctx,
			"INSERT INTO audiences (id, gender, country, age_group, social_hours, purchases, description) VALUES ($1,$2,$3,$4,$5,$6,$7)",
			a.ID, a.Gender, a.Country, a.AgeGroup, a.SocialHours, a.Purchases, a.Description,
		)
		if err != nil {
			log.Println("Failed to insert audience:", err)
			return
		}

		if err = tx.Commit(ctx); err != nil {
			log.Println("Failed to commit audience transaction:", err)
		}
	}
}

func (p *PostgresStore) Get(userID uuid.UUID) []models.Asset {
	ctx := context.Background()
	rows, err := p.pool.Query(ctx, "SELECT asset_id, asset_type FROM assets WHERE user_id=$1", userID)
	if err != nil {
		log.Println("Failed to get assets:", err)
		return nil
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var assetID, assetType string
		if err := rows.Scan(&assetID, &assetType); err != nil {
			log.Println("Failed to scan asset row:", err)
			continue
		}

		var asset models.Asset
		switch assetType {
		case "chart":
			var c models.Chart
			err := p.pool.QueryRow(ctx, `
				SELECT id, title, description, x_axis_title, y_axis_title 
				FROM charts WHERE id=$1`, assetID).Scan(&c.ID, &c.Title, &c.Description, &c.XAxisTitle, &c.YAxisTitle)
			if err != nil {
				log.Println("Failed to fetch chart:", err)
				continue
			}

			// Fetch chart data
			dataRows, err := p.pool.Query(ctx, `SELECT datapoint_code, value FROM chart_data WHERE chart_id=$1`, assetID)
			if err != nil {
				log.Println("Failed to fetch chart data:", err)
			} else {
				defer dataRows.Close()
				for dataRows.Next() {
					var dp models.ChartData
					if err := dataRows.Scan(&dp.DatapointCode, &dp.Value); err != nil {
						log.Println("Failed to scan chart data row:", err)
						continue
					}
					c.Data = append(c.Data, dp)
				}
			}
			asset = &c

		case "insight":
			var i models.Insight
			err := p.pool.QueryRow(ctx, `SELECT id, description FROM insights WHERE id=$1`, assetID).
				Scan(&i.ID, &i.Description)
			if err != nil {
				log.Println("Failed to fetch insight:", err)
				continue
			}
			asset = &i

		case "audience":
			var a models.Audience
			err := p.pool.QueryRow(ctx, `
				SELECT id, gender, country, age_group, social_hours, purchases, description 
				FROM audiences WHERE id=$1`, assetID).Scan(&a.ID, &a.Gender, &a.Country, &a.AgeGroup, &a.SocialHours, &a.Purchases, &a.Description)
			if err != nil {
				log.Println("Failed to fetch audience:", err)
				continue
			}
			asset = &a
		}
		assets = append(assets, asset)
	}
	return assets
}

func (p *PostgresStore) Remove(userID uuid.UUID, assetID string) bool {
	ctx := context.Background()
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		log.Println("Failed to start remove transaction:", err)
		return false
	}
	defer tx.Rollback(ctx)

	statements := []struct {
		query string
		args  []interface{}
	}{
		{"DELETE FROM chart_data WHERE chart_id=$1", []interface{}{assetID}},
		{"DELETE FROM charts WHERE id=$1", []interface{}{assetID}},
		{"DELETE FROM insights WHERE id=$1", []interface{}{assetID}},
		{"DELETE FROM audiences WHERE id=$1", []interface{}{assetID}},
		{"DELETE FROM assets WHERE asset_id=$1", []interface{}{assetID}},
		{"DELETE FROM favourites WHERE asset_id=$1 AND user_id=$2", []interface{}{assetID, userID}},
	}

	for _, stmt := range statements {
		if _, err := tx.Exec(ctx, stmt.query, stmt.args...); err != nil {
			log.Println("Failed to execute remove statement:", err)
			return false
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit remove transaction:", err)
		return false
	}

	return true
}

func (p *PostgresStore) EditDescription(userID uuid.UUID, assetID, newDesc string) bool {
	ctx := context.Background()
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		log.Println("Failed to start edit transaction:", err)
		return false
	}
	defer tx.Rollback(ctx)

	// First, get the asset type
	var assetType string
	err = tx.QueryRow(ctx, "SELECT asset_type FROM assets WHERE asset_id=$1 AND user_id=$2", assetID, userID).Scan(&assetType)
	if err != nil {
		log.Println("Failed to find asset for user:", err)
		return false
	}

	// Update the specific asset table
	var stmt string
	switch assetType {
	case "chart":
		stmt = "UPDATE charts SET description=$1 WHERE id=$2"
	case "insight":
		stmt = "UPDATE insights SET description=$1 WHERE id=$2"
	case "audience":
		stmt = "UPDATE audiences SET description=$1 WHERE id=$2"
	default:
		log.Println("Unknown asset type:", assetType)
		return false
	}

	if _, err := tx.Exec(ctx, stmt, newDesc, assetID); err != nil {
		log.Println("Failed to update description in specific table:", err)
		return false
	}

	// Also update the main assets table
	if _, err := tx.Exec(ctx, "UPDATE assets SET description=$1 WHERE asset_id=$2", newDesc, assetID); err != nil {
		log.Println("Failed to update description in assets table:", err)
		return false
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit edit transaction:", err)
		return false
	}

	return true
}

// ----------------- Favourite Methods -----------------

func (p *PostgresStore) AddFavourite(userID uuid.UUID, assetID, _ string) bool {
	ctx := context.Background()

	// Ensure user exists
	_, err := p.pool.Exec(ctx,
		"INSERT INTO users (id, name) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING",
		userID, "Unknown",
	)
	if err != nil {
		log.Println("Failed to ensure user exists:", err)
		return false
	}

	// Fetch the asset type from assets table
	var assetType string
	err = p.pool.QueryRow(ctx, "SELECT asset_type FROM assets WHERE asset_id=$1", assetID).Scan(&assetType)
	if err != nil {
		log.Println("Failed to fetch asset type:", err)
		return false
	}

	_, err = p.pool.Exec(ctx,
		"INSERT INTO favourites (user_id, asset_id, asset_type) VALUES ($1, $2, $3) ON CONFLICT (user_id, asset_id) DO NOTHING",
		userID, assetID, assetType,
	)
	if err != nil {
		log.Println("Failed to add favourite:", err)
		return false
	}

	log.Printf("Favourite added: user=%v, asset=%s, type=%s", userID, assetID, assetType)
	return true
}

func (p *PostgresStore) RemoveFavourite(userID uuid.UUID, assetID string) bool {
	_, err := p.pool.Exec(context.Background(),
		"DELETE FROM favourites WHERE user_id=$1 AND asset_id=$2",
		userID, assetID,
	)
	if err != nil {
		log.Println("Failed to remove favourite:", err)
		return false
	}
	return true
}

func (p *PostgresStore) GetFavourites(userID uuid.UUID) []models.Favourite {
	ctx := context.Background()

	rows, err := p.pool.Query(ctx,
		"SELECT asset_id, asset_type FROM favourites WHERE user_id=$1", userID)
	if err != nil {
		log.Println("Failed to get favourites:", err)
		return nil
	}
	defer rows.Close()

	log.Printf("Query executed for user %v", userID)

	var favs []models.Favourite
	for rows.Next() {
		var assetID, assetType string
		if err := rows.Scan(&assetID, &assetType); err != nil {
			log.Println("Failed to scan favourite row:", err)
			continue
		}
		log.Printf("Processing favourite: assetID=%s, assetType=%s", assetID, assetType)

		var asset models.Asset

		switch assetType {
		case "chart":
			var c models.Chart
			err := p.pool.QueryRow(ctx, `
				SELECT id, title, description, x_axis_title, y_axis_title
				FROM charts WHERE id=$1`, assetID).Scan(&c.ID, &c.Title, &c.Description, &c.XAxisTitle, &c.YAxisTitle)
			if err != nil {
				log.Println("Failed to fetch chart:", err)
				continue
			}

			// Fetch chart data
			dataRows, err := p.pool.Query(ctx, `SELECT datapoint_code, value FROM chart_data WHERE chart_id=$1`, assetID)
			if err != nil {
				log.Println("Failed to fetch chart data:", err)
			} else {
				defer dataRows.Close()
				for dataRows.Next() {
					var dp models.ChartData
					if err := dataRows.Scan(&dp.DatapointCode, &dp.Value); err != nil {
						log.Println("Failed to scan chart data row:", err)
						continue
					}
					c.Data = append(c.Data, dp)
				}
				log.Printf("Fetched %d data points for chart %s", len(c.Data), c.ID)
			}

			asset = &c
			log.Printf("Chart fetched successfully: %s", c.Title)

		case "insight":
			var i models.Insight
			err := p.pool.QueryRow(ctx, `SELECT id, description FROM insights WHERE id=$1`, assetID).
				Scan(&i.ID, &i.Description)
			if err != nil {
				log.Println("Failed to fetch insight:", err)
				continue
			}
			asset = &i
			log.Printf("Insight fetched successfully: %s", i.ID)

		case "audience":
			var a models.Audience
			err := p.pool.QueryRow(ctx, `
				SELECT id, gender, country, age_group, social_hours, purchases, description
				FROM audiences WHERE id=$1`, assetID).Scan(&a.ID, &a.Gender, &a.Country, &a.AgeGroup, &a.SocialHours, &a.Purchases, &a.Description)
			if err != nil {
				log.Println("Failed to fetch audience:", err)
				continue
			}
			asset = &a
			log.Printf("Audience fetched successfully: %s", a.ID)

		default:
			log.Printf("Unknown asset type: %s", assetType)
			continue
		}

		favs = append(favs, models.Favourite{
			UserID: userID,
			Asset:  asset,
		})
		log.Printf("Favourite appended for user %v: assetID=%s", userID, assetID)
	}

	if err := rows.Err(); err != nil {
		log.Println("Row iteration error:", err)
	}

	log.Printf("GetFavourites finished for user %v, total favourites: %d", userID, len(favs))
	return favs
}
