package service

import (
	port "github.com/Kercyn/crud_template/internal/core/port/outbound"
)

type UserPatchRequestBuilder struct {
	data port.PatchData
}

func (b *UserPatchRequestBuilder) Reset() {
	b.data = port.PatchData{}
}

func (b *UserPatchRequestBuilder) WithUserID(userID port.DataSourceID) {
	b.data.ID = userID
}

func (b *UserPatchRequestBuilder) SetName(name string) {
	b.data.Fields["name"] = name
}

func (b *UserPatchRequestBuilder) SetSurname(surname string) {
	b.data.Fields["surname"] = surname
}

func (b *UserPatchRequestBuilder) SetAssetDescription(
	assetID port.DataSourceID,
	description string,
) {
	b.data.Fields["assets"].(map[string]interface{})[assetID.String()].(map[string]interface{})["description"] = description
}

func (b *UserPatchRequestBuilder) SetAssetIsFavourite(
	assetID port.DataSourceID,
	isFavourite bool,
) {
	b.data.Fields["assets"].(map[string]interface{})[assetID.String()].(map[string]interface{})["is_favourite"] = isFavourite
}

func (b *UserPatchRequestBuilder) Build() port.PatchData {
	return b.data
}
