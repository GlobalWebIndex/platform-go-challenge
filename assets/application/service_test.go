package application

import (
	"context"
	"testing"

	"github.com/manos/favourites/assets/domain"
)

type assetRepoStub struct {
	insights  map[string]domain.InsightAsset
	audiences map[string]domain.AudienceAsset
	charts    map[string]domain.ChartAsset
}

func (r *assetRepoStub) GetInsight(_ context.Context, userId string, id string) (domain.InsightAsset, bool) {
	a, ok := r.insights[id]
	return a, ok && a.UserId == userId
}
func (r *assetRepoStub) GetAudience(_ context.Context, userId string, id string) (domain.AudienceAsset, bool) {
	a, ok := r.audiences[id]
	return a, ok && a.UserId == userId
}
func (r *assetRepoStub) GetChart(_ context.Context, userId string, id string) (domain.ChartAsset, bool) {
	a, ok := r.charts[id]
	return a, ok && a.UserId == userId
}
func (r *assetRepoStub) GetInsights(_ context.Context, userId string, ids []string) map[string]domain.InsightAsset {
	out := map[string]domain.InsightAsset{}
	for _, id := range ids {
		if a, ok := r.insights[id]; ok && a.UserId == userId {
			out[id] = a
		}
	}
	return out
}
func (r *assetRepoStub) GetAudiences(_ context.Context, userId string, ids []string) map[string]domain.AudienceAsset {
	out := map[string]domain.AudienceAsset{}
	for _, id := range ids {
		if a, ok := r.audiences[id]; ok && a.UserId == userId {
			out[id] = a
		}
	}
	return out
}
func (r *assetRepoStub) GetCharts(_ context.Context, userId string, ids []string) map[string]domain.ChartAsset {
	out := map[string]domain.ChartAsset{}
	for _, id := range ids {
		if a, ok := r.charts[id]; ok && a.UserId == userId {
			out[id] = a
		}
	}
	return out
}

func TestAssetServiceBatchMethods(t *testing.T) {
	repo := &assetRepoStub{
		insights: map[string]domain.InsightAsset{
			"ins-1": {Asset: domain.Asset{Id: "ins-1", UserId: "user-1"}},
		},
		audiences: map[string]domain.AudienceAsset{
			"aud-1": {Asset: domain.Asset{Id: "aud-1", UserId: "user-1"}},
		},
		charts: map[string]domain.ChartAsset{
			"chr-1": {Asset: domain.Asset{Id: "chr-1", UserId: "user-1"}},
		},
	}

	svc := NewAssetService(repo)

	if m, err := svc.GetInsights(context.Background(), "user-1", []string{"ins-1"}); err != nil || len(m) != 1 {
		t.Fatalf("GetInsights unexpected result len=%d err=%v", len(m), err)
	}
	if _, err := svc.GetInsights(context.Background(), "user-1", []string{"ins-missing"}); err == nil {
		t.Fatalf("expected error when missing insights")
	}

	if m, err := svc.GetAudiences(context.Background(), "user-1", []string{"aud-1"}); err != nil || len(m) != 1 {
		t.Fatalf("GetAudiences unexpected result len=%d err=%v", len(m), err)
	}
	if _, err := svc.GetAudiences(context.Background(), "user-1", []string{"aud-missing"}); err == nil {
		t.Fatalf("expected error when missing audiences")
	}

	if m, err := svc.GetCharts(context.Background(), "user-1", []string{"chr-1"}); err != nil || len(m) != 1 {
		t.Fatalf("GetCharts unexpected result len=%d err=%v", len(m), err)
	}
	if _, err := svc.GetCharts(context.Background(), "user-1", []string{"chr-missing"}); err == nil {
		t.Fatalf("expected error when missing charts")
	}
}
