package utils

import (
	"ownify_api/internal/dto"
	"reflect"
)

func ConvertToEntity[T dto.SQLable](data *T) ([]string, []interface{}) {
	entity := reflect.ValueOf(data).Elem()
	cols := []string{}
	values := []interface{}{}
	for i := 0; i < entity.NumField(); i++ {
		field := entity.Type().Field(i).Name
		value := entity.FieldByName(field)
		if IsZeroOfUnderlyingType(value) {
			continue
		}
		cols = append(cols, ToSnakeCase(field))
		values = append(values, value)
	}
	return cols, values
}

func IsZeroOfUnderlyingType(v reflect.Value) bool {
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
