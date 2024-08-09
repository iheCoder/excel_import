package excel_import

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
