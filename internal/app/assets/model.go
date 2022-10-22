package assets

type Chart struct {
	ID         uint32 `json:"id"`
	Title      string `json:"title"`
	AxisXTitle string `json:"axis_x_title"`
	AxisYTitle string `json:"axis_y_title"`
	Data       []byte `json:"data"`
}

type Insight struct {
	ID    uint32 `json:"id"`
	Title string `json:"title"`
}

type Audience struct {
	ID                         uint32       `json:"id"`
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
)
