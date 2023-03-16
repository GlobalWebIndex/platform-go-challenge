package repository

import (
	//"fmt"

	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
)

type UserQuery interface {
	CreateUser(
		user dto.BriefUser,
	) error
	ValidUser(pubKey string, idFingerprint string) (*dto.BriefUser, error)
	DeleteUser(pubKey string) error
	VerifyUser(userId string, pubKey string) (*dto.BriefUser, error)
}

type userQuery struct{}

func (u *userQuery) CreateUser(
	user dto.BriefUser) error {
	if err := user.Valid(); err != nil {
		return err
	}

	cols, values := utils.ConvertToEntity(&user)
	sqlBuilder := utils.NewSqlBuilder()
	query, err := sqlBuilder.Insert(UserTableName, cols, values)
	if err != nil {
		return err
	}
	_, err = DB.Exec(*query)
	if err != nil {
		return err
	}
	return nil
}

func (u *userQuery) DeleteUser(pubKey string) error {
	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Delete(UserTableName, []utils.Tuple{{Key: "pub_addr", Val: pubKey}}, "OR")
	if err != nil {
		return err
	}
	_, err = DB.Exec(*sql)
	if err != nil {
		return err
	}
	return nil
}

func (u *userQuery) ValidUser(pubKey string, idFingerprint string) (*dto.BriefUser, error) {

	var user dto.BriefUser
	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Select(UserTableName, []string{
		"first_name",
		"last_name",
		"birth_day",
		"gender",
		"nationality",
	}, []utils.Tuple{{Key: "pub_addr", Val: pubKey}, {Key: "id_fingerprint", Val: idFingerprint}}, "=", "OR")
	if err != nil {
		return nil, err
	}
	err = DB.QueryRow(*sql).Scan(
		&user.FirstName,
		&user.LastName,
		&user.BirthDay,
		&user.Gender,
		&user.Nationality,
	)
	if err != nil {
		return nil, err
	}
	user.PubAddr = pubKey
	return &user, nil
}

func (b *userQuery) VerifyUser(userId string, pubKey string) (*dto.BriefUser, error) {
	var user dto.BriefUser
	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Select(UserTableName, []string{
		"first_name",
		"last_name",
		"birth_day",
		"gender",
		"nationality",
	}, []utils.Tuple{{Key: "pub_addr", Val: pubKey}, {Key: "user_id", Val: userId}}, "=", "AND")
	if err != nil {
		return nil, err
	}
	err = DB.QueryRow(*sql).Scan(
		&user.FirstName,
		&user.LastName,
		&user.BirthDay,
		&user.Gender,
		&user.Nationality,
	)
	if err != nil {
		return nil, err
	}
	user.PubAddr = pubKey
	return &user, nil
}
