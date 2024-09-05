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

type GormTag struct {
	// The column name of the database.
	// tagName: column
	Column string

	// The primary key of the database.
	// tagName: primary_key
	PrimaryKey bool

	// The auto increment of the database.
	// tagName: auto_increment
	AutoIncrement bool

	// The default value of the database.
	// tagName: default
	Default string

	// The not null of the database.
	// tagName: not null
	NotNull bool

	// The size of the database.
	// tagName: size
	Size int

	// The type of the database.
	// tagName: type
	Type string
}
