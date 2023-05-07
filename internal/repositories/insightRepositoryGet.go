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

func (repo *InsightRepository) GetInsight(
	ctx context.Context,
	uid uint32,
) (*domain.Insight, error) {
	var err error
	var modelInsight *Insight

	err = repo.db.WithContext(ctx).
		Model(Insight{}).
		Where("id = ?", uid).
		Take(&modelInsight).Error

	if err == gorm.ErrRecordNotFound {
		return &domain.Insight{}, apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(uid)) + " not found"),
		}
	}

	if err != nil {
		return &domain.Insight{}, err
	}

	return &domain.Insight{
		Text:        modelInsight.Text,
		Description: modelInsight.Description,
	}, err
}
