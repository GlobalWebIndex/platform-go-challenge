package user_model

// User represents a user in the system with an ID and an Email.
type User struct {
	ID    int    `json:"id"`    // Unique identifier for the user.
	Email string `json:"email"` // Email is the user's email address.
}
