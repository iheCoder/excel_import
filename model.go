package excel_import

type CheckMode string
type TreeFlag string

const (
	CheckModeOn = "on"

	TreeFlagParentID TreeFlag = "parent_id"
	TreeFlagKey      TreeFlag = "key"
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
	// tree flag
	// tagName: tree
	Tree TreeFlag
}

func CheckChkKeyMatch(cm CheckMode, key string) bool {
	return cm == CheckMode(key)
}
