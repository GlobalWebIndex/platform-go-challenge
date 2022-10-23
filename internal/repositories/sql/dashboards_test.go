package sql

import (
	"context"
	"testing"

	"platform-go-challenge/internal/app/assets"
	"platform-go-challenge/internal/app/dashboards"
	"platform-go-challenge/internal/pagination"
	"platform-go-challenge/test"

	"github.com/stretchr/testify/assert"
)

func TestGetUserDashboardIDSuccess(t *testing.T) {
	db, down := test.DBInitQueriesWithTest(
		t,
		[]string{dashboardsUp},
		[]string{dashboardsDown},
	)
	defer down()

	repo := NewDashboardsRepository(db)
	got, err := repo.GetUserDashboardID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), got)

	got, err = repo.GetUserDashboardID(context.Background(), 11)
	assert.Error(t, err)
	assert.Equal(t, uint32(0), got)
}

func TestGetUserDashboardSuccess(t *testing.T) {
	db, down := test.DBInitQueriesWithTest(
		t,
		[]string{chartsUp, audiencesUp, insightsUp, dashboardsUp, d2aUp},
		[]string{chartsDown, dashboardsDown, audiencesDown, insightsDown, d2aDown},
	)
	defer down()
	expectedStarredCharts := dashboards.StarredCharts{
		dashboards.StarredChart{
			Chart: assets.Chart{
				ID:         1,
				Title:      "chart_1",
				AxisXTitle: "x axis",
				AxisYTitle: "y axis",
				Data:       []byte{0x7b, 0x22, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3a, 0x20, 0x31, 0x7d},
			},
			Description: "chart description 1",
		},
		dashboards.StarredChart{
			Chart: assets.Chart{
				ID:         2,
				Title:      "chart_2",
				AxisXTitle: "x axis",
				AxisYTitle: "y axis",
				Data:       []byte{0x7b, 0x22, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3a, 0x20, 0x32, 0x7d},
			},
			Description: "chart description 2",
		},
	}
	expectedStarredInsights := dashboards.StarredInsights{
		dashboards.StarredInsight{
			Insight: assets.Insight{
				ID:    1,
				Title: "40% of millenials spend more than 3hours on social media daily",
			},
			Description: "insight description 1",
		},
	}

	expectedStarredAudiences := dashboards.StarredAudiences{
		dashboards.StarredAudience{
			Audience: assets.Audience{
				ID:                         1,
				Gender:                     assets.GenderMale,
				BirthCountry:               "gr",
				AgeGroup:                   assets.AgeGroupYoungAdult,
				HoursSpentOnline:           10,
				NumberOfPerchasesLastMonth: 5,
			},
			Description: "audience description 1",
		},
	}
	expected := dashboards.Dashboard{
		ID:        1,
		UserID:    1,
		Charts:    expectedStarredCharts,
		Insights:  expectedStarredInsights,
		Audiences: expectedStarredAudiences,
	}

	repo := NewDashboardsRepository(db)
	got, err := repo.GetUserDashboard(context.Background(), 1, pagination.Pagination{Page: 1, PerPage: 10})
	assert.NoError(t, err)
	assert.Len(t, got.Charts, 2)
	assert.Equal(t, expected, got)
}

func TestAddToDashboardSuccess(t *testing.T) {
	db, down := test.DBInitQueriesWithTest(
		t,
		[]string{chartsUp, audiencesUp, insightsUp, dashboardsUp, d2aUp},
		[]string{chartsDown, dashboardsDown, audiencesDown, insightsDown, d2aDown},
	)
	defer down()
	repo := NewDashboardsRepository(db)
	got, err := repo.GetUserDashboard(context.Background(), 1, pagination.Pagination{Page: 1, PerPage: 10})
	assert.NoError(t, err)
	assert.Len(t, got.Audiences, 1)

	action := dashboards.DashboardActionParams{
		AssetID:     2,
		AssetType:   assets.AssetTypeAudience,
		DashboardID: 1,
		Description: "audience decription for new star",
	}

	err = repo.AddToDashboard(context.Background(), action)
	assert.NoError(t, err)
	got, err = repo.GetUserDashboard(context.Background(), 1, pagination.Pagination{Page: 1, PerPage: 10})
	assert.NoError(t, err)
	assert.Len(t, got.Audiences, 2)
}

func TestRemoveFromDashboardSuccess(t *testing.T) {
	db, down := test.DBInitQueriesWithTest(
		t,
		[]string{chartsUp, audiencesUp, insightsUp, dashboardsUp, d2aUp},
		[]string{chartsDown, dashboardsDown, audiencesDown, insightsDown, d2aDown},
	)
	defer down()
	repo := NewDashboardsRepository(db)
	got, err := repo.GetUserDashboard(context.Background(), 1, pagination.Pagination{Page: 1, PerPage: 10})
	assert.NoError(t, err)
	assert.Len(t, got.Charts, 2)

	action := dashboards.DashboardActionParams{
		AssetID:     1,
		AssetType:   assets.AssetTypeChart,
		DashboardID: 1,
	}

	err = repo.RemoveFromDashboard(context.Background(), action)
	assert.NoError(t, err)

	expectedStarredCharts := dashboards.StarredCharts{
		dashboards.StarredChart{
			Chart: assets.Chart{
				ID:         2,
				Title:      "chart_2",
				AxisXTitle: "x axis",
				AxisYTitle: "y axis",
				Data:       []byte{0x7b, 0x22, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3a, 0x20, 0x32, 0x7d},
			},
			Description: "chart description 2",
		},
	}

	got, err = repo.GetUserDashboard(context.Background(), 1, pagination.Pagination{Page: 1, PerPage: 10})
	assert.NoError(t, err)
	assert.Len(t, got.Charts, 1)
	assert.Equal(t, expectedStarredCharts, got.Charts)
}

func TestEditDescriptionSuccess(t *testing.T) {
	db, down := test.DBInitQueriesWithTest(
		t,
		[]string{chartsUp, audiencesUp, insightsUp, dashboardsUp, d2aUp},
		[]string{chartsDown, dashboardsDown, audiencesDown, insightsDown, d2aDown},
	)
	defer down()
	repo := NewDashboardsRepository(db)
	got, err := repo.GetUserDashboard(context.Background(), 1, pagination.Pagination{Page: 1, PerPage: 10})
	assert.NoError(t, err)
	assert.Equal(t, got.Audiences[0].Description, "audience description 1")

	action := dashboards.DashboardActionParams{
		AssetID:     1,
		AssetType:   assets.AssetTypeAudience,
		DashboardID: 1,
		Description: "audience description 1 updated",
	}

	err = repo.EditDescription(context.Background(), action)
	assert.NoError(t, err)

	got, err = repo.GetUserDashboard(context.Background(), 1, pagination.Pagination{Page: 1, PerPage: 10})
	assert.NoError(t, err)
	assert.Equal(t, action.Description, got.Audiences[0].Description)
}
