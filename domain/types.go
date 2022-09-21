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
	X []float64 `json:"x"`
	Y []float64 `json:"y"`
}

type Chart struct {
	Title       string `validate:"required" json:"title"`
	XTitle      string `validate:"required" json:"xTitle"`
	YTitle      string `validate:"required" json:"yTitle"`
	Description string `validate:"required" json:"description"`
	Data        XYData `json:"data"`
}

type Insight struct {
	Text        string `validate:"required" json:"text"`
	Description string `validate:"required" json:"description"`
}

type Audience struct {
	AgeMax            int        `validate:"required,gte=1,lte=102" json:"ageMax"`
	AgeMin            int        `validate:"required,gte=1,lte=102" json:"ageMin"`
	Gender            GenderType `validate:"required" json:"gender"`
	Country           string     `validate:"required" json:"country"`
	HoursSpent        int        `validate:"required,gte=1,lte=24" json:"hoursSpent"`
	NumberOfPurchases int        `validate:"required,gte=1,lte=100" json:"numberOfPurchases"`
	Description       string     `validate:"required" json:"description"`
}

type GenderType string

const (
	MaleGenderType   = GenderType("male")
	FemaleGenderType = GenderType("female")
)

type AssetType string

const (
	InsightAssetType  = AssetType("insights")
	AudienceAssetType = AssetType("audiences")
	ChartAssetType    = AssetType("charts")
)

type Asset struct {
	ID          uint        `json:"id"`
	IsFavourite *bool       `json:"isFavourite,omitempty"`
	Data        interface{} `json:"data"`
}

func (a *Asset) GetData() interface{} {
	return a.Data
}

type InputAsset struct {
	Data interface{}
}

func (ia *InputAsset) GetData() interface{} {
	return ia.Data
}

type QueryFavouriteAssets struct {
	FromUserID uint `json:"fromUserID"`
	OnlyFav    bool `json:"onlyFavourite"`
}

type QueryAssets struct {
	Limit  int       `validate:"required,gte=1" json:"limit"`
	LastID uint      `validate:"gte=0" json:"lastID"`
	Type   AssetType `validate:"required" json:"type"`
	IsDesc bool      `json:"isDesc"`
}

type ListedAssets struct {
	Limit   int       `json:"limit"`
	FirstID uint      `json:"firstID"`
	LastID  uint      `json:"lastID"`
	Type    AssetType `json:"type"`
	Assets  []Asset   `json:"assets"`
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

type IAsset interface {
	GetData() interface{}
}

type IDomain interface {
	GetAsset(ctx context.Context, user *User, assetID uint, assetType AssetType) (*Asset, error)
	AddAsset(ctx context.Context, user *User, asset InputAsset) (*Asset, error)
	DeleteAsset(ctx context.Context, user *User, assetID uint, assetType AssetType) error
	UpdateAsset(ctx context.Context, user *User, assetID uint, asset InputAsset) (*Asset, error)
	ListAssets(ctx context.Context, user *User, query QueryAssets, favQuery *QueryFavouriteAssets) (*ListedAssets, error)
	FavouriteAsset(ctx context.Context, uuser *User, assetID uint, assetType AssetType, isFavourite bool) error
	CreateUser(ctx context.Context, user User) (*User, error)
	LoginUser(ctx context.Context, cred LoginCredentials) (*User, error)
}

type IDBRepository interface {
	AddAsset(ctx context.Context, asset InputAsset) (*Asset, error)
	DeleteAsset(ctx context.Context, at AssetType, assetID uint) error
	UpdateAsset(ctx context.Context, assetID uint, asset InputAsset) (*Asset, error)
	GetAsset(ctx context.Context, at AssetType, assetID uint) (*Asset, error)
	ListAssets(ctx context.Context, query QueryAssets) (*ListedAssets, error)
	RemoveFavouriteAssetFromEveryone(ctx context.Context, assetID uint, at AssetType) error
	FavouriteAsset(ctx context.Context, userID, assetID uint, at AssetType, isFavourite bool) (uint, error)
	ListFavouriteAssets(ctx context.Context, userID uint, onlyFav bool, query QueryAssets) (*ListedAssets, error)
	AddUser(ctx context.Context, user User) (*User, error)
	FindUser(ctx context.Context, username string) (*User, error)
	UserExists(ctx context.Context, username string) (bool, error)
	GetUser(ctx context.Context, userID uint) (*User, error)
}
