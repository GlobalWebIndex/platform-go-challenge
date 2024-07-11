package domain

type Insight struct {
	Id          uint   `json:"id" gorm:"primaryKey"`
	TopicId     uint   `json:"topicId"`
	Description string `json:"description"`
}
