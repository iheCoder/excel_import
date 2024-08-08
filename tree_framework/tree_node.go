package tree_framework

import "fmt"

var (
	defaultKeyGen = genNodeKey
)

type rawCellContent struct {
	val    string
	isLeaf bool
}

type rawCellWhole struct {
	contents     [][]string
	cellContents [][]rawCellContent
}

type treeNode struct {
	value    string
	parent   *treeNode
	rank     int
	children []*treeNode
	item     any
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
	genKeyFunc generateNodeKey
	// the start row of the content
	startRow int
	// the end condition of the function
	ef endFunc
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
