package general_framework

import (
	"excel_import"
	util "excel_import/utils"
)

var (
	tcfFuncMap = map[excel_import.FormatCheckFunc]excel_import.FormatChecker{
		excel_import.FormatCheckFuncInt:      util.CheckIsInt,
		excel_import.FormatCheckFuncFloat:    util.CheckIsFloat,
		excel_import.FormatCheckFuncChinese:  util.CheckIsContainsChinese,
		excel_import.FormatCheckFuncEnglish:  util.CheckIsContainsEnglish,
		excel_import.FormatCheckFuncPinyin:   util.CheckIsPinyin,
		excel_import.FormatCheckFuncUrl:      util.CheckIsUrl,
		excel_import.FormatCheckFuncImageUrl: util.CheckIsImageUrl,
		excel_import.FormatCheckFuncHash:     util.CheckIsHash,
	}
)

type TagCommonTypeCheck struct {
	tcfFuncMap map[excel_import.FormatCheckFunc]excel_import.FormatChecker
}

func NewTagCommonTypeCheck() *TagCommonTypeCheck {
	return &TagCommonTypeCheck{
		tcfFuncMap: tcfFuncMap,
	}
}

func (t *TagCommonTypeCheck) CheckContents(contents [][]string, tags []*excel_import.ExcelImportTagAttr) bool {
	contents = util.ReverseMatrix(contents)

	for _, tag := range tags {
		if !t.checkContents(contents[tag.ColumnIndex], tag) {
			return false
		}
	}
	return true
}

func (t *TagCommonTypeCheck) checkContents(content []string, tag *excel_import.ExcelImportTagAttr) bool {
	for _, c := range content {
		if !t.checkTypeFunc(tag, c) {
			return false
		}
	}
	return true
}

// checkTypeFunc checks the type of the given string.
func (t *TagCommonTypeCheck) checkTypeFunc(tag *excel_import.ExcelImportTagAttr, str string) bool {
	// if tcf is not set or str is empty, return true
	if len(tag.FCF) == 0 || len(str) == 0 {
		return true
	}

	// use the tcf function to check the type
	if f, ok := t.tcfFuncMap[tag.FCF]; ok {
		return f(str)
	}

	return true
}
