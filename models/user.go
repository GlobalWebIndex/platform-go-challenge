package models

// User represents a user with an ID and a list of favourite assets
// @Description User represents a user with an ID and a list of favourite assets
type User struct {
	ID         string  `json:"id" example:"1"`
	Favourites []Asset `json:"favourites"`
}
