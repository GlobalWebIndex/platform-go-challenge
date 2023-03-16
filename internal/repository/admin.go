package repository

import (
	//"fmt"

	"ownify_api/internal/utils"
)

type AdminQuery interface {
	GrantBusiness(email string, isApproved bool) error
}

type adminQuery struct{}

func (u *adminQuery) GrantBusiness(email string, isApproved bool) error {
	err := utils.IsEmail(email)
	if err != nil {
		return err
	}

	tableName := BusinessTableName
	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Update(tableName, []utils.Tuple{{Key: "is_approved", Val: isApproved}}, []utils.Tuple{{Key: "email", Val: email}}, "AND")
	if err != nil {
		return err
	}
	_, err = DB.Exec(*sql)
	return err
}
