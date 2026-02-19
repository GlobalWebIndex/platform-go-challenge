package pagination

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCursor_RoundTrip(t *testing.T) {
	tm := time.Date(2026, 2, 18, 10, 0, 0, 123, time.UTC)
	cur := EncodeCursor(tm, "asset-1")

	parsed, err := ParseCursor(cur)
	require.NoError(t, err)
	require.True(t, parsed.FavouritedAt.Equal(tm))
	require.Equal(t, "asset-1", parsed.AssetID)
}

func TestCursor_Parse_Invalid(t *testing.T) {
	_, err := ParseCursor("no-separator")
	require.Error(t, err)

	_, err = ParseCursor("bad-time|asset-1")
	require.Error(t, err)

	_, err = ParseCursor(time.Now().UTC().Format(time.RFC3339Nano) + "|")
	require.Error(t, err)
}
