package middleware

import (
	"excel_import"
	"excel_import/tree_framework"
	util "excel_import/utils"
	"gorm.io/gorm"
)

type ExcelRewriterMiddleware struct {
	path     string
	contents map[int][]string
	attrs    []*excel_import.ExcelImportTagAttr
}

func NewExcelRewriterPostHandler(path string) *ExcelRewriterMiddleware {
	return &ExcelRewriterMiddleware{
		path:     path,
		contents: make(map[int][]string),
	}
}

// PreImportHandle init excel import tag attr
func (e *ExcelRewriterMiddleware) PreImportHandle(tx *gorm.DB, info tree_framework.TreeInfo) error {
	if len(info.GetModels()) == 0 {
		return nil
	}

	// get model excel import tag attr
	model := info.GetModels()[0]
	if model == nil {
		return nil
	}
	attrs := util.ParseTag(model)

	// set attrs
	e.attrs = attrs

	return nil
}
