package application

import (
	"context"
	"testing"

	"github.com/manos/favourites/favourites/domain"
)

type favRepoStub struct {
	items []domain.FavouriteEntity
}

func (r *favRepoStub) FindByUser(_ context.Context, _ string) ([]domain.FavouriteEntity, error) {
	return r.items, nil
}

func (r *favRepoStub) FindByUserPage(_ context.Context, _ string, offset int, limit int) ([]domain.FavouriteEntity, int, error) {
	end := offset + limit
	if end > len(r.items) {
		end = len(r.items)
	}
	if offset > len(r.items) {
		offset = len(r.items)
	}
	return r.items[offset:end], len(r.items), nil
}

func (r *favRepoStub) Save(_ context.Context, _ domain.FavouriteEntity) error { return nil }
func (r *favRepoStub) FindByID(_ context.Context, _ string) (domain.FavouriteEntity, error) {
	return domain.FavouriteEntity{}, nil
}
func (r *favRepoStub) Delete(_ context.Context, _ string) error                      { return nil }
func (r *favRepoStub) UpdateDescription(_ context.Context, _ string, _ string) error { return nil }

type assetClientStub struct {
	insights  map[string]AssetDTO
	audiences map[string]AssetDTO
	charts    map[string]AssetDTO
}

func (c *assetClientStub) GetInsight(_ context.Context, _ string, _ string) (AssetDTO, error) {
	return AssetDTO{}, nil
}
func (c *assetClientStub) GetAudience(_ context.Context, _ string, _ string) (AssetDTO, error) {
	return AssetDTO{}, nil
}
func (c *assetClientStub) GetChart(_ context.Context, _ string, _ string) (AssetDTO, error) {
	return AssetDTO{}, nil
}
func (c *assetClientStub) GetInsights(_ context.Context, _ string, _ []string) (map[string]AssetDTO, error) {
	return c.insights, nil
}
func (c *assetClientStub) GetAudiences(_ context.Context, _ string, _ []string) (map[string]AssetDTO, error) {
	return c.audiences, nil
}
func (c *assetClientStub) GetCharts(_ context.Context, _ string, _ []string) (map[string]AssetDTO, error) {
	return c.charts, nil
}

// Deprecated: the original monolithic test was split into focused tests below.
// setupTestService returns a service wired with a repo and asset client
// containing three favourites (ins-1, aud-1, chr-1).
func setupTestService() *FavouriteService {
	repo := &favRepoStub{
		items: []domain.FavouriteEntity{
			{ID: "fav-1", UserID: "user-1", AssetID: "ins-1", Type: domain.FavouriteInsight, Description: "one"},
			{ID: "fav-2", UserID: "user-1", AssetID: "aud-1", Type: domain.FavouriteAudience, Description: "two"},
			{ID: "fav-3", UserID: "user-1", AssetID: "chr-1", Type: domain.FavouriteChart, Description: "three"},
		},
	}
	assets := &assetClientStub{
		insights:  map[string]AssetDTO{"ins-1": {ID: "ins-1", Name: "insight", Type: domain.FavouriteInsight}},
		audiences: map[string]AssetDTO{"aud-1": {ID: "aud-1", Name: "audience", Type: domain.FavouriteAudience}},
		charts:    map[string]AssetDTO{"chr-1": {ID: "chr-1", Name: "chart", Type: domain.FavouriteChart}},
	}
	return NewFavouriteService(repo, assets)
}

func TestGetFavouritesForUser_TotalMatchesRepositoryLength(t *testing.T) {
	svc := setupTestService()
	page, err := svc.GetFavouritesForUser(context.Background(), "user-1", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := page.Total, 3; got != want {
		t.Fatalf("total = %d, want %d", got, want)
	}
}

func TestGetFavouritesForUser_RespectsLimit(t *testing.T) {
	svc := setupTestService()
	page, err := svc.GetFavouritesForUser(context.Background(), "user-1", 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := len(page.Items), 2; got != want {
		t.Fatalf("items length = %d, want %d", got, want)
	}
}

func TestGetFavouritesForUser_ReturnsCorrectAssetIDsForPage(t *testing.T) {
	svc := setupTestService()
	page, err := svc.GetFavouritesForUser(context.Background(), "user-1", 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := []string{page.Items[0].Asset.ID, page.Items[1].Asset.ID}
	want := []string{"ins-1", "aud-1"}
	if len(got) != len(want) {
		t.Fatalf("unexpected items length: got=%d want=%d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("asset id at index %d = %s, want %s", i, got[i], want[i])
		}
	}
}

func TestGetFavouritesForUser_PageOutOfRangeReturnsEmptyItems(t *testing.T) {
	svc := setupTestService()
	page, err := svc.GetFavouritesForUser(context.Background(), "user-1", 10, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := len(page.Items), 0; got != want {
		t.Fatalf("items length for out-of-range page = %d, want %d", got, want)
	}
}
