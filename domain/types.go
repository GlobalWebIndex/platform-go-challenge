package domain

import "context"

type Domain struct {
	repo IDBRepository
}

type Chart struct {
	ID          string
	XTitle      string
	YTitle      string
	Description string
	DataJson    string
}

type Insight struct {
	ID          string
	Text        string
	Description string
}

type Audience struct {
	AgeMax            int
	AgeMin            int
	Gender            string
	Country           string
	HoursSpent        int
	NumberOfPurchases int
	Description       string
}

type AssetType string

const (
	ChartAssetType    = AssetType("chart")
	InsightAssetType  = AssetType("insight")
	AudienceAssetType = AssetType("audience")
)

type Asset struct {
	AssetType AssetType
	Data      interface{}
}
type QueryAssets struct {
}

type User struct {
	Username string
	Password string
	IsAdmin  bool
}

type LoginCredentials struct {
	Username string
	Password string
}

type IDomain interface {
	AddAsset(ctx context.Context, asset Asset) error
	DeleteAsset(ctx context.Context, assetID string) error
	UpdateAsset(ctx context.Context, assetID string, asset Asset) error
	ListAssets(ctx context.Context, userID string, query QueryAssets) error
	FavourAnAsset(ctx context.Context, userID, assetID string) error
	CreateUser(ctx context.Context, user User) error
	LoginUser(ctx context.Context, cred LoginCredentials) error
}

type IDBRepository interface {
	AddAsset(ctx context.Context, asset Asset) (*Asset, error)
	DeleteAsset(ctx context.Context, assetID string) error
	UpdateAsset(ctx context.Context, asset Asset) (*Asset, error)
	ListAssets(ctx context.Context, userID string, query QueryAssets) error
	FavourAnAsset(ctx context.Context, userID, assetID string) (string, error)
	CreateUser(ctx context.Context, user User) (*User, error)
	FindUser(ctx context.Context, cred LoginCredentials) (*User, error)
}
