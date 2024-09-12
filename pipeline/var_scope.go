package pipeline

type scope struct {
	key             string
	parent          *scope
	children        []*scope
	parentScopeKeys []string
}
