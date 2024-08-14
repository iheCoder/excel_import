package tree_framework

import (
	"gorm.io/gorm"
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
		si.leafs = append(si.leafs, node.GetItem().(*rawModel))
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
