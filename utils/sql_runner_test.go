package util

import "testing"

type testStruct struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func TestGenerateInsertSQLWithValues(t *testing.T) {
	type testData struct {
		ts  *testStruct
		sql string
	}

	tests := []testData{
		{
			ts: &testStruct{
				Name: "hello",
				Age:  10,
			},
			sql: "INSERT INTO test_struct (name, age) VALUES ('hello', 10);",
		},
		{
			ts: &testStruct{
				Name: "world",
			},
			sql: "INSERT INTO test_struct (name) VALUES ('world');",
		},
		{
			ts: &testStruct{
				Age: 10,
			},
			sql: "INSERT INTO test_struct (age) VALUES (10);",
		},
	}

	for _, tt := range tests {
		t.Run("TestGenerateInsertSQLWithValues", func(t *testing.T) {
			if got := GenerateInsertSQLWithValues("test_struct", tt.ts); got != tt.sql {
				t.Errorf("GenerateInsertSQLWithValues() = %v, want %v", got, tt.sql)
			}
		})
	}
}
