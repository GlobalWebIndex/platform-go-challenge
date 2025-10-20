package db

import (
	"github.com/google/uuid"

	"github.com/mr-polymath/platform-go-challenge/internal/model"
)

// SeedMemoryStore populates the given MemoryStore with demo users and assets.
func SeedMemoryStore(store *MemoryStore) {
	user1 := uuid.MustParse("431ff79b-4e65-4625-9f7f-c8125e2be057")
	user2 := uuid.MustParse("d61733f6-910e-4e26-bb31-0cac3a7d4279")

	// User 1 Assets
	chart1 := model.NewAsset(model.ChartType, "Monthly Revenue", model.ChartData{
		Title: "Revenue by Month",
		XAxis: "Month",
		YAxis: "Revenue ($)",
		Values: []model.ChartPoint{
			{X: "2025-01", Y: 10000},
			{X: "2025-02", Y: 12000},
			{X: "2025-03", Y: 15000},
			{X: "2025-04", Y: 18000},
			{X: "2025-05", Y: 17000},
		},
	})

	insight1 := model.NewAsset(model.InsightType, "Social Media Usage", model.InsightData{
		Text: "40% of millennials spend more than 3 hours on social media daily.",
	})

	audience1 := model.NewAsset(model.AudienceType, "Frequent Buyers", model.AudienceData{
		Gender:             "Female",
		BirthCountry:       "UK",
		AgeGroup:           "25-34",
		HoursSocialMedia:   2.5,
		PurchasesLastMonth: 8,
	})

	store.AddFavourite(user1, &chart1)
	store.AddFavourite(user1, &insight1)
	store.AddFavourite(user1, &audience1)

	// User 2 Assets
	chart2 := model.NewAsset(model.ChartType, "Quarterly Expenses", model.ChartData{
		Title: "Expenses by Quarter",
		XAxis: "Quarter",
		YAxis: "Expenses ($)",
		Values: []model.ChartPoint{
			{X: "Q1", Y: 5000},
			{X: "Q2", Y: 7000},
			{X: "Q3", Y: 6000},
			{X: "Q4", Y: 8000},
		},
	})

	insight2 := model.NewAsset(model.InsightType, "E-commerce Trend", model.InsightData{
		Text: "Online purchases increased by 25% last quarter.",
	})

	store.AddFavourite(user2, &chart2)
	store.AddFavourite(user2, &insight2)

	// Log users so you can test easily
	println("Seeded users:")
	println("User 1:", user1.String())
	println("User 2:", user2.String())
}
