package repository

import (
	//"fmt"

	"fmt"
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"

	//"ownify_api/internal/dto"

	sq "github.com/Masterminds/squirrel"
)

type UserQuery interface {
	CreateUser(
		user dto.BriefUser,
	) error
	GetUser(pubKey string, idFingerprint string) (*dto.BriefUser, error)
	DeleteUser(pubKey string) error
	GetLastUserId(walletType string) (*int64, error)
	VerifyUser(userId string, pubKey string) (*interface{}, error)
}

type userQuery struct{}

// VerifyUser implements UserQuery
func (*userQuery) VerifyUser(userId string, pubKey string) (*interface{}, error) {
	panic("unimplemented")
}

// GetLastUserId implements UserQuery
func (*userQuery) GetLastUserId(walletType string) (*int64, error) {
	panic("unimplemented")
}

func (u *userQuery) CreateUser(
	user dto.BriefUser) error {
	if !user.Valid() {
		return fmt.Errorf("[ERR] invalid Info: %v", user.PubKey)
	}

	cols, values := utils.ConvertToEntity(&user)
	sqlBuilder := utils.NewSqlBuilder()

	query, err := sqlBuilder.Insert(domain.UserTableName, cols, values)
	fmt.Println(query)

	if err != nil {
		return err
	}
	_, err = DB.Exec(*query)
	if err != nil {
		return err
	}
	return nil
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

func (u *userQuery) DeleteUser(pubKey string) error {
	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Select(domain.UserTableName, []string{}, []utils.Tuple{{Key: "pub_addr", Val: pubKey}}, "=", "OR")
	if err != nil {
		return err
	}
	_, err = DB.Exec(*sql)
	if err != nil {
		return err
	}
	return nil
}

func (u *userQuery) GetUser(pubKey string, idFingerprint string) (*dto.BriefUser, error) {
	var user dto.BriefUser
	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Select(domain.UserTableName, []string{
		"first_name",
		"last_name",
		"birth_day",
		"gender",
		"nationality",
		"created_time",
	}, []utils.Tuple{{Key: "pub_addr", Val: pubKey}, {Key: "id_fingerprint", Val: idFingerprint}}, "=", "OR")
	if err != nil {
		return nil, err
	}
	err = DB.QueryRow(*sql).Scan(
		&user.FirstName,
		&user.LastName,
		&user.BirthDay,
		&user.Gender,
		&user.BirthDay,
		&user.UserId,
	)
	if err != nil {
		return nil, err
	}
	user.PubKey = pubKey
	return &user, nil
}

func (u *userQuery) GetBusiness(email string) (*interface{}, error) {

	var business interface{}
	err := pgQb().Select("*").Where(sq.Eq{"email": email}).From(domain.BusinessTableName).QueryRow().Scan(&business)
	if err != nil {
		return nil, err
	}
	return &business, err
}

func (b *businessQuery) VerifyUser(userId string, pubKey string) (*interface{}, error) {
	var user interface{}
	sqlBuilder := utils.NewSqlBuilder()
	conditions := []utils.Tuple{{Key: "user_id", Val: userId}, {Key: "pub_addr", Val: pubKey}}
	sql, err := sqlBuilder.Select(domain.BusinessTableName, []string{}, conditions, "=", "AND")
	fmt.Println(*sql)
	if err != nil {
		return nil, err
	}

	user, err = DB.Exec(*sql)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
