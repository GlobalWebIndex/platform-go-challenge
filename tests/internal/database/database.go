// Package database is a wrapper for mongoDB functions
package database

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CountDocuments(
	client *mongo.Client,
	databaseName string,
	collection string,
	filter any,
) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(databaseName)
	dbCollection := db.Collection(collection)

	totalDocuments, err := dbCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return totalDocuments, nil 
}

func UpdateOne(
	client *mongo.Client, 
	dbName, collection string, 
	filter bson.D,
	update bson.M,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("db Name = ", dbName)
	log.Println("db Collection = ", collection)

	db := client.Database(dbName)
	dbCollection := db.Collection(collection)


	res, err := dbCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("no matching document found")
	} else if res.ModifiedCount == 0 {
		return errors.New("nothing was changed")
	} else {
		return nil
	}

}

func DeleteOne(
	client *mongo.Client, 
	dbName, collection string, 
	filter bson.D,
	item string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("db Name = ", dbName)
	log.Println("db Collection = ", collection)

	db := client.Database(dbName)
	dbCollection := db.Collection(collection)

	res, err := dbCollection.DeleteOne(ctx, filter)
	if res.DeletedCount == 0 {
		return errors.New("no favourite asset with the combination provided")
	}
	if err != nil {
		return err
	}

	return nil
}

func InsertOne(
	client *mongo.Client,
	dbName, collection string,
	insertElement any, 
	item string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("db = ", dbName)
	log.Println("dbCollection = ", collection)

	db := client.Database(dbName)
	dbCollection := db.Collection(collection)

	if _, err := dbCollection.InsertOne(ctx, insertElement); err != nil {
		log.Printf("Error inserting %v: %v \n", item, err)
		return err
	}

	return nil
}

func FetchUsingFilter(
	client *mongo.Client,
	dbName, collection string,
	filter bson.D,
	result any,
	item string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("db Name = ", dbName)
	log.Println("db Collection = ", collection)

	db := client.Database(dbName)
	dbCollection := db.Collection(collection)

	println("dbCollection struct = ", dbCollection)


	if err := dbCollection.FindOne(ctx, filter).Decode(result); err != nil {
		return err	
	}

	return nil
}

func FetchAssetsByIDs(
	client *mongo.Client,
	dbName, collection string,
	assetIDs []string,
	result any,
) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("db Name = ", dbName)
	log.Println("db Collection = ", collection)

	filter := bson.M{
		"_id": bson.M{"$in": assetIDs},
	}

	db := client.Database(dbName)
	dbCollection := db.Collection(collection)

	println("dbCollection struct = ", dbCollection)

	cursor, err := dbCollection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return err
	}


	return nil
}

func FetchSkipLimit(
	client *mongo.Client,
	dbName, collection string, 
	offset, limit int,
	filter bson.D,
	result any,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(dbName)
	dbCollection := db.Collection(collection)
	var opts = options.Find()

	opts = options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "_id", Value: 1}}) 
	cursor, err := dbCollection.Find(ctx, filter, opts)

	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}















