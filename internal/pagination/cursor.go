package pagination

import (
	"fmt"
	"strings"
	"time"
)

type Cursor struct {
	FavouritedAt time.Time
	AssetID      string
}

// Cursor token from favourited_at + asset_id.
func EncodeCursor(t time.Time, assetID string) string {
	return fmt.Sprintf("%s|%s", t.UTC().Format(time.RFC3339Nano), assetID)
}

// Parses a cursor token: "<rfc3339nano>|<asset_id>"
func ParseCursor(s string) (Cursor, error) {
	parts := strings.SplitN(s, "|", 2)
	if len(parts) != 2 {
		return Cursor{}, fmt.Errorf("cursor must be in format <favourited_at>|<asset_id>")
	}

	ts, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return Cursor{}, fmt.Errorf("cursor favourited_at must be RFC3339Nano")
	}

	assetID := strings.TrimSpace(parts[1])
	if assetID == "" {
		return Cursor{}, fmt.Errorf("cursor asset_id must be non-empty")
	}

	return Cursor{FavouritedAt: ts, AssetID: assetID}, nil
}
