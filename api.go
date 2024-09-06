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

// CellFormatter format the cell content.
type CellFormatter func(s string) string

// EndFunc is the end condition of the excel contents
type EndFunc func(s []string) bool

// RowFilter filter the row
type RowFilter func(s []string) bool

// FormatChecker check the type of the cell content
type FormatChecker func(s string) bool

type CorrectnessChecker interface {
	// PreCollect pre collect the data.
	PreCollect(tx *gorm.DB) error
	// CheckCorrect check the correctness of the import.
	CheckCorrect(tx *gorm.DB) error
}
