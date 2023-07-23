package dto

import "gwi_api/internal/domain"

type AssetDto struct {
	CreatedBy   uint64
	Type        domain.AssetType
	Description string
	Data        interface{}
}
