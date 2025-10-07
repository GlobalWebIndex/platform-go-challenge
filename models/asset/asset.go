package asset_model

import "encoding/json"

// Asset represents a digital asset with an ID, associated user ID, and a description.
type Asset struct {
	ID          int             `json:"id"`          // Unique identifier for the asset
	UserID      int             `json:"user_id"`     // Identifier for the user who owns the asset
	Description json.RawMessage `json:"description"` // Description of the asset
}
