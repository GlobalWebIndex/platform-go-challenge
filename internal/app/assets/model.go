package assets

type Assets struct {
	Charts    Charts    `json:"charts"`
	Insights  Insights  `json:"insights"`
	Audiences Audiences `json:"audience"`
}

type Chart struct {
	ID         uint32                 `json:"id"`
	Title      string                 `json:"title"`
	AxisXTitle string                 `json:"axis_x_title"`
	AxisYTitle string                 `json:"axis_y_title"`
	Data       map[string]interface{} `json:"data"`
}

type Charts []Chart

type Insight struct {
	ID    uint32 `json:"id"`
	Title string `json:"title"`
}

type Insights []Insight

type Audience struct {
	ID                         uint32       `json:"id"`
	Gender                     GenderOption `json:"gender"`
	BirthCountry               string       `json:"birth_country"`
	AgeGroup                   AgeGroup     `json:"age_group"`
	HoursSpentOnline           uint32       `json:"hours_spent_online"`
	NumberOfPerchasesLastMonth uint32       `json:"number_of_perchases_last_month"`
}

type Audiences []Audience

type GenderOption string

var (
	GenderMale   = GenderOption("male")
	GenderFemale = GenderOption("female")
	GenderOther  = GenderOption("other")
)

type AgeGroup string

var (
	AgeGroupTeenager   = AgeGroup("teenagers")
	AgeGroupYoungAdult = AgeGroup("young-adults")
	AgeGroupMiddleAged = AgeGroup("middle-aged")
	AgeGroupSenior     = AgeGroup("seniors")
)

type AssetType string

var (
	AssetTypeChart    = AssetType("chart")
	AssetTypeInsight  = AssetType("insight")
	AssetTypeAudience = AssetType("audience")

	AssetTypes = map[AssetType]struct{}{
		AssetTypeChart:    {},
		AssetTypeInsight:  {},
		AssetTypeAudience: {},
	}
)
