package util

import "reflect"

type SimpleModelFactory struct {
	maxColumnCount int
	elemType       reflect.Type
}

func NewSimpleModelFactory(model any) *SimpleModelFactory {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		panic("input is not a pointer to a struct")
	}

	tags := ParseTag(model)
	maxColumnCount := 0
	for _, tag := range tags {
		if tag.ColumnIndex > maxColumnCount {
			maxColumnCount = tag.ColumnIndex
		}
	}
	maxColumnCount++

	elemType := v.Elem().Type()

	return &SimpleModelFactory{
		maxColumnCount: maxColumnCount,
		elemType:       elemType,
	}
}

func (s *SimpleModelFactory) GetModel() any {
	return reflect.New(s.elemType).Interface()
}

func (s *SimpleModelFactory) MinColumnCount() int {
	return s.maxColumnCount
}
