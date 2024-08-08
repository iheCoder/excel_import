package excel_import

import "gorm.io/gorm"

type LevelImporter interface {
	// import level tree node
	ImportLevelNode(tx *gorm.DB, parent *treeNode) error
}
