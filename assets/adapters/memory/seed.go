package memory

import (
	"time"

	"github.com/gwi/platform-go-challenge/assets/domain"
)

// Seed returns pre-seeded insights, charts, and audiences for the in-memory catalog.
// IDs (ins-001, chr-001, aud-001, ...) are stable so add-by-reference keeps working.
func Seed(now time.Time) (
	insights map[string]domain.Insight,
	charts map[string]domain.Chart,
	audiences map[string]domain.Audience,
) {
	base := func(id, name, summary, ownerID string, createdDeltaDays int) domain.BaseAsset {
		return domain.BaseAsset{
			ID:        id,
			Name:      name,
			Summary:   summary,
			OwnerID:   ownerID,
			CreatedAt: now.AddDate(0, 0, -createdDeltaDays),
			UpdatedAt: now,
		}
	}

	insights = map[string]domain.Insight{
		"ins-001": {BaseAsset: base("ins-001", "Q3 rainfall and crop yield", "Correlation between precipitation and harvest volume", "user-1", 10), Body: "Regions with 15–20% above-average rainfall saw a 8% lift in cereal yield; irrigation cannot fully offset drought in the south."},
		"ins-002": {BaseAsset: base("ins-002", "Office occupancy post-pandemic", "Return-to-office patterns by city", "user-2", 9), Body: "Mid-week occupancy has stabilised at 62% of 2019 levels; Tuesday and Wednesday peak."},
		"ins-003": {BaseAsset: base("ins-003", "Solar adoption by region", "Residential PV installation growth", "user-1", 8), Body: "Sunbelt states lead; payback period under 7 years in 40% of zip codes."},
		"ins-004": {BaseAsset: base("ins-004", "Supply chain delays in timber", "Lead times for structural lumber", "user-3", 7), Body: "Average lead time remains 3–4 weeks above pre-2020 baseline; softwood most affected."},
		"ins-005": {BaseAsset: base("ins-005", "Local high-street footfall", "Weekly footfall index vs same week last year", "user-4", 6), Body: "Footfall up 12% YoY in town centres; weekend peaks exceed weekday by 40%."},
	}

	charts = map[string]domain.Chart{
		"chr-001": {BaseAsset: base("chr-001", "Monthly rainfall index", "Deviation from 10-year mean by month", "user-1", 10), AxisX: "Month", AxisY: "Index (100 = mean)", Series: []float64{82, 94, 118, 105, 91, 76, 88, 102, 115, 124, 98, 86}},
		"chr-002": {BaseAsset: base("chr-002", "Occupancy by weekday", "Share of capacity used", "user-2", 9), AxisX: "Day", AxisY: "% capacity", Series: []float64{48, 61, 64, 63, 58, 22, 18}},
		"chr-003": {BaseAsset: base("chr-003", "PV installations by quarter", "Cumulative residential (MW)", "user-1", 8), AxisX: "Quarter", AxisY: "MW added", Series: []float64{420, 485, 510, 498, 530, 555}},
		"chr-004": {BaseAsset: base("chr-004", "Lumber lead time (days)", "From order to delivery", "user-2", 7), AxisX: "Week", AxisY: "Days", Series: []float64{28, 31, 29, 33, 30, 32, 28}},
		"chr-005": {BaseAsset: base("chr-005", "Footfall index by day", "Index vs same day last year", "user-2", 6), AxisX: "Day", AxisY: "Index", Series: []float64{98, 104, 108, 106, 105, 118, 122}},
	}

	audiences = map[string]domain.Audience{
		"aud-001": {BaseAsset: base("aud-001", "Farm operators", "Active farm decision-makers", "user-1", 12), Reach: 2100000, SharePct: 4.2, SampleSize: 18400, Gender: []string{"all"}, AgeGroups: []string{"35-44", "45-54", "55-64"}, BirthCountry: []string{"US"}},
		"aud-002": {BaseAsset: base("aud-002", "Office workers hybrid", "At least 2 days in office", "user-2", 11), Reach: 28500000, SharePct: 8.6, SampleSize: 52000, AgeGroups: []string{"25-34", "35-44"}, HoursSocialMedia: []string{"0-1", "1-2"}},
		"aud-003": {BaseAsset: base("aud-003", "Homeowners planning solar", "Considering PV in next 24 months", "user-1", 10), Reach: 7200000, SharePct: 5.4, SampleSize: 31000, BirthCountry: []string{"US"}, AgeGroups: []string{"35-44", "45-54"}, Purchases: []string{"home improvement"}},
		"aud-004": {BaseAsset: base("aud-004", "Construction procurement", "Specify or buy building materials", "user-2", 9), Reach: 890000, SharePct: 1.2, SampleSize: 4200, Gender: []string{"male", "female"}, AgeGroups: []string{"25-34", "35-44", "45-54"}},
		"aud-005": {BaseAsset: base("aud-005", "High-street shoppers", "Visited town centre in last 7 days", "user-2", 8), Reach: 44200000, SharePct: 13.1, SampleSize: 68000, Gender: []string{"all"}, AgeGroups: []string{"18-24", "25-34", "35-44"}, HoursSocialMedia: []string{"1-2", "2-3", "3+"}},
	}

	return insights, charts, audiences
}
