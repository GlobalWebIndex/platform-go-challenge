package repo_inmemory

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/manos/favourites/assets/domain"
)

// todo add mutex for concurrency in real app
type AssetRepository struct {
	mu        sync.RWMutex
	insights  map[string]domain.InsightAsset
	audiences map[string]domain.AudienceAsset
	charts    map[string]domain.ChartAsset
}

func NewInMemoryAssetRepository() *AssetRepository {
	now := time.Now()

	base := func(id, name, desc, user string, createdDeltaDays int) domain.Asset {
		created := now.AddDate(0, 0, -createdDeltaDays)
		return domain.Asset{
			Id:          id,
			Name:        name,
			Description: desc,
			UserId:      user,
			CreatedAt:   created,
			UpdatedAt:   now,
		}
	}

	insights := map[string]domain.InsightAsset{
		"ins-001": {Asset: base("ins-001", "Gen Z snack trends", "Summary of snack behaviors", "user-1", 10), Text: "Gen Z over-indexes on salty snacks and late-night snacking."},
		"ins-002": {Asset: base("ins-002", "Streaming churn drivers", "Why users cancel subscriptions", "user-2", 9), Text: "Price increases and content fatigue are the top churn triggers."},
		"ins-003": {Asset: base("ins-003", "Fitness motivations", "Drivers behind workout habits", "user-1", 8), Text: "Consistency is tied to social accountability and short routines."},
		"ins-004": {Asset: base("ins-004", "Travel intent 2025", "Where people want to go next", "user-3", 7), Text: "City breaks and nearby destinations dominate due to cost sensitivity."},
		"ins-005": {Asset: base("ins-005", "Brand trust signals", "What increases consumer trust", "user-4", 6), Text: "Transparency, reviews, and customer support responsiveness matter most."},
	}

	audiences := map[string]domain.AudienceAsset{
		"aud-001": {Asset: base("aud-001", "Core Audience", "Baseline audience definition", "user-1", 12), SampleSize: 328865, TotalRespondents: 976245, EstimatedReach: 995490000, PopulationPercent: 36.2},
		"aud-002": {Asset: base("aud-002", "Frequent Travelers", "Travelers with 3+ trips/year", "user-2", 11), SampleSize: 120540, TotalRespondents: 976245, EstimatedReach: 210000000, PopulationPercent: 7.8},
		"aud-003": {Asset: base("aud-003", "Mobile Gamers", "People who game on mobile weekly", "user-1", 10), SampleSize: 244901, TotalRespondents: 976245, EstimatedReach: 450000000, PopulationPercent: 16.4},
		"aud-004": {Asset: base("aud-004", "Eco-conscious Shoppers", "Prefer sustainable brands", "user-2", 9), SampleSize: 98211, TotalRespondents: 976245, EstimatedReach: 175000000, PopulationPercent: 6.1},
		"aud-005": {Asset: base("aud-005", "Premium Streamers", "Paying for 2+ subscriptions", "user-2", 8), SampleSize: 150332, TotalRespondents: 976245, EstimatedReach: 260000000, PopulationPercent: 9.5},
	}

	charts := map[string]domain.ChartAsset{
		"chr-001": {Asset: base("chr-001", "Snack Frequency", "Snacking per day distribution", "user-1", 10), XAxisTitle: "Snacks/day", YAxisTitle: "% of audience", Data: []float64{5, 12, 25, 30, 18, 10}},
		"chr-002": {Asset: base("chr-002", "Churn Reasons", "Top reasons for subscription cancellation", "user-2", 9), XAxisTitle: "Reason index", YAxisTitle: "Score", Data: []float64{78, 64, 52, 41, 33}},
		"chr-003": {Asset: base("chr-003", "Workout Days", "Days worked out per week", "user-1", 8), XAxisTitle: "Days/week", YAxisTitle: "% of audience", Data: []float64{8, 14, 22, 26, 18, 12}},
		"chr-004": {Asset: base("chr-004", "Travel Intent", "Intent to travel next 12 months", "user-2", 7), XAxisTitle: "Intent bucket", YAxisTitle: "% of audience", Data: []float64{15, 28, 34, 23}},
		"chr-005": {Asset: base("chr-005", "Trust Signals", "Impact of trust signals on conversion", "user-2", 6), XAxisTitle: "Signal index", YAxisTitle: "Lift", Data: []float64{1.05, 1.12, 1.18, 1.09, 1.15}},
	}

	return &AssetRepository{
		insights:  insights,
		audiences: audiences,
		charts:    charts,
	}
}

func (r *AssetRepository) GetInsight(_ context.Context, userId string, id string) (domain.InsightAsset, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.insights[id]
	if a.UserId != userId {
		log.Printf("repo assets: insight id=%s user=%s forbidden owner=%s", id, userId, a.UserId)
		return domain.InsightAsset{}, false
	}
	log.Printf("repo assets: insight id=%s user=%s hit=%t", id, userId, ok)
	return a, ok
}

func (r *AssetRepository) GetInsights(_ context.Context, userId string, ids []string) map[string]domain.InsightAsset {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]domain.InsightAsset, len(ids))
	for _, id := range ids {
		if a, ok := r.insights[id]; ok && a.UserId == userId {
			out[id] = a
		}
	}
	log.Printf("repo assets: insights batch user=%s requested=%d returned=%d", userId, len(ids), len(out))
	return out
}

func (r *AssetRepository) GetAudience(_ context.Context, userId string, id string) (domain.AudienceAsset, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.audiences[id]
	if a.UserId != userId {
		log.Printf("repo assets: audience id=%s user=%s forbidden owner=%s", id, userId, a.UserId)
		return domain.AudienceAsset{}, false
	}
	log.Printf("repo assets: audience id=%s user=%s hit=%t", id, userId, ok)
	return a, ok
}

func (r *AssetRepository) GetAudiences(_ context.Context, userId string, ids []string) map[string]domain.AudienceAsset {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]domain.AudienceAsset, len(ids))
	for _, id := range ids {
		if a, ok := r.audiences[id]; ok && a.UserId == userId {
			out[id] = a
		}
	}
	log.Printf("repo assets: audiences batch user=%s requested=%d returned=%d", userId, len(ids), len(out))
	return out
}

func (r *AssetRepository) GetChart(_ context.Context, userId string, id string) (domain.ChartAsset, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.charts[id]
	if a.UserId != userId {
		log.Printf("repo assets: chart id=%s user=%s forbidden owner=%s", id, userId, a.UserId)
		return domain.ChartAsset{}, false
	}
	log.Printf("repo assets: chart id=%s user=%s hit=%t", id, userId, ok)
	return a, ok
}

func (r *AssetRepository) GetCharts(_ context.Context, userId string, ids []string) map[string]domain.ChartAsset {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]domain.ChartAsset, len(ids))
	for _, id := range ids {
		if a, ok := r.charts[id]; ok && a.UserId == userId {
			out[id] = a
		}
	}
	log.Printf("repo assets: charts batch user=%s requested=%d returned=%d", userId, len(ids), len(out))
	return out
}
