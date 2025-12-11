package favourites

import (
	"errors"
	"time"
)

// AssetType describes the kind of asset stored as a favourite.
type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

// Asset is a tagged union. Type acts as the tag; exactly one of Chart,
// Insight, or Audience should be non-nil.
type Asset struct {
	ID          string
	Type        AssetType
	Description string

	Chart    *Chart
	Insight  *Insight
	Audience *Audience
}

// Chart contains chart-specific fields.
type Chart struct {
	Title string
	XAxis string
	YAxis string
	Data  []float64
}

// Insight contains insight-specific fields.
type Insight struct {
	Text string
}

// Audience contains simplified audience targeting fields.
type Audience struct {
	Gender                   string
	BirthCountry             string
	AgeGroup                 string
	HoursOnSocialMediaPerDay int
	PurchasesLastMonth       int
}

// Favourite links a user to an asset. Asset details are stored separately.
type Favourite struct {
	UserID    string
	AssetID   string
	CreatedAt time.Time
}

// FavouriteEntry combines a Favourite link with its Asset. Used when listing
// favourites to avoid leaking storage joins into higher layers.
type FavouriteEntry struct {
	Favourite
	Asset Asset
}

// ListFilter captures optional filters for listing favourites. Can be extended
// with pagination/sorting without changing the repository interface.
type ListFilter struct {
	// Type filters by asset type when non-nil.
	Type   *AssetType
	Limit  int
	Offset int
}

// Domain errors used by the favourites package.
var (
	ErrInvalidAssetType    = errors.New("invalid asset type")
	ErrInvalidAssetPayload = errors.New("invalid asset payload")

	ErrFavouriteNotFound = errors.New("favourite not found")
)

// Validate checks that exactly one payload matches the declared Type. Business
// rules (max lengths, value ranges) should be enforced in service/handler layers.
func (a Asset) Validate() error {
	switch a.Type {
	case AssetTypeChart:
		if a.Chart == nil || a.Insight != nil || a.Audience != nil {
			return ErrInvalidAssetPayload
		}
	case AssetTypeInsight:
		if a.Insight == nil || a.Chart != nil || a.Audience != nil {
			return ErrInvalidAssetPayload
		}
	case AssetTypeAudience:
		if a.Audience == nil || a.Chart != nil || a.Insight != nil {
			return ErrInvalidAssetPayload
		}
	default:
		return ErrInvalidAssetType
	}

	return nil
}

// IsValidAssetType reports whether t is one of the supported asset types.
func IsValidAssetType(t AssetType) bool {
	switch t {
	case AssetTypeChart, AssetTypeInsight, AssetTypeAudience:
		return true
	default:
		return false
	}
}
