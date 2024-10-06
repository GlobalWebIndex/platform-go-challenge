package chart

import (
	"fmt"
	"math/rand/v2"
	"svenvermeulen/platform-go-challenge/pkg/model"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetCharts(ids []uuid.UUID) map[uuid.UUID]*model.Chart {
	// example implementation
	// returns a number of random charts with the given id's
	// This simulates an expensive service call so it will delay relatively long for every item in the ids slice
	results := make(map[uuid.UUID]*model.Chart, len(ids))

	for _, id := range ids {
		resource := []string{"sales", "clicks", "whatever"}[rand.IntN(3)]
		numPoints := rand.IntN(100)
		points := make([]model.DataPoint, numPoints)
		for p := range numPoints {
			points[p].X = float64(p)
			points[p].Y = float64(2 * p)
		}

		results[id] = &model.Chart{
			Id:         id,
			Title:      fmt.Sprintf("%s chart number %d", resource, rand.IntN(10)),
			XAxisTitle: "time",
			YAxisTitle: fmt.Sprintf("number of %s", resource),
			DataPoints: points,
		}
		time.Sleep(50 * time.Millisecond)
	}
	return results
}
