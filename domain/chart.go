package domain

type Chart struct {
	Id    uint        `json:"id" gorm:"primaryKey"`
	Title string      `json:"title"`
	AxisX string      `json:"axisX"`
	AxisY string      `json:"axisY"`
	Data  []ChartData `json:"data" gorm:"foreignKey:ChartId;constraint:onDelete:CASCADE;"`
}

type ChartData struct {
	Id      uint    `json:"id" gorm:"primaryKey"`
	ChartId uint    `json:"chartId"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
}
