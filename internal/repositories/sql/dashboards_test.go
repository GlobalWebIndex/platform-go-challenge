package sql

import (
	"context"
	"testing"

	"platform-go-challenge/internal/app/assets"
	"platform-go-challenge/internal/app/dashboards"
	"platform-go-challenge/test"

	"github.com/stretchr/testify/assert"
)

var (
	chartsUp = `INSERT INTO challenge.charts
	(id,	title, 		x_axis, 	y_axis, 	data) VALUES
	(1,		"chart_1",	"x axis",	"y axis", 	'{"data":1}'),
	(2,		"chart_2",	"x axis",	"y axis", 	'{"data":2}'),
	(3,		"chart_3",	"x axis",	"y axis", 	'{"data":3}');`
	chartsDown = "DELETE FROM  challenge.charts;"

	audiencesUp = `INSERT INTO challenge.audiences
	(id, gender,   country_of_birth, 	age_group, 			hours_spent_online, number_of_purchases_last_month) VALUES
	(1,	 'male',   'gr', 				'young-adults', 	10, 				5),
	(2,	 'female', 'gr', 				'young-adults', 	10, 				5),
	(3,	 'male',   'de', 				'teenagers', 		25, 				2);
	`
	audiencesDown = "DELETE FROM  challenge.audiences;"

	insightsUp = `INSERT INTO challenge.insights
	(id, title) VALUES
	(1,  '40% of millenials spend more than 3hours on social media daily'),
	(2,  '60% of teenagers spend more than 6hours on social media daily'),
	(3,  '10% of seniors spend less than 3hours on social media weekly');
	`
	insightsDown = `DELETE FROM  challenge.insights;`

	dashboardsUp = `INSERT INTO challenge.dashboards
	(id,user_id) VALUES
	(1,	1),
	(2,	2),
	(3,	3);`
	dashboardsDown = "DELETE FROM  challenge.dashboards;"

	d2aUp = `INSERT INTO challenge.dashboards2assets
	(dashboard_id, 	asset_id, 	asset_type, description) VALUES
	(1, 			1, 			'chart', 		'chart description 1'),
	(1, 			2, 			'chart', 		'chart description 2'),
	(1, 			1, 			'audience', 	'audience description 1'),
	(1, 			1, 			'insight', 		'insight description 1'),
	(2, 			1, 			'chart', 		'chart description 1 alt');
	`
	d2aDown = "DELETE FROM  challenge.dashboards2assets;"
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
	got, err := repo.GetUserDashboard(context.Background(), 1)
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
	got, err := repo.GetUserDashboard(context.Background(), 1)
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
	got, err = repo.GetUserDashboard(context.Background(), 1)
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
	got, err := repo.GetUserDashboard(context.Background(), 1)
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

	got, err = repo.GetUserDashboard(context.Background(), 1)
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
	got, err := repo.GetUserDashboard(context.Background(), 1)
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

	got, err = repo.GetUserDashboard(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, action.Description, got.Audiences[0].Description)
}
