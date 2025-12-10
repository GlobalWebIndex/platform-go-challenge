package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/api/v1/dto"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/domain/favourites"
)

// stubService implements favourites.Service for handler tests.
type stubService struct {
	listResult []favourites.FavouriteEntry
	listTotal  int
	listErr    error
	lastList   struct {
		ctx    context.Context
		userID string
		filter favourites.ListFilter
	}

	addResult favourites.FavouriteEntry
	addErr    error
	lastAdd   struct {
		ctx    context.Context
		userID string
		asset  favourites.Asset
	}

	removeErr  error
	lastRemove struct {
		ctx     context.Context
		userID  string
		assetID string
	}

	editResult favourites.FavouriteEntry
	editErr    error
	lastEdit   struct {
		ctx         context.Context
		userID      string
		assetID     string
		description string
	}
}

func (s *stubService) ListUserFavourites(ctx context.Context, userID string, filter favourites.ListFilter) ([]favourites.FavouriteEntry, int, error) {
	s.lastList.ctx = ctx
	s.lastList.userID = userID
	s.lastList.filter = filter
	return s.listResult, s.listTotal, s.listErr
}

func (s *stubService) AddFavourite(ctx context.Context, userID string, asset favourites.Asset) (favourites.FavouriteEntry, error) {
	s.lastAdd.ctx = ctx
	s.lastAdd.userID = userID
	s.lastAdd.asset = asset
	return s.addResult, s.addErr
}

func (s *stubService) RemoveFavourite(ctx context.Context, userID, assetID string) error {
	s.lastRemove.ctx = ctx
	s.lastRemove.userID = userID
	s.lastRemove.assetID = assetID
	return s.removeErr
}

func (s *stubService) EditFavouriteDescription(ctx context.Context, userID, assetID, description string) (favourites.FavouriteEntry, error) {
	s.lastEdit.ctx = ctx
	s.lastEdit.userID = userID
	s.lastEdit.assetID = assetID
	s.lastEdit.description = description
	return s.editResult, s.editErr
}

func TestListUserFavourites_Success(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	stub := &stubService{
		listResult: []favourites.FavouriteEntry{
			{
				Favourite: favourites.Favourite{
					UserID:    "user-1",
					AssetID:   "asset-1",
					CreatedAt: now,
				},
				Asset: favourites.Asset{
					ID:          "asset-1",
					Type:        favourites.AssetTypeInsight,
					Description: "desc",
					Insight: &favourites.Insight{
						Text: "hello",
					},
				},
			},
		},
		listTotal: 1,
	}
	h := NewFavouritesHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/v1/users/user-1/favourites", nil)
	// chi needs the route context with URL params populated.
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userId", "user-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.ListUserFavourites(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp dto.ListFavouritesResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 favourite, got %d", len(resp.Items))
	}
	if resp.Items[0].Asset.ID != "asset-1" {
		t.Fatalf("expected asset ID asset-1, got %s", resp.Items[0].Asset.ID)
	}
	if resp.Items[0].UserID != "user-1" {
		t.Fatalf("expected user ID user-1, got %s", resp.Items[0].UserID)
	}
	if resp.Items[0].AssetID != "asset-1" {
		t.Fatalf("expected asset ID asset-1, got %s", resp.Items[0].AssetID)
	}
	if resp.Meta.Total != 1 {
		t.Fatalf("expected meta total 1, got %d", resp.Meta.Total)
	}
	if resp.Meta.Count != 1 {
		t.Fatalf("expected meta count 1, got %d", resp.Meta.Count)
	}
}

func TestAddFavourite_BadJSON(t *testing.T) {
	t.Parallel()

	h := NewFavouritesHandler(&stubService{})
	body := bytes.NewBufferString(`{ invalid json }`)
	req := httptest.NewRequest(http.MethodPost, "/v1/users/user-1/favourites", body)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userId", "user-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.AddFavourite(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestAddFavourite_Success(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	stub := &stubService{
		addResult: favourites.FavouriteEntry{
			Favourite: favourites.Favourite{
				UserID:    "user-1",
				AssetID:   "asset-1",
				CreatedAt: now,
			},
			Asset: favourites.Asset{
				ID:          "asset-1",
				Type:        favourites.AssetTypeChart,
				Description: "chart-desc",
				Chart: &favourites.Chart{
					Title: "t",
					XAxis: "x",
					YAxis: "y",
					Data:  []float64{1, 2},
				},
			},
		},
	}
	h := NewFavouritesHandler(stub)

	reqBody := dto.CreateFavouriteRequest{
		Asset: dto.CreateAssetRequest{
			Type:        dto.AssetTypeChart,
			Description: "chart-desc",
			Chart: &dto.ChartPayload{
				Title: "t",
				XAxis: "x",
				YAxis: "y",
				Data:  []float64{1, 2},
			},
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		t.Fatalf("encode body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/users/user-1/favourites", &buf)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userId", "user-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.AddFavourite(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp dto.CreateFavouriteResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.AssetID != "asset-1" {
		t.Fatalf("expected asset ID asset-1, got %s", resp.AssetID)
	}
	if resp.UserID != "user-1" {
		t.Fatalf("expected user ID user-1, got %s", resp.UserID)
	}
}

func TestRemoveFavourite_Success(t *testing.T) {
	t.Parallel()

	stub := &stubService{}
	h := NewFavouritesHandler(stub)

	req := httptest.NewRequest(http.MethodDelete, "/v1/users/user-1/favourites/asset-1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userId", "user-1")
	rctx.URLParams.Add("assetId", "asset-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.RemoveFavourite(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rr.Code)
	}
	if stub.lastRemove.userID != "user-1" || stub.lastRemove.assetID != "asset-1" {
		t.Fatalf("RemoveFavourite called with (%s, %s), want (user-1, asset-1)",
			stub.lastRemove.userID, stub.lastRemove.assetID)
	}
}

func TestEditFavouriteDescription_Success(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	stub := &stubService{
		editResult: favourites.FavouriteEntry{
			Favourite: favourites.Favourite{
				UserID:    "user-1",
				AssetID:   "asset-1",
				CreatedAt: now,
			},
			Asset: favourites.Asset{
				ID:          "asset-1",
				Type:        favourites.AssetTypeChart,
				Description: "updated-desc",
				Chart: &favourites.Chart{
					Title: "t",
					XAxis: "x",
					YAxis: "y",
					Data:  []float64{1, 2},
				},
			},
		},
	}
	h := NewFavouritesHandler(stub)

	body := dto.EditFavouriteDescriptionRequest{Description: "updated-desc"}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		t.Fatalf("encode body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/v1/users/user-1/favourites/asset-1", &buf)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userId", "user-1")
	rctx.URLParams.Add("assetId", "asset-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.EditFavouriteDescription(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp dto.FavouriteItemResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.UserID != "user-1" {
		t.Fatalf("expected user ID user-1, got %s", resp.UserID)
	}
	if resp.AssetID != "asset-1" {
		t.Fatalf("expected asset ID asset-1, got %s", resp.AssetID)
	}
	if resp.Asset.Description != "updated-desc" {
		t.Fatalf("expected description updated-desc, got %s", resp.Asset.Description)
	}
}

func TestEditFavouriteDescription_NotFound(t *testing.T) {
	t.Parallel()

	stub := &stubService{
		editErr: favourites.ErrFavouriteNotFound,
	}
	h := NewFavouritesHandler(stub)

	body := dto.EditFavouriteDescriptionRequest{Description: "new"}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		t.Fatalf("encode body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/v1/users/user-1/favourites/asset-1", &buf)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userId", "user-1")
	rctx.URLParams.Add("assetId", "asset-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.EditFavouriteDescription(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}
}

func TestListUserFavourites_PropagatesServiceError(t *testing.T) {
	t.Parallel()

	svcErr := errors.New("boom")
	stub := &stubService{
		listErr: svcErr,
	}
	h := NewFavouritesHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/v1/users/user-1/favourites", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userId", "user-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.ListUserFavourites(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}
}
