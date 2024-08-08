package util

import (
	"fmt"
	"testing"
)

func TestWriteExcelContent(t *testing.T) {
	path := "excel_test_data.xlsx"
	content := map[int][]string{
		1: {"1", "2", "3", "4", "5"},
	}
	err := WriteExcelColumnContent(path, content)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDivideSheetsIntoTables(t *testing.T) {
	path := "excel_test_sheets.xlsx"
	paths, err := DivideSheetsIntoTables(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range paths {
		fmt.Println(p)
	}
}

func TestCombineTablesIntoOne(t *testing.T) {
	paths := []string{
		"excel_test_sheets_Sheet1.xlsx",
		"excel_test_sheets_Sheet2.xlsx",
		"excel_test_sheets_Sheet3.xlsx",
	}

	err := CombineTablesIntoOne(paths...)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("done")
}
