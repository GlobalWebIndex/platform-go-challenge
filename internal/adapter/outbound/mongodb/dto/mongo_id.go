package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoID struct {
	data primitive.ObjectID
}

func NewMongoID(id string) (MongoID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return MongoID{}, err
	}

	return MongoID{
		data: objectID,
	}, nil
}

func (m MongoID) String() string {
	return m.data.String()
}

func (m MongoID) AsDataSourceID() any {
	return m.data
}
