package sqldb

import (
	"encoding/json"
	"platform-go-challenge/domain"
)

func (c *Chart) FromDomain(asset *domain.Chart) {
	dataJson, _ := json.Marshal(asset.Data)
	c.Title = asset.Title
	c.XTitle = asset.XTitle
	c.YTitle = asset.YTitle
	c.Description = asset.Description
	c.Data = dataJson
}

func (c *Chart) ToDomain() *domain.Chart {
	asset := domain.Chart{
		Title:       c.Title,
		XTitle:      c.XTitle,
		YTitle:      c.YTitle,
		Description: c.Description,
	}
	data := domain.XYData{}
	json.Unmarshal(c.Data, &data)
	asset.Data = data
	return &asset
}

func (in *Insight) FromDomain(asset *domain.Insight) {
	in.Text = asset.Text
	in.Description = asset.Description
}

func (in *Insight) ToDomain() *domain.Insight {
	return &domain.Insight{
		Text:        in.Text,
		Description: in.Description,
	}
}

func (au *Audience) FromDomain(asset *domain.Audience) {
	au.AgeMax = asset.AgeMax
	au.AgeMin = asset.AgeMin
	au.Gender = string(asset.Gender)
	au.Country = asset.Country
	au.HoursSpent = asset.HoursSpent
	au.NumberOfPurchases = asset.NumberOfPurchases
	au.Description = asset.Description
}

func (au *Audience) ToDomain() *domain.Audience {
	return &domain.Audience{
		AgeMax:            au.AgeMax,
		AgeMin:            au.AgeMin,
		Gender:            domain.GenderType(au.Gender),
		Country:           au.Country,
		HoursSpent:        au.HoursSpent,
		NumberOfPurchases: au.NumberOfPurchases,
		Description:       au.Description,
	}
}
