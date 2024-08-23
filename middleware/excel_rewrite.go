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

func (e *ExcelRewriterMiddleware) PostLevelImportHandle(tx *gorm.DB, node *tree_framework.TreeNode) error {
	// get models
	models := node.GetItems()

	// iterate models and write to content
	for i, attr := range e.attrs {
		if !attr.Rewrite || attr.ColumnIndex < 0 {
			continue
		}

		for _, model := range models {
			s, err := util.GetFieldString(model, i)
			if err != nil {
				return err
			}

			e.contents[attr.ColumnIndex] = append(e.contents[attr.ColumnIndex], s)
		}
	}

	return nil
}
