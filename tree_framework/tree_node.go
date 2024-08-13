package tree_framework

import (
	"excel_import"
	"fmt"
)

var (
	defaultKeyGen = genNodeKey
	defaultOptCfg = &treeImportOptionalCfg{
		genKeyFunc:     defaultKeyGen,
		startRow:       1,
		treeColEndFunc: defaultTreeColEndFunc,
		ef:             defaultRawEndFunc,
	}
)

type rawCellContent struct {
	val    string
	isLeaf bool
}

type rawCellWhole struct {
	contents        [][]string
	cellContents    [][]rawCellContent
	root            *TreeNode
	totalModelCount int
}

func (r *rawCellWhole) GetRoot() *TreeNode {
	return r.root
}

func (r *rawCellWhole) GetNodeCount() int {
	return r.totalModelCount
}

func (r *rawCellWhole) GetLeafCount() int {
	return len(r.contents)
}

type TreeNode struct {
	id       int64
	row      int
	value    string
	parent   *TreeNode
	rank     int
	children []*TreeNode
	item     any
}

func (t *TreeNode) GetItem() any {
	return t.item
}

func (t *TreeNode) GetValue() string {
	return t.value
}

func (t *TreeNode) GetParent() *TreeNode {
	return t.parent
}

func (t *TreeNode) GetChildren() []*TreeNode {
	return t.children
}

func (t *TreeNode) CheckIsLeaf() bool {
	return len(t.children) == 0
}

func (t *TreeNode) CheckIsRoot() bool {
	return t.parent == nil
}

func (t *TreeNode) GetID() int64 {
	return t.id
}

// SetID set the id of the tree node
// should be called after import the tree node
func (t *TreeNode) SetID(id int64) {
	t.id = id
}

func (t *TreeNode) GetRank() int {
	return t.rank
}

func (t *TreeNode) GetRow() int {
	return t.row
}

func constructLevelNode(s string, parent *TreeNode, level int, row int) *TreeNode {
	node := &TreeNode{
		value:  s,
		parent: parent,
		rank:   level,
		row:    row,
	}
	parent.children = append(parent.children, node)
	return node
}

type TreeImportCfg struct {
	LevelOrder []int
	// the TreeBoundary of the tree node. the last index that belongs to the tree
	TreeBoundary int
	// the column count of the raw model
	ColumnCount int
	// the rawModel factory
	ModelFac excel_import.RowModelFactory
}

type treeImportOptionalCfg struct {
	genKeyFunc GenerateNodeKey
	// the start row of the content
	startRow int
	// the end condition of the function
	ef RowEndFunc
	// the end condition of the tree column
	treeColEndFunc ColEndFunc
}

func genNodeKey(s []string, level int) string {
	return fmt.Sprintf("%s_%d", s[len(s)-1], level)
}

func genPrefixNodeKey(s []string, level int) string {
	var key string
	for _, x := range s {
		key += x + "_"
	}
	key += fmt.Sprintf("%d", level)
	return key
}

func defaultTreeColEndFunc(next string) bool {
	return len(next) == 0
}

func defaultRawEndFunc(s []string) bool {
	return len(s) == 0 || len(s[0]) == 0
}
