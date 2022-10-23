package dashboards

import (
	"context"

	"platform-go-challenge/internal/pagination"
)

type Repository interface {
	GetUserDashboard(ctx context.Context, userID uint32, pgn pagination.Pagination) (Dashboard, error)
	GetUserDashboardID(ctx context.Context, userID uint32) (uint32, error)

	AddToDashboard(ctx context.Context, actionParams DashboardActionParams) error
	RemoveFromDashboard(ctx context.Context, actionParams DashboardActionParams) error
	EditDescription(ctx context.Context, actionParams DashboardActionParams) error
}
