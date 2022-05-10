package gwi

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type MysqlRepository struct {
	db *sqlx.DB
}

func NewMysqlRepository(db *sqlx.DB) (*MysqlRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &MysqlRepository{db: db}, nil
}

// sqlContextGetter is an interface provided both by transaction and standard db connection
// type sqlContextGetter interface {
// 	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
// }

func (myr MysqlRepository) Exist(ctx context.Context, userid string, assetid string) (bool, error) {
	return false, nil
}

func (myr MysqlRepository) RetrieveFavourites(ctx context.Context, userid string) ([]Asset, error) {
	return []Asset{}, nil
}

func (myr MysqlRepository) AddAssetToFavourites(ctx context.Context, userid string, asset Asset) (bool, error) {
	return false, nil
}

func (myr MysqlRepository) UpdateFavourite(ctx context.Context, userid string, asset Asset) (bool, error) {
	return false, nil
}

func (myr MysqlRepository) RemoveFavourite(ctx context.Context, userid string, asset Asset) (bool, error) {
	return false, nil
}
