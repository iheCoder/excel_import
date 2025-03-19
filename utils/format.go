package util

import "strings"

func FormatCell(cell string) string {
	return strings.TrimSpace(removeSpecialChars(cell))
}

// removeSpecialChars remove all characters with Unicode code points greater than \U00100000
// and remove zero-width characters
func removeSpecialChars(s string) string {
	// zeroWidthChars is a list of zero-width characters
	zeroWidthChars := []rune{
		0x200B, // ZERO WIDTH SPACE
		0x200C, // ZERO WIDTH NON-JOINER
		0x200D, // ZERO WIDTH JOINER
		0xFEFF, // ZERO WIDTH NO-BREAK SPACE
	}

	result := make([]rune, 0, len(s))
	for _, r := range s {
		// check if the character is greater than 0x100000
		if r >= 0x100000 {
			continue
		}

		// check if the character is a zero-width character
		isZeroWidth := false
		for _, zw := range zeroWidthChars {
			if r == zw {
				isZeroWidth = true
				break
			}
		}
		// if the character is not a zero-width character, add it to the result
		if !isZeroWidth {
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
