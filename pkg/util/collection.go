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

func PageSlice(collectionLength, pageNo, pageSize int) (int, int) {
	startIndex := (pageNo - 1) * pageSize
	endIndex := pageNo * pageSize

	if startIndex >= collectionLength {
		return -1, -1
	}

	if endIndex > collectionLength {
		endIndex = collectionLength
	}

	return startIndex, endIndex
}
