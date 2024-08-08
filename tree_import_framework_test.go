package excel_import

import "testing"

func TestConstructTree(t *testing.T) {
	contents := [][]string{
		{"a1", "b2", "c3", "d4", "e5"},
		{"a1", "b2", "c2", "d4", "e5"},
		{"a1", "b2", "c2", "d1", "e5"},
		{"a1", "b2", "c2", "d1", "e7"},
	}
	cfg := &treeImportCfg{
		levelOrder: []int{0, 1, 2, 3, 4},
	}
	tif := NewTreeImportFramework(nil, cfg)

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
