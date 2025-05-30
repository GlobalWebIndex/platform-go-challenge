package postgres

import "github.com/uptrace/bun"

// User represents a user entity in the database.
// It includes fields for user identification, authentication, and relations to other entities.
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"` // Specifies the table name and alias for the User model.

	ID             int    `bun:"id,autoincrement,pk"`     // Primary key for the User table.
	Username       string `bun:"username,unique,notnull"` // Unique username for the user.
	HashedPassword string `bun:"hashed_password,notnull"` // Hashed password for user authentication.

	Favorites []*Favorite `bun:"rel:has-many,join:id=user_id"` // One-to-many relation with the Favorite table.
}
