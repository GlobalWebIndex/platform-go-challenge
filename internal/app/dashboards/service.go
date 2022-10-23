package dashboards

import (
	"context"

	"platform-go-challenge/internal/app/assets"
	"platform-go-challenge/internal/pagination"
)

type Service struct {
	repository Repository
}

func NewDashboardService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetUserDashboard(ctx context.Context, userID uint32, pgn pagination.Pagination) (Dashboard, error) {
	return s.repository.GetUserDashboard(ctx, userID, pgn)
}

func (s *Service) GetUserDashboardID(ctx context.Context, userID uint32) (uint32, error) {
	return s.repository.GetUserDashboardID(ctx, userID)
}

func (s *Service) AddToDashboard(ctx context.Context, userID uint32, assetID uint32, assetType assets.AssetType, description string) error {
	dashboardID, err := s.GetUserDashboardID(ctx, userID)
	if err != nil {
		return err
	}
	actionParams := DashboardActionParams{
		DashboardID: dashboardID,
		AssetID:     assetID,
		AssetType:   assetType,
		Description: description,
	}
	return s.repository.AddToDashboard(ctx, actionParams)
}

func (s *Service) RemoveFromDashboard(ctx context.Context, userID uint32, assetID uint32, assetType assets.AssetType) error {
	dashboardID, err := s.GetUserDashboardID(ctx, userID)
	if err != nil {
		return err
	}
	actionParams := DashboardActionParams{
		DashboardID: dashboardID,
		AssetID:     assetID,
		AssetType:   assetType,
	}
	return s.repository.RemoveFromDashboard(ctx, actionParams)
}

func (s *Service) EditDescription(ctx context.Context, userID uint32, assetID uint32, assetType assets.AssetType, description string) error {
	dashboardID, err := s.GetUserDashboardID(ctx, userID)
	if err != nil {
		return err
	}
	actionParams := DashboardActionParams{
		DashboardID: dashboardID,
		AssetID:     assetID,
		AssetType:   assetType,
		Description: description,
	}
	return s.repository.EditDescription(ctx, actionParams)
}
