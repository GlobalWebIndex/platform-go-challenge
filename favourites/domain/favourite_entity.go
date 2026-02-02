package domain

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type FavouriteType string

const (
	FavouriteInsight  FavouriteType = "INSIGHT"
	FavouriteAudience FavouriteType = "AUDIENCE"
	FavouriteChart    FavouriteType = "CHART"
)

type FavouriteEntity struct {
	ID          string
	UserID      string
	AssetID     string
	Description string
	Type        FavouriteType
	CreatedAt   time.Time
}

func NewFavourite(userID string, favType FavouriteType, assetID string, description string) (FavouriteEntity, error) {
	// Here you can add validation for favType if needed
	return FavouriteEntity{
		ID:          generateID(), // Assume generateID() generates a unique ID
		UserID:      userID,
		AssetID:     assetID,
		Description: description,
		Type:        favType,
		CreatedAt:   time.Now(),
	}, nil
}

func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("fav-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
