package entities

import (
	"errors"
	"fmt"
)

type AudienceEntity struct {
	AssetBaseEntity
	Gender          string  `db:"gender"`
	BirthCountry    string  `db:"birth_country"`
	AgeGroup        string  `db:"age_group"`
	HoursSocial     float64 `db:"hours_social"`
	PurchasesLastMo int     `db:"purchases_last_mo"`
}

// Validate Data Consistency Validation
func (a *AudienceEntity) Validate() error {
	if err := a.AssetBaseEntity.Validate(); err != nil {
		return err
	}

	if a.Gender != "" && !isValidGender(a.Gender) {
		return fmt.Errorf("invalid gender: %s", a.Gender)
	}
	if a.AgeGroup != "" && !isValidAgeGroup(a.AgeGroup) {
		return fmt.Errorf("invalid age group: %s", a.AgeGroup)
	}
	if a.HoursSocial < 0 {
		return errors.New("hours social cannot be negative")
	}
	if a.PurchasesLastMo < 0 {
		return errors.New("purchases last month cannot be negative")
	}
	return nil
}

func isValidGender(gender string) bool {
	validGenders := []string{"male", "female", "non-binary", "other", ""}
	for _, g := range validGenders {
		if gender == g {
			return true
		}
	}
	return false
}

func isValidAgeGroup(ageGroup string) bool {
	validAgeGroups := []string{"18-24", "25-34", "35-44", "45-54", "55-64", "65+", ""}
	for _, ag := range validAgeGroups {
		if ageGroup == ag {
			return true
		}
	}
	return false
}
