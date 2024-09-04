package util

import (
	"fmt"
	"testing"
)

func TestWriteExcelContent(t *testing.T) {
	path := "../testdata/excel_test_data.xlsx"
	content := map[int][]string{
		1: {"1", "2", "3", "4", "5"},
	}
	err := WriteExcelColumnContent(path, content)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDivideSheetsIntoTables(t *testing.T) {
	path := "../testdata/excel_test_sheets.xlsx"
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
		"../testdata/excel_test_sheets_Sheet1.xlsx",
		"../testdata/excel_test_sheets_Sheet2.xlsx",
		"../testdata/excel_test_sheets_Sheet3.xlsx",
	}

	err := CombineTablesIntoOne(paths...)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("done")
}

func TestDivideExcelContent(t *testing.T) {
	path := "../testdata/excel_test_data.xlsx"
	contents, err := DivideExcelContent(path, 2)
	if err != nil {
		t.Fatal(err)
	}

	for _, content := range contents {
		fmt.Println(content)
	}
}

func TestDivideMultipleTreesIntoMultipleTables(t *testing.T) {
	path := "../testdata/excel_tree_mul_tree_data.xlsx"
	ukColIndex := []int{0, 1, 2, 3}
	paths, err := DivideMultipleTreesIntoMultipleTables(path, "../testdata/trees_output", ukColIndex)
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 8 {
		t.Fatalf("expected 8, got %d", len(paths))
	}

	t.Log("done")
}

func TestSetHyperlinksInColumn(t *testing.T) {
	path := "../testdata/excel_test_data.xlsx"
	urls := []string{
		"https://www.baidu.com",
		"https://www.google.com",
		"https://www.bing.com",
		"https://www.yahoo.com",
		"https://www.sogou.com",
	}
	index := 2

	if err := SetHyperlinksInColumn(path, urls, index); err != nil {
		t.Fatal(err)
	}

	t.Log("done")
}

func TestDivideSheetsIntoTablesBySuffixKey(t *testing.T) {
	path := "../testdata/excel_test_divide.xlsx"
	paths, err := DivideSheetsIntoTablesByDefaultSuffixKey(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 3 {
		t.Fatalf("expected 3, got %d", len(paths))
	}

	for _, p := range paths {
		fmt.Println(p)
	}
}
