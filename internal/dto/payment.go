package dto

type Subscription struct {
	Email          string `json:"email"`
	CustomerId     string `json:"customer_id"`
	PriceId        string `json:"price_id"`
	SubscriptionId string `json:"subscription_id"`
	EndAt          int64  `json:"email_at"`
}
