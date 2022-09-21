package httpapi

import (
	"platform-go-challenge/domain"

	"github.com/golang-jwt/jwt"
)

type StatusType string

const (
	FailureStatus = "failure"
	SuccessStatus = "success"
)

const (
	AssetTypeCharts    = "charts"
	AssetTypeInsights  = "insights"
	AssetTypeAudiences = "audiences"
)

type JwtUserClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

type RequestUserCreation struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

type RequestUserLogin struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	ExpiresInMinutes int    `json:"expiresInMinutes"`
}

type ResponseStatus struct {
	Status StatusType `json:"status"`
	Error  string     `json:"error,omitempty"`
}

type ResponseLogin struct {
	Status    StatusType `json:"status"`
	Error     *error     `json:"error,omitempty"`
	Token     *string    `json:"token,omitempty"`
	Username  *string    `json:"username,omitempty"`
	ExpiresAt *int64     `json:"expiresAt,omitempty"`
}

type UserJson struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
}

type QueryAssets struct {
	domain.QueryAssets
	Who *domain.QueryFavouriteAssets `json:"who"`
}
