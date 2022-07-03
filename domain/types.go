package domain

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type Domain struct {
	validate *validator.Validate
	repo     IDBRepository
}

type XYData struct {
	X []interface{}
	Y []interface{}
}

type Chart struct {
	Title       string `validate:"required"`
	XTitle      string `validate:"required"`
	YTitle      string `validate:"required"`
	Description string `validate:"required"`
	Data        XYData
}

type Insight struct {
	Text        string `validate:"required"`
	Description string `validate:"required"`
}

type Audience struct {
	AgeMax            int        `validate:"required,gte=1,lte=102"`
	AgeMin            int        `validate:"required,gte=1,lte=102"`
	Gender            GenderType `validate:"required"`
	Country           string     `validate:"required"`
	HoursSpent        int        `validate:"required,gte=1,lte=24"`
	NumberOfPurchases int        `validate:"required,gte=1,lte=100"`
	Description       string     `validate:"required"`
}

type GenderType string

const (
	MaleGenderType   = GenderType("male")
	FemaleGenderType = GenderType("female")
)

type AssetType string

const (
	InsightAssetType  = AssetType("insight")
	AudienceAssetType = AssetType("audience")
	ChartAssetType    = AssetType("chart")
)

type Asset struct {
	ID          uint
	IsFavourite bool
	Data        interface{}
}

type QueryFavouriteAssets struct {
	FromUserID uint
	OnlyFav    bool
}

type QueryAssets struct {
	Limit  int       `validate:"required,gte=1"`
	LastID uint      `validate:"gte=0"`
	Type   AssetType `validate:"required"`
	IsDesc bool
}

type ListedAssets struct {
	Limit   int
	FirstID uint
	LastID  uint
	Type    AssetType
	Assets  []Asset
}

type User struct {
	ID       uint
	Username string `validate:"required"`
	Password string `validate:"required"`
	IsAdmin  bool
}

type LoginCredentials struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type IDomain interface {
	GetAsset(ctx context.Context, user *User, assetID uint, assetType AssetType) (*Asset, error)
	AddAsset(ctx context.Context, user *User, asset Asset) (*Asset, error)
	DeleteAsset(ctx context.Context, user *User, assetID uint, assetType AssetType) error
	UpdateAsset(ctx context.Context, user *User, asset Asset) (*Asset, error)
	ListAssets(ctx context.Context, user *User, query QueryAssets, favQuery *QueryFavouriteAssets) (*ListedAssets, error)
	FavouriteAsset(ctx context.Context, uuser *User, assetID uint, assetType AssetType, isFavourite bool) error
	CreateUser(ctx context.Context, user User) (*User, error)
	LoginUser(ctx context.Context, cred LoginCredentials) (*User, error)
}

type IDBRepository interface {
	AddAsset(ctx context.Context, asset Asset) (*Asset, error)
	DeleteAsset(ctx context.Context, at AssetType, assetID uint) error
	UpdateAsset(ctx context.Context, asset Asset) (*Asset, error)
	GetAsset(ctx context.Context, at AssetType, assetID uint) (*Asset, error)
	ListAssets(ctx context.Context, query QueryAssets) (*ListedAssets, error)
	FavouriteAsset(ctx context.Context, userID, assetID uint, at AssetType, isFavourite bool) (uint, error)
	ListFavouriteAssets(ctx context.Context, userID uint, onlyFav bool, query QueryAssets) (*ListedAssets, error)
	AddUser(ctx context.Context, user User) (*User, error)
	FindUser(ctx context.Context, username string) (*User, error)
	UserExists(ctx context.Context, username string) (bool, error)
	GetUser(ctx context.Context, userID uint) (*User, error)
}
