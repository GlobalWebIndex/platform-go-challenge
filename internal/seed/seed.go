package seed

import (
	"context"
	"fmt"
	"log"
	"time"

	"platform-go-challenge/internal/domain/entities"
	"platform-go-challenge/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Seed populates the repository with dummy data for testing.
// In a real application, this would be guarded by a flag or only run in a test environment.
func Seed(ctx context.Context, repo repository.Repository, logger *zap.Logger) {
	logger.Info("seeding repository with data")

	user1, err := entities.NewUser()
	if err != nil {
		log.Fatal("failed to create user 1", zap.Error(err))
	}
	user2, err := entities.NewUser()
	if err != nil {
		log.Fatal("failed to create user 2", zap.Error(err))
	}

	repo.SaveUser(user1)
	repo.SaveUser(user2)

	// --- 2. Create Assets ---
	// Create an Insight asset
	insightAsset := &entities.Insight{
		AssetBase: entities.AssetBase{
			ID:          uuid.New().String(),
			Type:        "insight",
			Description: "LalalInsight",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Text: "40% of millennials spend more than 3 hours on social media daily.",
	}
	// Create a Chart asset
	chartAsset := &entities.Chart{
		AssetBase: entities.AssetBase{
			ID:          uuid.New().String(),
			Type:        "chart",
			Description: "LalalChart",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Title:      "Q3 Sales Overview",
		AxesTitles: []string{"Product", "Sales"},
		Data:       []int{150, 200, 180, 220},
	}
	// Create an Audience asset
	audienceAsset := &entities.Audience{
		AssetBase: entities.AssetBase{
			ID:          uuid.New().String(),
			Type:        "audience",
			Description: "LalaAudience",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Gender:       "Males",
		BirthCountry: "USA",
		AgeGroup:     "24-35",
	}

	// Add all assets to the repository
	repo.SaveAsset(insightAsset)
	repo.SaveAsset(chartAsset)
	repo.SaveAsset(audienceAsset)

	// --- 3. Add Favourites for Users ---
	// User 1 favourites: Insight and Chart
	repo.AddFavourite(user1.GetID(), insightAsset.GetID())
	repo.AddFavourite(user1.GetID(), chartAsset.GetID())

	// User 2 favourites: Insight and Audience
	repo.AddFavourite(user2.GetID(), insightAsset.GetID())
	repo.AddFavourite(user2.GetID(), audienceAsset.GetID())

	logger.Info("seeding complete", zap.Int("users", 2), zap.Int("assets", 3))
	logger.Info(fmt.Sprintf("Test User 1 ID: %s", user1.GetID()))
	logger.Info(fmt.Sprintf("Test User 2 ID: %s", user2.GetID()))
}
