package database

import (
	"context"
	"errors"
	"favourites/models"
	"favourites/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	// "go.mongodb.org/mongo-driver/mongo/options" // TODO
)

type InsightService interface {
	GetAll(ctx context.Context) ([]models.Insight, error)
	GetByID(ctx context.Context, id string) (models.Insight, error)
	Create(ctx context.Context, m *models.Insight) error
	CreateAll(ctx context.Context, result []*models.Insight) error
	Update(ctx context.Context, id string, m models.Insight) error
	Delete(ctx context.Context, id string) error
}

type insightService struct {
	C *mongo.Collection
}

var _ InsightService = (*insightService)(nil)

func NewInsightService(collection *mongo.Collection) InsightService {

	return &insightService{C: collection}
}

func (s *insightService) GetAll(ctx context.Context) ([]models.Insight, error) {
	cur, err := s.C.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.Insight

	for cur.Next(ctx) {
		if err = cur.Err(); err != nil {
			return nil, err
		}

		var elem models.Insight
		err = cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}

	return results, nil
}

func (s *insightService) GetByID(ctx context.Context, id string) (models.Insight, error) {
	var insight models.Insight
	filter, err := utils.MatchID(id)
	if err != nil {
		return insight, err
	}

	err = s.C.FindOne(ctx, filter).Decode(&insight)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return insight, errors.New(utils.ErrorNotFound)
	}
	return insight, err
}

func (s *insightService) Create(ctx context.Context, m *models.Insight) error {
	m.CreatedAt = time.Now()
	m.ModifiedAt = time.Now()
	_, err := s.C.InsertOne(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func (s *insightService) CreateAll(ctx context.Context, insights []*models.Insight) error {

	var insightsI []interface{}
	for _, i := range insights {
		i.CreatedAt = time.Now()
		i.ModifiedAt = time.Now()

		insightsI = append(insightsI, i)
	}
	_, err := s.C.InsertMany(ctx, insightsI)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (s *insightService) Update(ctx context.Context, id string, m models.Insight) error {
	filter, err := utils.MatchID(id)
	if err != nil {
		return err
	}

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

func (s *insightService) Delete(ctx context.Context, id string) error {
	filter, err := utils.MatchID(id)
	if err != nil {
		return err
	}
	_, err = s.C.DeleteOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(utils.ErrorNotFound)
		}
		return err
	}

	return nil
}
