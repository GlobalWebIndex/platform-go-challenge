package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Audience struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" validate:"required"`
	Gender       Gender             `json:"gender,omitempty"`
	BirthCountry string             `json:"birth_country,omitempty" bson:"birth_country,omitempty" validate:"country_code"`
	AgeGroups    Range              `json:"age_groups,omitempty" bson:"age_groups,omitempty"`
	Attributes   []Attribute        `json:"social_commonalities,omitempty" bson:"social_commonalities,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt   time.Time          `json:"modified_at" bson:"modified_at"`
}

func (a Audience) Description() string {
	return a.Name
}
func (a Audience) GetId() primitive.ObjectID {
	return a.ID
}
func (a Audience) GetAssetType() AssetInterface {
	return a
}
