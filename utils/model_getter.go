package util

import (
	"errors"
	"reflect"
	"strconv"
)

// GetFieldString get the string value of a field in a struct
func GetFieldString(m any, i int) (string, error) {
	if m == nil {
		return "", errors.New("model is nil")
	}
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return "", errors.New("input is not a pointer to a struct")
	}
	v = v.Elem()

	if i < 0 || i >= v.NumField() {
		return "", errors.New("field index out of range")
	}

	field := v.Field(i)

	switch field.Kind() {
	case reflect.String:
		return field.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'g', -1, 64), nil
	default:
		return "", errors.New("unsupported field type")
	}
}
