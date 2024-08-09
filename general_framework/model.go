package general_framework

type RowType string
type ColumnType int

type FieldsOrder int

type optionFunc func(*importFramework)
type endFunc func(s []string) bool

type rawWhole struct {
	rawContents []*rawContent
}

type rawContent struct {
	row         int
	sectionType RowType
	content     []string
	model       any
}

type importControl struct {
	// the start row of the content
	startRow int
	// the end condition of the function
	ef endFunc
	// enable type check
	enableTypeCheck bool
	// enable import parallel
	enableParallel bool
	// the max parallel number
	maxParallel int
}

var defaultImportControl = importControl{
	startRow: 1,
}
