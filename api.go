package excel_import

import "gorm.io/gorm"

type ExcelImporter interface {
	// Import imports the excel file.
	Import(path string) error
}

type RowModelFactory interface {
	// MinColumnCount the min row count to construct raw model
	MinColumnCount() int
	// GetModel get the model
	GetModel() any
}

type PostHandler interface {
	// PostHandle post handle the section.
	PostHandle(tx *gorm.DB) error
}
