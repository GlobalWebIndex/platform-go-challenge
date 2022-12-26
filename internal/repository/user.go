package repository

import (
	//"fmt"

	"ownify_api/internal/domain"
	"ownify_api/internal/dto"

	//"ownify_api/internal/dto"

	"github.com/Masterminds/squirrel"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
)

type UserQuery[T domain.Userable] interface {
	CreateUser(
		user dto.BriefUser,
	) (*int64, error)
	GetUser(id int64) (*T, error)
	GetUserByBriefInfo(user dto.BriefUser) (*int64, error)
	DeleteUser(userID int64) error
}

type userQuery[T domain.Userable] struct{}

func (u *userQuery[T]) CreateUser(
	user dto.BriefUser) (*int64, error) {
	qb := pgQb().
		Insert(domain.PersonTableName).
		Columns("first_name", "email", "password", "last_name",
			"role", "verified", "email_code", "balance", "phone_number").
		// Values(user.FirstName, user.Email, user.Password, user.LastName,
		// 	user.Role, user.Verified, user.EmailCode, user.Balance, user.PhoneNumber).
		Suffix("RETURNING id")

	var id int64
	err := qb.QueryRow().Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (u *userQuery[T]) UpdateUser(user T) (*int64, error) {
	qb := pgQb().
		Insert(domain.PersonTableName).
		Columns("first_name", "email", "password", "last_name",
			"role", "verified", "email_code", "balance", "phone_number").
		// Values(user.FirstName, user.Email, user.Password, user.LastName,
		// 	user.Role, user.Verified, user.EmailCode, user.Balance, user.PhoneNumber).
		Suffix("RETURNING id")

	var id int64
	err := qb.QueryRow().Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (u *userQuery[T]) DeleteUser(userID int64) error {
	qb := pgQb().
		Delete(domain.PersonTableName).
		From(domain.PersonTableName).
		Where(squirrel.Eq{"id": userID})

	_, err := qb.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (u *userQuery[T]) GetUser(id int64) (*T, error) {
	// qb := pgQb().
	// 	Delete(domain.PersonTableName).
	// 	From(domain.PersonTableName).
	// 	Where(squirrel.Eq{"id": id})

	// _, err := qb.Exec()
	// if err != nil {
	// 	return err
	// }
	return nil, nil
}

func (u *userQuery[T]) GetUserByBriefInfo(user dto.BriefUser) (*int64, error) {
	// qb := pgQb().
	// 	Delete(domain.PersonTableName).
	// 	From(domain.PersonTableName).
	// 	Where(squirrel.Eq{"id": id})

	// _, err := qb.Exec()
	// if err != nil {
	// 	return err
	// }
	return nil, nil
}
