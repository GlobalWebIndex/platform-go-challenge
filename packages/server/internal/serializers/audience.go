package serializers

import "github.com/alexandros-georgantas/platform-go-challenge/internal/models"

type Audience struct {
	ID                         int    `json:"id"`
	Gender                     string `json:"gender"`
	CountryOfBirth             string `json:"countryIfBirth"`
	AgeGroupe                  string `json:"ageGroup"`
	DailyHoursOnSocialMedia    int    `json:"dailyHoursOnSocialMedia"`
	LastMonthNumberOfPurchases int    `json:"lastMonthNumberOfPurchases"`
}

type AudienceResponse struct {
	Audiences  *[]models.Audience `json:"items"`
	TotalCount int                `json:"totalCount"`
}
