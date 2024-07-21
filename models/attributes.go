package models

type Gender string

const (
	FEMALE Gender = "Female"
	MALE   Gender = "Male"
	OTHER  Gender = "Other"
)

type AttrType string

const (
	SOCIAL   AttrType = "social"
	PHYSICAL AttrType = "physical"
)

type AgeGroup struct {
	Name  string `json:"name,omitempty"`
	Range Range  `json:"range" validate:"required"`
}

type Attribute struct {
	Name         string  `json:"name" validate:"required"`
	Type         string  `json:"type" validate:"required"`
	Range        Range   `json:"range,omitempty" validate:"required_without_all=ExactInt ExactDecimal"`
	ExactInt     int     `json:"value" validate:"required_without_all=Range ExactDecimal"`
	ExactDecimal float64 `json:"decimal_value" validate:"required_without_all=Range ExactInt"`
	Unit         Unit    `json:"unit_of_measure,omitempty"`
}

/*type AttrType interface {
	Name() string
}

type SocialAttrType string

func (a *SocialAttrType) Name() string {
	return string(*a)
}

type PhysicalAttrType string

func (a *PhysicalAttrType) Name() string {
	return string(*a)
}
*/
