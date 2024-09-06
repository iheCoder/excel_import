package util

import "excel_import"

func DefaultRowEndFunc(s []string) bool {
	return len(s) == 0 || len(s[0]) == 0
}

func DefaultRowFilter(s []string) bool {
	return false
}

func FilterRows(contents [][]string, filter excel_import.RowFilter) [][]string {
	res := make([][]string, 0, len(contents))
	for _, row := range contents {
		if !filter(row) {
			res = append(res, row)
		}
	}

	return res
}
