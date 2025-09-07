package entities

import "time"

// AssetType is an enum for the different asset types.
type AssetType string

type Asset interface {
	GetID() string
	GetType() string
	GetDescription() string
	SetDescription(newDescription string)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

type AssetBase struct {
	ID          string    `json:"id" binding:"required"`
	Type        AssetType `json:"type" binding:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Chart struct {
	AssetBase
	Title      string   `json:"title" binding:"required"`
	AxesTitles []string `json:"axesTitles" binding:"required"`
	Data       []int    `json:"data" binding:"required"`
}

func (c *Chart) GetID() string {
	return c.ID
}
func (c *Chart) GetType() string {
	return "chart"
}
func (c *Chart) GetDescription() string {
	return c.Description
}
func (c *Chart) SetDescription(newDescription string) {
	c.Description = newDescription
}
func (c *Chart) GetCreatedAt() time.Time {
	return c.CreatedAt
}
func (c *Chart) GetUpdatedAt() time.Time {
	return c.UpdatedAt
}
func (c *Chart) GetTitle() string {
	return c.Title
}
func (c *Chart) GetAxesTitles() []string {
	return c.AxesTitles
}
func (c *Chart) GetData() []int {
	return c.Data
}

// Insight is a text-based asset.
type Insight struct {
	AssetBase
	Text string `json:"text" binding:"required" `
}

func (i *Insight) GetID() string {
	return i.ID
}
func (i *Insight) GetType() string {
	return "insight"
}
func (i *Insight) GetDescription() string {
	return i.Description
}
func (i *Insight) SetDescription(newDescription string) {
	i.Description = newDescription
}
func (i *Insight) GetCreatedAt() time.Time {
	return i.CreatedAt
}
func (i *Insight) GetUpdatedAt() time.Time {
	return i.UpdatedAt
}
func (i *Insight) GetText() string {
	return i.Text
}

// Audience is an asset with specific user characteristics.
type Audience struct {
	AssetBase
	Gender         string `json:"gender"`
	BirthCountry   string `json:"birthCountry"`
	AgeGroup       string `json:"ageGroup"`
	HoursOnSocial  int    `json:"hoursOnSocial"`
	PurchasesMonth int    `json:"purchasesMonth"`
}

func (a *Audience) GetID() string {
	return a.ID
}
func (a *Audience) GetType() string {
	return "audience"
}
func (a *Audience) GetDescription() string {
	return a.Description
}
func (a *Audience) SetDescription(newDescription string) {
	a.Description = newDescription
}
func (a *Audience) GetCreatedAt() time.Time {
	return a.CreatedAt
}
func (a *Audience) GetUpdatedAt() time.Time {
	return a.UpdatedAt
}
func (a *Audience) GetGender() string {
	return a.Gender
}
func (a *Audience) GetBirthCountry() string {
	return a.BirthCountry
}
func (a *Audience) GetAgeGroup() string {
	return a.AgeGroup
}
func (a *Audience) GetHoursOnSocial() int {
	return a.HoursOnSocial
}
func (a *Audience) GetPurchasesMonth() int {
	return a.PurchasesMonth
}
