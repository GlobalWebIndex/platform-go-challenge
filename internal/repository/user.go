package repository

import (
	//"fmt"

	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
	"strings"

	//"ownify_api/internal/dto"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
)

type UserQuery interface {
	CreateUser(
		user dto.BriefUser,
	) (*int64, error)
	GetUser(id int64, walletType string) (*interface{}, error)
	GetUserByBriefInfo(user dto.BriefUser) (*int64, error)
	DeleteUser(userID int64, walletType string) error
	GetLastUserId(walletType string) (*int64, error)
}

type userQuery struct{}

func (u *userQuery) CreateUser(
	user dto.BriefUser) (*int64, error) {
	tableName := domain.PersonTableName
	if user.WalletType == domain.BusinessWallet {
		tableName = domain.BusinessTableName
	}

	var user_id int64 = 0
	err := pgQb().Select("user_id").OrderBy("user_id DESC").From(tableName).QueryRow().Scan(&user_id)
	if !strings.Contains(err.Error(), domain.NoRows) {
		return nil, err
	}
	cols := []string{"user_id", "chain_id", "wallet_address"}
	values := []interface{}{user_id + 1, user.ChainId, user.PubKey}
	sqlBuilder := utils.NewSqlBuilder()

	query, err := sqlBuilder.Insert(tableName, cols, values)
	if err != nil {
		return nil, err
	}
	_, err = DB.Exec(*query)
	if err != nil {
		return nil, err
	}
	user_id += 1
	return &user_id, nil
}



// func (u *userQuery) UpdateUser(user T) (*int64, error) {
// 	qb := pgQb().
// 		Insert(domain.PersonTableName).
// 		Columns("first_name", "email", "password", "last_name",
// 			"role", "verified", "email_code", "balance", "phone_number").
// 		// Values(user.FirstName, user.Email, user.Password, user.LastName,
// 		// 	user.Role, user.Verified, user.EmailCode, user.Balance, user.PhoneNumber).
// 		Suffix("RETURNING id")

// 	var id int64
// 	err := qb.QueryRow().Scan(&id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &id, nil
// }

func (u *userQuery) DeleteUser(userID int64, walletType string) error {
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

func (u *userQuery) GetUser(id int64, walletType string) (*interface{}, error) {
	tableName := domain.PersonTableName
	if walletType == domain.BusinessWallet {
		tableName = domain.BusinessTableName
	}

	var user interface{}
	err := pgQb().Select("*").Where(sq.Eq{"id": id}).From(tableName).QueryRow().Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (u *userQuery) GetBusiness(email string) (*interface{}, error) {

	var business interface{}
	err := pgQb().Select("*").Where(sq.Eq{"email": email}).From(domain.BusinessTableName).QueryRow().Scan(&business)
	if err != nil {
		return nil, err
	}
	return &business, err
}

func (u *userQuery) GetUserByBriefInfo(user dto.BriefUser) (*int64, error) {
	var user_id int64
	err := pgQb().
		Select("user_id").
		From(domain.PersonTableName).
		Where(squirrel.Eq{"chain_id": user.ChainId, "wallet_addres": user.PubKey}).QueryRow().Scan(&user_id)

	if err != nil {
		return nil, err
	}
	return &user_id, nil
}

func (u *userQuery) GetLastUserId(walletType string) (*int64, error) {
	tableName := domain.PersonTableName
	if walletType == domain.BusinessWallet {
		tableName = domain.BusinessTableName
	}
	var user_id int64 = 0
	err := pgQb().Select("user_id").OrderBy("user_id DESC").From(tableName).QueryRow().Scan(&user_id)
	if err != nil {
		return nil, err
	}
	return &user_id, nil
}
