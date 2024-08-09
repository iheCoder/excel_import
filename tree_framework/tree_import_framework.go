package tree_framework

import (
	"excel_import/utils"
	"fmt"
	"gorm.io/gorm"
)

type treeImportFramework struct {
	db            *gorm.DB
	recorder      *util.UnexpectedRecorder
	cfg           *treeImportCfg
	nodes         map[string]*TreeNode
	levelImporter []LevelImporter
	ocfg          *treeImportOptionalCfg
}

func NewTreeImportFramework(db *gorm.DB, cfg *treeImportCfg, levelImporter []LevelImporter, options ...optionFunc) *treeImportFramework {
	if cfg == nil {
		panic("cfg should not nil")
	}

	tif := &treeImportFramework{
		db:            db,
		cfg:           cfg,
		nodes:         make(map[string]*TreeNode),
		levelImporter: levelImporter,
		ocfg:          defaultOptCfg,
	}

	for _, option := range options {
		option(tif)
	}

	if tif.cfg.modelFac == nil {
		panic("model factory should not nil")
	}

	return tif
}

func WithGenKeyFunc(gkf generateNodeKey) optionFunc {
	return func(framework *treeImportFramework) {
		framework.ocfg.genKeyFunc = gkf
	}
}

func WithStartRow(sr int) optionFunc {
	return func(framework *treeImportFramework) {
		framework.ocfg.startRow = sr
	}
}

func WithEndFunc(ef rowEndFunc) optionFunc {
	return func(framework *treeImportFramework) {
		framework.ocfg.ef = ef
	}
}

func WithColEndFunc(cf colEndFunc) optionFunc {
	return func(framework *treeImportFramework) {
		framework.ocfg.cf = cf
	}
}

func (t *treeImportFramework) parseRawWhole(content [][]string) (*rawCellWhole, error) {
	// construct the tree
	root, err := t.constructTree(content)
	if err != nil {
		return nil, err
	}

	cellContents := make([][]rawCellContent, len(content))
	models := make([]any, len(content))
	for i, row := range content {
		// parse the cell content
		cellContents[i] = make([]rawCellContent, len(row))
		for j, cell := range row {
			cellContents[i][j] = rawCellContent{val: cell, isLeaf: t.checkIsLeaf(i, row)}
		}

		// parse the content into models
		var model any
		if t.cfg.modelFac != nil {
			model = t.cfg.modelFac.GetModel()
			if err := util.FillModelOrder(model, row); err != nil {
				return nil, err
			}
		}
		models[i] = model
	}

	// fill the model into the leaf tree node
	t.fillModelIntoLeafNode(root, models)

	return &rawCellWhole{
		contents:     content,
		cellContents: cellContents,
		root:         root,
	}, nil
}

func (t *treeImportFramework) fillModelIntoLeafNode(root *TreeNode, models []any) {
	if root == nil {
		return
	}

	if root.CheckIsLeaf() {
		root.item = models[root.row-t.ocfg.startRow]
		return
	}

	for _, child := range root.children {
		t.fillModelIntoLeafNode(child, models)
	}
}

func (t *treeImportFramework) checkIsLeaf(i int, row []string) bool {
	if i == t.cfg.boundary {
		return true
	}

	var next string
	if i+1 < len(row) {
		next = row[i+1]
	}
	return t.ocfg.cf(next)
}

func (t *treeImportFramework) importTree(whole *rawCellWhole) error {
	root := whole.root

	// import the tree
	children := root.children
	for _, importer := range t.levelImporter {
		nextChildren := make([]*TreeNode, 0)
		for _, child := range children {
			if err := importer.ImportLevelNode(t.db, child); err != nil {
				return err
			}
			nextChildren = append(nextChildren, child.children...)
		}

		children = nextChildren
	}

	return nil
}

func (t *treeImportFramework) constructTree(rcContents [][]string) (*TreeNode, error) {
	// remove the column which is not belongs to the tree
	rcContents = rcContents[:t.cfg.boundary+1]

	// reverse the matrix
	contents := reverseMatrix(rcContents)

	// construct the tree
	root := &TreeNode{}
	parent := root
	for level, i := range t.cfg.levelOrder {
		for j, s := range contents[i] {
			// if current node has been constructed, skip it
			curKey := t.ocfg.genKeyFunc(rcContents[j][:i+1], level+1)
			if _, ok := t.nodes[curKey]; ok {
				continue
			}

			// find the parent node
			if level > 0 {
				porder := t.cfg.levelOrder[level-1]
				parent = t.findParent(rcContents[j][:porder+1], level)
			}
			if parent == nil {
				return nil, fmt.Errorf("parent not found for %s", s)
			}

			// construct the node
			node := constructLevelNode(s, parent, level+1, j+t.ocfg.startRow)
			t.nodes[curKey] = node
		}
	}

	return root, nil
}

func (t *treeImportFramework) findParent(s []string, level int) *TreeNode {
	key := t.ocfg.genKeyFunc(s, level)
	if node, ok := t.nodes[key]; ok {
		return node
	}
	return nil
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
