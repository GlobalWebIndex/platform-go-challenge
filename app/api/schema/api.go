package schema

// import "errors"
// import "github.com/xdg-go/stringprep"

// "reflect"
// "fmt"

const (
	CHARACTERISTIC_AGE_GROUP                = "age group"
	CHARACTERISTIC_BIRTH_COUNTRY            = "birth country"
	CHARACTERISTIC_GENDER                   = "gender"
	CHARACTERISTIC_PURCHASES_LAST_MONTH     = "purchases made last month"
	CHARACTERISTIC_SOCIAL_MEDIA_DAILY_HOURS = "hours spent on social media per day" // todo: should be treated like age groups
	GENDER_FEMALE                           = "female"
	GENDER_MALE                             = "male"
)

type AssetUpdateRequest struct {
	Description string `json:"description" binding:"required"`
}

type UserFavoriteRequest struct {
	AssetID string `json:"assetId" binding:"required"`
}

type UserDto struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type AssetDto struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type ChartAssetDto struct {
	Title      string `json:"title"`
	AxesTitles string `json:"axesTitles"`
	PlotData   string `json:"plotData"`
}

type InsightAssetDto struct {
	Text string `json:"text"`
}

type AudienceAssetDto struct {
	Characteristics []string `json:"characteristics"`
}

type AssetDetailsDto struct {
	Asset           AssetDto          `json:"asset"`
	ChartDetails    *ChartAssetDto    `json:"chart,omitempty"`
	InsightDetails  *InsightAssetDto  `json:"insight,omitempty"`
	AudienceDetails *AudienceAssetDto `json:"audience,omitempty"`
}

type UserFavoriteDto struct {
	ID      string          `json:"id"`
	Details AssetDetailsDto `json:"details"`
}

/* error */
type HttpError interface {
	Error() string
	Status() int
}

type ApiError struct {
	Code   int   `json:"code"`
	Reason error `json:"reason"`
}

func NewApiError(code int, reason error) error {
	return &ApiError{
		Code:   code,
		Reason: reason,
	}
}

func (e *ApiError) Error() string {
	if e.Reason == nil {
		return "something went wrong"
	}
	return e.Reason.Error()
}

func (e *ApiError) Status() int {
	return e.Code
}

type ValidationError struct {
	Key string
}

/* internal */

type AssetFilter struct {
	AssetID string
}