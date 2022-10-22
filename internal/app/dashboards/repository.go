package dashboards

import (
	"context"
)

type Repository interface {
	GetUserDashboard(ctx context.Context, userID uint32) (Dashboard, error)
	GetUserDashboardID(ctx context.Context, userID uint32) (uint32, error)

	AddToDashboard(ctx context.Context, actionParams DashboardActionParams) error
	RemoveFromDashboard(ctx context.Context, actionParams DashboardActionParams) error
	EditDescription(ctx context.Context, actionParams DashboardActionParams) error
}
