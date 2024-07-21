package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Insight struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Text       string             `json:"text"` // TODO: check possible length for string. Might need to change.
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt time.Time          `json:"modified_at" bson:"modified_at"`
}

func (i Insight) Description() string {
	return i.Text
}
func (i Insight) GetId() primitive.ObjectID {
	return i.ID
}
func (i Insight) GetAssetType() AssetInterface {
	return i
}
