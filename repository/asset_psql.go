package repository

import (
	"platform2.0-go-challenge/models"
)

type AssetRepository struct{}

func (a AssetRepository) GetUserAssetsPagination(user_id, limit, offset int) (*models.AssetResponse, error) {
	response, errs := getAllPagination(user_id, limit, offset)

	return response, errs
}

func (a AssetRepository) GetUserAssets(user_id int) (*models.AssetResponse, error) {
	response, errs := getAll(user_id)

	return response, errs
}

func getAll(user_id int) (*models.AssetResponse, error) {
	var response *models.AssetResponse = new(models.AssetResponse)

	chartrepo := ChartRepository{}
	charts, err := chartrepo.GetCharts(user_id)
	if err != nil {
		return nil, err
	}
	response.Charts = charts

	insightsrepo := Insightrepository{}
	insights, err := insightsrepo.GetInsights(user_id)
	if err != nil {
		return nil, err
	}
	response.Insights = insights

	audiencesrepo := AudienceRepository{}
	audiences, err := audiencesrepo.GetAudiences(user_id)
	if err != nil {
		return nil, err
	}
	response.Audiences = audiences

	return response, err
}

func getAllPagination(user_id, limit, offset int) (*models.AssetResponse, error) {
	var response *models.AssetResponse = new(models.AssetResponse)

	chartrepo := ChartRepository{}
	charts, err := chartrepo.GetChartsPagination(user_id, limit, offset)
	if err != nil {
		return nil, err
	}
	response.Charts = charts

	insightsrepo := Insightrepository{}
	insights, err := insightsrepo.GetInsightsPagination(user_id, limit, offset)
	if err != nil {
		return nil, err
	}
	response.Insights = insights

	audiencesrepo := AudienceRepository{}
	audiences, err := audiencesrepo.GetAudiencesPagination(user_id, limit, offset)
	if err != nil {
		return nil, err
	}
	response.Audiences = audiences

	return response, err
}
