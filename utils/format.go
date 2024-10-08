package util

import "strings"

func FormatCell(cell string) string {
	return strings.TrimSpace(removeLargeUnicodeChars(cell))
}

// removeLargeUnicodeChars remove all characters with Unicode code points greater than \U00100000
func removeLargeUnicodeChars(s string) string {
	result := make([]rune, 0, len(s))
	for _, r := range s {
		if r < 0x100000 {
			result = append(result, r)
		}
	}
	return string(result)
}

func ReverseMatrix(contents [][]string) [][]string {
	if len(contents) == 0 {
		return contents
	}

	n := getMinColCount(contents)
	m := len(contents)
	res := make([][]string, n)
	for i := 0; i < n; i++ {
		res[i] = make([]string, m)
		for j := 0; j < m; j++ {
			res[i][j] = contents[j][i]
		}
	}

	return res
}

func getMinColCount(contents [][]string) int {
	result := len(contents[0])
	for _, v := range contents {
		if len(v) < result {
			result = len(v)
		}
	}

	return result
}
