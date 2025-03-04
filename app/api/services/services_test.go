package services

import (
	"app/api/schema"
	"app/models"

	"errors"
	"net/http"
	"testing"
)

const (
	ASSET_CHART_ID = "asset-chart-id"
	UNKNOWN        = "unknown"
)

type mockGService struct{}

func (s *mockGService) FindAssets(filter schema.AssetFilter) []schema.AssetDetailsDto {
	return make([]schema.AssetDetailsDto, 0)
}

func (s *mockGService) FindAsset(assetID string) (schema.AssetDetailsDto, error) {
	if assetID == ASSET_CHART_ID {
		return schema.AssetDetailsDto{
			Asset:        schema.AssetDto{ID: ASSET_CHART_ID, Description: "desc"},
			ChartDetails: &schema.ChartAssetDto{Title: "title", AxesTitles: "axesTitles", PlotData: "plotData"},
		}, nil
	}
	return schema.AssetDetailsDto{}, errors.New("asset not found")
}

func (s *mockGService) UpdateAsset(assetID string, description string) error {
	//TODO return static
	return nil
}

func (s *mockGService) FindUser(userID string) (models.User, error) {
	//TODO return static
	return models.User{}, nil
}

func (s *mockGService) FindUserFavorites(userID string) []models.UserFavorite {
	//TODO return static
	return make([]models.UserFavorite, 0)
}

func (s *mockGService) CreateFavoriteAsset(assetID string, userID string) (string, error) {
	//TODO return static
	return "", nil
}

func (s *mockGService) DeleteUserFavorite(userFavoriteID string) error {
	//TODO return static
	return nil
}

func Test_getAsset_OK(t *testing.T) {
	gService := &mockGService{}
	apiService := NewXService(gService)

	asset, err := apiService.getAsset(ASSET_CHART_ID)
	if err != nil {
		t.Fatalf(`unhandled case with mocked service, err="%v"`, err)
	}
	if asset.Asset.ID != ASSET_CHART_ID {
		t.Fatal("something went wrong when mapping to response struct")
	}
}

func Test_getAsset_NotOK(t *testing.T) {
	gService := &mockGService{}
	apiService := NewXService(gService)

	_, err := apiService.getAsset(UNKNOWN)
	if err == nil {
		t.Fatal("expected error for not existing asset")
	}
	httpError, ok := err.(schema.HttpError)
	if !ok {
		t.Fatalf(`expected HttpError, but got "%s"`, err)
	}
	if httpError.Status() != http.StatusNotFound {
		t.Fatalf(`expected 404, but got %d`, httpError.Status())
	}
}
