package favorites

import (
	"fmt"
)

// AudienceRules defines the complex rules for an Audience.
type AudienceRules struct {
	Gender  string `json:"gender,omitzero"`
	Country string `json:"country,omitzero"`
	AgeMin  int    `json:"age_min,omitzero"`
	AgeMax  int    `json:"age_max,omitzero"`
}

// Audience represents a target audience asset.
type Audience struct {
	BaseAsset
	Rules AudienceRules `json:"rules"`
}

// Validate implements the Asset interface.
func (a Audience) Validate() error {
	if err := a.ValidateCommon(AssetTypeAudience); err != nil {
		return err
	}
	// Rules validation
	if a.Rules.Country == "" {
		return fmt.Errorf("%w: country rule is required for Audience", ErrValidation)
	}
	if a.Rules.AgeMin < 0 || a.Rules.AgeMax < 0 {
		return fmt.Errorf("%w: age limits cannot be negative", ErrValidation)
	}
	if a.Rules.AgeMax > 0 && a.Rules.AgeMin > a.Rules.AgeMax {
		return fmt.Errorf("%w: age_min cannot be greater than age_max", ErrValidation)
	}
	return nil
}
