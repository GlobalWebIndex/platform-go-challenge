package domain

type AssetType string

type Asset struct {
	ID          uint64
	CreatedBy   uint64
	Type        AssetType
	Description string
	Data        interface{}
}
