package dto

type BriefProduct struct {
	ChainId        int    `json:"chain_id"`
	AssetId        int32  `json:"asset_id"`
	Owner          string `json:"owner"`
	Barcode        string `json:"bar_code"`
	ItemName       string `json:"item_name"`
	BrandName      string `json:"brand_name"`
	AdditionalData string `json:"additional_data"`
	Location       string `json:"location"`
	IssueDate      string `json:"issue_date"`
}
