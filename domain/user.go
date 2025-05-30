package domain

// User represents a user account in the system.
// It stores basic user identification and authentication information.
type User struct {
	ID             string // Unique identifier for the user
	Username       string // Username chosen by the user for login purposes
	HashedPassword string // Securely hashed version of the user's password
}
