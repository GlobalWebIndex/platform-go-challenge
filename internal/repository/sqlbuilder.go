package repository

import (
	"fmt"
	"reflect"
	"strings"
)

type SqlBuilder interface {
	Insert(tableName string, cols []string, values []interface{}) (*string, error)
}

type sqlBuilder struct{}

func NewSqlBuilder() SqlBuilder {
	return &sqlBuilder{}
}

func (builder sqlBuilder) Insert(tableName string, cols []string, values []interface{}) (*string, error) {
	if len(cols) != len(values) {
		return nil, fmt.Errorf("[ERR] doesn't match length! cols: %s, values: %s", cols, values)
	}
	valueS := []string{}
	for _, val := range values {
		if reflect.ValueOf(val).Kind() == reflect.String {
			valueS = append(valueS, fmt.Sprintf("'%v'", val))
			continue
		}
		valueS = append(valueS, fmt.Sprintf("%v", val))
	}
	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", tableName, strings.Join(cols, ","), strings.Join(valueS, ","))
	return &query, nil
}
