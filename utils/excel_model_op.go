package util

import (
	"fmt"
	"strconv"
	"strings"
)

type Field struct {
	Name    string
	Type    string
	Comment string
}

type StructInfo struct {
	Name   string
	Fields []Field
}

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

// GenerateStructString generates a string representation of a struct.
// The struct is defined by the structName and the contents.
// the field name is english word ascending order
// the field type is check by the contents
// the field comment is the fieldComment combined with number ascending order
func GenerateStructString(structName string, fieldComment []string, contents [][]string) string {
	info := StructInfo{
		Name: structName,
	}

	for i, v := range contents {
		info.Fields = append(info.Fields, Field{
			Name:    getExcelColIndex(i),
			Type:    detectType(v),
			Comment: fmt.Sprintf("%s\t%d", fieldComment[i], i+1),
		})
	}

	return info.String()
}

// getExcelColIndex returns the excel col index for the given number.
// e.g. 0 -> "A", 1 -> "B", 25 -> "Z", 26 -> "AA", 27 -> "AB", 5201314 -> "HELLO"
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

// String implements the fmt.Stringer interface for StructInfo.
func (s StructInfo) String() string {
	var sb strings.Builder
	// Write the struct definition line
	sb.WriteString(fmt.Sprintf("type %s struct {\n", s.Name))

	// Write each field line with type and comment
	for _, field := range s.Fields {
		// Use format: "\tName Type // Comment\n"
		sb.WriteString(fmt.Sprintf("\t%s %s // %s\n", field.Name, field.Type, field.Comment))
	}

	// Close the struct definition
	sb.WriteString("}\n")
	return sb.String()
}
