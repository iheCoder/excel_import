package pipeline

type scope struct {
	key      string
	parent   *scope
	children []*scope
	vars     map[string]bool
}

func newScope(key string, parent *scope) *scope {
	return &scope{
		key:      key,
		parent:   parent,
		children: make([]*scope, 0),
		vars:     make(map[string]bool),
	}
}

// addChild adds a child to the scope.
func (s *scope) addChild(child *scope) {
	child.parent = s
	s.children = append(s.children, child)
}

// addVar adds a var to the scope.
func (s *scope) addVar(varName string) bool {
	if s.checkVarConflict(varName) {
		return false
	}

	s.vars[varName] = true
	return true
}

// checkVarConflict checks if the var name conflicts with the root to current node path and current node's children.
func (s *scope) checkVarConflict(varName string) bool {
	if _, ok := s.vars[varName]; ok {
		return true
	}

	parent := s.parent
	for parent != nil {
		if _, ok := parent.vars[varName]; ok {
			return true
		}
		parent = parent.parent
	}

	for _, child := range s.children {
		if child.checkVarConflict(varName) {
			return true
		}
	}

	return false
}
