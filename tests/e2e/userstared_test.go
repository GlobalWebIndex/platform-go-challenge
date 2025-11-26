package e2e

import (
	"encoding/json"
	"platform-go-challenge/models"
	"testing"
)

// GraphQL response types (IDs are strings in GraphQL)
type gqlAudience struct {
	ID            string `json:"id"`
	Gender        string `json:"gender"`
	Birthcountry  string `json:"birthcountry"`
	Agegroup      string `json:"agegroup"`
	Dailyhours    int    `json:"dailyhours"`
	Noofpurchases int    `json:"noofpurchases"`
}

type gqlChart struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Xaxistitle string `json:"xaxistitle"`
	Yaxistitle string `json:"yaxistitle"`
}

type gqlInsight struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// TestUserStared_EmptyFavourites tests when user has no favourites
func TestUserStared_EmptyFavourites(t *testing.T) {
	CleanupTestData(testDB)

	query := `
		query GetUserStared($userID: ID!) {
			userstared(userID: $userID) {
				userid
				audience {
					id
					gender
				}
				chart {
					id
					title
				}
				insight {
					id
					text
				}
			}
		}
	`

	variables := map[string]interface{}{
		"userID": "999",
	}

	resp := ExecuteGraphQL(t, query, variables)

	if len(resp.Errors) > 0 {
		t.Fatalf("expected no errors, got: %v", resp.Errors)
	}

	var result struct {
		Userstars struct {
			Userid   int           `json:"userid"`
			Audience []gqlAudience `json:"audience"`
			Chart    []gqlChart    `json:"chart"`
			Insight  []gqlInsight  `json:"insight"`
		} `json:"userstared"`
	}

	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.Userstars.Userid != 999 {
		t.Errorf("expected userid 999, got %d", result.Userstars.Userid)
	}

	if len(result.Userstars.Audience) != 0 {
		t.Errorf("expected 0 audiences, got %d", len(result.Userstars.Audience))
	}

	if len(result.Userstars.Chart) != 0 {
		t.Errorf("expected 0 charts, got %d", len(result.Userstars.Chart))
	}

	if len(result.Userstars.Insight) != 0 {
		t.Errorf("expected 0 insights, got %d", len(result.Userstars.Insight))
	}
}

// TestUserStared_WithFavourites tests fetching user favourites with all asset types
func TestUserStared_WithFavourites(t *testing.T) {
	CleanupTestData(testDB)

	// Seed test data
	audienceID, chartID, insightID := SeedTestData(t, testDB)

	// Create favourites for user 1
	favourites := []models.UserStar{
		{UserID: 1, Type: models.AssetTypeAudience, AssetID: audienceID},
		{UserID: 1, Type: models.AssetTypeChart, AssetID: chartID},
		{UserID: 1, Type: models.AssetTypeInsight, AssetID: insightID},
	}

	for _, fav := range favourites {
		if err := testDB.Create(&fav).Error; err != nil {
			t.Fatalf("failed to create favourite: %v", err)
		}
	}

	query := `
		query GetUserStared($userID: ID!) {
			userstared(userID: $userID) {
				userid
				audience {
					id
					gender
					birthcountry
					agegroup
					dailyhours
					noofpurchases
				}
				chart {
					id
					title
					xaxistitle
					yaxistitle
				}
				insight {
					id
					text
				}
			}
		}
	`

	variables := map[string]interface{}{
		"userID": "1",
	}

	resp := ExecuteGraphQL(t, query, variables)

	if len(resp.Errors) > 0 {
		t.Fatalf("expected no errors, got: %v", resp.Errors)
	}

	var result struct {
		Userstars struct {
			Userid   int           `json:"userid"`
			Audience []gqlAudience `json:"audience"`
			Chart    []gqlChart    `json:"chart"`
			Insight  []gqlInsight  `json:"insight"`
		} `json:"userstared"`
	}

	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Verify user ID
	if result.Userstars.Userid != 1 {
		t.Errorf("expected userid 1, got %d", result.Userstars.Userid)
	}

	// Verify audiences
	if len(result.Userstars.Audience) != 1 {
		t.Fatalf("expected 1 audience, got %d", len(result.Userstars.Audience))
	}
	if result.Userstars.Audience[0].Gender != "Male" {
		t.Errorf("expected gender 'Male', got '%s'", result.Userstars.Audience[0].Gender)
	}
	if result.Userstars.Audience[0].Birthcountry != "USA" {
		t.Errorf("expected birthcountry 'USA', got '%s'", result.Userstars.Audience[0].Birthcountry)
	}

	// Verify charts
	if len(result.Userstars.Chart) != 1 {
		t.Fatalf("expected 1 chart, got %d", len(result.Userstars.Chart))
	}
	if result.Userstars.Chart[0].Title != "Sales Chart" {
		t.Errorf("expected title 'Sales Chart', got '%s'", result.Userstars.Chart[0].Title)
	}

	// Verify insights
	if len(result.Userstars.Insight) != 1 {
		t.Fatalf("expected 1 insight, got %d", len(result.Userstars.Insight))
	}
	if result.Userstars.Insight[0].Text != "Revenue increased by 20% this quarter" {
		t.Errorf("expected specific insight text, got '%s'", result.Userstars.Insight[0].Text)
	}
}

// TestUserStared_MultipleFavouritesOfSameType tests multiple favourites of the same type
func TestUserStared_MultipleFavouritesOfSameType(t *testing.T) {
	CleanupTestData(testDB)

	// Create multiple audiences
	audience1 := models.Audience{
		Gender:        "Male",
		BirthCountry:  "USA",
		AgeGroup:      "25-34",
		DailyHours:    5,
		NoOfPurchases: 10,
	}
	testDB.Create(&audience1)

	audience2 := models.Audience{
		Gender:        "Female",
		BirthCountry:  "Canada",
		AgeGroup:      "35-44",
		DailyHours:    3,
		NoOfPurchases: 7,
	}
	testDB.Create(&audience2)

	// Create favourites
	testDB.Create(&models.UserStar{UserID: 2, Type: models.AssetTypeAudience, AssetID: audience1.ID})
	testDB.Create(&models.UserStar{UserID: 2, Type: models.AssetTypeAudience, AssetID: audience2.ID})

	query := `
		query GetUserStared($userID: ID!) {
			userstared(userID: $userID) {
				userid
				audience {
					id
					gender
					birthcountry
				}
			}
		}
	`

	variables := map[string]interface{}{
		"userID": "2",
	}

	resp := ExecuteGraphQL(t, query, variables)

	if len(resp.Errors) > 0 {
		t.Fatalf("expected no errors, got: %v", resp.Errors)
	}

	var result struct {
		Userstars struct {
			Userid   int           `json:"userid"`
			Audience []gqlAudience `json:"audience"`
		} `json:"userstared"`
	}

	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(result.Userstars.Audience) != 2 {
		t.Fatalf("expected 2 audiences, got %d", len(result.Userstars.Audience))
	}

	// Verify both audiences are returned
	genders := make(map[string]bool)
	for _, aud := range result.Userstars.Audience {
		genders[aud.Gender] = true
	}

	if !genders["Male"] || !genders["Female"] {
		t.Errorf("expected both Male and Female audiences")
	}
}

// TestUserStared_OnlySpecificUser tests that only the requested user's favourites are returned
func TestUserStared_OnlySpecificUser(t *testing.T) {
	CleanupTestData(testDB)

	// Create test data
	audienceID, chartID, _ := SeedTestData(t, testDB)

	// Create favourites for user 1
	testDB.Create(&models.UserStar{UserID: 1, Type: models.AssetTypeAudience, AssetID: audienceID})

	// Create favourites for user 2
	testDB.Create(&models.UserStar{UserID: 2, Type: models.AssetTypeChart, AssetID: chartID})

	query := `
		query GetUserStared($userID: ID!) {
			userstared(userID: $userID) {
				userid
				audience { id }
				chart { id }
			}
		}
	`

	// Query for user 1
	variables := map[string]interface{}{
		"userID": "1",
	}

	resp := ExecuteGraphQL(t, query, variables)

	if len(resp.Errors) > 0 {
		t.Fatalf("expected no errors, got: %v", resp.Errors)
	}

	var result struct {
		Userstars struct {
			Userid   int           `json:"userid"`
			Audience []gqlAudience `json:"audience"`
			Chart    []gqlChart    `json:"chart"`
		} `json:"userstared"`
	}

	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// User 1 should have 1 audience and 0 charts
	if len(result.Userstars.Audience) != 1 {
		t.Errorf("expected 1 audience for user 1, got %d", len(result.Userstars.Audience))
	}

	if len(result.Userstars.Chart) != 0 {
		t.Errorf("expected 0 charts for user 1, got %d", len(result.Userstars.Chart))
	}
}

// TestUserStared_InvalidUserID tests error handling for invalid user ID
func TestUserStared_InvalidUserID(t *testing.T) {
	query := `
		query GetUserStared($userID: ID!) {
			userstared(userID: $userID) {
				userid
			}
		}
	`

	variables := map[string]interface{}{
		"userID": "invalid",
	}

	resp := ExecuteGraphQL(t, query, variables)

	// Should have an error for invalid user ID
	if len(resp.Errors) == 0 {
		t.Fatalf("expected error for invalid user ID, got none")
	}

	if resp.Errors[0].Message == "" {
		t.Errorf("expected error message, got empty string")
	}
}
