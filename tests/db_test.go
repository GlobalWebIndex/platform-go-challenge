package tests

import (
	"platform-go-challenge/db"
	"platform-go-challenge/models"
	"testing"
)

func TestAddAndGetFavorites(test *testing.T) {
	database := db.NewDB()
	asset := models.Asset{ID: "12", Description: "Testing an asset"}

	//Add the asset
	err := database.AddFavorite("user123", asset)
	if err != nil {
		test.Errorf("unexpected error: %v", err)
	}

	//Get the asset
	assets, err := database.GetFavorites("user123")
	if err != nil {
		test.Errorf("unexpected error: %v", err)
	}
	if len(assets) != 1 {
		test.Errorf("expected 1 asset, got %d", len(assets))
	}
	if assets[0].ID != "12" {
		test.Errorf("expected 12, got %s", assets[0].ID)
	}
}

func TestRemoveFavorite(test *testing.T) {
	database := db.NewDB()
	asset := models.Asset{ID: "12", Description: "Testing an asset to be removed"}
	_ = database.AddFavorite("user123", asset)

	err := database.RemoveFavorite("user123", "12")
	if err != nil {
		test.Errorf("unexpected error: %v", err)
	}

	assets, _ := database.GetFavorites("user123")
	if len(assets) != 0 {
		test.Errorf("expected 0 assets, got %d", len(assets))
	}
}

func TestUpdateFavorite(test *testing.T) {
	database := db.NewDB()
	asset := models.Asset{ID: "12", Description: "Old description"}
	_ = database.AddFavorite("user123", asset)

	updatedAsset := models.Asset{ID: "12", Description: "New description"}
	err := database.UpdateFavorites("user123", updatedAsset)
	if err != nil {
		test.Errorf("unexpected error: %v", err)
	}

	assets, _ := database.GetFavorites("user123")
	if assets[0].Description != "New description" {
		test.Errorf("expected new description, got %s", assets[0].Description)
	}
}

func TestInvalidUserID(test *testing.T) {
	database := db.NewDB()
	asset := models.Asset{ID: "12", Description: "Just an asset"}

	err := database.AddFavorite("", asset)
	if err == nil {
		test.Errorf("expected error for empty userId, got nil")
	}

	//Get favorites should fail for invalid userId
	_, err = database.GetFavorites("")
	if err == nil {
		test.Errorf("expected error for empty userId, got nil")
	}

}
