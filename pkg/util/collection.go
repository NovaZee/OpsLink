package util

import (
	"reflect"
	"strconv"
)

func StructToStringSlice(entity interface{}) []string {
	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)

	if entityType.Kind() != reflect.Struct {
		return nil
	}

	numFields := entityType.NumField()
	values := make([]string, numFields)

	for i := 0; i < numFields; i++ {
		field := entityValue.Field(i)
		// 使用 strconv 包将属性值转换为字符串
		values[i] = strconv.FormatInt(field.Int(), 10)
	}
	return values
}
