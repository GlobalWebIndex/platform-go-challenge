package assets

type Chart struct {
	ID         uint32      `json:"id"`
	Title      string      `json:"title"`
	AxisXTitle string      `json:"axis_x_title"`
	AxisYTitle string      `json:"axis_y_title"`
	Data       interface{} `json:"data"`
}

type Insight struct {
	ID    uint32 `json:"id"`
	Title string `json:"title"`
}

type Audience struct {
	Gender                     GenderOption `json:"gender"`
	BirthCountry               string       `json:"birth_country"`
	AgeGroup                   AgeGroup     `json:"age_group"`
	HoursSpentOnline           uint32       `json:"hours_spent_online"`
	NumberOfPerchasesLastMonth uint32       `json:"number_of_perchases_last_month"`
}

type GenderOption string

var (
	GenderMale   = GenderOption("male")
	GenderFemale = GenderOption("female")
	GenderOther  = GenderOption("other")
)

type AgeGroup string

var (
	AgeGroupTeenager   = AgeGroup("teenager")
	AgeGroupYoungAdult = AgeGroup("young-adult")
	AgeGroupMiddleAged = AgeGroup("middle-aged")
	AgeGroupSenior     = AgeGroup("senior")
)

type AssetType string

var (
	AssetTypeChart    = AssetType("chart")
	AssetTypeInsight  = AssetType("insight")
	AssetTypeAudience = AssetType("audience")
)
