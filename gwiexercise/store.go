package gwiexercise

// Store defines the persistence operations for user assets.
// Implementations may store data in-memory, on disk, or using external services.
type Store interface {
	List(userID string) []Asset
	Add(userID string, a Asset) (Asset, error)
	Remove(userID, assetID string) error
	UpdateDescription(userID, assetID, description string) (Asset, error)
}
