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

const (
	eq   = "="
	like = "LIKE"
)

type CondOperator string

func (s CondOperator) valid() bool {
	if s == eq || s == like {
		return true
	}
	return false
}

type SqlBuilder interface {
	Insert(tableName string, cols []string, values []interface{}) (*string, error)
	Delete(tableName string, conditions []Tuple, joinKey string) (*string, error)
	Select(tableName string, targets []string, conditions []Tuple, condOperator CondOperator, joinKey string) (*string, error)
	Update(tableName string, values []Tuple, conditions []Tuple, joinKey string) (*string, error)
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
		if reflect.ValueOf(val).Kind() == reflect.Struct || reflect.ValueOf(val).Kind() == reflect.String {
			valueS = append(valueS, strings.TrimSpace(fmt.Sprintf("\"%v\"", val)))
			continue
		}
		valueS = append(valueS, fmt.Sprintf("%v", val))
	}
	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", tableName, strings.Join(cols, ","), strings.Join(valueS, ","))
	return &query, nil
}

func (builder sqlBuilder) Select(tableName string, targets []string, conditions []Tuple, condOperator CondOperator, joinKey string) (*string, error) {
	if !condOperator.valid() {
		return nil, fmt.Errorf("[ERR] invalid operator: %s", condOperator)
	}
	targetStr := "*"
	if len(targets) != 0 {
		targetStr = strings.Join(targets, ",")
	}
	if len(conditions) == 0 {
		sql := fmt.Sprintf("SELECT %s FROM %s", targetStr, tableName)
		return &sql, nil
	}
	condStr := *conditionBuilder(conditions, condOperator, joinKey)
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", targetStr, tableName, condStr)
	return &sql, nil
}

func (builder sqlBuilder) Delete(tableName string, conditions []Tuple, joinKey string) (*string, error) {
	condStr := *conditionBuilder(conditions, eq, joinKey)
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, condStr)
	return &sql, nil
}

func conditionBuilder(conditions []Tuple, operator CondOperator, joinKey string) *string {

	condSize := len(conditions)
	if condSize == 0 {
		return nil
	}
	condStr := ""
	prefix := ""
	suffix := ""
	if operator == like {
		prefix = "%"
		suffix = "%"
	}
	for index, pair := range conditions {
		_, ok := pair.Val.(string)
		if ok {
			condStr += pair.Key + fmt.Sprintf(" %s \"%s%v%s\"", operator, prefix, pair.Val, suffix)
		} else {
			condStr += pair.Key + fmt.Sprintf(" %s %s%v%s", operator, prefix, pair.Val, suffix)
		}

		if index < condSize-1 {
			condStr += " " + joinKey + " "
		}
	}
	return &condStr
}

func GenerateCond(keys []string, values []interface{}) []Tuple {
	if len(keys) != len(values) {
		return []Tuple{}
	}
	conds := []Tuple{}
	for i := 0; i < len(keys); i++ {
		conds = append(conds, Tuple{Key: keys[i], Val: values[i]})
	}
	return conds
}


func (*sqlBuilder) Update(tableName string, values []Tuple, conditions []Tuple, joinKey string) (*string, error) {
	valueQ := conditionBuilder(values, eq, ",")
	condQ := conditionBuilder(conditions, eq, joinKey)
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, *valueQ, *condQ)
	return &sql, nil
}
