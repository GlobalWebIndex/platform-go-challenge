package repository

import (
	"fmt"
	"go-challenge/models"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// func (r *InMemoryFavouriteRepository) seedMockAssets() {
// 	now := time.Now()

// 	// Mock insight
// 	r.assets["asset-001"] = models.Asset{
// 		ID:          "asset-001",
// 		Type:        models.AssetTypeInsight,
// 		Description: "Statistical insight about millennial social media usage",
// 		Insight: &models.Insight{
// 			Text: "40% of millennials spend more than 3 hours on social media daily",
// 		},
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}

// 	// Mock chart
// 	r.assets["asset-002"] = models.Asset{
// 		ID:          "asset-002",
// 		Type:        models.AssetTypeChart,
// 		Description: "Weekly engagement metrics showing daily active users",
// 		Chart: &models.Chart{
// 			Title:      "Daily Active Users",
// 			XAxisTitle: "Day of Week",
// 			YAxisTitle: "Number of Users",
// 			Data: []models.DataPoint{
// 				{X: "Monday", Y: 1000},
// 				{X: "Tuesday", Y: 1200},
// 				{X: "Wednesday", Y: 1100},
// 				{X: "Thursday", Y: 1300},
// 				{X: "Friday", Y: 1500},
// 			},
// 		},
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}

// 	// Mock audience
// 	r.assets["asset-003"] = models.Asset{
// 		ID:          "asset-003",
// 		Type:        models.AssetTypeAudience,
// 		Description: "Target demographic for marketing campaign",
// 		Audience: &models.Audience{
// 			Characteristics: []models.Characteristic{
// 				{Type: models.CharacteristicGender, Value: models.GenderMale},
// 				{Type: models.CharacteristicAgeGroup, Value: models.AgeGroup{Min: 24, Max: 35}},
// 				{Type: models.CharacteristicSocialMediaHours, Value: models.SocialMediaHours{Min: 3.0}},
// 			},
// 		},
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}
// }

// func (r *InMemoryFavouriteRepository) seedMockFavourites() {
// 	now := time.Now()
// 	// User-001 has favourited asset-001 and asset-002
// 	r.favourites["user-001"] = []models.Favourite{
// 		{
// 			ID:        uuid.New().String(),
// 			UserID:    "user-001",
// 			AssetID:   "asset-001",
// 			CreatedAt: now,
// 		},
// 		{
// 			ID:        uuid.New().String(),
// 			UserID:    "user-001",
// 			AssetID:   "asset-002",
// 			CreatedAt: now,
// 		},
// 	}

// 	// User-002 has favourited asset-003
// 	r.favourites["user-002"] = []models.Favourite{
// 		{
// 			ID:        uuid.New().String(),
// 			UserID:    "user-002",
// 			AssetID:   "asset-003",
// 			CreatedAt: now,
// 		},
// 	}
// }

// New function to seed both assets and a large set of random favourites
func (r *InMemoryFavouriteRepository) seedMockFavourites() {
	// 1. Generate Assets
	r.SeedMockAssets()

	// Initialize the favourites map
	if r.favourites == nil {
		r.favourites = make(map[string][]models.Favourite)
	}

	// 2. Define Mock Users and Asset IDs
	const totalAssets = 40

	// Create a single slice of all 100 asset IDs
	allAssetIDs := make([]string, totalAssets)
	for i := 0; i < totalAssets; i++ {
		// match the asset-001 format
		allAssetIDs[i] = fmt.Sprintf("asset-%03d", i+1)
	}

	// Shuffle the asset IDs do the split
	// covers different types (Insight/Chart/Audience) randomly
	rand.Shuffle(len(allAssetIDs), func(i, j int) {
		allAssetIDs[i], allAssetIDs[j] = allAssetIDs[j], allAssetIDs[i]
	})

	now := time.Now()

	// create a new Favourite
	createFavourite := func(userID, assetID string) models.Favourite {
		return models.Favourite{
			ID:        uuid.New().String(),
			UserID:    userID,
			AssetID:   assetID,
			CreatedAt: now,
		}
	}

	// --- User 1: 20 Favorites ---
	userID1 := "user-001"

	// Slice the first 50 shuffled asset IDs
	assetsForUser1 := allAssetIDs[:20]
	r.favourites[userID1] = make([]models.Favourite, 0, 20)

	for _, assetID := range assetsForUser1 {
		r.favourites[userID1] = append(r.favourites[userID1], createFavourite(userID1, assetID))
	}

	// --- User 2: 25 Favorites (Indices 50 to 74) ---
	userID2 := "user-002"

	// Slice the next 25 shuffled asset IDs
	assetsForUser2 := allAssetIDs[20:30]
	r.favourites[userID2] = make([]models.Favourite, 0, 10)

	for _, assetID := range assetsForUser2 {
		r.favourites[userID2] = append(r.favourites[userID2], createFavourite(userID2, assetID))
	}

	// --- User 3: 25 Favorites (Indices 75 to 99) ---
	userID3 := "user-003"

	// Slice the last 25 shuffled asset IDs
	assetsForUser3 := allAssetIDs[30:40]
	r.favourites[userID3] = make([]models.Favourite, 0, 10)

	for _, assetID := range assetsForUser3 {
		r.favourites[userID3] = append(r.favourites[userID3], createFavourite(userID3, assetID))
	}
}

// Function to generate  mock assets
func (r *InMemoryFavouriteRepository) SeedMockAssets() {
	assentNum := 40
	// Initialize the map if it's nil
	if r.assets == nil {
		r.assets = make(map[string]models.Asset)
	}

	// Seed the random number generator
	rand.NewSource(time.Now().UnixNano())

	// Generate  assets, based on three types (Insight, Chart, Audience)
	for i := 1; i <= assentNum; i++ {
		assetID := fmt.Sprintf("asset-%03d", i)
		assetType := i % 3

		var newAsset models.Asset
		now := time.Now()

		switch assetType {
		case 0: // Insight (i=3, 6, 9, ...)
			newAsset = generateRandomInsight(assetID, now)
		case 1: // Chart (i=1, 4, 7, ...)
			newAsset = generateRandomChart(assetID, now)
		case 2: // Audience (i=2, 5, 8, ...)
			newAsset = generateRandomAudience(assetID, now)
		}

		r.assets[assetID] = newAsset
	}
}

// Helper Functions to Generate Specific Asset Types.
func generateRandomInsight(id string, now time.Time) models.Asset {
	// Generate a random percentage and metric
	percentage := rand.Intn(50) + 10 // 10% to 59%
	metrics := []string{"engagement rate", "click-through rate", "conversion rate", "bounce rate"}
	metric := metrics[rand.Intn(len(metrics))]

	description := fmt.Sprintf("Statistical insight for %s data.", metric)
	text := fmt.Sprintf("A recent study showed that the average %s increased by %d%% last month.", metric, percentage)

	return models.Asset{
		ID:          id,
		Type:        models.AssetTypeInsight,
		Description: description,
		Insight: &models.Insight{
			Text: text,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func generateRandomChart(id string, now time.Time) models.Asset {
	// Generate random data points for 7 days
	dataPoints := make([]models.DataPoint, 7)
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

	baseValue := rand.Intn(500) + 1000 // Base users between 1000 and 1499

	for i, day := range days {
		// Random fluctuation around the base value
		yValue := float64(baseValue + rand.Intn(500) - 250)
		dataPoints[i] = models.DataPoint{X: day, Y: yValue}
	}

	titles := []string{"Daily User Signups", "Weekly Revenue", "Server Load Metrics", "Mobile App Sessions"}
	title := titles[rand.Intn(len(titles))]

	description := fmt.Sprintf("Chart showing the past week's data for: %s", title)

	return models.Asset{
		ID:          id,
		Type:        models.AssetTypeChart,
		Description: description,
		Chart: &models.Chart{
			Title:      title,
			XAxisTitle: "Day",
			YAxisTitle: "Value",
			Data:       dataPoints,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func generateRandomAudience(id string, now time.Time) models.Asset {
	// Generate random characteristics
	characteristics := make([]models.Characteristic, 0)

	// 1. Random Age Group
	minAge := rand.Intn(10) + 18         // 18 to 27
	maxAge := minAge + rand.Intn(15) + 5 // 5 to 19 years older

	characteristics = append(characteristics, models.Characteristic{
		Type:  models.CharacteristicAgeGroup,
		Value: models.AgeGroup{Min: minAge, Max: maxAge},
	})

	// 2. Random Gender (Male, Female, or Any)
	genders := []models.Gender{models.GenderMale, models.GenderFemale}
	gender := genders[rand.Intn(len(genders))]

	characteristics = append(characteristics, models.Characteristic{
		Type:  models.CharacteristicGender,
		Value: gender,
	})

	// 3. Random Social Media Hours
	minHours := float64(rand.Intn(4) + 1) // 1.0 to 4.0
	characteristics = append(characteristics, models.Characteristic{
		Type:  models.CharacteristicSocialMediaHours,
		Value: models.SocialMediaHours{Min: minHours},
	})

	description := fmt.Sprintf("Target audience between %d-%d years old, spending over %.1f hrs on social media.", minAge, maxAge, minHours)

	return models.Asset{
		ID:          id,
		Type:        models.AssetTypeAudience,
		Description: description,
		Audience: &models.Audience{
			Characteristics: characteristics,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
