package models

type User struct {
	Id       string `json:"user_id"`  // user id is actually the keycloak user id
	Username string `json:"username"` // this should be the keycloak username
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Asset[T any] struct {
	Id          int    `json:"asset_id"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Data        T      `json:"data"`
}

type Chart struct {
	XAxisTitle string `json:"X_axis_title"`
	YAxisTitle string `json:"Y_axis_title"`
	Data       string `json:"data"`
}

type Insight struct {
	Text string `json:"text"`
}

type Audience struct {
	Gender             string `json:"gender"`
	BirthCountry       string `json:"birth_country"`
	AgeGroups          string `json:"age_groups"`
	DailyHoursOnSocial int8   `json:"daily_hours_on_social"`
	PurchasesLastMonth int8   `json:"purchases_last_month"`
}

type Favorites struct {
	UsedId  string `json:"user_id"`
	AssetId int    `json:"asset_id"`
}
