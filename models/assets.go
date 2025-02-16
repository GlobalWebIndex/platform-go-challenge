package models

type AssetInterface interface {
	GetID() string
	GetType() AssetType
	GetDescription() string
}

type AssetType string

const (
	ChartType    AssetType = "chart"
	InsightType  AssetType = "insight"
	AudienceType AssetType = "audience"
)

type Asset struct {
	ID          string    `json:"id"`
	Type        AssetType `json:"type"`
	Description string    `json:"description"`
}

func (a Asset) GetID() string          { return a.ID }
func (a Asset) GetType() AssetType     { return a.Type }
func (a Asset) GetDescription() string { return a.Description }

type Chart struct {
	Asset
	Title      string   `json:"title"`
	AxesTitles []string `json:"axes_titles"`
	Data       []int    `json:"data"`
}

type Insight struct {
	Asset
}

type Audience struct {
	Asset
	Gender                       string `json:"gender"`
	BirthCountry                 string `json:"birth_country"`
	AgeGroup                     string `json:"age_group"`
	HoursSpentDailyOnSocialMedia int    `json:"hours_spent_daily_on_social_media"`
	PurchasesLastMonth           int    `json:"purchases_last_month"`
}

func (c Chart) GetID() string          { return c.ID }
func (c Chart) GetType() AssetType     { return c.Type }
func (c Chart) GetDescription() string { return c.Description }

func (i Insight) GetID() string          { return i.ID }
func (i Insight) GetType() AssetType     { return i.Type }
func (i Insight) GetDescription() string { return i.Description }

func (a Audience) GetID() string          { return a.ID }
func (a Audience) GetType() AssetType     { return a.Type }
func (a Audience) GetDescription() string { return a.Description }
