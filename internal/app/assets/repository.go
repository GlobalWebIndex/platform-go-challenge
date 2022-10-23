package assets

import "context"

type Repository interface {
	List(ctx context.Context) (Assets, error)
	ListCharts(ctx context.Context) (Charts, error)
	ListInsights(ctx context.Context) (Insights, error)
	ListAudiences(ctx context.Context) (Audiences, error)
}
