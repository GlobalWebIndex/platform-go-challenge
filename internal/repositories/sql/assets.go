package sql

import (
	"context"
	"database/sql"
	"encoding/json"

	"platform-go-challenge/internal/app/assets"

	sq "github.com/Masterminds/squirrel"
)

// Assets sql repository struct.
type Assets struct {
	dbClient BasicConnectionWithTransactions
}

// NewAssetsRepository constructor.
func NewAssetsRepository(dbClient BasicConnectionWithTransactions) *Assets {
	return &Assets{dbClient: dbClient}
}

func (d *Assets) List(ctx context.Context) (assets.Assets, error) {
	assetList := assets.Assets{}

	charts, err := d.ListCharts(ctx)
	if err != nil {
		return assets.Assets{}, err
	}

	insights, err := d.ListInsights(ctx)
	if err != nil {
		return assets.Assets{}, err
	}

	audiences, err := d.ListAudiences(ctx)
	if err != nil {
		return assets.Assets{}, err
	}

	assetList.Charts = charts
	assetList.Insights = insights
	assetList.Audiences = audiences

	return assetList, nil
}

func (d *Assets) ListCharts(ctx context.Context) (assets.Charts, error) {
	rows, err := sq.Select("ch.id", "ch.title", "ch.x_axis", "ch.y_axis", "ch.data").
		From(chartsTable + " AS ch").
		RunWith(d.dbClient).QueryContext(ctx)
	if err != nil {
		return assets.Charts{}, err
	}
	defer rows.Close()
	return scanCharts(rows)
}

func scanCharts(rows *sql.Rows) (assets.Charts, error) {
	var err error
	charts := make([]assets.Chart, 0, 5)
	for rows.Next() {
		data := []byte{}
		chart := assets.Chart{}
		err = rows.Scan(
			&chart.ID,
			&chart.Title,
			&chart.AxisXTitle,
			&chart.AxisYTitle,
			&data)
		if err != nil {
			return assets.Charts{}, err
		}
		err = json.Unmarshal(data, &chart.Data)
		if err != nil {
			return assets.Charts{}, err
		}
		charts = append(charts, chart)
	}
	return charts, nil
}

func (d *Assets) ListInsights(ctx context.Context) (assets.Insights, error) {
	rows, err := sq.Select("ins.id", "ins.title").
		From(insightsTable + " AS ins").
		RunWith(d.dbClient).QueryContext(ctx)
	if err != nil {
		return assets.Insights{}, err
	}
	defer rows.Close()
	return scanInsights(rows)
}

func scanInsights(rows *sql.Rows) (assets.Insights, error) {
	var err error
	insights := make(assets.Insights, 0, 5)
	for rows.Next() {
		insight := assets.Insight{}
		err = rows.Scan(
			&insight.ID,
			&insight.Title)
		if err != nil {
			return assets.Insights{}, err
		}
		insights = append(insights, insight)
	}
	return insights, nil
}

func (d *Assets) ListAudiences(ctx context.Context) (assets.Audiences, error) {
	rows, err := sq.Select(
		"au.id",
		"au.gender",
		"au.country_of_birth",
		"au.age_group",
		"au.hours_spent_online",
		"au.number_of_purchases_last_month").
		From(audiencesTable + " AS au").
		RunWith(d.dbClient).QueryContext(ctx)
	if err != nil {
		return assets.Audiences{}, err
	}
	defer rows.Close()
	return scanAudiences(rows)
}

func scanAudiences(rows *sql.Rows) (assets.Audiences, error) {
	var err error
	audiences := make(assets.Audiences, 0, 5)
	for rows.Next() {
		audience := assets.Audience{}
		err = rows.Scan(
			&audience.ID,
			&audience.Gender,
			&audience.BirthCountry,
			&audience.AgeGroup,
			&audience.HoursSpentOnline,
			&audience.NumberOfPerchasesLastMonth)
		if err != nil {
			return assets.Audiences{}, err
		}
		audiences = append(audiences, audience)
	}
	return audiences, nil
}
