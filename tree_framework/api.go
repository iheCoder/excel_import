package tree_framework

import "gorm.io/gorm"

type LevelImporter interface {
	// import level tree node
	ImportLevelNode(tx *gorm.DB, parent *TreeNode) error
}

type GenerateNodeKey func(s []string, level int) string
type RowEndFunc func(s []string) bool
type ColEndFunc func(next string) bool
type OptionFunc func(*TreeImportFramework)
