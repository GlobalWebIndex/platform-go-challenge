package service

import (
	"encoding/json"
	"testing"

	"github.com/KostasDasios/platform-go-challenge/internal/repo"
	"github.com/KostasDasios/platform-go-challenge/internal/models"
)

func mustRaw(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func TestService_CreateListUpdateDelete(t *testing.T) {
	repo := repo.NewInMemoryRepo()
	svc := NewService(repo)

	user := "kostas"

	// create: insight
	insight := models.Insight{
		AssetBase: models.AssetBase{Type: models.AssetInsight, Description: "baseline"},
		Text:      "40% of users…",
	}
	f1, err := svc.CreateFavourite(user, mustRaw(insight))
	if err != nil {
		t.Fatalf("create insight: %v", err)
	}
	if f1.ID == "" || f1.Type != models.AssetInsight {
		t.Fatalf("unexpected favourite: %+v", f1)
	}

	// create: chart (valid)
	chart := models.Chart{
		AssetBase:  models.AssetBase{Type: models.AssetChart, Description: "c1"},
		Title:      "Sales",
		AxisXTitle: "Month",
		AxisYTitle: "€",
		Data:       []float64{1, 2, 3},
	}
	_, err = svc.CreateFavourite(user, mustRaw(chart))
	if err != nil {
		t.Fatalf("create chart: %v", err)
	}

	// list
	list, err := svc.ListFavourites(user)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 favourites, got %d", len(list))
	}

	// update description
	upd, err := svc.UpdateFavouriteDescription(user, f1.ID, "updated")
	if err != nil {
		t.Fatalf("update desc: %v", err)
	}
	if upd.Description != "updated" {
		t.Fatalf("desc not updated: %+v", upd)
	}

	// delete
	if err := svc.DeleteFavourite(user, f1.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	list, _ = svc.ListFavourites(user)
	if len(list) != 1 {
		t.Fatalf("expected 1 favourite after delete, got %d", len(list))
	}
}

func TestService_ValidationErrors(t *testing.T) {
	repo := repo.NewInMemoryRepo()
	svc := NewService(repo)

	// invalid user
	if _, err := svc.ListFavourites("!!!"); err == nil {
		t.Fatalf("expected invalid user id")
	}

	// invalid asset: unknown type
	raw := mustRaw(struct {
		Type string `json:"type"`
	}{Type: "unknown"})
	if _, err := svc.CreateFavourite("ok_user", raw); err == nil {
		t.Fatalf("expected error for unknown asset type")
	}

	// invalid chart: missing title/data
	badChart := models.Chart{
		AssetBase: models.AssetBase{Type: models.AssetChart},
	}
	if _, err := svc.CreateFavourite("ok_user", mustRaw(badChart)); err == nil {
		t.Fatalf("expected chart validation error")
	}
}
