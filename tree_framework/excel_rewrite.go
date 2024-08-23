package tree_framework

import (
	"excel_import"
	util "excel_import/utils"
	"gorm.io/gorm"
)

type ExcelRewriterTreeMiddleware struct {
	path     string
	contents map[int][]string
	attrs    []*excel_import.ExcelImportTagAttr
	startRow int
}

func NewExcelRewriterTreeMiddleware(path string) *ExcelRewriterTreeMiddleware {
	return &ExcelRewriterTreeMiddleware{
		path:     path,
		contents: make(map[int][]string),
		startRow: 1,
	}
}

func (e *ExcelRewriterTreeMiddleware) SetStartRow(startRow int) {
	e.startRow = startRow
}

// PreImportHandle init excel import tag attr
func (e *ExcelRewriterTreeMiddleware) PreImportHandle(tx *gorm.DB, info TreeInfo) error {
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

func (e *ExcelRewriterTreeMiddleware) PostLevelImportHandle(tx *gorm.DB, node *TreeNode) error {
	// check if node is leaf
	if !node.CheckIsLeaf() {
		return nil
	}

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

func (e *ExcelRewriterTreeMiddleware) PostHandle(tx *gorm.DB) error {
	return util.WriteExcelColumnContentByStartRow(e.path, e.contents, e.startRow)
}
