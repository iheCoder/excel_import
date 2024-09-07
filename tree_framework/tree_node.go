package tree_framework

import (
	"excel_import"
	util "excel_import/utils"
	"fmt"
)

var (
	defaultKeyGen = genNodeKey
	defaultOptCfg = &treeImportOptionalCfg{
		genKeyFunc:     defaultKeyGen,
		startRow:       1,
		treeColEndFunc: defaultTreeColEndFunc,
		ef:             util.DefaultRowEndFunc,
	}
)

type rawCellContent struct {
	val    string
	isLeaf bool
}

type rawCellWhole struct {
	contents       [][]string
	cellContents   [][]rawCellContent
	root           *TreeNode
	totalNodeCount int
	models         []any
	modelAttrs     []*excel_import.ExcelImportTagAttr
}

func (r *rawCellWhole) GetModelTags() []*excel_import.ExcelImportTagAttr {
	return r.modelAttrs
}

func (r *rawCellWhole) GetRoot() *TreeNode {
	return r.root
}

func (r *rawCellWhole) GetNodeCount() int {
	return r.totalNodeCount
}

func (r *rawCellWhole) GetLeafCount() int {
	return len(r.contents)
}

func (r *rawCellWhole) GetModels() []any {
	return r.models
}

type TreeNode struct {
	id       int64
	key      string
	value    string
	parent   *TreeNode
	rank     int
	children []*TreeNode
	extra    *TreeNodeExtra
	whole    *rawCellWhole
}

type TreeNodeExtra struct {
	items []*TreeNodeItem
}

type TreeNodeItem struct {
	item any
	row  int
}

func (t *TreeNode) GetKey() string {
	return t.key
}

func (t *TreeNode) GetModelAttrs() []*excel_import.ExcelImportTagAttr {
	return t.whole.GetModelTags()
}

// SetKey set the key of the tree node
// the key is the unique key of the tree node
func (t *TreeNode) SetKey(key string) {
	t.key = key
}

// GetItem get the item of the tree node
// the item is the raw model of the tree node, usually the leaf node has the item
func (t *TreeNode) GetItem() any {
	return t.extra.items[0].item
}

// GetItems get the items of the tree node
// the items is the raw model of the tree node, usually the leaf node has the item
func (t *TreeNode) GetItems() []any {
	items := make([]any, len(t.extra.items))
	for i, item := range t.extra.items {
		items[i] = item.item
	}

	return items
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

func (t *TreeNode) GetRows() []int {
	rows := make([]int, len(t.extra.items))
	for i, item := range t.extra.items {
		rows[i] = item.row
	}

	return rows
}

func (t *TreeNode) CheckIsRoot() bool {
	return t.parent == nil
}

func (t *TreeNode) GetID() int64 {
	return t.id
}

// SetID set the id of the tree node.
// should be called after import the tree node
func (t *TreeNode) SetID(id int64) {
	t.id = id
}

func (t *TreeNode) GetRank() int {
	return t.rank
}

func constructLevelNode(s string, parent *TreeNode, level int) *TreeNode {
	node := &TreeNode{
		value:  s,
		parent: parent,
		rank:   level,
		extra: &TreeNodeExtra{
			items: make([]*TreeNodeItem, 0),
		},
	}
	parent.children = append(parent.children, node)
	return node
}

type TreeImportCfg struct {
	// the tree level order of the tree node
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
	ef excel_import.EndFunc
	// the end condition of the tree column
	treeColEndFunc ColEndFunc
	// the cell format function
	cellFormatFunc excel_import.CellFormatter
	// the row filter function
	rowFilterFunc excel_import.RowFilter
	// enable format checker
	enableFormatChecker bool
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
