package dto

import (
	"github.com/Kercyn/crud_template/internal/core/domain"
	"time"
)

type Audience struct {
	AssetBase
	Gender                  string  `bson:"gender,omitempty"`
	AgeGroup                string  `bson:"age_group,omitempty"`
	NumOfPurchasesPerMonth  int     `bson:"num_of_purchases_per_month,omitempty"`
	HoursSpentDailyOnSocial float64 `bson:"hours_spent_daily_on_social,omitempty"`
}

func (a Audience) ToDomain() (domain.Asset, error) {
	domainAsset, err := a.AssetBase.ToDomain()
	if err != nil {
		return domain.Audience{}, err
	}

	//todo gender and agegroup from string to their type
	return domain.Audience{
		Asset:                     domainAsset,
		Gender:                    domain.GenderMale,
		AgeGroup:                  domain.AgeGroup19to24,
		NumberOfPurchasesPerMonth: a.NumOfPurchasesPerMonth,
		HoursSpendDailyOnSocial:   time.Duration(a.HoursSpentDailyOnSocial * float64(time.Hour)),
	}, nil
}
