package domain

type Audience struct {
	AssetBase
	Gender          string
	BirthCountry    string
	AgeGroup        string
	HoursSocial     float64
	PurchasesLastMo int
}

var _ Asset = (*Audience)(nil)

// Domain validation
func (a *Audience) Validate() error {
	if err := a.AssetBase.Validate(); err != nil {
		return err
	}
	return nil
}
