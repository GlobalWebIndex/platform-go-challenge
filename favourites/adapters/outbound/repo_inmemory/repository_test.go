package repo_inmemory

import (
	"context"
	"testing"

	"github.com/manos/favourites/favourites/domain"
)

func TestFindByUserPage(t *testing.T) {
	repo := NewInMemoryFavouriteRepository()

	page, total, err := repo.FindByUserPage(context.Background(), "user-1", 0, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Fatalf("total = %d, want 3", total)
	}
	if len(page) != 2 {
		t.Fatalf("len(page) = %d, want 2", len(page))
	}
	if page[0].ID != "fav-1" || page[1].ID != "fav-2" {
		t.Fatalf("unexpected order: %#v", page)
	}
}

func TestSaveFindDeleteUpdate(t *testing.T) {
	repo := NewInMemoryFavouriteRepository()
	ctx := context.Background()

	newFav := domain.FavouriteEntity{ID: "fav-new", UserID: "user-3", AssetID: "ins-9", Type: domain.FavouriteInsight, Description: "desc"}
	if err := repo.Save(ctx, newFav); err != nil {
		t.Fatalf("save error: %v", err)
	}

	// duplicate save should fail
	if err := repo.Save(ctx, domain.FavouriteEntity{ID: "other-id", UserID: "user-3", AssetID: "ins-9", Type: domain.FavouriteInsight}); err == nil {
		t.Fatalf("expected duplicate save error")
	}

	got, err := repo.FindByID(ctx, "fav-new")
	if err != nil {
		t.Fatalf("find error: %v", err)
	}
	if got.ID != "fav-new" {
		t.Fatalf("unexpected favourite: %#v", got)
	}

	if err := repo.UpdateDescription(ctx, "fav-new", "updated"); err != nil {
		t.Fatalf("update error: %v", err)
	}
	updated, _ := repo.FindByID(ctx, "fav-new")
	if updated.Description != "updated" {
		t.Fatalf("description not updated: %#v", updated)
	}

	if err := repo.Delete(ctx, "fav-new"); err != nil {
		t.Fatalf("delete error: %v", err)
	}
	if _, err := repo.FindByID(ctx, "fav-new"); err == nil {
		t.Fatalf("expected error after delete")
	}
}
