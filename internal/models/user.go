package models

import "github.com/google/uuid"

// Favourite links a user to an Asset
type Favourite struct {
	UserID uuid.UUID `json:"user_id"`
	Asset  Asset     `json:"asset"`
}
