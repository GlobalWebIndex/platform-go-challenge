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

func (repo *ChartRepository) GetChart(
	ctx context.Context,
	uid uint32,
) (*domain.Chart, error) {
	var err error
	var modelChart *Chart

	err = repo.db.WithContext(ctx).
		Model(Chart{}).
		Where("id = ?", uid).
		Take(&modelChart).Error

	if err == gorm.ErrRecordNotFound {
		return &domain.Chart{}, apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(uid)) + " not found"),
		}
	}

	if err != nil {
		return &domain.Chart{}, err
	}

	return &domain.Chart{
		Title:       modelChart.Title,
		XAxisTitle:  modelChart.XTitle,
		YAxisTitle:  modelChart.YTitle,
		Data:        modelChart.Data,
		Description: modelChart.Description,
	}, err
}
