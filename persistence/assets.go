package persistence

import (
	"context"
	"errors"
	"log"
	"platform-go-challenge/consts"
	"platform-go-challenge/domain"

	"github.com/arangodb/go-driver"
)

type AssetRepository struct {
	collection driver.Collection
}

func NewAssetRepository(collection driver.Collection) *AssetRepository {
	return &AssetRepository{collection: collection}
}

func (ar AssetRepository) CreateAssets() ([]domain.Asset, error) {
	numOfDocs, err := ar.collection.Count(context.Background())
	if err != nil {
		log.Println("Failed to get document count for asset collection")
	}
	if numOfDocs > 0 {
		return []domain.Asset{}, errors.New("there are already documents in asset collections")
	}

	assets := append(append(consts.Charts, consts.Insights...), consts.Audience...)

	arangoAssets := []domain.Asset{}
	for _, asset := range assets {
		meta, err := ar.collection.CreateDocument(context.Background(), &asset)
		if err != nil {
			log.Printf("Couldn't create new document")
			return []domain.Asset{}, err
		}
		asset.Key = meta.Key
		arangoAssets = append(arangoAssets, asset)
	}
	return arangoAssets, nil
}

func (ar AssetRepository) GetAsset(assetID string) (string, error) {
	asset, err := ar.collection.ReadDocument(context.Background(), assetID, nil)
	if err != nil {
		return "", err
	}
	return string(asset.ID), nil
}
