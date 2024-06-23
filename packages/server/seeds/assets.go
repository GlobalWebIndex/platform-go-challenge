package seeds

import (
	"errors"
	"strconv"

	"github.com/alexandros-georgantas/platform-go-challenge/internal/models"
	"github.com/bxcodec/faker/v3"

	"github.com/lib/pq"

	"gorm.io/gorm"
)

func SeedAssets(db *gorm.DB) error {
	var a []models.Asset

	db.Find(&a)

	if len(a) == 0 {
		// Creating charts and associate them with assets
		for i := 0; i < 1000; i++ {
			chart := models.Chart{
				Title:               faker.Word(),
				HorizontalAxisLabel: faker.Word(),
				VerticalAxisLabel:   faker.Word(),
				Data:                pq.Float64Array{1.1, 2.2, 3.3},
				Asset: models.Asset{
					Description: faker.Sentence(),
				},
			}

			if err := db.Create(&chart).Error; err != nil {
				return errors.New("something went wrong during seeding phase of charts")
			}
		}

		// Creating insights and associate them with assets
		for i := 0; i < 1000; i++ {
			insight := models.Insight{
				Text: faker.Sentence(),
				Asset: models.Asset{
					Description: faker.Sentence(),
				},
			}

			if err := db.Create(&insight).Error; err != nil {
				return errors.New("something went wrong during seeding phase of insights")
			}
		}

		// Creating audiences and associate them with assets
		for i := 0; i < 1000; i++ {
			i, _ := strconv.Atoi(faker.DayOfMonth())
			j, _ := strconv.Atoi(faker.DayOfMonth())
			audience := models.Audience{
				Gender:                     models.FEMALE,
				CountryOfBirth:             faker.Timezone(),
				AgeGroup:                   "ADULTS",
				DailyHoursOnSocialMedia:    i,
				LastMonthNumberOfPurchases: j,
				Asset: models.Asset{
					Description: faker.Sentence(),
				},
			}

			if err := db.Create(&audience).Error; err != nil {
				return errors.New("something went wrong during seeding phase of audiences")
			}
		}
	}
	return nil
}
