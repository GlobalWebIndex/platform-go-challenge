package users

import (
	"context"

	"platform-go-challenge/internal/app/assets"
	"platform-go-challenge/internal/app/dashboards"
	"platform-go-challenge/internal/pagination"
)

type Service struct {
	dashboardService dashboards.Service
}

func NewUsersService(dashboardService *dashboards.Service) *Service {
	return &Service{dashboardService: *dashboardService}
}

func (s *Service) GetDashboard(ctx context.Context, userID uint32, pgn pagination.Pagination) (dashboards.Dashboard, error) {
	return s.dashboardService.GetUserDashboard(ctx, userID, pgn)
}

func (s *Service) AddToDashboard(ctx context.Context, userID uint32, assetID uint32, assetType assets.AssetType, description string) error {
	return s.dashboardService.AddToDashboard(ctx, userID, assetID, assetType, description)
}

func (s *Service) RemoveFromDashboard(ctx context.Context, userID uint32, assetID uint32, assetType assets.AssetType) error {
	return s.dashboardService.RemoveFromDashboard(ctx, userID, assetID, assetType)
}

func (s *Service) EditDescription(ctx context.Context, userID uint32, assetID uint32, assetType assets.AssetType, description string) error {
	return s.dashboardService.EditDescription(ctx, userID, assetID, assetType, description)
}
