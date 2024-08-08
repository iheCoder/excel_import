package excel_import

import (
	"fmt"
	"gorm.io/gorm"
)

type treeImportFramework struct {
	db            *gorm.DB
	recorder      *unexpectedRecorder
	cfg           *orderLevelCfg
	nodes         map[string]*treeNode
	levelImporter []LevelImporter
}

func NewTreeImportFramework(db *gorm.DB, cfg *orderLevelCfg) *treeImportFramework {
	return &treeImportFramework{
		db:    db,
		cfg:   cfg,
		nodes: make(map[string]*treeNode),
	}
}

func (t *treeImportFramework) importTree(contents [][]string) error {
	root, err := t.constructTree(contents)
	if err != nil {
		return err
	}

	children := root.children
	for _, importer := range t.levelImporter {
		nextChildren := make([]*treeNode, 0)
		for _, child := range children {
			if err = importer.ImportLevelNode(t.db, child); err != nil {
				return err
			}
			nextChildren = append(nextChildren, child.children...)
		}

		children = nextChildren
	}

	return nil
}

func (t *treeImportFramework) constructTree(contents [][]string) (*treeNode, error) {
	// reverse the matrix
	contents = reverseMatrix(contents)

	// construct the tree
	root := &treeNode{}
	parent := root
	for level, i := range t.cfg.levelOrder {
		for j, s := range contents[i] {
			curKey := genNodeKey(s, level+1)
			if _, ok := t.nodes[curKey]; ok {
				continue
			}

			if level > 0 {
				porder := t.cfg.levelOrder[level-1]
				parent = t.findParent(contents[porder][j], level)
			}

			if parent == nil {
				return nil, fmt.Errorf("parent not found for %s", s)
			}

			node := constructLevelNode(s, parent, level+1)
			t.nodes[genNodeKey(s, level+1)] = node
		}
	}

	return root, nil
}

func (t *treeImportFramework) findParent(s string, level int) *treeNode {
	key := genNodeKey(s, level)
	if node, ok := t.nodes[key]; ok {
		return node
	}
	return nil
}

func genNodeKey(s string, level int) string {
	return fmt.Sprintf("%s_%d", s, level)
}

func reverseMatrix(contents [][]string) [][]string {
	if len(contents) == 0 {
		return contents
	}

	n := len(contents[0])
	m := len(contents)
	res := make([][]string, n)
	for i := 0; i < n; i++ {
		res[i] = make([]string, m)
		for j := 0; j < m; j++ {
			res[i][j] = contents[j][i]
		}
	}

	return res
}
