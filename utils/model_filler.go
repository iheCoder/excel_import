package util

import (
	"errors"
	"excel_import"
	"fmt"
	"reflect"
	"strconv"
)

func CheckModelOrder(model any, values []string) error {
	fieldOrders := make([]int, len(values))
	for i := 0; i < len(values); i++ {
		fieldOrders[i] = i
	}
	return CheckModel(model, values, fieldOrders)
}

// CheckModel 检查输入的值是否符合模型的定义.允许值为空
func CheckModel(model interface{}, values []string, fieldOrders []int) error {
	v := reflect.ValueOf(model)

	// 检查输入是否是指向结构体的指针
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("input is not a pointer to a struct")
	}

	// 获取结构体的实际值
	v = v.Elem()

	// 检查字段数量是否匹配
	n := v.NumField()
	if n != len(fieldOrders) {
		return errors.New("field count does not match")
	}

	// 检查字段顺序是否正确
	for _, order := range fieldOrders {
		if order < 0 || order >= n || order >= len(values) {
			return errors.New("field order is out of range")
		}
	}

	// 根据字段信息设置字段值
	for i := 0; i < n; i++ {
		if err := checkField(v, i, fieldOrders[i], values[fieldOrders[i]]); err != nil {
			return err
		}
	}

	return nil
}

func checkField(v reflect.Value, i, colIndex int, value string) error {
	// regard as valid if the value is empty
	if len(value) == 0 {
		return nil
	}

	field := v.Field(i)

	switch field.Kind() {
	case reflect.String:
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("第%d列不为整数: %v", colIndex+1, err)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		_, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("第%d列不为无符号整数: %v", colIndex+1, err)
		}

	case reflect.Float32, reflect.Float64:
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("第%d列不为浮点数: %v", colIndex+1, err)
		}

	default:
		return errors.New("unsupported field type")
	}

	return nil
}

func FillModelOrder(model any, values []string) error {
	fieldOrders := make([]int, len(values))
	for i := 0; i < len(values); i++ {
		fieldOrders[i] = i
	}
	return FillModel(model, values, fieldOrders)
}

// FillModelByTag fill model by tag
func FillModelByTag(model any, values []string) error {
	attrs := ParseTag(model)
	fieldOrders := make([]int, len(attrs))
	for i, attr := range attrs {
		fieldOrders[i] = attr.ColumnIndex
	}

	return FillModel(model, values, fieldOrders)
}

// FillModelByTags fill model by tags
func FillModelByTags(tags []*excel_import.ExcelImportTagAttr, model any, values []string) error {
	fieldOrders := make([]int, len(tags))
	for i, tag := range tags {
		fieldOrders[i] = tag.ColumnIndex
	}

	return FillModel(model, values, fieldOrders)
}

func FillModel(model interface{}, values []string, fieldOrders []int) error {
	v := reflect.ValueOf(model)

	// 检查输入是否是指向结构体的指针
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("input is not a pointer to a struct")
	}

	// 获取结构体的实际值
	v = v.Elem()

	// 检查字段数量是否匹配
	fieldNum := v.NumField()

	// 检查字段顺序是否正确
	valNums := len(values)
	for _, order := range fieldOrders {
		if order < 0 || order >= valNums {
			return errors.New("field order is out of range")
		}
	}

	n := min(fieldNum, valNums)
	// 根据字段信息设置字段值
	for i := 0; i < n; i++ {
		if err := setField(v, i, values[fieldOrders[i]]); err != nil {
			return err
		}
	}

	return nil
}

func setField(v reflect.Value, i int, value string) error {
	field := v.Field(i)
	if !field.CanSet() {
		return errors.New("field is unexported")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil && len(value) != 0 {
			return err
		}
		field.SetInt(fieldValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil && len(value) != 0 {
			return err
		}
		field.SetUint(fieldValue)
	case reflect.Float32, reflect.Float64:
		fieldValue, err := strconv.ParseFloat(value, 64)
		if err != nil && len(value) != 0 {
			return err
		}
		field.SetFloat(fieldValue)
	default:
		return errors.New("unsupported field type")
	}

	return nil
}

func NewModel(model any) any {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil
	}

	return reflect.New(v.Elem().Type()).Interface()
}
