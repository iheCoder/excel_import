package excel_import

type CheckMode string

const (
	CheckModeOn = "on"
)

type ExcelImportTagAttr struct {
	// The column index of the excel file.
	// -1 means not set.
	// tagName: index
	ColumnIndex int
	// weather the column is rewrite.
	// tagName: rewrite
	Rewrite bool
	// check model
	// tagName: chk
	Check CheckMode
}
