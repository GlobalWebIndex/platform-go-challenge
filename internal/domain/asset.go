package domain

type AssetType string

type Asset struct {
    ID          string
    Type        AssetType
    Description string
    Data        interface{}
}
