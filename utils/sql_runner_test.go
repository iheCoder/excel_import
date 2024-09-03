package util

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

type testStruct struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func initDB() *gorm.DB {
	source := "account:password@tcp(host:3306)/english_agent?loc=PRC&charset=utf8mb4&parseTime=True&multiStatements=true"
	db, err := gorm.Open(mysql.Open(source), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func TestSqlRunnerBatchExec(t *testing.T) {
	db := initDB()
	tx := db.Begin()
	sqlPath := "../testdata/sql_runner_test.sql"
	runner := NewSqlSentencesRunner(sqlPath, tx, "test_table")

	datas := []*testStruct{
		{
			Name: "1",
			Age:  10,
		},
		{
			Name: "2",
			Age:  20,
		},
		{
			Name: "3",
			Age:  30,
		},
		{
			Name: "4",
			Age:  40,
		},
		{
			Name: "5",
			Age:  50,
		},
	}

	for _, data := range datas {
		if err := runner.GenerateSqlInsertSentences(data); err != nil {
			tx.Rollback()
			t.Fatal(err)
		}
	}

	if err := runner.RunSqlSentencesWithBatch(2); err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	if err := tx.Commit().Error; err != nil {
		t.Fatal(err)
	}

	t.Log("done")
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

func TestGenerateUpdateSQLWithValues(t *testing.T) {
	type testData struct {
		updates map[string]any
		where   map[string]any
		sql     string
	}

	tests := []testData{
		{
			updates: map[string]any{
				"name": "hello",
				"age":  10,
			},
			where: map[string]any{
				"name": "world",
			},
			sql: "UPDATE test_struct SET name = 'hello', age = 10 WHERE name = 'world';",
		},
		{
			updates: map[string]any{
				"name": "hello",
			},
			where: map[string]any{
				"age": 10,
			},
			sql: "UPDATE test_struct SET name = 'hello' WHERE age = 10;",
		},
		{
			updates: map[string]any{
				"age": 10,
			},
			where: map[string]any{
				"name": "world",
			},
			sql: "UPDATE test_struct SET age = 10 WHERE name = 'world';",
		},
	}

	for _, tt := range tests {
		t.Run("TestGenerateUpdateSQLWithValues", func(t *testing.T) {
			if got := GenerateUpdateSQLWithValues("test_struct", tt.updates, tt.where); got != tt.sql {
				t.Errorf("GenerateUpdateSQLWithValues() = %v, want %v", got, tt.sql)
			}
		})
	}
}

func TestFormatString(t *testing.T) {
	type testData struct {
		s        string
		expected string
	}

	tests := []testData{
		{
			s:        "hello",
			expected: "'hello'",
		},
		{
			s:        "O'Reilly",
			expected: "'O''Reilly'",
		},
		{
			s:        "hello\"x'",
			expected: "'hello\"x'''",
		},
		{
			s:        "hello'x\"",
			expected: "'hello''x\"'",
		},
		{
			s:        "hello#",
			expected: "'hello#'",
		},
	}

	for _, tt := range tests {
		t.Run("TestFormatString", func(t *testing.T) {
			if got := formatValueString(tt.s); got != tt.expected {
				t.Errorf("formatString() = %v, want %v", got, tt.expected)
			}
		})
	}
}
