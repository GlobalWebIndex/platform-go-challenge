package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manos/favourites/favourites/application"
	"github.com/manos/favourites/favourites/domain"
	"github.com/manos/favourites/http/auth"
)

func TestFavouritesHandler_ListWithPagination(t *testing.T) {
	repo := &favouriteTestRepo{
		data: map[string]domain.FavouriteEntity{
			"fav-1": {ID: "fav-1", UserID: "user-1", AssetID: "ins-1", Type: domain.FavouriteInsight},
		},
	}
	assets := &assetClientTestStub{
		assets: map[string]application.AssetDTO{"ins-1": {ID: "ins-1", Name: "n", Type: domain.FavouriteInsight}},
	}
	svc := application.NewFavouriteService(repo, assets)
	h := NewFavouritesHandler(svc, 1, 20, 100)

	req := httptest.NewRequest(http.MethodGet, "/v1/favourites?page=2&limit=1", nil)
	ctx := context.WithValue(req.Context(), auth.ContextKeySubject, "user-1")
	rec := httptest.NewRecorder()

	params := GetV1FavouritesParams{
		Page:  intPtr(2),
		Limit: intPtr(1),
	}

	h.GetV1Favourites(rec, req.WithContext(ctx), params)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp FavouritePage
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Page == nil || *resp.Page != 2 || resp.Limit == nil || *resp.Limit != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestFavouritesHandler_UpdateDescription(t *testing.T) {
	repo := &favouriteTestRepo{
		data: map[string]domain.FavouriteEntity{
			"fav-1": {ID: "fav-1", UserID: "user-1", AssetID: "ins-1", Type: domain.FavouriteInsight},
		},
	}
	assets := &assetClientTestStub{}
	svc := application.NewFavouriteService(repo, assets)
	h := NewFavouritesHandler(svc, 1, 20, 100)

	body := bytes.NewBufferString(`{"description":"new desc"}`)
	req := httptest.NewRequest(http.MethodPatch, "/v1/favourites/fav-1", body)
	ctx := context.WithValue(req.Context(), auth.ContextKeySubject, "user-1")
	rec := httptest.NewRecorder()

	h.PatchV1FavouritesId(rec, req.WithContext(ctx), "fav-1")

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if !repo.updateCalled {
		t.Fatalf("expected update description to be called")
	}
}

func TestFavouritesHandler_RejectsLimitOverMax(t *testing.T) {
	repo := &favouriteTestRepo{
		data: map[string]domain.FavouriteEntity{
			"fav-1": {ID: "fav-1", UserID: "user-1", AssetID: "ins-1", Type: domain.FavouriteInsight},
		},
	}
	assets := &assetClientTestStub{
		assets: map[string]application.AssetDTO{"ins-1": {ID: "ins-1", Name: "n", Type: domain.FavouriteInsight}},
	}
	svc := application.NewFavouriteService(repo, assets)
	h := NewFavouritesHandler(svc, 1, 1, 2)

	req := httptest.NewRequest(http.MethodGet, "/v1/favourites?limit=10", nil)
	ctx := context.WithValue(req.Context(), auth.ContextKeySubject, "user-1")
	rec := httptest.NewRecorder()

	params := GetV1FavouritesParams{
		Limit: intPtr(10),
	}

	h.GetV1Favourites(rec, req.WithContext(ctx), params)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

type favouriteTestRepo struct {
	data         map[string]domain.FavouriteEntity
	updateCalled bool
}

func (r *favouriteTestRepo) FindByUser(_ context.Context, userID string) ([]domain.FavouriteEntity, error) {
	out := make([]domain.FavouriteEntity, 0)
	for _, f := range r.data {
		if f.UserID == userID {
			out = append(out, f)
		}
	}
	return out, nil
}
func (r *favouriteTestRepo) FindByUserPage(_ context.Context, userID string, _, _ int) ([]domain.FavouriteEntity, int, error) {
	list, _ := r.FindByUser(context.Background(), userID)
	return list, len(list), nil
}
func (r *favouriteTestRepo) Save(_ context.Context, fav domain.FavouriteEntity) error {
	r.data[fav.ID] = fav
	return nil
}
func (r *favouriteTestRepo) FindByID(_ context.Context, id string) (domain.FavouriteEntity, error) {
	if fav, ok := r.data[id]; ok {
		return fav, nil
	}
	return domain.FavouriteEntity{}, application.ErrFavouriteNotFound
}
func (r *favouriteTestRepo) Delete(_ context.Context, id string) error {
	delete(r.data, id)
	return nil
}
func (r *favouriteTestRepo) UpdateDescription(_ context.Context, id string, desc string) error {
	if fav, ok := r.data[id]; ok {
		fav.Description = desc
		r.data[id] = fav
		r.updateCalled = true
		return nil
	}
	return application.ErrFavouriteNotFound
}

type assetClientTestStub struct {
	assets map[string]application.AssetDTO
}

func (c *assetClientTestStub) GetInsight(_ context.Context, _ string, _ string) (application.AssetDTO, error) {
	return application.AssetDTO{}, nil
}
func (c *assetClientTestStub) GetAudience(_ context.Context, _ string, _ string) (application.AssetDTO, error) {
	return application.AssetDTO{}, nil
}
func (c *assetClientTestStub) GetChart(_ context.Context, _ string, _ string) (application.AssetDTO, error) {
	return application.AssetDTO{}, nil
}
func (c *assetClientTestStub) GetInsights(_ context.Context, _ string, _ []string) (map[string]application.AssetDTO, error) {
	return c.assets, nil
}
func (c *assetClientTestStub) GetAudiences(_ context.Context, _ string, _ []string) (map[string]application.AssetDTO, error) {
	return c.assets, nil
}
func (c *assetClientTestStub) GetCharts(_ context.Context, _ string, _ []string) (map[string]application.AssetDTO, error) {
	return c.assets, nil
}
