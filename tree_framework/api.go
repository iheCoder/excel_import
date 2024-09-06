package tree_framework

import (
	"excel_import"
	"gorm.io/gorm"
)

type LevelImporter interface {
	// ImportLevelNode import level tree node
	ImportLevelNode(tx *gorm.DB, node *TreeNode) error
}

type LevelImportPostHandler interface {
	// PostLevelImportHandle post handle the level
	PostLevelImportHandle(tx *gorm.DB, node *TreeNode) error
}

type TreePreHandler interface {
	// PreImportHandle pre handle before import
	PreImportHandle(tx *gorm.DB, info TreeInfo) error
}

type TreeInfo interface {
	// GetRoot return the root node of the tree
	GetRoot() *TreeNode
	// GetNodeCount return the node count of the tree
	GetNodeCount() int
	// GetLeafCount return the leaf count of the tree
	GetLeafCount() int
	// GetModels return the models of the tree
	GetModels() []any
}

type TreeMiddleware interface {
	TreePreHandler
	LevelImportPostHandler
	excel_import.PostHandler
}

type GenerateNodeKey func(s []string, level int) string
type ColEndFunc func(next string) bool
type OptionFunc func(*TreeImportFramework)
