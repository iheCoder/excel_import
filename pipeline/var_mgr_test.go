package pipeline

import "testing"

func TestGenerateVarNameByUpperCase(t *testing.T) {
	type testData struct {
		typeName string
		expected string
	}

	tests := []testData{
		{
			typeName: "VarMgr",
			expected: "vm",
		},
		{
			typeName: "ModelGraph",
			expected: "mg",
		},
		{
			typeName: "Field",
			expected: "f",
		},
		{
			typeName: "X",
			expected: "x",
		},
		{
			typeName: "private",
			expected: "private",
		},
		{
			typeName: "",
			expected: "",
		},
	}

	for _, test := range tests {
		if got := GenerateVarNameByUpperCase(test.typeName); got != test.expected {
			t.Fatalf("expect %s, but got %s", test.expected, got)
		}
	}
}

func TestGenerateVarNameByLastWord(t *testing.T) {
	type testData struct {
		typeName string
		expected string
	}

	tests := []testData{
		{
			typeName: "VarMgr",
			expected: "mgr",
		},
		{
			typeName: "ModelGraph",
			expected: "graph",
		},
		{
			typeName: "Field",
			expected: "field",
		},
		{
			typeName: "X",
			expected: "x",
		},
		{
			typeName: "private",
			expected: "private",
		},
		{
			typeName: "",
			expected: "",
		},
	}

	for _, test := range tests {
		if got := GenerateVarNameByLastWord(test.typeName); got != test.expected {
			t.Fatalf("expect %s, but got %s", test.expected, got)
		}
	}
}
