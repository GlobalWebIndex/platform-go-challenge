package entities

import (
	"fmt"
)

type ChartEntity struct {
	AssetBaseEntity
	AxesTitles string `db:"axes_titles"` // JSON serialized
	Data       string `db:"data"`        // JSON serialized
}

// Validate Data Consistency Validation
func (c *ChartEntity) Validate() error {
	if err := c.AssetBaseEntity.Validate(); err != nil {
		return err
	}
	if c.Data != "" {
		if len(c.Data) > 1000 {
			return fmt.Errorf("data string too long")
		}
	}

	if len(c.Data) > 0 {
		for i, r := range c.Data {
			if r == 0 {
				return fmt.Errorf("invalid rune at position %d", i)
			}
		}
	}

	return nil
}
