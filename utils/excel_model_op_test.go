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

func TestGenerateStructString(t *testing.T) {
	type testData struct {
		structName   string
		fieldComment []string
		contents     [][]string
		expected     string
	}

	tests := []testData{
		{
			structName: "test",
			fieldComment: []string{
				"id",
				"name",
			},
			contents: [][]string{
				{"1", "2"},
				{"a", "b"},
			},
			expected: "type test struct {\n\tA int // id\t0\n\tB string // name\t1\n}\n",
		},
		{
			structName: "test",
			fieldComment: []string{
				"你好",
			},
			contents: [][]string{
				{"1"},
			},
			expected: "type test struct {\n\tA int // 你好\t0\n}\n",
		},
	}

	for _, test := range tests {
		if GenerateStructString(test.structName, test.fieldComment, test.contents) != test.expected {
			t.Errorf("expected:\n %s, got:\n %s", test.expected, GenerateStructString(test.structName, test.fieldComment, test.contents))
		}
	}
}

func TestGenerateExcelModelString(t *testing.T) {
	path := "../testdata/excel_test_resource.xlsx"
	structName := "resourceTest"
	expected := "type resourceTest struct {\n\tA string // 名称\t0\n\tB int // 类型\t1\n\tC string // 饮料品牌\t2\n\tD string // 鞋品牌\t3\n\tE int // 尺码\t4\n\tF string // 建筑额外信息\t5\n}\n"
	structStr, err := GenerateExcelModelString(path, structName)
	if err != nil {
		t.Fatal(err)
	}

	if structStr != expected {
		t.Errorf("expected:\n %s, got:\n %s", expected, structStr)
	}
}
