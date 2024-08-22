package util

import (
	"excel_import"
	"reflect"
	"testing"
)

type ParseTagTest1 struct {
	A string `excel:"index:1"`
	B string `excel:"index:3"`
	C string `excel:"index:5"`
}
type ParseTagTest2 struct {
	A string `excel:"index:1"`
	B string
	C string `excel:"index:3"`
}
type ParseTagTest3 struct {
	A string
	B string
	C string
}

func TestParseTag(t *testing.T) {
	type testData struct {
		st       any
		expected []*excel_import.ExcelImportTagAttr
	}

	tests := []testData{
		{
			st: &ParseTagTest1{},
			expected: []*excel_import.ExcelImportTagAttr{
				{
					ColumnIndex: 1,
				},
				{
					ColumnIndex: 3,
				},
				{
					ColumnIndex: 5,
				},
			},
		},
		{
			st: &ParseTagTest2{},
			expected: []*excel_import.ExcelImportTagAttr{
				{
					ColumnIndex: 1,
				},
				{
					ColumnIndex: 2,
				},
				{
					ColumnIndex: 3,
				},
			},
		},
		{
			st: &ParseTagTest3{},
			expected: []*excel_import.ExcelImportTagAttr{
				{
					ColumnIndex: 0,
				},
				{
					ColumnIndex: 1,
				},
				{
					ColumnIndex: 2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("TestParseTag", func(t *testing.T) {
			if got := ParseTag(tt.st); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseTag() = %v, want %v", got, tt.expected)
			}
		})
	}
}
