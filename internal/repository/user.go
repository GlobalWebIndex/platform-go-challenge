package repository

import (
	//"fmt"

	"fmt"
	"gwi_api/internal/dto"
)

type UserQuery interface {
	CreateUser(
		user dto.UserDto) (*uint64, error)

	DeleteUser(userId uint64) error
	GetUserPasswordByEmail(email string) (*string, error)
	GetUserIdByEmail(email string) (*uint64, error)
}

type userQuery struct{}

// GetUserIdByEmail implements UserQuery.
func (u *userQuery) GetUserIdByEmail(email string) (*uint64, error) {
	userId, found := DB.userEmails[email]
	if !found {
		return nil, fmt.Errorf("%s", "did not exist user")
	}
	return &userId, nil
}

// GetUserPasswordByEmail implements UserQuery.
func (u *userQuery) GetUserPasswordByEmail(email string) (*string, error) {
	userId, found := DB.userEmails[email]
	if !found {
		return nil, fmt.Errorf("%s", "did not exist user")
	}

	user, found := DB.users[userId]
	if !found {
		return nil, fmt.Errorf("%s", "did not exist user")
	}
	return &user.Password, nil
}

func (u *userQuery) CreateUser(
	user dto.UserDto) (*uint64, error) {
	return DB.RegisterUser(user)
}

func (u *userQuery) DeleteUser(userId uint64) error {
	return DB.DeleteUser(userId)
}

func (u *userQuery) ValidUser(pubKey string, idFingerprint string) (*dto.UserDto, error) {

	var user dto.UserDto
	// sqlBuilder := utils.NewSqlBuilder()
	// // sql, err := sqlBuilder.Select(UserTableName, []string{
	// // 	"first_name",
	// // 	"last_name",
	// // 	"birth_day",
	// // 	"gender",
	// // 	"nationality",
	// // }, []utils.Tuple{{Key: "pub_addr", Val: pubKey}, {Key: "id_fingerprint", Val: idFingerprint}}, "=", "OR")
	// if err != nil {
	// 	return nil, err
	// }

	// if err != nil {
	// 	return nil, err
	// }

	return &user, nil
}
