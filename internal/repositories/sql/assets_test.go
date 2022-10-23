package sql

import (
	"context"
	"testing"

	"platform-go-challenge/test"

	"github.com/stretchr/testify/assert"
)

func TestAssets_List(t *testing.T) {
	db, down := test.DBInitQueriesWithTest(
		t,
		[]string{chartsUp, audiencesUp, insightsUp},
		[]string{chartsDown, audiencesDown, insightsDown},
	)
	defer down()

	repo := NewAssetsRepository(db)
	assetsList, err := repo.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, assetsList.Charts, 3)
	assert.Len(t, assetsList.Insights, 3)
	assert.Len(t, assetsList.Audiences, 3)
}

func TestAssets_ListCharts(t *testing.T) {
	db, down := test.DBInitQueriesWithTest(
		t,
		[]string{chartsUp, audiencesUp, insightsUp},
		[]string{chartsDown, audiencesDown, insightsDown},
	)
	defer down()

	repo := NewAssetsRepository(db)
	charts, err := repo.ListCharts(context.Background())
	assert.NoError(t, err)
	assert.Len(t, charts, 3)
}

func TestAssets_ListInsights(t *testing.T) {
	db, down := test.DBInitQueriesWithTest(
		t,
		[]string{chartsUp, audiencesUp, insightsUp},
		[]string{chartsDown, audiencesDown, insightsDown},
	)
	defer down()

	repo := NewAssetsRepository(db)
	insights, err := repo.ListInsights(context.Background())
	assert.NoError(t, err)
	assert.Len(t, insights, 3)
}

func TestAssets_ListAudiences(t *testing.T) {
	db, down := test.DBInitQueriesWithTest(
		t,
		[]string{chartsUp, audiencesUp, insightsUp},
		[]string{chartsDown, audiencesDown, insightsDown},
	)
	defer down()

	repo := NewAssetsRepository(db)
	audiences, err := repo.ListAudiences(context.Background())
	assert.NoError(t, err)
	assert.Len(t, audiences, 3)
}
