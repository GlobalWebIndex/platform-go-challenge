package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Chart is one of the available assets
//
// swagger:model Chart
type Chart struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title      string             `json:"title" validate:"required"`
	XAxis      Axis               `json:"x_axis" validate:"required"`
	YAxis      Axis               `json:"y_axis" validate:"required"`
	Points     []Point            `json:"points" validate:"required"`
	CreatedAt  time.Time          `json:"created_at"  bson:"created_at"`
	ModifiedAt time.Time          `json:"modified_at" bson:"modified_at"`
}

func (c Chart) Description() string {
	return c.Title
}
func (c Chart) GetId() primitive.ObjectID {
	return c.ID
}
func (c Chart) GetAssetType() AssetInterface {
	return c
}
