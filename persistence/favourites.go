package persistence

import (
	"context"
	"log"
	"platform-go-challenge/domain"
	"time"

	"github.com/arangodb/go-driver"
)

type FavouriteRepository struct {
	db         driver.Database
	collection map[string]driver.Collection
}

func NewFavouriteRepository(db driver.Database, collection map[string]driver.Collection) *FavouriteRepository {
	return &FavouriteRepository{db: db, collection: collection}
}

func (fr FavouriteRepository) AddToFavourites(userID, assetID string) (domain.Favourite, error) {
	favourite := domain.Favourite{
		From:      userID,
		To:        assetID,
		CreatedAt: time.Now(),
	}
	_, err := fr.collection["favourite"].CreateDocument(context.Background(), &favourite)
	if err != nil {
		log.Println("Couldn't create new document to favourites")
		return domain.Favourite{}, err
	}
	return favourite, nil
}

func (fr FavouriteRepository) GetFavourites(userID string) ([]domain.Asset, error) {
	bindVars := map[string]interface{}{
		"userID": userID,
	}

	query := `FOR v, e IN 1..1 OUTBOUND @userID favourite RETURN v`
	cursor, err := fr.db.Query(context.Background(), query, bindVars)
	if err != nil {
		log.Printf("Failed to execute query")
		return []domain.Asset{}, err
	}
	defer cursor.Close()

	var assets []domain.Asset
	for {
		var asset domain.Asset
		meta, err := cursor.ReadDocument(context.Background(), &asset)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Printf("Failed to read document")
			return []domain.Asset{}, err
		}

		asset = domain.Asset{
			Key:         meta.Key,
			Type:        asset.Type,
			Description: asset.Description,
			CreatedAt:   asset.CreatedAt,
			Data:        asset.Data,
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

func (fr FavouriteRepository) EditFavourites(userID, assetID, description string) (domain.Asset, error) {
	query := `FOR f IN favourite
	          FILTER f._from==@userID AND f._to==@assetID
			  LET asset=DOCUMENT(f._to)
              RETURN asset`
	bindVars := map[string]interface{}{
		"userID":  userID,
		"assetID": assetID,
	}
	cursor, err := fr.db.Query(context.Background(), query, bindVars)
	if err != nil {
		log.Printf("Failed to execute query")
		return domain.Asset{}, err
	}
	defer cursor.Close()

	var asset domain.Asset
	meta, err := cursor.ReadDocument(context.Background(), &asset)
	if err != nil {
		log.Printf("Failed to read document")
		return domain.Asset{}, err
	}

	asset.Description = description

	_, err = fr.collection["asset"].UpdateDocument(context.Background(), meta.Key, &asset)
	if err != nil {
		log.Println("Couldn't update document")
		return domain.Asset{}, err
	}
	return asset, nil
}

func (fr FavouriteRepository) DeleteFavourite(userID, assetID string) (domain.Favourite, error) {
	query := `FOR f IN favourite
	          FILTER f._from==@userID AND f._to==@assetID
			  RETURN f`
	bindVars := map[string]interface{}{
		"userID":  userID,
		"assetID": assetID,
	}
	cursor, err := fr.db.Query(context.Background(), query, bindVars)
	if err != nil {
		log.Printf("Failed to execute query")
		return domain.Favourite{}, err
	}
	defer cursor.Close()

	var favourite domain.Favourite
	meta, err := cursor.ReadDocument(context.Background(), &favourite)
	if err != nil {
		log.Printf("Failed to read document")
		return domain.Favourite{}, err
	}

	_, err = fr.collection["favourite"].RemoveDocument(context.Background(), meta.Key)
	if err != nil {
		log.Println("Couldn't delete document from favourites")
		return domain.Favourite{}, err
	}
	return favourite, nil
}
