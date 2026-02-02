package repo_inmemory

import (
	"context"
	"testing"
)

func TestAssetRepositoryOwnershipAndBatch(t *testing.T) {
	repo := NewInMemoryAssetRepository()
	ctx := context.Background()

	if _, ok := repo.GetInsight(ctx, "user-1", "ins-001"); !ok {
		t.Fatalf("expected insight for owner")
	}
	if _, ok := repo.GetInsight(ctx, "user-2", "ins-1"); ok {
		t.Fatalf("expected forbidden insight for other user")
	}

	batch := repo.GetInsights(ctx, "user-1", []string{"ins-001", "ins-999"})
	if len(batch) != 1 {
		t.Fatalf("batch insights len=%d, want 1", len(batch))
	}

	if _, ok := repo.GetAudience(ctx, "user-1", "aud-001"); !ok {
		t.Fatalf("expected audience for owner")
	}
	if _, ok := repo.GetAudience(ctx, "user-2", "aud-001"); ok {
		t.Fatalf("expected forbidden audience for other user")
	}

	audBatch := repo.GetAudiences(ctx, "user-1", []string{"aud-001", "aud-999"})
	if len(audBatch) != 1 {
		t.Fatalf("batch audiences len=%d, want 1", len(audBatch))
	}

	if _, ok := repo.GetChart(ctx, "user-1", "chr-001"); !ok {
		t.Fatalf("expected chart for owner")
	}
	if _, ok := repo.GetChart(ctx, "user-2", "chr-001"); ok {
		t.Fatalf("expected forbidden chart for other user")
	}

	chartBatch := repo.GetCharts(ctx, "user-1", []string{"chr-001", "chr-999"})
	if len(chartBatch) != 1 {
		t.Fatalf("batch charts len=%d, want 1", len(chartBatch))
	}
}
