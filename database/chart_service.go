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
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type ChartService interface {
	GetAll(ctx context.Context) ([]models.Chart, error)
	GetByID(ctx context.Context, id string) (models.Chart, error)
	Create(ctx context.Context, m *models.Chart) error
	CreateAll(ctx context.Context, m []*models.Chart) error
	Update(ctx context.Context, id string, m models.Chart) error
	Delete(ctx context.Context, id string) error
}

type chartService struct {
	C *mongo.Collection
}

var _ ChartService = (*chartService)(nil)

func NewChartService(collection *mongo.Collection) ChartService {

	return &chartService{C: collection}
}

func (s *chartService) GetAll(ctx context.Context) ([]models.Chart, error) {
	cur, err := s.C.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.Chart

	for cur.Next(ctx) {
		if err = cur.Err(); err != nil {
			return nil, err
		}

		var elem models.Chart
		err = cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}
	return results, nil
}

func (s *chartService) GetByID(ctx context.Context, id string) (models.Chart, error) {
	var chart models.Chart
	filter, err := utils.MatchID(id)
	if err != nil {
		return chart, err
	}

	err = s.C.FindOne(ctx, filter).Decode(&chart)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return chart, errors.New(utils.ErrorNotFound)
	}
	return chart, err
}

func (s *chartService) Create(ctx context.Context, m *models.Chart) error {

	m.CreatedAt = time.Now()
	m.ModifiedAt = time.Now()

	_, err := s.C.InsertOne(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func (s *chartService) CreateAll(ctx context.Context, charts []*models.Chart) error {

	var chartsI []interface{}
	for _, i := range charts {
		i.CreatedAt = time.Now()
		i.ModifiedAt = time.Now()

		chartsI = append(chartsI, i)
	}
	_, err := s.C.InsertMany(context.TODO(), chartsI)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (s *chartService) Update(ctx context.Context, id string, m models.Chart) error {
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

func (s *chartService) Delete(ctx context.Context, id string) error {
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
