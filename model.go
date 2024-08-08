package excel_import

type RowType string
type ColumnType int

type FieldsOrder int

type optionFunc func(*importFramework)

type rawWhole struct {
	rawContents []*rawContent
}

type rawContent struct {
	sectionType RowType
	content     []string
	model       any
}
