package general_framework

import "excel_import"

type RowType string
type ColumnType int

type FieldsOrder int

type OptionFunc func(*ImportFramework)
type EndFunc func(s []string) bool

type RawWhole struct {
	rawContents []*RawContent
}

type RawContent struct {
	Row         int
	SectionType RowType
	Content     []string
	Model       any
	effect      importEffect
}

func (r *RawContent) SetInsertModel(model any) {
	r.effect.insertedModel = model
}

func (r *RawContent) SetUpdateCond(updates, wheres map[string]any) {
	r.effect.updates = updates
	r.effect.wheres = wheres
}

// the effect of the import
type importEffect struct {
	// the inserted model
	insertedModel any
	// the updated and where condition
	updates, wheres map[string]any
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
	// the cell format function
	CellFormatFunc excel_import.CellFormatter
}

var defaultImportControl = ImportControl{
	StartRow: 1,
	Ef:       defaultRawEndFunc,
}

func defaultRawEndFunc(s []string) bool {
	return len(s) == 0 || len(s[0]) == 0
}
