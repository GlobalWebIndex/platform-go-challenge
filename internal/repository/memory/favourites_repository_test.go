package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/domain/favourites"
)

func TestCreateFavourite_VariousScenarios(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		asset     favourites.Asset
		userID    string
		wantErr   error
		wantIDSet bool
	}{
		{
			name:   "success_assigns_id_and_links_favourite",
			userID: "user-1",
			asset: favourites.Asset{
				Type: favourites.AssetTypeChart,
				Chart: &favourites.Chart{
					Title: "Example",
					XAxis: "Age",
					YAxis: "Score",
					Data:  []float64{1, 2, 3},
				},
				Description: "My chart",
			},
			wantErr:   nil,
			wantIDSet: true,
		},
		{
			name:   "non_empty_asset_id_returns_error",
			userID: "user-1",
			asset: favourites.Asset{
				ID:   "pre-set",
				Type: favourites.AssetTypeInsight,
				Insight: &favourites.Insight{
					Text: "Some insight",
				},
			},
			wantErr:   favourites.ErrInvalidAssetPayload,
			wantIDSet: false,
		},
		{
			name:   "invalid_asset_payload_returns_error",
			userID: "user-1",
			asset: favourites.Asset{
				// Mismatched type/payload: type says chart but Chart is nil,
				// which should fail Asset.Validate().
				Type:        favourites.AssetTypeChart,
				Description: "Invalid chart",
				Chart:       nil,
			},
			wantErr:   favourites.ErrInvalidAssetPayload,
			wantIDSet: false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repo := NewFavouritesRepository()

			entry, err := repo.CreateFavourite(ctx, tt.userID, tt.asset)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("CreateFavourite error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				// On error we don't care about entry fields.
				return
			}

			if entry.Asset.ID == "" && tt.wantIDSet {
				t.Fatalf("expected asset ID to be set, got empty string")
			}

			if entry.Favourite.UserID != tt.userID {
				t.Fatalf("expected Favourite.UserID %q, got %q", tt.userID, entry.Favourite.UserID)
			}

			if entry.Favourite.AssetID != entry.Asset.ID {
				t.Fatalf("expected Favourite.AssetID %q to match Asset.ID %q",
					entry.Favourite.AssetID, entry.Asset.ID)
			}
		})
	}
}

func TestListFavourites_ByUserAndTypeFilter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := NewFavouritesRepository()

	// User 1: chart + insight
	chart := favourites.Asset{
		Type: favourites.AssetTypeChart,
		Chart: &favourites.Chart{
			Title: "Chart",
		},
	}
	insight := favourites.Asset{
		Type: favourites.AssetTypeInsight,
		Insight: &favourites.Insight{
			Text: "Insight",
		},
	}
	if _, err := repo.CreateFavourite(ctx, "user-1", chart); err != nil {
		t.Fatalf("setup CreateFavourite user-1 chart: %v", err)
	}
	if _, err := repo.CreateFavourite(ctx, "user-1", insight); err != nil {
		t.Fatalf("setup CreateFavourite user-1 insight: %v", err)
	}

	// User 2: one audience
	audience := favourites.Asset{
		Type: favourites.AssetTypeAudience,
		Audience: &favourites.Audience{
			Gender: "Male",
		},
	}
	if _, err := repo.CreateFavourite(ctx, "user-2", audience); err != nil {
		t.Fatalf("setup CreateFavourite user-2 audience: %v", err)
	}

	// List all for user-1
	entries, total, err := repo.ListFavourites(ctx, "user-1", favourites.ListFilter{})
	if err != nil {
		t.Fatalf("ListFavourites returned error: %v", err)
	}
	if got, want := len(entries), 2; got != want {
		t.Fatalf("expected %d entries for user-1, got %d", want, got)
	}
	if total != 2 {
		t.Fatalf("expected total 2 for user-1, got %d", total)
	}

	// List only charts for user-1
	tChart := favourites.AssetTypeChart
	filter := favourites.ListFilter{Type: &tChart}

	chartEntries, totalCharts, err := repo.ListFavourites(ctx, "user-1", filter)
	if err != nil {
		t.Fatalf("ListFavourites with filter returned error: %v", err)
	}
	if got, want := len(chartEntries), 1; got != want {
		t.Fatalf("expected %d chart entry for user-1, got %d", want, got)
	}
	if chartEntries[0].Asset.Type != favourites.AssetTypeChart {
		t.Fatalf("expected AssetTypeChart, got %v", chartEntries[0].Asset.Type)
	}
	if totalCharts != 1 {
		t.Fatalf("expected totalCharts=1, got %d", totalCharts)
	}

	// User-2 should not see user-1's favourites
	entries2, total2, err := repo.ListFavourites(ctx, "user-2", favourites.ListFilter{})
	if err != nil {
		t.Fatalf("ListFavourites user-2 returned error: %v", err)
	}
	if got, want := len(entries2), 1; got != want {
		t.Fatalf("expected %d entry for user-2, got %d", want, got)
	}
	if entries2[0].Asset.Type != favourites.AssetTypeAudience {
		t.Fatalf("expected AssetTypeAudience, got %v", entries2[0].Asset.Type)
	}
	if total2 != 1 {
		t.Fatalf("expected total2=1 for user-2, got %d", total2)
	}
}

func TestListFavourites_Pagination(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := NewFavouritesRepository()

	userID := "user-1"

	// Create 5 favourites for the same user in sequence so CreatedAt ordering is stable.
	for i := 0; i < 5; i++ {
		asset := favourites.Asset{
			Type: favourites.AssetTypeInsight,
			Insight: &favourites.Insight{
				Text: fmt.Sprintf("insight-%d", i),
			},
			Description: "for pagination test",
		}
		if _, err := repo.CreateFavourite(ctx, userID, asset); err != nil {
			t.Fatalf("CreateFavourite #%d error: %v", i, err)
		}
	}

	// Request page with limit=2, offset=1
	filter := favourites.ListFilter{
		Limit:  2,
		Offset: 1,
	}

	page, total, err := repo.ListFavourites(ctx, userID, filter)
	if err != nil {
		t.Fatalf("ListFavourites returned error: %v", err)
	}

	if total != 5 {
		t.Fatalf("expected total=5, got %d", total)
	}
	if len(page) != 2 {
		t.Fatalf("expected page length 2, got %d", len(page))
	}

	// Sanity: ordering by CreatedAt ascending means we should see items 1 and 2
	// (0-based index) in this page. We only check that their CreatedAt is within
	// the range (i.e. page is not empty and sorted) rather than exact IDs/text.
	if !page[0].Favourite.CreatedAt.Before(page[1].Favourite.CreatedAt) &&
		!page[0].Favourite.CreatedAt.Equal(page[1].Favourite.CreatedAt) {
		t.Fatalf("expected page to be ordered by CreatedAt ascending")
	}
}

func TestListFavourites_ContextCancelled(t *testing.T) {
	t.Parallel()

	repo := NewFavouritesRepository()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, _, err := repo.ListFavourites(ctx, "user-1", favourites.ListFilter{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRemoveFavourite_IsIdempotent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := NewFavouritesRepository()

	asset := favourites.Asset{
		Type: favourites.AssetTypeChart,
		Chart: &favourites.Chart{
			Title: "To be removed",
		},
	}
	entry, err := repo.CreateFavourite(ctx, "user-1", asset)
	if err != nil {
		t.Fatalf("CreateFavourite returned error: %v", err)
	}

	// First removal
	if err := repo.RemoveFavourite(ctx, "user-1", entry.Asset.ID); err != nil {
		t.Fatalf("RemoveFavourite returned error: %v", err)
	}

	// Second removal should be a no-op, not an error
	if err := repo.RemoveFavourite(ctx, "user-1", entry.Asset.ID); err != nil {
		t.Fatalf("expected no error on second RemoveFavourite, got %v", err)
	}
}

func TestUpdateDescriptionForUser_VariousScenarios(t *testing.T) {
	t.Parallel()

	type setupFn func(t *testing.T, repo favourites.Repository) (userID, assetID string)

	ctx := context.Background()

	tests := []struct {
		name        string
		setup       setupFn
		newDesc     string
		wantErr     error
		wantDescSet bool
	}{
		{
			name: "succeeds_when_favourite_exists",
			setup: func(t *testing.T, repo favourites.Repository) (string, string) {
				asset := favourites.Asset{
					Type: favourites.AssetTypeInsight,
					Insight: &favourites.Insight{
						Text: "Old",
					},
					Description: "Old description",
				}
				entry, err := repo.CreateFavourite(ctx, "user-1", asset)
				if err != nil {
					t.Fatalf("setup CreateFavourite: %v", err)
				}
				return "user-1", entry.Asset.ID
			},
			newDesc:     "New description",
			wantErr:     nil,
			wantDescSet: true,
		},
		{
			name: "returns_err_favourite_not_found_when_missing",
			setup: func(t *testing.T, repo favourites.Repository) (string, string) {
				return "user-1", "non-existent"
			},
			newDesc:     "desc",
			wantErr:     favourites.ErrFavouriteNotFound,
			wantDescSet: false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewFavouritesRepository()

			userID, assetID := tt.setup(t, repo)

			updated, err := repo.UpdateDescriptionForUser(ctx, userID, assetID, tt.newDesc)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("UpdateDescriptionForUser error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				return
			}

			if updated.Asset.Description != tt.newDesc {
				t.Fatalf("expected description %q, got %q", tt.newDesc, updated.Asset.Description)
			}

			// Ensure subsequent listing shows updated description.
			entries, total, err := repo.ListFavourites(ctx, userID, favourites.ListFilter{})
			if err != nil {
				t.Fatalf("ListFavourites returned error: %v", err)
			}
			if got, want := len(entries), 1; got != want {
				t.Fatalf("expected %d entry, got %d", want, got)
			}
			if total != 1 {
				t.Fatalf("expected total=1, got %d", total)
			}
			if entries[0].Asset.Description != tt.newDesc {
				t.Fatalf("expected description %q in listed asset, got %q",
					tt.newDesc, entries[0].Asset.Description)
			}
		})
	}
}

func TestFavouritesRepository_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	repo := NewFavouritesRepository()
	ctx := context.Background()

	const n = 50
	userID := "user-1"

	// 1) Concurrently create N favourites for the same user.
	var wg sync.WaitGroup
	idsCh := make(chan string, n)

	makeAsset := func(i int) favourites.Asset {
		return favourites.Asset{
			Type: favourites.AssetTypeInsight,
			Insight: &favourites.Insight{
				Text: fmt.Sprintf("insight-%d", i),
			},
			Description: "concurrent test asset",
		}
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			asset := makeAsset(i)
			entry, err := repo.CreateFavourite(ctx, userID, asset)
			if err != nil {
				t.Errorf("CreateFavourite error: %v", err)
				return
			}

			idsCh <- entry.Asset.ID
		}(i)
	}

	wg.Wait()
	close(idsCh)

	ids := make([]string, 0, n)
	for id := range idsCh {
		ids = append(ids, id)
	}

	if len(ids) != n {
		t.Fatalf("expected %d created favourites, got %d", n, len(ids))
	}

	// 2) Concurrently remove all favourites, while also running some readers.
	wg = sync.WaitGroup{}

	// Concurrent removers: one goroutine per favourite ID.
	for _, assetID := range ids {
		wg.Add(1)
		go func(assetID string) {
			defer wg.Done()

			if err := repo.RemoveFavourite(ctx, userID, assetID); err != nil {
				t.Errorf("RemoveFavourite(%s) error: %v", assetID, err)
			}
		}(assetID)
	}

	// Concurrent readers: list favourites while removals are happening.
	const readers = 10
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if _, _, err := repo.ListFavourites(ctx, userID, favourites.ListFilter{}); err != nil {
				t.Errorf("ListFavourites error: %v", err)
			}
		}()
	}

	wg.Wait()

	// 3) After all concurrent operations, user-1 should have no favourites.
	entries, total, err := repo.ListFavourites(ctx, userID, favourites.ListFilter{})
	if err != nil {
		t.Fatalf("final ListFavourites returned error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 favourites after concurrent removals, got %d", len(entries))
	}
	if total != 0 {
		t.Fatalf("expected total=0 after concurrent removals, got %d", total)
	}
}
