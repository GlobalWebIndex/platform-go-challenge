package domain

import (
	"fmt"
	"strings"
)

type Insight struct {
	AssetBase
	Text string
}

var _ Asset = (*Insight)(nil)

// Domain validation
func (i *Insight) Validate() error {
	if err := i.AssetBase.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(i.Text) == "" {
		return fmt.Errorf("insight text is required")
	}

	return nil
}
