package repositories

import (
	"context"
	"errors"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (repo *InsightRepository) UpdateInsightDescription(
	ctx context.Context,
	uid uint32,
	description string,
) error {
	err := repo.db.WithContext(ctx).Model(&Insight{}).
		Where("id = ?", uid).
		Update("description", description).Error

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(uid)) + " not found"),
		}
	}

	return err
}
