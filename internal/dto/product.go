package dto

import (
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

type BriefProduct struct {
	ChainId        int    `json:"chain_id"`
	AssetId        int64  `json:"asset_id"`
	Owner          string `json:"owner"`
	Barcode        string `json:"bar_code"`
	ItemName       string `json:"item_name"`
	BrandName      string `json:"brand_name"`
	AdditionalData string `json:"additional_data"`
	Location       string `json:"location"`
	IssueDate      int32  `json:"issue_date"`
}

func (p *BriefProduct) Valid() bool {
	_, err := types.DecodeAddress(p.Owner)
	return !(p.AssetId == 0 || strings.TrimSpace(p.Owner) == "" || strings.TrimSpace(p.Barcode) == "" || err != nil)
}
