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
