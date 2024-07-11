package domain

type Favourite struct {
	Id      uint `json:"id" gorm:"primaryKey"`
	UserId  uint `json:"userId"`
	AssetId uint `json:"assetId"`
}
