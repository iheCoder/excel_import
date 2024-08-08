package excel_import

type ExcelImporter interface {
	// Import imports the excel file.
	Import(path string) error
}
