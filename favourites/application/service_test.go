package application

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gwi/platform-go-challenge/favourites/domain"
)

type fakeRepo struct {
	users   map[string][]domain.Asset
	counter int
}

func (f *fakeRepo) ListPaginated(userID string, limit, offset int) ([]domain.Asset, int) {
	list := f.users[userID]
	total := len(list)
	if offset >= total {
		return []domain.Asset{}, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	out := make([]domain.Asset, end-offset)
	copy(out, list[offset:end])
	return out, total
}

func (f *fakeRepo) Get(userID string, assetID string) (*domain.Asset, bool) {
	for i := range f.users[userID] {
		if f.users[userID][i].ID == assetID {
			a := f.users[userID][i]
			return &a, true
		}
	}
	return nil, false
}

func (f *fakeRepo) Add(userID string, asset domain.Asset) domain.Asset {
	if asset.ID == "" {
		f.counter++
		asset.ID = fmt.Sprintf("fake-id-%d", f.counter)
	}
	f.users[userID] = append(f.users[userID], asset)
	return asset
}

func (f *fakeRepo) Remove(userID string, assetID string) bool {
	list := f.users[userID]
	for i, a := range list {
		if a.ID == assetID {
			f.users[userID] = append(list[:i], list[i+1:]...)
			return true
		}
	}
	return false
}

func (f *fakeRepo) UpdateDescription(userID string, assetID string, description string) (*domain.Asset, bool) {
	for i := range f.users[userID] {
		if f.users[userID][i].ID == assetID {
			f.users[userID][i].Description = description
			a := f.users[userID][i]
			return &a, true
		}
	}
	return nil, false
}

func (f *fakeRepo) RemoveFavouritesBySourceAssetID(sourceAssetID string) int {
	var n int
	for userID, list := range f.users {
		for i := len(list) - 1; i >= 0; i-- {
			if list[i].SourceAssetID == sourceAssetID {
				f.users[userID] = append(list[:i], list[i+1:]...)
				list = f.users[userID]
				n++
			}
		}
	}
	return n
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{users: make(map[string][]domain.Asset), counter: 0}
}

var _ FavouritesRepository = (*fakeRepo)(nil)

func TestFavouritesService_Add_and_ListPaginated(t *testing.T) {
	svc := NewFavouritesService(newFakeRepo(), nil)
	userID := "u1"
	asset := domain.Asset{Type: domain.AssetTypeChart, Description: "Chart", Chart: &domain.ChartData{Title: "T"}}
	added, err := svc.AddFavourite(userID, asset)
	if err != nil {
		t.Fatal(err)
	}
	if added.ID == "" {
		t.Fatal("expected generated ID")
	}
	items, total := svc.ListFavouritesPaginated(userID, 10, 0)
	if total != 1 || len(items) != 1 {
		t.Errorf("expected 1 item total=%d, got total=%d len=%d", 1, total, len(items))
	}
}

func TestFavouritesService_GetFavourite(t *testing.T) {
	svc := NewFavouritesService(newFakeRepo(), nil)
	userID := "u2"
	added, err := svc.AddFavourite(userID, domain.Asset{Type: domain.AssetTypeInsight, Insight: &domain.InsightData{Text: "insight"}})
	if err != nil {
		t.Fatal(err)
	}
	got, ok := svc.GetFavourite(userID, added.ID)
	if !ok {
		t.Fatal("expected found")
	}
	if got.ID != added.ID {
		t.Errorf("expected ID %s, got %s", added.ID, got.ID)
	}
	_, ok = svc.GetFavourite(userID, "nonexistent")
	if ok {
		t.Error("expected not found")
	}
}

func TestFavouritesService_RemoveFavourite(t *testing.T) {
	svc := NewFavouritesService(newFakeRepo(), nil)
	userID := "u3"
	added, err := svc.AddFavourite(userID, domain.Asset{Type: domain.AssetTypeChart, Chart: &domain.ChartData{}})
	if err != nil {
		t.Fatal(err)
	}
	ok := svc.RemoveFavourite(userID, added.ID)
	if !ok {
		t.Fatal("expected Remove to succeed")
	}
	_, total := svc.ListFavouritesPaginated(userID, 10, 0)
	if total != 0 {
		t.Errorf("expected 0 after remove, got %d", total)
	}
}

func TestFavouritesService_UpdateFavouriteDescription(t *testing.T) {
	svc := NewFavouritesService(newFakeRepo(), nil)
	userID := "u4"
	added, err := svc.AddFavourite(userID, domain.Asset{Type: domain.AssetTypeAudience, Description: "old", Audience: &domain.AudienceData{}})
	if err != nil {
		t.Fatal(err)
	}
	updated, ok := svc.UpdateFavouriteDescription(userID, added.ID, "new desc")
	if !ok {
		t.Fatal("expected Update to succeed")
	}
	if updated.Description != "new desc" {
		t.Errorf("expected 'new desc', got %s", updated.Description)
	}
}

// fakeAssetCatalog: returns ErrAssetNotFound for "missing-asset"; for "chart-1" returns a seeded chart (add-by-reference test).
type fakeAssetCatalog struct {
	byKey map[string]*domain.Asset // "type:id" -> asset
}

func newFakeAssetCatalog() *fakeAssetCatalog {
	return &fakeAssetCatalog{
		byKey: map[string]*domain.Asset{
			"chart:chart-1": {
				Type:          domain.AssetTypeChart,
				SourceAssetID: "chart-1",
				Description:   "Catalog chart",
				Chart:         &domain.ChartData{Title: "Sales", XAxisTitle: "Month", YAxisTitle: "Revenue"},
			},
		},
	}
}

func (f *fakeAssetCatalog) GetAsset(userID string, sourceAssetID string, assetType string) (*domain.Asset, error) {
	if sourceAssetID == "missing-asset" {
		return nil, ErrAssetNotFound
	}
	key := assetType + ":" + sourceAssetID
	if a, ok := f.byKey[key]; ok {
		out := *a
		return &out, nil
	}
	return nil, ErrAssetNotFound
}

var _ AssetCatalogClient = (*fakeAssetCatalog)(nil)

func TestFavouritesService_AddFavourite_rejects_unknown_source_asset_id(t *testing.T) {
	svc := NewFavouritesService(newFakeRepo(), newFakeAssetCatalog())
	asset := domain.Asset{
		Type:          domain.AssetTypeChart,
		SourceAssetID: "missing-asset",
		Chart:         &domain.ChartData{Title: "T"},
	}
	_, err := svc.AddFavourite("u1", asset)
	if err == nil {
		t.Fatal("expected ErrAssetNotFound")
	}
	if !errors.Is(err, ErrAssetNotFound) {
		t.Errorf("expected ErrAssetNotFound, got %v", err)
	}
}

func TestFavouritesService_AddFavourite_add_by_reference(t *testing.T) {
	svc := NewFavouritesService(newFakeRepo(), newFakeAssetCatalog())
	// Add with only source_asset_id and type (add-by-reference); no chart body.
	asset := domain.Asset{Type: domain.AssetTypeChart, SourceAssetID: "chart-1"}
	added, err := svc.AddFavourite("u1", asset)
	if err != nil {
		t.Fatal(err)
	}
	if added.SourceAssetID != "chart-1" {
		t.Errorf("expected source_asset_id chart-1, got %s", added.SourceAssetID)
	}
	if added.Chart == nil {
		t.Fatal("expected chart data from catalog")
	}
	if added.Chart.Title != "Sales" {
		t.Errorf("expected chart title Sales, got %s", added.Chart.Title)
	}
	if added.Description != "Catalog chart" {
		t.Errorf("expected description from catalog, got %s", added.Description)
	}
}
