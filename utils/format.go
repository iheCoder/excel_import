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
