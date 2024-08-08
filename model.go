package excel_import

type RowType string
type ColumnType int

type FieldsOrder int

type optionFunc func(*importFramework)

// used for recognize row section
type sectionRecognizer func(s []string) RowType

type rawContent struct {
	sectionTypes []RowType
	content      [][]string
}
