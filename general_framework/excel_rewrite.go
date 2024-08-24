package general_framework

import (
	"excel_import"
	util "excel_import/utils"
	"gorm.io/gorm"
)

type ExcelRewriterMiddleware struct {
	path     string
	contents map[int][]string
	attrs    []*excel_import.ExcelImportTagAttr
	startRow int
}

func NewExcelRewriterMiddleware(path string) *ExcelRewriterMiddleware {
	return &ExcelRewriterMiddleware{
		path:     path,
		contents: make(map[int][]string),
		startRow: 1,
	}
}

func (e *ExcelRewriterMiddleware) SetStartRow(startRow int) {
	e.startRow = startRow
}

func (e *ExcelRewriterMiddleware) PreImportHandle(tx *gorm.DB, content *RawContent) error {
	if content.Model == nil {
		return nil
	}

	// get model excel import tag attr
	attrs := util.ParseTag(content.Model)

	// set attrs
	e.attrs = attrs

	return nil
}

func (e *ExcelRewriterMiddleware) PostImportSectionHandle(tx *gorm.DB, s *RawContent) error {
	// get models
	model := s.Model

	// iterate models and write to content
	for i, attr := range e.attrs {
		if !attr.Rewrite || attr.ColumnIndex < 0 {
			continue
		}

		// write to content
		c, err := util.GetFieldString(model, attr.ColumnIndex)
		if err != nil {
			return err
		}
		e.contents[i] = append(e.contents[i], c)
	}

	return nil
}

func (e *ExcelRewriterMiddleware) PostImportHandle(tx *gorm.DB) error {
	// write to excel
	if err := util.WriteExcelColumnContentByStartRow(e.path, e.contents, e.startRow); err != nil {
		return err
	}

	return nil
}
