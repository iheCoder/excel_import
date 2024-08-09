package tree_framework

import "fmt"

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
	root         *treeNode
}

type treeNode struct {
	value    string
	parent   *treeNode
	rank     int
	children []*treeNode
	item     any
}

func (t *treeNode) CheckIsLeaf() bool {
	return len(t.children) == 0
}

func constructLevelNode(s string, parent *treeNode, level int) *treeNode {
	node := &treeNode{
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
