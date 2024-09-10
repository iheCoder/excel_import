package util

import (
	"excel_import"
	"fmt"
	"strconv"
	"strings"
)

type typeCode struct {
	Type string
	Is   bool
}

var (
	strTypeCode = typeCode{
		Type: "string",
		Is:   true,
	}
	intTypeCode = typeCode{
		Type: "int",
		Is:   true,
	}
	floatTypeCode = typeCode{
		Type: "float",
		Is:   true,
	}
)

func GenerateExcelModelString(path string, structName string) (string, error) {
	// read excel content
	contents, err := ReadExcelContent(path)
	if err != nil {
		return "", err
	}

	// get field comments
	fieldComments := make([]string, 0)
	for _, v := range contents[0] {
		v = FormatCell(v)
		if len(v) == 0 {
			break
		}

		fieldComments = append(fieldComments, v)
	}

	// end row
	n := len(fieldComments)
	contents = contents[1:]
	for i, v := range contents {
		if DefaultRowEndFunc(v) {
			contents = contents[:i]
			break
		}
	}

	// revise contents column length as n
	for i, v := range contents {
		if len(v) > n {
			contents[i] = v[:n]
		}
		if len(v) < n {
			contents[i] = append(v, make([]string, n-len(v))...)
		}
	}

	// inverse the row and column
	contents = (contents)

	return GenerateStructString(structName, fieldComments, ReverseMatrix(contents)), nil
}

// GenerateStructString generates a string representation of a struct.
// The struct is defined by the structName and the contents.
// the field name is english word ascending order
// the field type is check by the contents
// the field comment is the fieldComment combined with number ascending order
func GenerateStructString(structName string, fieldComment []string, contents [][]string) string {
	if len(contents) == 0 {
		return ""
	}

	info := excel_import.StructInfo{
		Name: structName,
	}
	for i, v := range contents {
		info.Fields = append(info.Fields, excel_import.Field{
			Name:    getExcelColIndex(i),
			Type:    detectType(v),
			Comment: fmt.Sprintf("%s\t%d", fieldComment[i], i),
		})
	}

	return info.String()
}

// getExcelColIndex returns the excel col index for the given number.
// e.g. 0 -> "A", 1 -> "B", 25 -> "Z", 26 -> "AA", 27 -> "AB", 3752126 -> "HELLO"
func getExcelColIndex(i int) string {
	if i < 0 {
		return ""
	}

	const base = 26
	const a = 65
	var sb strings.Builder

	for i >= 0 {
		sb.WriteByte(byte(i%base) + a)
		i /= base
		i--
	}

	return reverse(sb.String())
}

// TranslateNumIndexByExcelColumn translates the excel column to the number index.
// e.g. "A" -> 0, "B" -> 1, "Z" -> 25, "AA" -> 26, "AB" -> 27, "HELLO" -> 3752126
func TranslateNumIndexByExcelColumn(s string) int {
	if len(s) == 0 {
		panic("empty string")
	}

	// check if the string is valid
	for i := 0; i < len(s); i++ {
		if s[i] < 'A' || s[i] > 'Z' {
			panic("invalid string")
		}
	}

	// translate the string to number
	const base = 26
	const a = 65
	var sum int
	for i := 0; i < len(s); i++ {
		sum = sum*base + int(s[i]-a) + 1
	}

	return sum - 1
}

func reverse(s string) string {
	var sb strings.Builder
	for i := len(s) - 1; i >= 0; i-- {
		sb.WriteByte(s[i])
	}
	return sb.String()
}

func resetTypeCodes() {
	strTypeCode.Is = true
	intTypeCode.Is = true
	floatTypeCode.Is = true
}

func checkType(s string) {
	if len(s) == 0 {
		return
	}

	if intTypeCode.Is {
		_, err := strconv.Atoi(s)
		if err != nil {
			intTypeCode.Is = false
		}
	}

	if floatTypeCode.Is {
		_, err := strconv.ParseFloat(s, 64)
		if err != nil {
			floatTypeCode.Is = false
		}
	}
}

func ensureType() string {
	defer resetTypeCodes()

	if intTypeCode.Is {
		return intTypeCode.Type
	}
	if floatTypeCode.Is {
		return floatTypeCode.Type
	}

	return strTypeCode.Type
}

func detectType(values []string) string {
	if len(values) == 0 {
		return strTypeCode.Type
	}

	for _, v := range values {
		checkType(v)
	}

	return ensureType()
}
