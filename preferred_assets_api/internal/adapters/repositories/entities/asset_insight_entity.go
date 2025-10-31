package entities

import (
	"fmt"
	"unicode/utf8"
)

type InsightEntity struct {
	AssetBaseEntity
	Text string `db:"text"`
}

// Validate Data Consistency Validation
func (i *InsightEntity) Validate() error {
	if err := i.AssetBaseEntity.Validate(); err != nil {
		return err
	}

	if utf8.RuneCountInString(i.Text) > 1000 {
		return fmt.Errorf("insight text cannot exceed 1000 characters")
	}
	return nil
}
