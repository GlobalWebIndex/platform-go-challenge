package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/giannis84/platform-go-challenge/internal/models"
)

const maxStringLength = 255

// AssetType represents the type of asset being favourited.
type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

// AddFavouriteRequest is the request payload for adding a favourite.
type AddFavouriteRequest struct {
	AssetType   AssetType       `json:"asset_type"`
	Description string          `json:"description"`
	AssetData   json.RawMessage `json:"asset_data"`
}

// UpdateDescriptionRequest is the request payload for updating a favourite's description.
type UpdateDescriptionRequest struct {
	Description string `json:"description"`
}

// ParseAddFavouriteRequest validates the request and returns the parsed asset.
// It handles asset type validation and type-specific unmarshaling.
func ParseAddFavouriteRequest(req *AddFavouriteRequest) (models.Asset, error) {
	switch req.AssetType {
	case AssetTypeChart:
		var chart models.Chart
		if err := json.Unmarshal(req.AssetData, &chart); err != nil {
			return nil, fmt.Errorf("invalid chart data: %w", err)
		}
		return &chart, nil
	case AssetTypeInsight:
		var insight models.Insight
		if err := json.Unmarshal(req.AssetData, &insight); err != nil {
			return nil, fmt.Errorf("invalid insight data: %w", err)
		}
		return &insight, nil
	case AssetTypeAudience:
		var audience models.Audience
		if err := json.Unmarshal(req.AssetData, &audience); err != nil {
			return nil, fmt.Errorf("invalid audience data: %w", err)
		}
		return &audience, nil
	default:
		return nil, fmt.Errorf("invalid asset_type: %q", req.AssetType)
	}
}

// ValidateAssetID validates that an asset ID is not empty.
func ValidateAssetID(assetID string) error {
	if strings.TrimSpace(assetID) == "" {
		return &ValidationError{Errors: []string{"asset_id is required"}}
	}
	return nil
}

// IsInvalidAssetType returns true if the error is due to invalid asset type.
func IsInvalidAssetType(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), "invalid asset_type:")
}

// IsInvalidAssetData returns true if the error is due to invalid asset data.
func IsInvalidAssetData(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), "invalid ")
}

var (
	validGenders          = []string{"Male", "Female"}
	validAgeGroups        = []string{"18-24", "25-34", "35-44", "45-54", "55+"}
	validSocialMediaHours = []string{"0-1", "1-3", "3-5", "5+"}
)

// ValidationError holds a list of field-level validation errors.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(e.Errors, "; "))
}

// validate collects errors and returns a *ValidationError if any exist.
func validate(checks ...func() string) error {
	var errs []string
	for _, check := range checks {
		if msg := check(); msg != "" {
			errs = append(errs, msg)
		}
	}
	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}

func requireNonEmpty(field, value string) string {
	if strings.TrimSpace(value) == "" {
		return fmt.Sprintf("%s is required", field)
	}
	return ""
}

func checkMaxLength(field, value string, max int) string {
	if len(value) > max {
		return fmt.Sprintf("%s exceeds maximum length of %d", field, max)
	}
	return ""
}

func checkInList(field, value string, allowed []string) string {
	for _, v := range allowed {
		if value == v {
			return ""
		}
	}
	return fmt.Sprintf("%s has invalid value %q (allowed: %s)", field, value, strings.Join(allowed, ", "))
}

func checkNonNegative(field string, value int) string {
	if value < 0 {
		return fmt.Sprintf("%s must not be negative", field)
	}
	return ""
}

// validateChart validates required fields and length constraints for a Chart asset.
func validateChart(c *models.Chart) error {
	return validate(
		func() string { return requireNonEmpty("id", c.ID) },
		func() string { return checkMaxLength("id", c.ID, maxStringLength) },
		func() string { return requireNonEmpty("title", c.Title) },
		func() string { return checkMaxLength("title", c.Title, maxStringLength) },
		func() string { return requireNonEmpty("x_axis_title", c.XAxisTitle) },
		func() string { return checkMaxLength("x_axis_title", c.XAxisTitle, maxStringLength) },
		func() string { return requireNonEmpty("y_axis_title", c.YAxisTitle) },
		func() string { return checkMaxLength("y_axis_title", c.YAxisTitle, maxStringLength) },
	)
}

// validateInsight validates required fields and length constraints for an Insight asset.
func validateInsight(i *models.Insight) error {
	return validate(
		func() string { return requireNonEmpty("id", i.ID) },
		func() string { return checkMaxLength("id", i.ID, maxStringLength) },
		func() string { return requireNonEmpty("text", i.Text) },
		func() string { return checkMaxLength("text", i.Text, maxStringLength) },
	)
}

// validateAudience validates an Audience asset. Only ID is required;
// other fields are optional but validated when provided.
func validateAudience(a *models.Audience) error {
	checks := []func() string{
		func() string { return requireNonEmpty("id", a.ID) },
		func() string { return checkMaxLength("id", a.ID, maxStringLength) },
		func() string { return checkNonNegative("purchases_last_month", a.PurchasesLastMonth) },
	}

	for i, g := range a.Gender {
		g := g
		i := i
		checks = append(checks, func() string {
			return checkInList(fmt.Sprintf("gender[%d]", i), g, validGenders)
		})
	}

	for i, c := range a.BirthCountry {
		c := c
		i := i
		checks = append(checks, func() string {
			return requireNonEmpty(fmt.Sprintf("birth_country[%d]", i), c)
		})
	}

	for i, ag := range a.AgeGroups {
		ag := ag
		i := i
		checks = append(checks, func() string {
			return checkInList(fmt.Sprintf("age_groups[%d]", i), ag, validAgeGroups)
		})
	}

	if a.SocialMediaHoursDaily != "" {
		checks = append(checks, func() string {
			return checkInList("social_media_hours_daily", a.SocialMediaHoursDaily, validSocialMediaHours)
		})
	}

	return validate(checks...)
}

// validateDescription validates the description field on update requests.
func validateDescription(description string) error {
	return validate(
		func() string { return requireNonEmpty("description", description) },
		func() string { return checkMaxLength("description", description, maxStringLength) },
	)
}
