package model

import (
	"encoding/json"
	"errors"
	"strings"
)

var ErrInvalidPayload = errors.New("invalid payload")

func ValidateAsset(a Asset) error {
	if a.ID == "" || a.Type == "" || len(a.Payload) == 0 {
		return ErrInvalidPayload
	}

	switch a.Type {
	case AssetTypeInsight:
		var p InsightPayload
		if err := json.Unmarshal(a.Payload, &p); err != nil {
			return ErrInvalidPayload
		}
		if strings.TrimSpace(p.Text) == "" {
			return ErrInvalidPayload
		}
		return nil

	case AssetTypeChart:
		var p ChartPayload
		if err := json.Unmarshal(a.Payload, &p); err != nil {
			return ErrInvalidPayload
		}
		if strings.TrimSpace(p.Title) == "" {
			return ErrInvalidPayload
		}
		if strings.TrimSpace(p.Axes.X) == "" || strings.TrimSpace(p.Axes.Y) == "" {
			return ErrInvalidPayload
		}
		// Data can be anything; just require it exists (optional)
		return nil

	case AssetTypeAudience:
		var p AudiencePayload
		if err := json.Unmarshal(a.Payload, &p); err != nil {
			return ErrInvalidPayload
		}
		// Require at least one dimension present
		if len(p.Gender)+len(p.BirthCountry)+len(p.AgeGroups)+len(p.HoursSocialMedia)+len(p.Purchases) == 0 {
			return ErrInvalidPayload
		}
		return nil

	default:
		return ErrInvalidPayload
	}
}
