package outbound

type UserPatchRequestBuilder interface {
	Reset()
	WithUserID(userID DataSourceID)
	SetName(name string)
	SetSurname(name string)
	SetAssetDescription(assetID DataSourceID, description string)
	SetAssetIsFavourite(assetID DataSourceID, isFavourite bool)
	Build() PatchData
}
