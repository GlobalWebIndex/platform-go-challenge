package favorite_model

import asset_model "gwi/models/asset"

// Favorite represents a user's favorite asset in the system.
type Favorite struct {
	ID      int                `json:"id"`              // Unique identifier for the favorite
	UserID  int                `json:"user_id"`         // Identifier for the user who owns the asset
	AssetID int                `json:"asset_id"`        // Identifier for the asset who belongs to the user
	Asset   *asset_model.Asset `json:"asset,omitempty"` // linked object (optional)
}
