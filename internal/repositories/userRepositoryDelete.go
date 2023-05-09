package repositories

import (
	"context"
	"errors"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (repo *UserRepository) DeleteUser(ctx context.Context, uid uint32) error {
	// I've chosen to do a soft delete on Users as supported by Gorm. On many platforms
	// when deleting a user, he is not actually deleted and can restore his account.
	// This is done by keeping user-asset relationships and delete by updating the deleted_at
	// field in the DB
	db := repo.db.WithContext(ctx).Model(&User{}).
		Where("id = ?", uid).
		Take(&User{}).
		Delete(&User{})

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
