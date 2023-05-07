package repositories

import (
	"context"
	"errors"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (repo *UserRepository) AddFavouriteAsset(
	ctx context.Context,
	userId uint32,
	assetId uint32,
	assetType string,
) error {
	switch assetType {
	case string(domain.ChartType):
		return repo.addFavouriteChart(ctx, userId, assetId)
	case string(domain.AudienceType):
		return repo.addFavouriteAudience(ctx, userId, assetId)
	case string(domain.InsightType):
		return repo.addFavouriteInsight(ctx, userId, assetId)
	default:
		return apierrors.UnknownAssetTypeErrorWrapper{
			ReturnedStatusCode: http.StatusInternalServerError,
			OriginalError:      errors.New("unknown asset type " + assetType),
		}
	}
}

func (repo *UserRepository) addFavouriteChart(
	ctx context.Context,
	userId uint32,
	assetId uint32,
) error {
	var err error
	var modelUser *User
	var modelChart *Chart

	err = repo.db.WithContext(ctx).
		Preload("FavouriteCharts").
		Model(User{}).
		Where("id = ?", userId).
		Take(&modelUser).Error

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(userId)) + " not found"),
		}
	}

	if err != nil {
		return err
	}

	err = repo.db.WithContext(ctx).
		Model(Chart{}).
		Where("id = ?", assetId).
		Take(&modelChart).Error

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("chartId " + strconv.Itoa(int(assetId)) + " not found"),
		}
	}

	if err != nil {
		return err
	}

	modelUser.FavouriteCharts = append(
		modelUser.FavouriteCharts,
		*modelChart,
	)

	return repo.db.WithContext(ctx).Omit("FavouriteCharts.*").Save(&modelUser).Error
}

func (repo *UserRepository) addFavouriteAudience(
	ctx context.Context,
	userId uint32,
	assetId uint32,
) error {
	var err error
	var modelUser *User
	var modelAudience *Audience

	err = repo.db.WithContext(ctx).
		Preload("FavouriteAudiences").
		Model(User{}).
		Where("id = ?", userId).
		Take(&modelUser).Error

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(userId)) + " not found"),
		}
	}

	if err != nil {
		return err
	}

	err = repo.db.WithContext(ctx).
		Model(Audience{}).
		Where("id = ?", assetId).
		Take(&modelAudience).Error

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("audienceId " + strconv.Itoa(int(assetId)) + " not found"),
		}
	}

	if err != nil {
		return err
	}

	modelUser.FavouriteAudiences = append(
		modelUser.FavouriteAudiences,
		*modelAudience,
	)

	return repo.db.WithContext(ctx).Omit("FavouriteAudiences.*").Save(&modelUser).Error
}

func (repo *UserRepository) addFavouriteInsight(
	ctx context.Context,
	userId uint32,
	assetId uint32,
) error {
	var err error
	var modelUser *User
	var modelInsight *Insight

	err = repo.db.WithContext(ctx).
		Preload("FavouriteInsights").
		Model(User{}).
		Where("id = ?", userId).
		Take(&modelUser).Error

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(userId)) + " not found"),
		}
	}

	if err != nil {
		return err
	}

	err = repo.db.WithContext(ctx).
		Model(Insight{}).
		Where("id = ?", assetId).
		Take(&modelInsight).Error

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("insightId " + strconv.Itoa(int(assetId)) + " not found"),
		}
	}

	if err != nil {
		return err
	}

	modelUser.FavouriteInsights = append(
		modelUser.FavouriteInsights,
		*modelInsight,
	)

	return repo.db.WithContext(ctx).Omit("FavouriteInsights.*").Save(&modelUser).Error
}

func (repo *UserRepository) RemoveFavouriteAsset(
	ctx context.Context,
	userId uint32,
	assetId uint32,
	assetType string,
) error {
	switch assetType {
	case string(domain.ChartType):
		return repo.removeFavouriteChart(ctx, userId, assetId)
	case string(domain.AudienceType):
		return repo.removeFavouriteAudience(ctx, userId, assetId)
	case string(domain.InsightType):
		return repo.removeFavouriteInsight(ctx, userId, assetId)
	default:
		return apierrors.UnknownAssetTypeErrorWrapper{
			ReturnedStatusCode: http.StatusInternalServerError,
			OriginalError:      errors.New("unknown asset type " + assetType),
		}
	}
}

func (repo *UserRepository) removeFavouriteChart(
	ctx context.Context,
	userId uint32,
	assetId uint32,
) error {
	var user User
	err := repo.db.WithContext(ctx).First(&user, userId).Error

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(userId)) + " not found"),
		}
	}

	if err != nil {
		return err
	}

	chartAssociation := repo.db.Model(&user).Association("FavouriteCharts")

	err = chartAssociation.Delete(
		&Chart{Model: gorm.Model{ID: uint(assetId)}},
	)

	return err
}

func (repo *UserRepository) removeFavouriteInsight(
	ctx context.Context,
	userId uint32,
	assetId uint32,
) error {
	var user User
	err := repo.db.WithContext(ctx).First(&user, userId).Error

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(userId)) + " not found"),
		}
	}

	if err != nil {
		return err
	}

	insightAssociation := repo.db.WithContext(ctx).
		Model(&user).Association("FavouriteInsights")

	err = insightAssociation.Delete(
		&Insight{Model: gorm.Model{ID: uint(assetId)}},
	)

	return err
}

func (repo *UserRepository) removeFavouriteAudience(
	ctx context.Context,
	userId uint32,
	assetId uint32,
) error {
	var user User
	err := repo.db.WithContext(ctx).First(&user, userId).Error

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(userId)) + " not found"),
		}
	}

	if err != nil {
		return err
	}

	audienceAssociation := repo.db.WithContext(ctx).
		Model(&user).Association("FavouriteAudiences")

	err = audienceAssociation.Delete(
		&Audience{Model: gorm.Model{ID: uint(assetId)}},
	)

	return err
}

func (repo *UserRepository) EditFavouriteAssetDescription(
	ctx context.Context,
	userId uint32,
	assetId uint32,
	assetType,
	description string,
) error {
	switch assetType {
	case string(domain.ChartType):
		return repo.editFavouriteChartDescription(ctx, userId, assetId, description)
	case string(domain.AudienceType):
		return repo.editFavouriteAudienceDescription(ctx, userId, assetId, description)
	case string(domain.InsightType):
		return repo.editFavouriteInsightDescription(ctx, userId, assetId, description)
	default:
		return apierrors.UnknownAssetTypeErrorWrapper{
			ReturnedStatusCode: http.StatusInternalServerError,
			OriginalError:      errors.New("unknown asset type " + assetType),
		}
	}
}

func (repo *UserRepository) editFavouriteChartDescription(
	ctx context.Context,
	userId uint32,
	assetId uint32,
	description string,
) error {
	var err error

	modelChart := Chart{
		Model: gorm.Model{ID: uint(assetId)},
	}
	err = repo.db.WithContext(ctx).
		Model(User{}).
		Where("id = ?", userId).
		Association("FavouriteCharts").
		Find(&modelChart)

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(userId)) + " not found"),
		}
	}

	if err != nil {
		return nil
	}

	modelChart.Description = description
	return repo.db.WithContext(ctx).Save(&modelChart).Error
}

func (repo *UserRepository) editFavouriteAudienceDescription(
	ctx context.Context,
	userId uint32,
	assetId uint32,
	description string,
) error {
	var err error

	modelAudience := Audience{
		Model: gorm.Model{ID: uint(assetId)},
	}
	err = repo.db.WithContext(ctx).
		Model(User{}).
		Where("id = ?", userId).
		Association("FavouriteAudiences").
		Find(&modelAudience)

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(userId)) + " not found"),
		}
	}

	if err != nil {
		return nil
	}

	modelAudience.Description = description
	return repo.db.WithContext(ctx).Save(&modelAudience).Error
}

func (repo *UserRepository) editFavouriteInsightDescription(
	ctx context.Context,
	userId uint32,
	assetId uint32,
	description string,
) error {
	var err error

	modelInsight := Insight{
		Model: gorm.Model{ID: uint(assetId)},
	}
	err = repo.db.WithContext(ctx).
		Model(User{}).
		Where("id = ?", userId).
		Association("FavouriteInsights").
		Find(&modelInsight)

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(userId)) + " not found"),
		}
	}

	if err != nil {
		return nil
	}

	modelInsight.Description = description
	return repo.db.WithContext(ctx).Save(&modelInsight).Error
}
