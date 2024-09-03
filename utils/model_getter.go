package util

import (
	"errors"
	"excel_import"
	"fmt"
	"reflect"
	"strconv"
)

var (
	unexpectedSentence = "expected model: %v, but got: %v"
)

// GetFieldString get the string value of a field in a struct
func GetFieldString(m any, i int) (string, error) {
	if m == nil {
		return "", nil
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

func CompareModel(real, expected any, attr []*excel_import.ExcelImportTagAttr, key string) error {
	// check if both are nil
	if real == nil || expected == nil {
		if real == nil && expected == nil {
			return nil
		}
		return errors.New(fmt.Sprintf(unexpectedSentence, expected, real))
	}

	// check if both are not pointers to structs
	vReal := reflect.ValueOf(real)
	vExpected := reflect.ValueOf(expected)

	if vReal.Kind() != reflect.Ptr || vReal.Elem().Kind() != reflect.Struct {
		return errors.New("input is not a pointer to a struct")
	}
	if vExpected.Kind() != reflect.Ptr || vExpected.Elem().Kind() != reflect.Struct {
		return errors.New("input is not a pointer to a struct")
	}

	// get the struct values
	vReal = vReal.Elem()
	vExpected = vExpected.Elem()

	// check if the number of fields are equal
	if vReal.NumField() != vExpected.NumField() {
		return errors.New("number of fields not equal")
	}

	for i, a := range attr {
		// check if the field is not checked
		if !excel_import.CheckChkKeyMatch(a.Check, key) {
			continue
		}

		fieldReal := vReal.Field(i)
		fieldExpected := vExpected.Field(i)

		if !reflect.DeepEqual(fieldReal.Interface(), fieldExpected.Interface()) {
			return errors.New(fmt.Sprintf(unexpectedSentence, fieldExpected.Interface(), fieldReal.Interface()))
		}
	}

	return nil
}
