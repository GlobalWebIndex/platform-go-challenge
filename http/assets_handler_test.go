package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manos/favourites/assets/application"
	"github.com/manos/favourites/assets/domain"
	"github.com/manos/favourites/http/auth"
)

type assetRepoTestStub struct {
	insights map[string]domain.InsightAsset
}

func (r *assetRepoTestStub) GetInsight(_ context.Context, userId string, id string) (domain.InsightAsset, bool) {
	a, ok := r.insights[id]
	return a, ok && a.UserId == userId
}
func (r *assetRepoTestStub) GetAudience(_ context.Context, _ string, _ string) (domain.AudienceAsset, bool) {
	return domain.AudienceAsset{}, false
}
func (r *assetRepoTestStub) GetChart(_ context.Context, _ string, _ string) (domain.ChartAsset, bool) {
	return domain.ChartAsset{}, false
}
func (r *assetRepoTestStub) GetInsights(_ context.Context, _ string, _ []string) map[string]domain.InsightAsset {
	return nil
}
func (r *assetRepoTestStub) GetAudiences(_ context.Context, _ string, _ []string) map[string]domain.AudienceAsset {
	return nil
}
func (r *assetRepoTestStub) GetCharts(_ context.Context, _ string, _ []string) map[string]domain.ChartAsset {
	return nil
}

func TestAssetsHandler_GetInsight(t *testing.T) {
	repo := &assetRepoTestStub{
		insights: map[string]domain.InsightAsset{
			"ins-1": {Asset: domain.Asset{Id: "ins-1", Name: "n", UserId: "user-1"}},
		},
	}
	svc := application.NewAssetService(repo)
	h := NewAssetsHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/v1/assets/insights/ins-1", nil)
	ctx := context.WithValue(req.Context(), auth.ContextKeySubject, "user-1")
	rec := httptest.NewRecorder()

	h.GetV1AssetsTypeId(rec, req.WithContext(ctx), Insights, "ins-1")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["id"] != "ins-1" {
		t.Fatalf("unexpected response body: %+v", resp)
	}
}
