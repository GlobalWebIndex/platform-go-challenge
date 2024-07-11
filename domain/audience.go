package domain

type Gender string

type Audience struct {
	Id                 uint   `json:"id" gorm:"primaryKey"`
	Gender             Gender `json:"gender"`
	Country            string `json:"country"`
	Age                uint   `json:"age"`
	SocialMediaTime    uint   `json:"s_media_time"` // time spent on social media in minutes
	LastMonthPurchases uint   `json:"last_month_purchases"`
}

const (
	Male   Gender = "male"
	Female Gender = "female"
	Other  Gender = "other"
)
