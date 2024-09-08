package pipeline

// FieldNode is a node in the graph
type FieldNode struct {
	StructName string
	FieldName  string
}

// RelationAdapter is a function adapter for relation models
type RelationAdapter struct {
	f func(from any) any
}

// ModelGraph is a graph of models
type ModelGraph struct {
	Edges map[FieldNode]map[FieldNode]RelationAdapter
}

func NewModelGraph() *ModelGraph {
	return &ModelGraph{
		Edges: make(map[FieldNode]map[FieldNode]RelationAdapter),
	}
}

// AddEdge adds an edge to the graph
func (mg *ModelGraph) AddEdge(from, to FieldNode, adapter RelationAdapter) {
	if mg.Edges[from] == nil {
		mg.Edges[from] = make(map[FieldNode]RelationAdapter)
	}
	mg.Edges[from][to] = adapter
}

// GetEdgeAdapter gets an edge from the graph
func (mg *ModelGraph) GetEdgeAdapter(from, to FieldNode) (RelationAdapter, bool) {
	adapter, ok := mg.Edges[from][to]
	return adapter, ok
}

// GetEdge gets all edges from a node
func (mg *ModelGraph) GetEdge(from FieldNode) map[FieldNode]RelationAdapter {
	return mg.Edges[from]
}

// GetOneEdge gets one edge from a node
func (mg *ModelGraph) GetOneEdge(from FieldNode) (FieldNode, RelationAdapter, bool) {
	for to, adapter := range mg.Edges[from] {
		return to, adapter, true
	}

	return FieldNode{}, RelationAdapter{}, false
}
