package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config Mongo Configuration.
type Config struct {
	Database string
	User     string
	Pass     string
	Host     string
	Port     string
}

// New returns a handle to a Mongo client.
func New(ctx context.Context, cfg Config) (*mongo.Database, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/", cfg.User, cfg.Pass, cfg.Host, cfg.Port)
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	return db.Database(cfg.Database), err
}
