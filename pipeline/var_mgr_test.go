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
			expected: "p",
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

func TestGenerateVarNameByLowerFirst(t *testing.T) {
	type testData struct {
		typeName string
		expected string
	}

	tests := []testData{
		{
			typeName: "VarMgr",
			expected: "varMgr",
		},
		{
			typeName: "ModelGraph",
			expected: "modelGraph",
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
		if got := GenerateVarNameByLowerFirst(test.typeName); got != test.expected {
			t.Fatalf("expect %s, but got %s", test.expected, got)
		}
	}
}

func TestCheckVarConflictInScope(t *testing.T) {
	v := NewVarMgr()
	childA := "childA"
	childB := "childB"
	grandChildA := "grandChildA"
	grandChildB := "grandChildB"
	grandChildC := "grandChildC"

	v.AddScopeAtRoot(childA)
	v.AddScopeAtRoot(childB)
	v.AddScope(grandChildA, childA)
	v.AddScope(grandChildB, childA)
	v.AddScope(grandChildC, childB)

	v.AddVarInScope("root", "")
	v.AddVarInScope("childA", "childA")
	v.AddVarInScope("childB", "childB")
	v.AddVarInScope("grandChildA", "grandChildA")
	v.AddVarInScope("grandChildB", "grandChildB")
	v.AddVarInScope("grandChildC", "grandChildC")

	// should be no conflict
	noConflictVarName := "noConflict"
	if got := v.checkVarConflictInScope(noConflictVarName, v.findScope("")); got {
		t.Fatalf("expect false, but got true")
	}

	// should be conflict with childA
	rootVarName := "root"
	if got := v.checkVarConflictInScope(rootVarName, v.findScope(childA)); !got {
		t.Fatalf("expect true, but got false")
	}

	// should be no conflict with childB
	childAVarName := "childA"
	if got := v.checkVarConflictInScope(childAVarName, v.findScope(childB)); got {
		t.Fatalf("expect false, but got true")
	}

	// should be conflict with childA
	grandChildAVarName := "grandChildA"
	if got := v.checkVarConflictInScope(grandChildAVarName, v.findScope(childA)); !got {
		t.Fatalf("expect true, but got false")
	}

	// should be conflict with root
	if got := v.checkVarConflictInScope(grandChildAVarName, v.findScope("")); !got {
		t.Fatalf("expect true, but got false")
	}
}
