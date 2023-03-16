package repository

import (
	"fmt"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
)

type LoggerQuery interface {
	LogUserActivity(userActivity *dto.UserActivity)
}

type loggerQuery struct{}

func (l *loggerQuery) LogUserActivity(userActivity *dto.UserActivity) {

	tableName := AssetVerifyTableName
	sqlBuilder := utils.NewSqlBuilder()
	cols,vals := utils.ConvertToEntity(userActivity)
	sql, err := sqlBuilder.Insert(tableName, cols, vals)
	if err != nil {
		fmt.Println("[ERR] loggin:", err.Error())
	}
	_, err = DB.Exec(*sql)
	if err != nil {
		fmt.Println("[ERR] loggin:", err.Error())
	}
}
