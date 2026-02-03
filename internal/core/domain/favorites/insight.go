package favorites

import (
	"fmt"
)

// Insight represents a textual insight asset.
type Insight struct {
	BaseAsset
	Content string `json:"content,omitzero"`
}

// Validate implements the Asset interface.
func (i Insight) Validate() error {
	if err := i.ValidateCommon(AssetTypeInsight); err != nil {
		return err
	}
	if i.Content == "" {
		return fmt.Errorf("%w: content is required for Insight", ErrValidation)
	}
	return nil
}
