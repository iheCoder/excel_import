package tree_framework

import (
	"excel_import"
	"excel_import/utils"
	"fmt"
	"gorm.io/gorm"
)

type TreeImportFramework struct {
	db            *gorm.DB
	recorder      *util.UnexpectedRecorder
	cfg           *TreeImportCfg
	nodes         map[string]*TreeNode
	levelImporter []LevelImporter
	// the root importer
	rootImporter     LevelImporter
	ocfg             *treeImportOptionalCfg
	progressReporter *util.ProgressReporter
	postHandler      excel_import.PostHandler
	preHandler       TreePreHandler
	middlewares      []TreeMiddleware
	correctCheckers  []excel_import.CorrectnessChecker
}

func NewTreeImportStrictOrderFramework(db *gorm.DB, treeBoundary, colCount int, modelFac excel_import.RowModelFactory, importer LevelImporter, options ...OptionFunc) *TreeImportFramework {
	treeLevel := treeBoundary + 1
	levelOrder := make([]int, treeLevel)
	for i := 0; i < treeLevel; i++ {
		levelOrder[i] = i
	}

	levelImporter := make([]LevelImporter, treeLevel)
	for i := 0; i < treeLevel; i++ {
		levelImporter[i] = importer
	}

	cfg := &TreeImportCfg{
		TreeBoundary: treeBoundary,
		ModelFac:     modelFac,
		LevelOrder:   levelOrder,
		ColumnCount:  colCount,
	}
	return NewTreeImportFramework(db, cfg, importer, levelImporter, options...)
}

func NewTreeImportFramework(db *gorm.DB, cfg *TreeImportCfg, rootImporter LevelImporter, levelImporter []LevelImporter, options ...OptionFunc) *TreeImportFramework {
	if cfg == nil {
		panic("cfg should not nil")
	}
	if len(cfg.LevelOrder) == 0 {
		panic("level order should not empty")
	}
	if len(levelImporter) == 0 {
		panic("level importer should not empty")
	}
	if len(levelImporter) != len(cfg.LevelOrder) {
		panic("level importer should be equal to level order")
	}

	tif := &TreeImportFramework{
		db:               db,
		cfg:              cfg,
		nodes:            make(map[string]*TreeNode),
		levelImporter:    levelImporter,
		rootImporter:     rootImporter,
		ocfg:             defaultOptCfg,
		recorder:         util.NewDefaultUnexpectedRecorder(),
		progressReporter: util.NewProgressReporter(true),
	}

	for _, option := range options {
		option(tif)
	}

	if tif.cfg.ModelFac == nil {
		panic("rawModel factory should not nil")
	}

	return tif
}

func WithGenKeyFunc(gkf GenerateNodeKey) OptionFunc {
	return func(framework *TreeImportFramework) {
		framework.ocfg.genKeyFunc = gkf
	}
}

func WithStartRow(sr int) OptionFunc {
	return func(framework *TreeImportFramework) {
		framework.ocfg.startRow = sr
	}
}

func WithEndFunc(ef excel_import.EndFunc) OptionFunc {
	return func(framework *TreeImportFramework) {
		framework.ocfg.ef = ef
	}
}

func WithColEndFunc(cf ColEndFunc) OptionFunc {
	return func(framework *TreeImportFramework) {
		framework.ocfg.treeColEndFunc = cf
	}
}

func WithPostHandler(ph excel_import.PostHandler) OptionFunc {
	return func(framework *TreeImportFramework) {
		framework.postHandler = ph
	}
}

func WithPreHandler(ph TreePreHandler) OptionFunc {
	return func(framework *TreeImportFramework) {
		framework.preHandler = ph
	}
}

func WithMiddlewares(middlewares ...TreeMiddleware) OptionFunc {
	return func(framework *TreeImportFramework) {
		framework.middlewares = middlewares
	}
}

func WithRowFilter(f excel_import.RowFilter) OptionFunc {
	return func(framework *TreeImportFramework) {
		framework.ocfg.rowFilterFunc = f
	}
}

func (t *TreeImportFramework) WithOption(option OptionFunc) *TreeImportFramework {
	option(t)
	return t
}

func (t *TreeImportFramework) Import(path string) error {
	defer t.recorder.Flush()
	defer t.progressReporter.Report()

	// parse the content
	whole, err := t.parseContent(path)
	if err != nil {
		fmt.Printf("parse file content failed: %v\n", err)
		return err
	}

	// pre handle the content
	if t.preHandler != nil {
		err = t.preHandler.PreImportHandle(t.db, whole)
		if err != nil {
			fmt.Printf("pre handler failed: %v\n", err)
			return err
		}
	}

	// middleware pre handle
	for _, middleware := range t.middlewares {
		if err = middleware.PreImportHandle(t.db, whole); err != nil {
			fmt.Printf("middleware pre handle failed: %v\n", err)
			return err
		}
	}

	// import the tree
	err = t.importTree(whole)
	if err != nil {
		fmt.Printf("import tree failed: %v\n", err)
		return err
	}

	// middleware post handle
	for _, middleware := range t.middlewares {
		if err = middleware.PostHandle(t.db); err != nil {
			fmt.Printf("middleware post handle failed: %v\n", err)
			return err
		}
	}

	// post handle
	if t.postHandler != nil {
		err = t.postHandler.PostHandle(t.db)
		if err != nil {
			fmt.Printf("post handler failed: %v\n", err)
			return err
		}
	}

	return nil
}

// EnableCorrectnessCheck enable the correctness check.
// must be called before Import.
func (t *TreeImportFramework) EnableCorrectnessCheck(correctnessCheckers ...excel_import.CorrectnessChecker) error {
	if len(correctnessCheckers) == 0 {
		correctnessCheckers = t.correctCheckers
	}

	for _, checker := range t.correctCheckers {
		if err := checker.PreCollect(t.db); err != nil {
			return err
		}
	}

	return nil
}

func (t *TreeImportFramework) parseContent(path string) (*rawCellWhole, error) {
	// read the excel content
	content, err := util.ReadExcelContent(path)
	if err != nil {
		return nil, err
	}

	// pre handle the raw content
	content = t.preHandleRawContent(content)

	// parse the raw content
	return t.parseRawWhole(content)
}

func (t *TreeImportFramework) preHandleRawContent(contents [][]string) [][]string {
	// skip the header default
	if t.ocfg.startRow > 0 {
		contents = contents[t.ocfg.startRow:]
	}

	// end row with func
	if t.ocfg.ef != nil {
		for i, row := range contents {
			if t.ocfg.ef(row) {
				contents = contents[:i]
				break
			}
		}
	}

	// filter the row
	if t.ocfg.rowFilterFunc != nil {
		contents = util.FilterRows(contents, t.ocfg.rowFilterFunc)
	}

	// format the content
	for i, row := range contents {
		// if the content is less than the min column count, complete it with empty string
		if len(row) < t.cfg.ColumnCount {
			row = append(row, make([]string, t.cfg.ColumnCount-len(row))...)
		}

		// format the cell
		fc := util.FormatCell
		if t.ocfg.cellFormatFunc != nil {
			fc = t.ocfg.cellFormatFunc
		}
		for j, cell := range row {
			row[j] = fc(cell)
		}

		contents[i] = row
	}

	return contents
}

func (t *TreeImportFramework) parseRawWhole(content [][]string) (*rawCellWhole, error) {
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
			if t.ocfg.treeColEndFunc(cell) {
				break
			}

			cellContents[i][j] = rawCellContent{val: cell, isLeaf: t.checkIsLeaf(i, row)}
		}

		// parse the content into models
		var model any
		if t.cfg.ModelFac != nil {
			model = t.cfg.ModelFac.GetModel()
			if err = util.FillModelByTag(model, row); err != nil {
				return nil, err
			}
		}
		models[i] = model
	}

	// fill the rawModel into the leaf tree node
	t.fillModelIntoLeafNode(root, models)

	// calculate the total node count
	totalNodeCount := t.calculateTotalNodeCount(root)

	return &rawCellWhole{
		contents:       content,
		cellContents:   cellContents,
		root:           root,
		totalNodeCount: totalNodeCount,
		models:         models,
	}, nil
}

func (t *TreeImportFramework) calculateTotalNodeCount(root *TreeNode) int {
	if root == nil {
		return 0
	}

	count := 1
	for _, child := range root.children {
		count += t.calculateTotalNodeCount(child)
	}

	return count
}

func (t *TreeImportFramework) fillModelIntoLeafNode(node *TreeNode, models []any) {
	if node == nil {
		return
	}

	if node.CheckIsLeaf() {
		for _, nodeItem := range node.extra.items {
			if nodeItem.row < t.ocfg.startRow || nodeItem.row >= t.ocfg.startRow+len(models) {
				fmt.Printf("row %d out of range\n", nodeItem.row)
				continue
			}

			nodeItem.item = models[nodeItem.row-t.ocfg.startRow]
		}
		return
	}

	for _, child := range node.children {
		t.fillModelIntoLeafNode(child, models)
	}
}

func (t *TreeImportFramework) checkIsLeaf(i int, row []string) bool {
	if i == t.cfg.TreeBoundary {
		return true
	}

	var next string
	if i+1 < len(row) {
		next = row[i+1]
	}
	return t.ocfg.treeColEndFunc(next)
}

func (t *TreeImportFramework) importTree(whole *rawCellWhole) error {
	t.progressReporter.StartProgress(whole.GetNodeCount())

	root := whole.root

	// import the root
	if err := t.importLevelNode(t.rootImporter, root); err != nil {
		return err
	}

	// import the tree
	nodes := root.GetChildren()
	for _, importer := range t.levelImporter {
		nextNodes := make([]*TreeNode, 0)
		for _, node := range nodes {
			if err := t.importLevelNode(importer, node); err != nil {
				return err
			}
			nextNodes = append(nextNodes, node.children...)
		}

		nodes = nextNodes
	}

	return nil
}

func (t *TreeImportFramework) importLevelNode(importer LevelImporter, node *TreeNode) error {
	status := util.ProgressStatusSuccess
	defer t.progressReporter.CommitProgress(1, status)

	if node == nil || importer == nil {
		return nil
	}

	if err := importer.ImportLevelNode(t.db, node); err != nil {
		fmt.Printf("import value %s section failed: %v\n", node.GetValue(), err)
		t.recorder.RecordImportError(util.CombineRowsErrors(node.GetRows(), err))
		status = util.ProgressStatusFailed
		return err
	}

	for _, middleware := range t.middlewares {
		if err := middleware.PostLevelImportHandle(t.db, node); err != nil {
			fmt.Printf("middleware post level import failed: %v\n", err)
			return err
		}
	}

	return nil
}

func (t *TreeImportFramework) constructTree(rcContents [][]string) (*TreeNode, error) {
	// reverse the matrix
	contents := util.ReverseMatrix(rcContents)

	// remove the column which is not belongs to the tree
	contents = contents[:t.cfg.TreeBoundary+1]

	// construct the tree
	root := &TreeNode{}
	parent := root
	for level, i := range t.cfg.LevelOrder {
		for j, s := range contents[i] {
			if t.ocfg.treeColEndFunc(s) {
				continue
			}

			// if current node has been constructed, skip it
			curKey := t.ocfg.genKeyFunc(rcContents[j][:i+1], level+1)
			node, ok := t.nodes[curKey]
			if !ok {
				// find the parent node
				if level > 0 {
					porder := t.cfg.LevelOrder[level-1]
					parent = t.findParent(rcContents[j][:porder+1], level)
				}
				if parent == nil {
					return nil, fmt.Errorf("parent not found for %s", s)
				}

				// construct the node
				node = constructLevelNode(s, parent, level+1)
				t.nodes[curKey] = node
			}

			// add the item into the node
			node.extra.items = append(node.extra.items, &TreeNodeItem{row: j + t.ocfg.startRow})
		}
	}

	return root, nil
}

func (t *TreeImportFramework) findParent(s []string, level int) *TreeNode {
	key := t.ocfg.genKeyFunc(s, level)
	if node, ok := t.nodes[key]; ok {
		return node
	}
	return nil
}

func (t *TreeImportFramework) CheckCorrect() error {
	for _, checker := range t.correctCheckers {
		if err := checker.CheckCorrect(t.db); err != nil {
			return err
		}
	}

	return nil
}
