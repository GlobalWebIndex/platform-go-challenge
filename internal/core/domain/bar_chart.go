package domain

type BarData struct {
	Label string
	Value float64
}

type BarChart struct {
	Asset
	Title      string
	YAxisTitle string
	Data       []BarData
}

func (b BarChart) GetType() AssetType {
	return b.Asset.GetType()
}

func (b BarChart) GetDescription() string {
	return b.Asset.GetDescription()
}

func (b BarChart) GetIsFavourite() bool {
	return b.Asset.GetIsFavourite()
}
