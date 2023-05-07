package domain

type Chart struct {
	Title       string
	XAxisTitle  string
	YAxisTitle  string
	Data        interface{}
	Description string
	AssetType   string
}

func (c Chart) GetDescription() string {
	return c.Description
}

func (c Chart) GetType() string {
	return string(ChartType)
}
