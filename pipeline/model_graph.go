package pipeline

import (
	"excel_import"
	util "excel_import/utils"
)

var (
	defaultRelationAdapter = RelationAdapter{
		f: func(from any) any {
			return from
		},
	}
)

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

// GetStructEdges gets all edges from a struct
func (mg *ModelGraph) GetStructEdges(structName string) map[FieldNode]RelationAdapter {
	edges := make(map[FieldNode]RelationAdapter)
	for from, tos := range mg.Edges {
		if from.StructName == structName {
			edges[from] = tos[from]
		}
	}

	return edges
}

// GetOneEdge gets one edge from a node
func (mg *ModelGraph) GetOneEdge(from FieldNode) (FieldNode, RelationAdapter, bool) {
	for to, adapter := range mg.Edges[from] {
		return to, adapter, true
	}

	return FieldNode{}, RelationAdapter{}, false
}

func NewModelGraphOneToMany(one any, many []any) *ModelGraph {
	// parse one, many tags
	oneTags := util.ParseTag(one)
	n := len(many)
	manyTags := make([][]*excel_import.ExcelImportTagAttr, 0, n)
	for _, m := range many {
		manyTags = append(manyTags, util.ParseTag(m))
	}

	// parse one, many struct info
	oneInfo := util.ParseStructInfo(one)
	manyInfos := make([]*excel_import.StructInfo, 0, n)
	for _, m := range many {
		manyInfos = append(manyInfos, util.ParseStructInfo(m))
	}

	// map one tag and related field
	m := make(map[string]excel_import.Field)
	for i, tag := range oneTags {
		m[tag.ID] = oneInfo.Fields[i]
	}

	// match same tag id, add to graph
	graph := NewModelGraph()
	for i, mTags := range manyTags {
		for j, mTag := range mTags {
			if field, ok := m[mTag.ID]; ok {
				oneField := FieldNode{StructName: oneInfo.Name, FieldName: field.Name}
				manyField := FieldNode{StructName: manyInfos[i].Name, FieldName: manyInfos[i].Fields[j].Name}
				graph.AddEdge(oneField, manyField, defaultRelationAdapter)
				graph.AddEdge(manyField, oneField, defaultRelationAdapter)
			}
		}
	}

	return graph
}
