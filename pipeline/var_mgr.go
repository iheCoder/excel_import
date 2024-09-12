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

	maxGenTries = 5
)

type VarInfo struct {
	varName string
}

// VarMgr used to manage var generation and check scope conflict
// for simplifying the variable name generation, we have some rules:
// 1. scope will only have two levels: file scope and func scope.
// 2. variable name will be generated by upper case, last word or the random name.
type VarMgr struct {
	// rootScope is the root scope of the var manager.
	// it's represented by the global scope.
	// global scope is the go file scope of the generated code.
	rootScope *scope
	// globalVarPool is the global variable pool.
	globalVarPool map[string]*VarInfo
}

// NewVarMgr creates a new VarMgr.
func NewVarMgr() *VarMgr {
	return &VarMgr{
		rootScope:     newScope("", nil),
		globalVarPool: make(map[string]*VarInfo),
	}
}

// AddScopeAtRoot creates a new scope at the root scope.
func (v *VarMgr) AddScopeAtRoot(key string) {
	newScope(key, v.rootScope)
}

// AddScope creates a new scope at the current scope.
func (v *VarMgr) AddScope(key, parentKey string) bool {
	s := v.findScope(parentKey)
	if s != nil {
		newScope(key, s)
		return true
	}

	return false
}

// GenerateVarNameInScope generates a variable name in the current scope.
func (v *VarMgr) GenerateVarNameInScope(typeName, scopeKey string) (varName string, success bool) {
	var tries int
	s := v.findScope(scopeKey)
	defer func() {
		if success {
			v.addVarInScope(varName, s)
		}
	}()

	// first try to generate the var name by the upper case of the type name
	varName = GenerateVarNameByUpperCase(typeName)
	if !v.checkVarConflictInScope(varName, s) {
		return varName, true
	}
	firstVarName := varName
	tries++

	// then try to generate the var name by the lower case of the first word of the type name
	varName = GenerateVarNameByLowerFirst(typeName)
	if !v.checkVarConflictInScope(varName, s) {
		return varName, true
	}
	tries++

	// later try to generate the var name by the last word of the type name
	varName = GenerateVarNameByLastWord(typeName)
	if !v.checkVarConflictInScope(varName, s) {
		return varName, true
	}
	tries++

	// if the above three methods failed, generate a suffix var name
	suffix := 'A'
	for tries < maxGenTries {
		varName = firstVarName + string(suffix)
		if !v.checkVarConflictInScope(varName, s) {
			return varName, true
		}
		tries++
		suffix++
	}

	return "", false
}

func (v *VarMgr) checkVarConflictInScope(varName string, s *scope) bool {
	// check the var conflict in the global var pool
	if _, ok := v.globalVarPool[varName]; !ok {
		return false
	}

	// check the var conflict in the current scope
	if s.checkVarConflict(varName) {
		return true
	}

	return false
}

func (v *VarMgr) addVarInScope(varName string, s *scope) {
	// add the var
	s.addVar(varName)

	// register the var
	v.globalVarPool[varName] = &VarInfo{
		varName: varName,
	}
}

// AddVarInScope adds a variable to the current scope.
func (v *VarMgr) AddVarInScope(varName, scopeKey string) bool {
	// check the keyword conflict
	if checkKeywordConflict(varName) {
		return false
	}

	// find the scope
	s := v.findScope(scopeKey)
	if s == nil {
		return false
	}

	// check the var name conflict
	if s.checkVarConflict(varName) {
		return false
	}

	// add the var
	v.addVarInScope(varName, s)

	return true
}

// findScope finds the scope by the key.
func (v *VarMgr) findScope(key string) *scope {
	return v.findScopeRecursive(v.rootScope, key)
}

// findScopeRecursive finds the scope by the key recursively.
func (v *VarMgr) findScopeRecursive(s *scope, key string) *scope {
	if s.key == key {
		return s
	}

	for _, c := range s.children {
		if r := v.findScopeRecursive(c, key); r != nil {
			return r
		}
	}

	return nil
}

// GenerateVarNameByLowerFirst generates a string by the lower case of the first word and combines the remains words.
func GenerateVarNameByLowerFirst(typeName string) string {
	if len(typeName) == 0 {
		return ""
	}

	words := splitCamelCase(typeName)
	var varName string
	if len(words) > 0 {
		varName = strings.ToLower(words[0])
		for i := 1; i < len(words); i++ {
			varName += words[i]
		}
	}

	return varName
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
