package excel_import

type treeNode struct {
	value    string
	parent   *treeNode
	rank     int
	children []*treeNode
	item     any
}

func constructLevelNode(s []string, parent *treeNode, level int) {
	m := make(map[string]bool)
	for _, v := range s {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = true
		parent.children = append(parent.children, &treeNode{parent: parent, rank: level, value: v})
	}

	return
}
