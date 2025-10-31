package inmemory

import (
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	//lru "github.com/hashicorp/golang-lru/v2"
)

type mockAsset struct {
	id          string
	t           entities.AssetType
	validateErr error
}

/*
func (m *mockAsset) Validate() error             { return m.validateErr }
func (m *mockAsset) GetID() string               { return m.id }
func (m *mockAsset) GetType() entities.AssetType { return m.t }

func newCache(t *testing.T, size int) *lru.Cache[string, entities.AssetEntity] {
	t.Helper()
	c, err := lru.New[string, entities.AssetEntity](size)
	if err != nil {
		t.Fatalf("failed to create lru cache: %v", err)
	}
	return c
}

func TestSaveGetExistsUpdateDelete(t *testing.T) {
	cache := newCache(t, 10)
	repo := NewAssetRepository(cache)

	a := &mockAsset{id: "asset-1", t: entities.AssetType(1)}
	var ae entities.AssetEntity = a

	// Save
	if err := repo.Save(ae); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Exists
	ok, err := repo.Exists("asset-1")
	if err != nil {
		t.Fatalf("Exists returned error: %v", err)
	}
	if !ok {
		t.Fatalf("expected asset to exist")
	}

	// GetByID
	got, err := repo.GetByID("asset-1")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.GetID() != "asset-1" {
		t.Fatalf("unexpected id: %s", got.GetID())
	}

	// Update (valid)
	a2 := &mockAsset{id: "asset-1", t: entities.AssetType(2)}
	if err := repo.Update(a2); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	updated, err := repo.GetByID("asset-1")
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}
	if updated.GetType() != entities.AssetType(2) {
		t.Fatalf("expected type 2, got %v", updated.GetType())
	}

	// Update (non-existing)
	if err := repo.Update(&mockAsset{id: "nope", t: 0}); err == nil {
		t.Fatalf("expected error when updating non-existing asset")
	}

	// Delete
	if err := repo.Delete("asset-1"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if exists, _ := repo.Exists("asset-1"); exists {
		t.Fatalf("asset should have been deleted")
	}

	// Delete non-existing
	if err := repo.Delete("asset-1"); err == nil {
		t.Fatalf("expected error deleting non-existing asset")
	}
}

func TestGetByIDsGetAllGetByType(t *testing.T) {
	cache := newCache(t, 10)
	repo := NewAssetRepository(cache)

	a1 := &mockAsset{id: "a1", t: entities.AssetType(1)}
	a2 := &mockAsset{id: "a2", t: entities.AssetType(2)}
	a3 := &mockAsset{id: "a3", t: entities.AssetType(1)}

	if err := repo.Save(a1); err != nil {
		t.Fatalf("Save a1: %v", err)
	}
	if err := repo.Save(a2); err != nil {
		t.Fatalf("Save a2: %v", err)
	}
	if err := repo.Save(a3); err != nil {
		t.Fatalf("Save a3: %v", err)
	}

	// GetByIDs (including a missing one)
	ids := []string{"a1", "a2", "missing", "a3"}
	res, err := repo.GetByIDs(ids)
	if err != nil {
		t.Fatalf("GetByIDs returned error: %v", err)
	}
	if len(res) != 3 {
		t.Fatalf("expected 3 assets, got %d", len(res))
	}

	// GetAll
	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll error: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 assets from GetAll, got %d", len(all))
	}

	// GetByType
	type1, err := repo.GetByType(entities.AssetType(1))
	if err != nil {
		t.Fatalf("GetByType error: %v", err)
	}
	if len(type1) != 2 {
		t.Fatalf("expected 2 assets of type 1, got %d", len(type1))
	}
}

func TestGetAudienceChartInsightByID(t *testing.T) {
	cache := newCache(t, 10)
	repo := NewAssetRepository(cache)

	// Insert concrete entities directly into the cache to control dynamic type.
	// Use zero-values; assume these concrete types are defined in entities package.
	aud := &entities.AudienceEntity{}
	chart := &entities.ChartEntity{}
	insight := &entities.InsightEntity{}

	cache.Add("aud-1", aud)
	cache.Add("chart-1", chart)
	cache.Add("insight-1", insight)

	// Successful retrievals
	if _, err := repo.GetAudienceByID("aud-1"); err != nil {
		t.Fatalf("GetAudienceByID failed: %v", err)
	}
	if _, err := repo.GetChartByID("chart-1"); err != nil {
		t.Fatalf("GetChartByID failed: %v", err)
	}
	if _, err := repo.GetInsightByID("insight-1"); err != nil {
		t.Fatalf("GetInsightByID failed: %v", err)
	}

	// Wrong type retrieval should error
	// e.g., attempt to get audience from chart id
	if _, err := repo.GetAudienceByID("chart-1"); err == nil {
		t.Fatalf("expected error when getting Audience from a Chart")
	}
	if _, err := repo.GetChartByID("insight-1"); err == nil {
		t.Fatalf("expected error when getting Chart from an Insight")
	}
	if _, err := repo.GetInsightByID("aud-1"); err == nil {
		t.Fatalf("expected error when getting Insight from an Audience")
	}

	// Non-existent id should return asset not found
	if _, err := repo.GetAudienceByID("nope"); err == nil {
		t.Fatalf("expected error for non-existent id")
	}
}
*/
