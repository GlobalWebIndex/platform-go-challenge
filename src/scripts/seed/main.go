package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"

	"platform-go-challenge/src/core"
	"platform-go-challenge/src/infra/repo"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required (e.g. postgres://yourapp:pass@localhost:5432/yourapp?sslmode=disable)")
	}

	// 1) run migrations + insert users (need users before favorites due to FK)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	migSQL, err := os.ReadFile("src/scripts/db/migrate.sql")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(string(migSQL)); err != nil {
		log.Fatal(err)
	}

	now := time.Now().UTC()

	// upsert two users (idempotent)
	upsertUser := func(id, email, fullName string, at time.Time) {
		_, err := db.Exec(`
INSERT INTO users (id, email, full_name, is_active, created_at, updated_at)
VALUES ($1,$2,$3, TRUE, $4, $4)
ON CONFLICT (id) DO UPDATE
SET email = EXCLUDED.email,
    full_name = EXCLUDED.full_name,
    is_active = EXCLUDED.is_active,
    updated_at = EXCLUDED.updated_at
`, id, email, fullName, at)
		if err != nil {
			log.Fatal(err)
		}
	}

	upsertUser("u1", "u1@example.com", "User One", now)
	upsertUser("u2", "u2@example.com", "User Two", now)

	_ = db.Close()

	// 2) open shared PG and split repos (NEW wiring)
	pg, err := repo.NewPG(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	assetsRepo := repo.NewAssetRepo(pg)
	favsRepo := repo.NewFavoriteRepo(pg)

	// 3) create 3 charts, 3 insights, 3 audiences
	// --- Charts ---
	ch1, _ := core.MakeChart(
		core.AssetID(uuid.NewString()),
		"Daily Active Users", "DAU last 7 days", "Day", "Users",
		[]core.ChartSeries{{Name: "DAU", Points: []core.ChartPoint{{X: 1, Y: 1200}, {X: 2, Y: 1300}, {X: 3, Y: 1500}}}},
		now,
	)
	ch2, _ := core.MakeChart(
		core.AssetID(uuid.NewString()),
		"Monthly Active Users", "MAU last 6 months", "Month", "Users",
		[]core.ChartSeries{{Name: "MAU", Points: []core.ChartPoint{{X: 1, Y: 10_000}, {X: 2, Y: 11_500}, {X: 3, Y: 12_200}}}},
		now,
	)
	ch3, _ := core.MakeChart(
		core.AssetID(uuid.NewString()),
		"Revenue by Day", "Revenue (last 7 days)", "Day", "Revenue",
		[]core.ChartSeries{{Name: "Revenue", Points: []core.ChartPoint{{X: 1, Y: 2200}, {X: 2, Y: 2800}, {X: 3, Y: 2600}}}},
		now,
	)

	// --- Insights ---
	in1, _ := core.MakeInsight(core.AssetID(uuid.NewString()),
		"Millennials Social Time", "Social usage stat",
		"40% of millennials spend more than 3 hours on social media daily",
		now,
	)
	in2, _ := core.MakeInsight(core.AssetID(uuid.NewString()),
		"Cart Abandonment", "E-commerce behavior",
		"~70% of online shopping carts are abandoned before purchase",
		now,
	)
	in3, _ := core.MakeInsight(core.AssetID(uuid.NewString()),
		"Push Opt-in", "Mobile growth",
		"Opt-in rates for push notifications average 60–70% on Android",
		now,
	)

	// --- Audiences ---
	gMale := core.GenderMale
	gFemale := core.GenderFemale
	minH3 := 3
	minH2 := 2
	a18 := 18
	a24 := 24
	a35 := 35
	a44 := 44
	aud1, _ := core.MakeAudience(core.AssetID(uuid.NewString()),
		"Heavy Social Males 24–35", "Segment for campaigns",
		&gMale, nil, &a24, &a35, &minH3, nil, now,
	)
	aud2, _ := core.MakeAudience(core.AssetID(uuid.NewString()),
		"Light Social Females 18–24", "Early brand adopters",
		&gFemale, nil, &a18, &a24, &minH2, nil, now,
	)
	country := "US"
	purchases := 3
	aud3, _ := core.MakeAudience(core.AssetID(uuid.NewString()),
		"High-Intent US 35–44", "Higher purchase signals",
		nil, &country, &a35, &a44, nil, &purchases, now,
	)

	assets := []core.AnyAsset{ch1, ch2, ch3, in1, in2, in3, aud1, aud2, aud3}

	// 4) upsert all assets (via AssetRepo)
	for _, a := range assets {
		if err := assetsRepo.Upsert(a); err != nil {
			log.Fatal(err)
		}
	}

	// 5) add favorites for u1 and u2 (via FavoriteRepo)
	addFav := func(user core.UserID, asset core.AnyAsset, offsetSeconds int) {
		if err := favsRepo.Add(user, asset.GetAssetBase().ID, now.Add(time.Duration(offsetSeconds)*time.Second)); err != nil {
			log.Fatal(err)
		}
	}
	// u1: 5 favorites
	addFav("u1", ch1, 0)
	addFav("u1", in1, 1)
	addFav("u1", aud1, 2)
	addFav("u1", ch2, 3)
	addFav("u1", in2, 4)

	// u2: 6 favorites
	addFav("u2", ch3, 0)
	addFav("u2", in3, 1)
	addFav("u2", aud2, 2)
	addFav("u2", aud3, 3)
	addFav("u2", ch1, 4)
	addFav("u2", in1, 5)

	log.Println("Seeded users, assets (3 of each type), and favorites for u1 & u2")
}
