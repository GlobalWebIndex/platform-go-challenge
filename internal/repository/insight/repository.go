package insight

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

func (r *Repository) GetInsights(ids []uuid.UUID) map[uuid.UUID]*model.Insight {
	// example implementation
	// returns a number of random insights with the given id's
	// Then delays 100μs for every item in the ids slice
	results := make(map[uuid.UUID]*model.Insight, len(ids))

	for _, id := range ids {
		results[id] = &model.Insight{
			Id:          id,
			Description: r.generateRandomDescription(),
		}
		time.Sleep(100 * time.Microsecond)
	}
	return results
}

func (r *Repository) generateRandomDescription() string {
	return fmt.Sprintf("%d%% of people between %d and %d spend more than %d hours per day online", rand.IntN(101), 15+rand.IntN(30), 20+rand.IntN(55), rand.IntN(8))
}
