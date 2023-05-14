package repository

import (
	"ownify_api/internal/utils"
)

type LicenseQuery interface {
	SaveAPIKey(email, userId, apiKey string) error
	GetAPIKey(email, userId string) ([]string, error)
}

type licenseQuery struct{}

func (l *licenseQuery) SaveAPIKey(email, userId, apiKey string) error {

	tableName := UserLicenseTableName
	sqlBuilder := utils.NewSqlBuilder()
	cols := []string{"user_id", "email", "api_key"}
	vals := []string{userId, email, apiKey}

	// Convert []string to []interface{}
	interfaceVals := make([]interface{}, len(vals))
	for i, v := range vals {
		interfaceVals[i] = v
	}

	sql, err := sqlBuilder.Insert(tableName, cols, interfaceVals)
	if err != nil {
		return err
	}
	_, err = DB.Exec(*sql)
	if err != nil {
		return err
	}
	return nil
}

func (l *licenseQuery) GetAPIKey(email, userId string) ([]string, error) {

	tableName := UserLicenseTableName
	sqlBuilder := utils.NewSqlBuilder()
	apiKeys := ""

	sql, err := sqlBuilder.Select(tableName, []string{"api_key"}, []utils.Tuple{
		{
			Key: "email",
			Val: email,
		},
		{
			Key: "user_id",
			Val: userId,
		},
	}, "=", "AND")

	if err != nil {
		return nil, err
	}
	err = DB.QueryRow(*sql).Scan(&apiKeys)
	if err != nil {
		return nil, err
	}
	return []string{apiKeys}, nil
}
