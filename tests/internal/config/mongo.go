package config

import (
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// InitMongoClient returns an initialized mongo client after successful connection with mongo
func InitMongoClient(mongoURI string) *mongo.Client {
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Fatalf("Error in instantiating mongo client %v", err)
	}

	return client
}
