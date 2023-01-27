package utils

import (
	"fmt"
	"reflect"
	"strings"
)

type Tuple struct {
	Key string
	Val interface{}
}

type SqlBuilder interface {
	Insert(tableName string, cols []string, values []interface{}) (*string, error)
	Delete(tableName string, conditions []Tuple, joinKey string) (*string, error)
	Select(tableName string, targets []string, conditions []Tuple, joinKey string) (*string, error)
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
		fmt.Println(reflect.ValueOf(val).Kind())
		if reflect.ValueOf(val).Kind() == reflect.Struct || reflect.ValueOf(val).Kind() == reflect.String {
			valueS = append(valueS, strings.TrimSpace(fmt.Sprintf("'%v'", val)))
			continue
		}
		valueS = append(valueS, fmt.Sprintf("%v", val))
	}
	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", tableName, strings.Join(cols, ","), strings.Join(valueS, ","))
	return &query, nil
}

func (builder sqlBuilder) Select(tableName string, targets []string, conditions []Tuple, joinKey string) (*string, error) {
	targetStr := "*"
	if len(targets) != 0 {
		targetStr = strings.Join(targets, ",")
	}
	condStr := *conditionBuilder(conditions, joinKey)
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", targetStr, tableName, condStr)
	return &sql, nil
}

func (builder sqlBuilder) Delete(tableName string, conditions []Tuple, joinKey string) (*string, error) {
	condStr := *conditionBuilder(conditions, joinKey)
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, condStr)
	return &sql, nil
}

func conditionBuilder(conditions []Tuple, joinKey string) *string {
	condSize := len(conditions)
	if condSize == 0 {
		return nil
	}
	condStr := ""
	for index, pair := range conditions {
		_, ok := pair.Val.(string)
		if ok {
			condStr += pair.Key + "=" + fmt.Sprintf("'%v'", pair.Val)
		} else {
			condStr += pair.Key + "=" + fmt.Sprintf("%v", pair.Val)
		}
		if index < condSize-1 {
			condStr += " " + joinKey + " "
		}
	}
	return &condStr
}
