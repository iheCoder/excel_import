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
