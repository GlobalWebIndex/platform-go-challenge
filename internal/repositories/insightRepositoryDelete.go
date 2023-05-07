package repositories

import (
	"context"
	"errors"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (repo *InsightRepository) DeleteInsight(ctx context.Context, uid uint32) error {
	db := repo.db.WithContext(ctx).Model(&Insight{}).
		Where("id = ?", uid).
		Take(&Insight{}).
		Delete(&Insight{})

	if db.Error == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(uid)) + " not found"),
		}
	}

	if db.Error != nil {
		return db.Error
	}
	return nil
}
