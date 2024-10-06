package audience

import (
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

func (r *Repository) GetAudiences(ids []uuid.UUID) map[uuid.UUID]*model.Audience {
	// example implementation
	// returns a subset of a number of random-generated audiences, with the provided ids

	// Does a fixed small delay. This is slow-changing data and would probably
	// be served from a in-memory cache anyway.
	results := make(map[uuid.UUID]*model.Audience, len(ids))

	for _, id := range ids {
		ageFrom := 15 + rand.IntN(15)
		ageTo := ageFrom + rand.IntN(60)

		hoursFrom := rand.IntN(5)
		hoursTo := hoursFrom + rand.IntN(8)

		purchasesFrom := rand.IntN(3)
		purchasesTo := rand.IntN(5)

		results[id] = &model.Audience{
			Id:             id,
			Gender:         []rune{'m', 'f', '?'}[rand.IntN(3)],
			BirthCountry:   []string{"United Kingdom", "Greece", "Netherlands"}[rand.IntN(3)],
			AgeFrom:        ageFrom,
			AgeTo:          ageTo,
			HoursSpentFrom: hoursFrom,
			HoursSpentTo:   hoursTo,
			PurchasesFrom:  purchasesFrom,
			PurchasesTo:    purchasesTo,
		}
	}
	time.Sleep(100 * time.Microsecond)
	return results
}
