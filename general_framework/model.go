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

// SetInsertModel set the inserted model
// used in sqlRunner middleware or batch feature
// model must be implemented schema.Tabler
func (r *RawContent) SetInsertModel(model any) {
	r.effect.insertedModel = model
}

// SetUpdateCond set the update and where condition
// used in sqlRunner middleware
func (r *RawContent) SetUpdateCond(updates, wheres map[string]any) {
	r.effect.updates = updates
	r.effect.wheres = wheres
}

// SetUpdateModelCond set the update model and condition
// used in batch feature
// model must be implemented schema.Tabler
func (r *RawContent) SetUpdateModelCond(model any, updates, wheres map[string]any) {
	r.effect.updateModel = model
	r.effect.updates = updates
	r.effect.wheres = wheres
}

func (r *RawContent) GetInsertModel() any {
	return r.effect.insertedModel
}

func (r *RawContent) GetUpdateCond() (any, map[string]any, map[string]any) {
	return r.effect.updateModel, r.effect.updates, r.effect.wheres
}

// the effect of the import
type importEffect struct {
	// the inserted model
	insertedModel any

	// the updated model
	updateModel any
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
	// enable import with batch
	// SetInsertModel or SetUpdateModelCond should be called in the ImportSection
	EnableBatch bool
	// the batch size
	BatchSize int
}

var defaultImportControl = ImportControl{
	StartRow:  1,
	Ef:        defaultRawEndFunc,
	BatchSize: defaultBatchSize,
}

func defaultRawEndFunc(s []string) bool {
	return len(s) == 0 || len(s[0]) == 0
}
