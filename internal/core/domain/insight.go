package domain

type Insight struct {
	Asset
	Data string
}

func (i Insight) GetType() AssetType {
	return i.Asset.GetType()
}

func (i Insight) GetDescription() string {
	return i.Asset.GetDescription()
}

func (i Insight) GetIsFavourite() bool {
	return i.Asset.GetIsFavourite()
}
