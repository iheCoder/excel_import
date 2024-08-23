package tree_framework

import (
	util "excel_import/utils"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

func TestConstructTree(t *testing.T) {
	contents := [][]string{
		{"a1", "b2", "c3", "d4", "e5"},
		{"a1", "b2", "c2", "d4", "e5"},
		{"a1", "b2", "c2", "d1", "e5"},
		{"a1", "b2", "c2", "d1", "e7"},
	}
	cfg := &TreeImportCfg{
		LevelOrder: []int{0, 1, 2, 3, 4},
	}
	tif := NewTreeImportFramework(nil, cfg, nil, nil)

	root, err := tif.constructTree(contents)
	if err != nil {
		t.Fatal(err)
	}

	if root == nil {
		t.Fatal("root is nil")
	}
	if root.children == nil {
		t.Fatal("root.children is nil")
	}
	if len(root.children) != 1 {
		t.Fatalf("root.children length is %d", len(root.children))
	}
	if root.children[0].children == nil {
		t.Fatal("root.children[0].children is nil")
	}
	if len(root.children[0].children) != 1 {
		t.Fatalf("root.children[0].children length is %d", len(root.children[0].children))
	}
	if root.children[0].children[0].children == nil {
		t.Fatal("root.children[0].children[0].children is nil")
	}
	if len(root.children[0].children[0].children) != 2 {
		t.Fatalf("root.children[0].children[0].children length is %d", len(root.children[0].children[0].children))
	}

	t.Log("done")
}

type modelFac struct {
}

func (mf *modelFac) GetModel() any {
	return &rawModel{}
}

func (mf *modelFac) MinColumnCount() int {
	return 4
}

type rawModel struct {
	L1, L2, L3 string
	Key        string
}

type simpleTestDataImporter struct {
	leafs []*rawModel
	msvs  []string
}

func (si *simpleTestDataImporter) ImportLevelNode(tx *gorm.DB, node *TreeNode) error {
	si.msvs = append(si.msvs, node.GetValue())
	if node.CheckIsLeaf() {
		si.leafs = append(si.leafs, &rawModel{})
	}

	return nil
}

func TestTreeImportFramework_Import(t *testing.T) {
	path := "../testdata/excel_tree_test_data.xlsx"

	mf := &modelFac{}
	cfg := &TreeImportCfg{
		LevelOrder:   []int{0, 1, 2},
		TreeBoundary: 2,
		ModelFac:     mf,
		ColumnCount:  4,
	}

	si := &simpleTestDataImporter{}
	levelImporters := []LevelImporter{si, si, si}

	tif := NewTreeImportFramework(nil, cfg, nil, levelImporters)
	err := tif.Import(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(si.leafs) != 5 || len(si.msvs) != 11 {
		t.Fatalf("leafs length is %d, models length is %d", len(si.leafs), len(si.msvs))
	}

	t.Log("done")
}

func TestTreeImportFramework_ImportStrictOrder(t *testing.T) {
	path := "../testdata/excel_tree_test_data.xlsx"
	mf := &modelFac{}
	si := &simpleTestDataImporter{}

	tif := NewTreeImportStrictOrderFramework(nil, 2, 4, mf, si)
	err := tif.Import(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(si.leafs) != 5 || len(si.msvs) != 12 {
		t.Fatalf("leafs length is %d, models length is %d", len(si.leafs), len(si.msvs))
	}

	t.Log("done")
}

func TestTreeImportFramework_ImportWithTag(t *testing.T) {
	path := "../testdata/excel_tree_test_tag.xlsx"
	mf := &modelTagFac{}
	si := &simpleTestDataTagImporter{}
	cfg := &TreeImportCfg{
		LevelOrder:   []int{1, 3, 4},
		TreeBoundary: 4,
		ModelFac:     mf,
		ColumnCount:  6,
	}
	levelImporters := []LevelImporter{si, si, si}

	tif := NewTreeImportFramework(nil, cfg, nil, levelImporters)
	err := tif.Import(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(si.leafs) != 5 || len(si.msvs) != 11 {
		t.Fatalf("leafs length is %d, models length is %d", len(si.leafs), len(si.msvs))
	}

	t.Log("done")
}

type modelTagFac struct {
}

func (mf *modelTagFac) GetModel() any {
	return &rawTagModel{}
}

func (mf *modelTagFac) MinColumnCount() int {
	return 6
}

type rawTagModel struct {
	L1  string `exi:"index:1"`
	L2  string `exi:"index:3"`
	L3  string `exi:"index:4"`
	Key string `exi:"index:5"`
}

type simpleTestDataTagImporter struct {
	leafs []*rawTagModel
	msvs  []string
}

func (si *simpleTestDataTagImporter) ImportLevelNode(tx *gorm.DB, node *TreeNode) error {
	si.msvs = append(si.msvs, node.GetValue())
	if node.CheckIsLeaf() {
		si.leafs = append(si.leafs, &rawTagModel{})
	}

	return nil
}

func TestTreeImportFramework_ImportWithExcelRewriteMiddleware(t *testing.T) {
	path := "../testdata/excel_tree_test_rewrite.xlsx"
	mf := &modelRewriteFac{}
	si := &simpleTestDataRewriteImporter{}
	cfg := &TreeImportCfg{
		LevelOrder:   []int{1, 3, 4},
		TreeBoundary: 4,
		ModelFac:     mf,
		ColumnCount:  7,
	}
	excelRewriteMiddleware := NewExcelRewriterTreeMiddleware(path)

	tif := NewTreeImportFramework(nil, cfg, nil, []LevelImporter{si, si, si}, WithMiddlewares(excelRewriteMiddleware))
	err := tif.Import(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(si.leafs) != 5 || len(si.msvs) != 11 {
		t.Fatalf("leafs length is %d, models length is %d", len(si.leafs), len(si.msvs))
	}

	// read leaf id from excel and check it
	leafIDColumnIndex := 6
	expectedLeafIDs := []int{0, 1, 2, 3, 4, 4}
	contents, err := util.ReadExcelContent(path)
	if err != nil {
		t.Fatal(err)
	}

	for i, d := range expectedLeafIDs {
		if contents[i+1][leafIDColumnIndex] != strconv.Itoa(d) {
			t.Fatalf("leaf id is %s, expect %d", contents[i+1][leafIDColumnIndex], d)
		}
	}

	t.Log("done")
}

type rawRewriteModel struct {
	L1     string `exi:"index:1"`
	L2     string `exi:"index:3"`
	L3     string `exi:"index:4"`
	Key    string `exi:"index:5"`
	LeafID int    `exi:"index:6,rewrite:true"`
}

type modelRewriteFac struct {
}

func (mf *modelRewriteFac) GetModel() any {
	return &rawRewriteModel{}
}

func (mf *modelRewriteFac) MinColumnCount() int {
	return 7
}

type simpleTestDataRewriteImporter struct {
	leafs  []*rawRewriteModel
	msvs   []string
	leafID int
}

func (si *simpleTestDataRewriteImporter) ImportLevelNode(tx *gorm.DB, node *TreeNode) error {
	si.msvs = append(si.msvs, node.GetValue())
	if node.CheckIsLeaf() {
		si.leafs = append(si.leafs, &rawRewriteModel{})
		models := node.GetItems()
		for _, model := range models {
			m := model.(*rawRewriteModel)
			m.LeafID = si.leafID
		}
		si.leafID++
	}

	return nil
}
