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

func (repo *UserRepository) UpdateUser(
	ctx context.Context,
	uid uint32,
	user *domain.User,
) error {
	err := repo.db.WithContext(ctx).Model(&User{}).
		Where("id = ?", uid).
		Updates(
			map[string]interface{}{
				"username": user.Username,
				"password": user.Password,
			},
		).Error

	if err == gorm.ErrRecordNotFound {
		return apierrors.UserNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("userId " + strconv.Itoa(int(uid)) + " not found"),
		}
	}

	return err
}
