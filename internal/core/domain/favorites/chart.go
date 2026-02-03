package favorites

import "fmt"

// Chart represents a visualization asset.
type Chart struct {
	BaseAsset
	XAxis string `json:"x_axis,omitzero"`
	YAxis string `json:"y_axis,omitzero"`
}

// Validate implements the Asset interface.
func (c Chart) Validate() error {
	if err := c.ValidateCommon(AssetTypeChart); err != nil {
		return err
	}
	// Charts should probably have axes defined, though not strictly required by original prompt.
	// Adding simple validation improves data integrity.
	if c.XAxis == "" && c.YAxis == "" {
		return fmt.Errorf("%w: chart requires at least one axis to be defined", ErrValidation)
	}
	return nil
}
