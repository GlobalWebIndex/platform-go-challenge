package persistence

import (
	"context"
	"log"

	"github.com/arangodb/go-driver"
)

type CollectionRepository struct {
	db driver.Database
}

func NewCollectionRepository(db driver.Database) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (cr CollectionRepository) GetOrCreate(name string, _type string) driver.Collection {
	collection, err := cr.db.Collection(context.Background(), name)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			if _type == "document" {
				collection, err = cr.db.CreateCollection(context.Background(), name, &driver.CreateCollectionOptions{
					Type: driver.CollectionTypeDocument,
				})
				if err != nil {
					log.Fatalf("Failed to create document collection: %v", err)
				}
			} else if _type == "edge" {
				collection, err = cr.db.CreateCollection(context.Background(), name, &driver.CreateCollectionOptions{
					Type: driver.CollectionTypeEdge,
				})
				if err != nil {
					log.Fatalf("Failed to create edge collection: %v", err)
				}
			} else {
				log.Fatalf("Unknown type of collection")
			}
		}
		log.Printf("Collection %s was created successfully", name)
	}
	return collection
}
