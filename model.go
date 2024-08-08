package excel_import

type RowType string
type ColumnType int

type FieldsOrder int

type optionFunc func(*importFramework)

// used for recognize row section
type sectionRecognizer func(s []string) RowType

type rawWhole struct {
	rawContents []*rawContent
}

type rawContent struct {
	sectionType RowType
	content     []string
	model       any
}
