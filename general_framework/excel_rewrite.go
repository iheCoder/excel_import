package general_framework

import "excel_import"

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
