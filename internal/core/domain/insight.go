package domain

type Insight struct {
	Text        string
	Description string
	AssetType   string
}

func (i Insight) GetDescription() string {
	return i.Description
}

func (i Insight) GetType() string {
	return string(InsightType)
}
