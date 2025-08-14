package domain

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserFavorite represents the relationship between a user and their favorite asset
type UserFavorite struct {
	UserID    string    `json:"user_id"`
	AssetID   string    `json:"asset_id"`
	Asset     Asset     `json:"asset"`
	AddedAt   time.Time `json:"added_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user
func NewUser(id, email, name string) *User {
	now := time.Now()
	return &User{
		ID:        id,
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewUserFavorite creates a new user favorite relationship
func NewUserFavorite(userID string, asset Asset) *UserFavorite {
	now := time.Now()
	return &UserFavorite{
		UserID:    userID,
		AssetID:   asset.GetID(),
		Asset:     asset,
		AddedAt:   now,
		UpdatedAt: now,
	}
}
