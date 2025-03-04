package configs

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"fmt"
	"log"
	"time"
)

type Configs struct {
	DbClient *mongo.Client
	DbName   string
}

func NewConfigs() *Configs {
	//TODO get from env params
	dbName := "godb"
	dbUser := "godbuser"
	dbPass := "godbpass"
	dbHost := "db"
	dbPort := "27017"
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=%s", dbUser, dbPass, dbHost, dbPort, dbName)))
	if err != nil {
		log.Fatal(err)
	}
	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return &Configs{DbClient: client, DbName: dbName}
}
