package sql

import (
	"context"
	"database/sql"

	"platform-go-challenge/internal/app/assets"
	"platform-go-challenge/internal/app/dashboards"

	sq "github.com/Masterminds/squirrel"
)

// Dashboards sql repository struct.
type Dashboards struct {
	dbClient BasicConnectionWithTransactions
}

// NewDashboardsRepository constructor.
func NewDashboardsRepository(dbClient BasicConnectionWithTransactions) *Dashboards {
	return &Dashboards{dbClient: dbClient}
}

// GetUserDashboardID returns id of a user's dashboard.
func (d *Dashboards) GetUserDashboardID(ctx context.Context, userID uint32) (uint32, error) {
	dashboardID := uint32(0)
	err := sq.Select("d.id").
		From(dashboardsTable + " AS d").
		Where(sq.Eq{"d.user_id": userID}).
		RunWith(d.dbClient).QueryRowContext(ctx).Scan(&dashboardID)
	if err != nil {
		return 0, err
	}

	return dashboardID, nil
}

// GetUserDashboard returns a user's dashboard.
func (d *Dashboards) GetUserDashboard(ctx context.Context, userID uint32) (dashboards.Dashboard, error) {
	dashboard := dashboards.Dashboard{}

	rows, err := sq.Select("d.id,d2a.asset_id,d2a.asset_type").
		From(dashboardsTable + " AS d").
		Join(dashboards2AssetsTable + " AS d2a ON d2a.dashboard_id = d.id").
		Where(sq.Eq{"d.user_id": userID}).
		RunWith(d.dbClient).QueryContext(ctx)
	if err != nil {
		return dashboards.Dashboard{}, err
	}
	defer rows.Close()
	dashboardID, assetsMap, err := scanAssetIDs(rows)
	if err != nil {
		return dashboards.Dashboard{}, err
	}

	starredCharts, err := d.getStarredCharts(ctx, dashboardID, assetsMap[assets.AssetTypeChart])
	if err != nil {
		return dashboards.Dashboard{}, err
	}
	starredInsights, err := d.getStarredInsights(ctx, dashboardID, assetsMap[assets.AssetTypeInsight])
	if err != nil {
		return dashboards.Dashboard{}, err
	}
	starredAudiences, err := d.getStarredAudiences(ctx, dashboardID, assetsMap[assets.AssetTypeAudience])
	if err != nil {
		return dashboards.Dashboard{}, err
	}
	dashboard.ID = dashboardID
	dashboard.UserID = userID
	dashboard.Charts = starredCharts
	dashboard.Insights = starredInsights
	dashboard.Audiences = starredAudiences

	return dashboard, nil
}

func (d *Dashboards) getStarredCharts(ctx context.Context, dashboardID uint32, chartIDs []uint32) (dashboards.StarredCharts, error) {
	rows, err := sq.Select("ch.id", "ch.title", "ch.x_axis", "ch.y_axis", "ch.data", "d2a.description").
		From(chartsTable + " AS ch").
		Join(dashboards2AssetsTable + " AS d2a ON d2a.asset_id = ch.id AND d2a.asset_type = '" + string(assets.AssetTypeChart) + "'").
		Where(sq.Eq{"ch.id": chartIDs, "d2a.dashboard_id": dashboardID}).
		RunWith(d.dbClient).QueryContext(ctx)
	if err != nil {
		return dashboards.StarredCharts{}, err
	}
	defer rows.Close()
	return scanStarredCharts(rows)
}

func scanStarredCharts(rows *sql.Rows) (dashboards.StarredCharts, error) {
	var err error
	starredCharts := make(dashboards.StarredCharts, 0, 5)
	for rows.Next() {
		starredChart := dashboards.StarredChart{}
		err = rows.Scan(
			&starredChart.ID,
			&starredChart.Title,
			&starredChart.AxisXTitle,
			&starredChart.AxisYTitle,
			&starredChart.Data,
			&starredChart.Description)
		if err != nil {
			return dashboards.StarredCharts{}, err
		}
		starredCharts = append(starredCharts, starredChart)
	}
	return starredCharts, nil
}

func (d *Dashboards) getStarredInsights(ctx context.Context, dashboardID uint32, insightIDs []uint32) (dashboards.StarredInsights, error) {
	rows, err := sq.Select("ins.id", "ins.title", "d2a.description").
		From(insightsTable + " AS ins").
		Join(dashboards2AssetsTable + " AS d2a ON d2a.asset_id = ins.id AND d2a.asset_type = '" + string(assets.AssetTypeInsight) + "'").
		Where(sq.Eq{"ins.id": insightIDs, "d2a.dashboard_id": dashboardID}).
		RunWith(d.dbClient).QueryContext(ctx)
	if err != nil {
		return dashboards.StarredInsights{}, err
	}
	defer rows.Close()
	return scanStarredInsights(rows)
}

func scanStarredInsights(rows *sql.Rows) (dashboards.StarredInsights, error) {
	var err error
	starredInsights := make(dashboards.StarredInsights, 0, 5)
	for rows.Next() {
		starredInsight := dashboards.StarredInsight{}
		err = rows.Scan(
			&starredInsight.ID,
			&starredInsight.Title,
			&starredInsight.Description)
		if err != nil {
			return dashboards.StarredInsights{}, err
		}
		starredInsights = append(starredInsights, starredInsight)
	}
	return starredInsights, nil
}

func (d *Dashboards) getStarredAudiences(ctx context.Context, dashboardID uint32, audienceIDs []uint32) (dashboards.StarredAudiences, error) {
	rows, err := sq.Select(
		"au.id",
		"au.gender",
		"au.country_of_birth",
		"au.age_group",
		"au.hours_spent_online",
		"au.number_of_purchases_last_month",
		"d2a.description").
		From(audiencesTable + " AS au").
		Join(dashboards2AssetsTable + " AS d2a ON d2a.asset_id = au.id AND d2a.asset_type = '" + string(assets.AssetTypeAudience) + "'").
		Where(sq.Eq{"au.id": audienceIDs, "d2a.dashboard_id": dashboardID}).
		RunWith(d.dbClient).QueryContext(ctx)
	if err != nil {
		return dashboards.StarredAudiences{}, err
	}
	defer rows.Close()
	return scanStarredAudiences(rows)
}

func scanStarredAudiences(rows *sql.Rows) (dashboards.StarredAudiences, error) {
	var err error
	starredAudiences := make(dashboards.StarredAudiences, 0, 5)
	for rows.Next() {
		starredAudience := dashboards.StarredAudience{}
		err = rows.Scan(
			&starredAudience.ID,
			&starredAudience.Gender,
			&starredAudience.BirthCountry,
			&starredAudience.AgeGroup,
			&starredAudience.HoursSpentOnline,
			&starredAudience.NumberOfPerchasesLastMonth,
			&starredAudience.Description)
		if err != nil {
			return dashboards.StarredAudiences{}, err
		}
		starredAudiences = append(starredAudiences, starredAudience)
	}
	return starredAudiences, nil
}

func scanAssetIDs(rows *sql.Rows) (uint32, map[assets.AssetType][]uint32, error) {
	assetsMap := make(map[assets.AssetType][]uint32, 3)
	chartsSlice := make([]uint32, 0, 5)
	audiencesSlice := make([]uint32, 0, 5)
	insightsSlice := make([]uint32, 0, 5)

	dashboardID := uint32(0)
	assetID := uint32(0)
	assetType := assets.AssetType("")
	var err error
	for rows.Next() {
		err = rows.Scan(&dashboardID, &assetID, &assetType)
		if err != nil {
			return 0, assetsMap, err
		}
		switch assetType {
		case assets.AssetTypeChart:
			chartsSlice = append(chartsSlice, assetID)
		case assets.AssetTypeAudience:
			audiencesSlice = append(audiencesSlice, assetID)
		case assets.AssetTypeInsight:
			insightsSlice = append(insightsSlice, assetID)
		}
	}
	assetsMap[assets.AssetTypeAudience] = audiencesSlice
	assetsMap[assets.AssetTypeChart] = chartsSlice
	assetsMap[assets.AssetTypeInsight] = insightsSlice

	return dashboardID, assetsMap, nil
}

// AddToDashboard adds a new asset to a user's dashboard.
func (d *Dashboards) AddToDashboard(ctx context.Context, actionParams dashboards.DashboardActionParams) error {
	setMap := map[string]interface{}{
		"dashboard_id": actionParams.DashboardID,
		"asset_id":     actionParams.AssetID,
		"asset_type":   actionParams.AssetType,
		"description":  actionParams.Description,
	}
	_, err := sq.Insert(dashboards2AssetsTable).SetMap(setMap).RunWith(d.dbClient).ExecContext(ctx)
	return err
}

// RemoveFromDashboard removes an asset from a user's dashboard.
func (d *Dashboards) RemoveFromDashboard(ctx context.Context, actionParams dashboards.DashboardActionParams) error {
	_, err := sq.Delete(dashboards2AssetsTable).Where(sq.Eq{
		"dashboard_id": actionParams.DashboardID,
		"asset_id":     actionParams.AssetID,
		"asset_type":   actionParams.AssetType,
	}).RunWith(d.dbClient).ExecContext(ctx)

	return err
}

// EditDescription edits the description of an asset in dashboard.
func (d *Dashboards) EditDescription(ctx context.Context, actionParams dashboards.DashboardActionParams) error {
	_, err := sq.Update(dashboards2AssetsTable).
		SetMap(map[string]interface{}{"description": actionParams.Description}).
		Where(sq.Eq{
			"dashboard_id": actionParams.DashboardID,
			"asset_id":     actionParams.AssetID,
			"asset_type":   actionParams.AssetType,
		}).RunWith(d.dbClient).ExecContext(ctx)

	return err
}
