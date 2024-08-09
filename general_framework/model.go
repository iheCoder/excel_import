package general_framework

type RowType string
type ColumnType int

type FieldsOrder int

type OptionFunc func(*ImportFramework)
type EndFunc func(s []string) bool

type rawWhole struct {
	rawContents []*RawContent
}

type RawContent struct {
	Row         int
	SectionType RowType
	Content     []string
	Model       any
}

type ImportControl struct {
	// the start row of the content
	StartRow int
	// the end condition of the function
	Ef EndFunc
	// enable type check
	EnableTypeCheck bool
	// enable import parallel
	EnableParallel bool
	// the max parallel number
	MaxParallel int
}

var defaultImportControl = ImportControl{
	StartRow: 1,
}
