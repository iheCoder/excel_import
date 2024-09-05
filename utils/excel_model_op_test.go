package util

import (
	"fmt"
	"testing"
)

func TestFmt(t *testing.T) {
	info := &StructInfo{
		Name: "test",
		Fields: []Field{
			{
				Name:    "ID",
				Type:    "int",
				Comment: "id",
			},
			{
				Name:    "Name",
				Type:    "string",
				Comment: "name",
			},
		},
	}
	fmt.Println(info.String())
}

func TestStructInfo_Fmt(t *testing.T) {
	type testData struct {
		info     *StructInfo
		expected string
	}

	tests := []testData{
		{
			info: &StructInfo{
				Name: "test",
				Fields: []Field{
					{
						Name:    "ID",
						Type:    "int",
						Comment: "id",
					},
					{
						Name:    "Name",
						Type:    "string",
						Comment: "name",
					},
				},
			},
			expected: "type test struct {\n\tID int // id\n\tName string // name\n}\n",
		},
		{
			info: &StructInfo{
				Name: "test",
				Fields: []Field{
					{
						Name:    "A",
						Type:    "int",
						Comment: "你好",
					},
				},
			},
			expected: "type test struct {\n\tA int // 你好\n}\n",
		},
	}

	for _, test := range tests {
		if test.info.String() != test.expected {
			t.Errorf("expected:\n %s, got:\n %s", test.expected, test.info.String())
		}
	}
}
