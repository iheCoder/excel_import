package excel_import

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

type orderLevelCfg struct {
	levelOrder []int
}
