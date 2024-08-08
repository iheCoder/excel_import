package excel_import

import (
	"fmt"
	"gorm.io/gorm"
)

type treeImportFramework struct {
	db       *gorm.DB
	recorder *unexpectedRecorder
	cfg      *orderLevelCfg
	nodes    map[string]*treeNode
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
