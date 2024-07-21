package database

import (
	"context"
	"errors"
	"favourites/models"
	"favourites/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type FavouriteService interface {
	GetAll(ctx context.Context, role string) ([]models.Favourite, error)
	GetByID(ctx context.Context, id string, role string) (models.Favourite, error)
	Create(ctx context.Context, m *models.Favourite) error
	Update(ctx context.Context, id string, m models.Favourite) error
	Delete(ctx context.Context, id string, role string) error
}

type favouriteService struct {
	C *mongo.Collection
}

var _ FavouriteService = (*favouriteService)(nil)

func NewFavouriteService(collection *mongo.Collection) FavouriteService {
	return &favouriteService{C: collection}
}

func (s *favouriteService) GetAll(ctx context.Context, role string) ([]models.Favourite, error) {
	cur, err := s.C.Aggregate(ctx, bson.A{
		bson.D{{"$match", bson.D{{"role", role}}}},
		bson.D{{"$lookup", bson.D{
			{"from", "insights"},
			{"localField", "asset_id"},
			{"foreignField", "_id"},
			{"as", "insight"},
		}}},
		bson.D{{"$lookup", bson.D{
			{"from", "audiences"},
			{"localField", "asset_id"},
			{"foreignField", "_id"},
			{"as", "audience"},
		}}},
		bson.D{{"$lookup", bson.D{
			{"from", "charts"},
			{"localField", "asset_id"},
			{"foreignField", "_id"},
			{"as", "chart"},
		}}},
		bson.D{{"$project", bson.D{
			{"_id", 1},
			{"title", 1},
			{"description", 1},
			{"asset_type", 1},
			{"asset_id", 1},
			{"role", 1},
			{"created_at", 1},
			{"modified_at", 1},
			{"asset",
				bson.D{{"$cond",
					bson.D{{"if",
						bson.D{{"$ne", bson.A{"$insight", bson.A{}}}}},
						{"then", bson.D{{"$first", "$insight"}}},
						{"else", bson.D{
							{"$cond", bson.D{
								{"if",
									bson.D{{"$ne", bson.A{"$chart", bson.A{}}}}},
								{"then", bson.D{{"$first", "$chart"}}},
								{"else", bson.D{{"$first", "$audience"}}},
							}}}}}}}}}}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.Favourite

	for cur.Next(ctx) {
		if err = cur.Err(); err != nil {
			return nil, err
		}

		//	elem := bson.D{}
		var elem models.Favourite
		err = cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}

	return results, nil
}

func (s *favouriteService) GetByID(ctx context.Context, id string, role string) (models.Favourite, error) {
	var f models.Favourite
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return f, err
	}

	filter := bson.D{{"_id", objectID}, {"role", role}}

	err = s.C.FindOne(context.TODO(), filter).Decode(&f)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return f, errors.New(utils.ErrorNotFound)
	}

	f.EvaluateAssetType()
	fmt.Println(f.GetAssetCollectionByType())

	err = utils.GetDB().Collection(f.GetAssetCollectionByType()).
		FindOne(nil, bson.D{{"_id", f.AssetId}}).Decode(f.Asset)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return f, errors.New(utils.ErrorNotFound)
	}
	return f, nil
}

func (s *favouriteService) Create(ctx context.Context, m *models.Favourite) error {
	m.CreatedAt = time.Now()
	m.ModifiedAt = time.Now()

	_, err := s.C.InsertOne(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func (s *favouriteService) Update(ctx context.Context, id string, m models.Favourite) error {
	filter, err := utils.MatchID(id)
	if err != nil {
		return err
	}
	m.ModifiedAt = time.Now()

	update := bson.D{
		{Key: "$set", Value: m},
	}

	_, err = s.C.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(utils.ErrorNotFound)
		}
		return err
	}

	return nil
}

func (s *favouriteService) Delete(ctx context.Context, id string, role string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objectID}, {"role", role}}
	result, err := s.C.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New(utils.ErrorNotFound)
	}
	return nil
}
