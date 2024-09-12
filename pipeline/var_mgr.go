package pipeline

import (
	"strings"
	"unicode"
)

// Keywords contains Go's reserved keywords that cannot be used as variable names.
var (
	goKeywords = map[string]struct{}{
		"break": {}, "default": {}, "func": {}, "interface": {}, "select": {},
		"case": {}, "defer": {}, "go": {}, "map": {}, "struct": {},
		"chan": {}, "else": {}, "goto": {}, "package": {}, "switch": {},
		"const": {}, "fallthrough": {}, "if": {}, "range": {}, "type": {},
		"continue": {}, "for": {}, "import": {}, "return": {}, "var": {},
	}
)

type VarMgr struct {
}

type scope struct {
	vars   map[string]bool
	parent *scope
}

func (s *scope) addVar(varName string) {
	s.vars[varName] = true
}

func (s *scope) hasVar(varName string) bool {
	_, ok := s.vars[varName]
	return ok
}

func (s *scope) removeVar(varName string) {
	delete(s.vars, varName)
}

// GenerateVarNameByUpperCase generates a string by the upper case of the input string.
func GenerateVarNameByUpperCase(typeName string) string {
	if len(typeName) == 0 {
		return ""
	}

	var result []rune
	for _, r := range typeName {
		if unicode.IsUpper(r) {
			result = append(result, unicode.ToLower(r))
		}
	}

	// return first letter if no upper case letter
	if len(result) == 0 {
		return string(unicode.ToLower(rune(typeName[0])))
	}

	return string(result)
}

// GenerateVarNameByLastWord generates a string by the last word of the input string.
// The last word is the last continuous upper case characters.
// If there is no upper case character, the input string is returned.
// the first character of the result is lower case.
func GenerateVarNameByLastWord(typeName string) string {
	if len(typeName) == 0 {
		return ""
	}

	words := splitCamelCase(typeName)
	if len(words) > 0 {
		return strings.ToLower(words[len(words)-1])
	}

	return ""
}

// splitCamelCase splits a camel case string into words.
func splitCamelCase(input string) []string {
	var words []string
	var lastPos int
	for i := 1; i < len(input); i++ {
		if unicode.IsUpper(rune(input[i])) {
			words = append(words, input[lastPos:i])
			lastPos = i
		}
	}
	words = append(words, input[lastPos:])
	return words
}

func checkKeywordConflict(varName string) bool {
	_, ok := goKeywords[varName]
	return ok
}
