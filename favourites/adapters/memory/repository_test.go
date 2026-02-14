package memory

import (
	"testing"

	"github.com/gwi/platform-go-challenge/favourites/domain"
)

func TestRepository_Add_and_ListPaginated(t *testing.T) {
	repo := NewRepository()
	userID := "user1"
	asset := domain.Asset{
		Type:        domain.AssetTypeChart,
		Description: "My chart",
		Chart:       &domain.ChartData{Title: "Test", XAxisTitle: "X", YAxisTitle: "Y", Data: []int{1, 2, 3}},
	}

	added := repo.Add(userID, asset)
	if added.ID == "" {
		t.Fatal("expected generated ID")
	}
	if added.Type != domain.AssetTypeChart {
		t.Errorf("expected type chart, got %s", added.Type)
	}

	items, total := repo.ListPaginated(userID, 10, 0)
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
	if items[0].ID != added.ID {
		t.Errorf("expected ID %s, got %s", added.ID, items[0].ID)
	}
}

func TestRepository_ListPaginated_empty(t *testing.T) {
	repo := NewRepository()
	items, total := repo.ListPaginated("nobody", 20, 0)
	if total != 0 || len(items) != 0 {
		t.Errorf("expected empty list, got total=%d items=%d", total, len(items))
	}
}

func TestRepository_ListPaginated_pagination(t *testing.T) {
	repo := NewRepository()
	userID := "user2"
	for i := 0; i < 5; i++ {
		repo.Add(userID, domain.Asset{Type: domain.AssetTypeInsight, Insight: &domain.InsightData{Text: "x"}})
	}

	page1, total := repo.ListPaginated(userID, 2, 0)
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(page1) != 2 {
		t.Errorf("expected 2 items on page 1, got %d", len(page1))
	}

	page2, _ := repo.ListPaginated(userID, 2, 2)
	if len(page2) != 2 {
		t.Errorf("expected 2 items on page 2, got %d", len(page2))
	}

	page3, _ := repo.ListPaginated(userID, 2, 4)
	if len(page3) != 1 {
		t.Errorf("expected 1 item on page 3, got %d", len(page3))
	}
}

func TestRepository_Get(t *testing.T) {
	repo := NewRepository()
	userID := "user3"
	added := repo.Add(userID, domain.Asset{Type: domain.AssetTypeChart, Chart: &domain.ChartData{Title: "T"}})

	got, ok := repo.Get(userID, added.ID)
	if !ok {
		t.Fatal("expected asset to be found")
	}
	if got.ID != added.ID {
		t.Errorf("expected ID %s, got %s", added.ID, got.ID)
	}

	_, ok = repo.Get(userID, "nonexistent")
	if ok {
		t.Error("expected not found for nonexistent ID")
	}
}

func TestRepository_Remove(t *testing.T) {
	repo := NewRepository()
	userID := "user4"
	added := repo.Add(userID, domain.Asset{Type: domain.AssetTypeInsight, Insight: &domain.InsightData{Text: "x"}})

	ok := repo.Remove(userID, added.ID)
	if !ok {
		t.Fatal("expected Remove to return true")
	}
	_, total := repo.ListPaginated(userID, 10, 0)
	if total != 0 {
		t.Errorf("expected 0 after remove, got %d", total)
	}

	ok = repo.Remove(userID, "nonexistent")
	if ok {
		t.Error("expected Remove to return false for nonexistent")
	}
}

func TestRepository_UpdateDescription(t *testing.T) {
	repo := NewRepository()
	userID := "user5"
	added := repo.Add(userID, domain.Asset{Type: domain.AssetTypeAudience, Description: "old", Audience: &domain.AudienceData{}})

	updated, ok := repo.UpdateDescription(userID, added.ID, "new description")
	if !ok {
		t.Fatal("expected UpdateDescription to succeed")
	}
	if updated.Description != "new description" {
		t.Errorf("expected description 'new description', got %s", updated.Description)
	}
	if updated.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	_, ok = repo.UpdateDescription(userID, "nonexistent", "x")
	if ok {
		t.Error("expected UpdateDescription to return false for nonexistent")
	}
}

func TestRepository_Add_with_existing_ID(t *testing.T) {
	repo := NewRepository()
	userID := "user6"
	asset := domain.Asset{
		ID:    "custom-id",
		Type:  domain.AssetTypeChart,
		Chart: &domain.ChartData{Title: "T"},
	}
	added := repo.Add(userID, asset)
	if added.ID != "custom-id" {
		t.Errorf("expected to keep custom ID, got %s", added.ID)
	}
	if added.CreatedAt.IsZero() || added.UpdatedAt.IsZero() {
		t.Error("expected CreatedAt and UpdatedAt to be set")
	}
}

func TestRepository_RemoveFavouritesBySourceAssetID(t *testing.T) {
	repo := NewRepository()
	sourceID := "catalog-asset-123"

	repo.Add("u1", domain.Asset{Type: domain.AssetTypeChart, SourceAssetID: sourceID, Chart: &domain.ChartData{Title: "A"}})
	repo.Add("u1", domain.Asset{Type: domain.AssetTypeInsight, SourceAssetID: sourceID, Insight: &domain.InsightData{Text: "B"}})
	repo.Add("u2", domain.Asset{Type: domain.AssetTypeChart, SourceAssetID: sourceID, Chart: &domain.ChartData{Title: "C"}})
	repo.Add("u1", domain.Asset{Type: domain.AssetTypeInsight, Insight: &domain.InsightData{Text: "no-source"}})

	removed := repo.RemoveFavouritesBySourceAssetID(sourceID)
	if removed != 3 {
		t.Errorf("expected 3 removed, got %d", removed)
	}

	_, total1 := repo.ListPaginated("u1", 10, 0)
	if total1 != 1 {
		t.Errorf("u1 should have 1 favourite left, got %d", total1)
	}
	_, total2 := repo.ListPaginated("u2", 10, 0)
	if total2 != 0 {
		t.Errorf("u2 should have 0 favourites, got %d", total2)
	}
}

func TestRepository_RemoveFavouritesBySourceAssetID_empty_returns_zero(t *testing.T) {
	repo := NewRepository()
	repo.Add("u1", domain.Asset{Type: domain.AssetTypeChart, Chart: &domain.ChartData{Title: "T"}})

	removed := repo.RemoveFavouritesBySourceAssetID("")
	if removed != 0 {
		t.Errorf("expected 0 when source_asset_id is empty, got %d", removed)
	}
	_, total := repo.ListPaginated("u1", 10, 0)
	if total != 1 {
		t.Errorf("favourites should be unchanged, got total %d", total)
	}
}
