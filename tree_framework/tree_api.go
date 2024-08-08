package tree_framework

import "gorm.io/gorm"

type LevelImporter interface {
	// import level tree node
	ImportLevelNode(tx *gorm.DB, parent *treeNode) error
}

type generateNodeKey func(s []string, level int) string
