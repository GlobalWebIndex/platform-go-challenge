package domain

import "fmt"

type Chart struct {
	AssetBase
	AxesTitles []string
	Data       [][]float64
}

var _ Asset = (*Chart)(nil)

// Validate performs domain validation for the Chart type.
func (c *Chart) Validate() error {
	if err := c.AssetBase.Validate(); err != nil {
		return err
	}

	if len(c.Data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	if len(c.AxesTitles) > 2 {
		return fmt.Errorf("maximum 2 axes titles allowed")
	}

	firstRowLength := len(c.Data[0])
	for i, row := range c.Data {
		if len(row) != firstRowLength {
			return fmt.Errorf("all data rows must have the same length, row %d has different length", i)
		}
	}

	return nil
}
