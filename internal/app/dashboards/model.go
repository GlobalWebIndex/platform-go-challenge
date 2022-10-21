package dashboards

import "platform-go-challenge/internal/app/assets"

type Dashboard struct {
	ID        uint32           `json:"id"`
	UserID    uint32           `json:"-"`
	Charts    StarredCharts    `json:"charts"`
	Insights  StarredInsights  `json:"insights"`
	Audiences StarredAudiences `json:"audiences"`
}

type StarredChart struct {
	assets.Chart
	Description string `json:"description"`
}

type StarredCharts []StarredChart

type StarredInsight struct {
	assets.Insight
	Description string `json:"description"`
}

type StarredInsights []StarredInsight

type StarredAudience struct {
	assets.Audience
	Description string `json:"description"`
}

type StarredAudiences []StarredAudience
