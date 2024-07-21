package models

type Range struct {
	From float64 `json:"from" validate:"required"`
	To   float64 `json:"to" validate:"required,gtfield=From"`
}

type Axis struct {
	Title string `json:"title"`
	Range Range  `json:"range" validate:"required"`
}
type Point struct {
	X float64 `json:"x" validate:"required"`
	Y float64 `json:"y" validate:"required"`
}

const (
	PERCENTAGE Unit = "%"
	HOUR       Unit = "hour(s)"
	DAY        Unit = "day(s)"
	MINUTE     Unit = "minute(s)"
	YEAR       Unit = "year(s)"
	SECOND     Unit = "second(s)"
)

type Unit string
