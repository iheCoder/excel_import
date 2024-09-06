package util

import "testing"

func TestFilterRows(t *testing.T) {
	contents := [][]string{
		{"1", "2", "3"},
		{"4", "5", "6"},
		{"7", "8", "9"},
	}
	filter := func(row []string) bool {
		return row[0] == "4"
	}
	res := FilterRows(contents, filter)
	if len(res) != 2 {
		t.Fatalf("expect 2, but got %d", len(res))
	}
	if res[0][0] != "1" || res[1][0] != "7" {
		t.Fatalf("expect 1 and 7, but got %s and %s", res[0][0], res[1][0])
	}
}
