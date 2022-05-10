package gwi

type Asseter interface {
	GetId() string
}

type Asset struct {
	AssetType AssetType `json:"type"`
	Chart     Chart     `json:"chart"`
	Insight   Insight   `json:"insight"`
	Audience  Audience  `json:"audience"`
}

type AssetType string

const (
	TypeChart    AssetType = "chart"
	TypeInsight  AssetType = "insight"
	TypeAudience AssetType = "audience"
)

type Chart struct {
	Id         string      `json:"id"`
	Title      string      `json:"title"`
	AxisYTitle string      `json:"axis_y_title"`
	AxisXTitle string      `json:"axis_x_title"`
	Data       interface{} `json:"data"`
}

type Insight struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

type Audience struct {
	Id                    string `json:"id"`
	Gender                string `json:"gender"`
	BornCountry           string `json:"born_country"`
	AgeGroup              string `json:"age_group"`
	DailyHoursSocialMedia string `json:"daily_hours_social_media"`
	PurchasesLastMonth    string `json:"purchases_last_month"`
}

func (c *Chart) GetId() string {
	return c.Id
}

func (i *Insight) GetId() string {
	return i.Id
}

func (a *Audience) GetId() string {
	return a.Id
}
