package models

// User struct defines the user model with ID and a map of favourite assets
type User struct {
	ID         int           `json:"id"`
	Favourites map[int]Asset `json:"favourites"` // map of Asset with asset id as key
}
