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
	tif := NewTreeImportFramework(nil, cfg, nil)

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

func (si *simpleTestDataImporter) ImportLevelNode(tx *gorm.DB, n *TreeNode) error {
	si.msvs = append(si.msvs, n.GetValue())
	if n.CheckIsLeaf() {
		si.leafs = append(si.leafs, n.GetItem().(*rawModel))
	}

	return nil
}

func TestTreeImportFramework_Import(t *testing.T) {
	path := "../testdata/excel_tree_test_data.xlsx"

	mf := &modelFac{}
	cfg := &TreeImportCfg{
		LevelOrder:   []int{0, 1, 2},
		TreeBoundary: 4,
		ModelFac:     mf,
	}

	ef := func(s []string) bool {
		return len(s) == 0 || len(s[0]) == 0
	}
	efo := WithEndFunc(ef)

	si := &simpleTestDataImporter{}
	levelImporters := []LevelImporter{si, si, si}

	tif := NewTreeImportFramework(nil, cfg, levelImporters, efo)
	err := tif.Import(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(si.leafs) != 5 || len(si.msvs) != 11 {
		t.Fatalf("leafs length is %d, models length is %d", len(si.leafs), len(si.msvs))
	}

	t.Log("done")
}
