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
	value    string
	parent   *TreeNode
	rank     int
	children []*TreeNode
	item     any
}

func (t *TreeNode) GetItem() any {
	return t.item
}

func (t *TreeNode) CheckIsLeaf() bool {
	return len(t.children) == 0
}

func constructLevelNode(s string, parent *TreeNode, level int) *TreeNode {
	node := &TreeNode{
		value:  s,
		parent: parent,
		rank:   level,
	}
	parent.children = append(parent.children, node)
	return node
}

type treeImportCfg struct {
	levelOrder []int
	// the boundary of the tree node
	boundary int
	// the model factory
	modelFac excel_import.RowModelFactory
}

type treeImportOptionalCfg struct {
	genKeyFunc generateNodeKey
	// the start row of the content
	startRow int
	// the end condition of the function
	ef rowEndFunc
	// the end condition of the column
	cf colEndFunc
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
