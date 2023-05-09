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

func (repo *AudienceRepository) GetAudience(
	ctx context.Context,
	uid uint32,
) (*domain.Audience, error) {
	var err error
	var modelAudience *Audience

	err = repo.db.WithContext(ctx).
		Model(Audience{}).
		Where("id = ?", uid).
		Take(&modelAudience).Error

	if err == gorm.ErrRecordNotFound {
		return &domain.Audience{}, apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(uid)) + " not found"),
		}
	}

	if err != nil {
		return &domain.Audience{}, err
	}

	return &domain.Audience{
		Gender:             modelAudience.Gender,
		BirthCountry:       modelAudience.BirthCountry,
		AgeGroup:           modelAudience.AgeGroup,
		HoursSpentDaily:    modelAudience.HoursSpentOnSocial,
		PurchasesLastMonth: modelAudience.PurchasesLastMonth,
		Description:        modelAudience.Description,
	}, err
}
