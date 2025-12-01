package models

type Chart struct {
	ID int `json:"id"`
	Title string `json:"title"`
	XAxis string `json:"x_axis"`
	YAxis string `json:"y_axis"`
	Data []Point `json:"data"`
	Description string `json:"description"`
}

func (c *Chart) GetID() int {
	return c.ID
}

func (c *Chart) GetType() AssetType {
	return AssetTypeChart
}

func (c *Chart) GetDescription() string {
	return c.Description
}

func (c *Chart) SetDescription(desc string) {
	c.Description = desc
}