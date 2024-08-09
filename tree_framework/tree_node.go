package tree_framework

import (
	"excel_import"
	"fmt"
)

var (
	defaultKeyGen = genNodeKey
	defaultOptCfg = &treeImportOptionalCfg{
		genKeyFunc: defaultKeyGen,
		startRow:   1,
		cf:         defaultRowEndFunc,
	}
)

type rawCellContent struct {
	val    string
	isLeaf bool
}

type rawCellWhole struct {
	contents     [][]string
	cellContents [][]rawCellContent
	root         *TreeNode
}

type TreeNode struct {
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
	// the TreeBoundary of the tree node
	TreeBoundary int
	// the rawModel factory
	ModelFac excel_import.RowModelFactory
}

type treeImportOptionalCfg struct {
	genKeyFunc GenerateNodeKey
	// the start row of the content
	startRow int
	// the end condition of the function
	ef RowEndFunc
	// the end condition of the column
	cf ColEndFunc
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

func defaultRowEndFunc(next string) bool {
	return len(next) == 0
}
