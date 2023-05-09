package repositories

import (
	"context"
	"errors"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"sync"
)

func (repo *UserRepository) GetUser(
	ctx context.Context,
	uid uint32,
) (*domain.User, error) {
	var err error
	var modelUser *User
	var mutex sync.Mutex
	var wg sync.WaitGroup

	err = repo.db.WithContext(ctx).
		Preload("FavouriteCharts").
		Preload("FavouriteAudiences").
		Preload("FavouriteInsights").
		Model(User{}).
		Where("id = ?", uid).
		Take(&modelUser).Error

	if err == gorm.ErrRecordNotFound {
		return &domain.User{}, apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(uid)) + " not found"),
		}
	}

	if err != nil {
		return &domain.User{}, err
	}

	var favouriteAssets []domain.Asset

	wg.Add(1)
	go func(
		charts *[]Chart,
		favouriteAssets *[]domain.Asset,
		mutex *sync.Mutex,
	) {
		defer wg.Done()

		for _, chart := range modelUser.FavouriteCharts {
			mutex.Lock()
			*favouriteAssets = append(*favouriteAssets, domain.Chart{
				Title:       chart.Title,
				XAxisTitle:  chart.XTitle,
				YAxisTitle:  chart.YTitle,
				Data:        chart.Data,
				Description: chart.Description,
				AssetType:   "chart",
			})
			mutex.Unlock()
		}
	}(&modelUser.FavouriteCharts, &favouriteAssets, &mutex)

	wg.Add(1)
	go func(
		audiences *[]Audience,
		favouriteAssets *[]domain.Asset,
		mutex *sync.Mutex,
	) {
		defer wg.Done()

		for _, audience := range modelUser.FavouriteAudiences {
			mutex.Lock()
			*favouriteAssets = append(*favouriteAssets, domain.Audience{
				Gender:             audience.Gender,
				BirthCountry:       audience.BirthCountry,
				AgeGroup:           audience.AgeGroup,
				HoursSpentDaily:    audience.HoursSpentOnSocial,
				PurchasesLastMonth: audience.PurchasesLastMonth,
				Description:        audience.Description,
				AssetType:          "audience",
			})
			mutex.Unlock()
		}

	}(&modelUser.FavouriteAudiences, &favouriteAssets, &mutex)

	wg.Add(1)
	go func(
		insights *[]Insight,
		favouriteAssets *[]domain.Asset,
		mutex *sync.Mutex,
	) {
		defer wg.Done()

		for _, insight := range modelUser.FavouriteInsights {
			mutex.Lock()
			*favouriteAssets = append(*favouriteAssets, domain.Insight{
				Text:        insight.Text,
				Description: insight.Description,
				AssetType:   "insight",
			})
			mutex.Unlock()
		}
	}(&modelUser.FavouriteInsights, &favouriteAssets, &mutex)

	wg.Wait()

	return &domain.User{
		Username:        modelUser.Username,
		Password:        modelUser.Password,
		FavouriteAssets: favouriteAssets,
	}, err
}
