package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"platform-go-challenge/src/core"
	apihttp "platform-go-challenge/src/infra/http"
	"platform-go-challenge/src/infra/repo"
)

func main() {
	addr := getenv("ADDR", ":8080")
	dsn := mustenv("DATABASE_URL")

	// one shared connection pool
	pg, err := repo.NewPG(dsn)
	if err != nil { log.Fatal(err) }
	defer pg.Close()

	// split repos -> interfaces
	var assetsRepo core.AssetRepository    = repo.NewAssetRepo(pg)
	var favoritesRepo core.FavoriteRepository = repo.NewFavoriteRepo(pg)

	assetSvc := core.NewAssetService(assetsRepo)
	favSvc   := core.NewFavoriteService(assetsRepo, favoritesRepo)
	now := func() time.Time { return time.Now().UTC() }

	h := apihttp.NewHandlers(assetSvc, favSvc, now)
	r := apihttp.SetupRouter(h)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string { if v := os.Getenv(k); v != "" { return v }; return def }
func mustenv(k string) string     { v := os.Getenv(k); if v == "" { log.Fatalf("%s is required", k) }; return v }
