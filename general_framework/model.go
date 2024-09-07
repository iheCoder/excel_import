package general_framework

import (
	"excel_import"
	util "excel_import/utils"
)

type RowType string
type ColumnType int

type FieldsOrder int

type OptionFunc func(*ImportFramework)

type RawWhole struct {
	rawContents []*RawContent

	modelInfo *ModelsInfo
}

type ModelsInfo struct {
	// the model tags of the rawContent Model
	excelModelTags []*excel_import.ExcelImportTagAttr
}

func (r *RawWhole) GetModelTags() []*excel_import.ExcelImportTagAttr {
	return r.modelInfo.excelModelTags
}

type RawContent struct {
	// the row number. keep the original row number
	Row int
	// the type of the row
	SectionType RowType
	// the content of the row
	Content []string
	// the model of the row
	Model any
	// the import effect.
	// for example, the inserted model or updated model
	effect importEffect
	// the whole operator
	whole *RawWhole
}

func (r *RawContent) GetModelTags() []*excel_import.ExcelImportTagAttr {
	return r.whole.GetModelTags()
}

func (r *RawContent) GetContent() []string {
	return r.Content
}

func (r *RawContent) GetModel() any {
	return r.Model
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
	Ef excel_import.EndFunc
	// enable tag format check
	// must be set model factory
	EnableTagFormatCheck bool
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
	// the row filter function
	RowFilter excel_import.RowFilter
}

var defaultImportControl = ImportControl{
	StartRow:  1,
	Ef:        util.DefaultRowEndFunc,
	BatchSize: defaultBatchSize,
}
