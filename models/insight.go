package models

type Insight struct {
	ID int `json:"id"`
	Text string `json:"text"`
	Description string `json:"description"`
}

func (i *Insight) GetID() int {
	return i.ID
}

func (i *Insight) GetType() AssetType {
	return AssetTypeInsight
}

func (i *Insight) GetDescription() string {
	return i.Description
}

func (i *Insight) SetDescription(desc string) {
	i.Description = desc
}