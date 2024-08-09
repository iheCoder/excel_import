package tree_framework

import "gorm.io/gorm"

type LevelImporter interface {
	// import level tree node
	ImportLevelNode(tx *gorm.DB, parent *TreeNode) error
}

type generateNodeKey func(s []string, level int) string
type rowEndFunc func(s []string) bool
type colEndFunc func(next string) bool
type optionFunc func(*treeImportFramework)
