package features

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

type TagFormatChecker struct {
	tcfFuncMap map[excel_import.FormatCheckFunc]excel_import.FormatChecker
}

func NewTagCommonFormatCheck() *TagFormatChecker {
	return &TagFormatChecker{
		tcfFuncMap: tcfFuncMap,
	}
}

func (t *TagFormatChecker) CheckContents(content []string, tags []*excel_import.ExcelImportTagAttr) error {
	errBuilder := util.NewErrBuilder()
	for i, c := range content {
		if err := t.checkFormatFunc(tags[i], c); err != nil {
			errBuilder.AddWithContent(c, err)
		}
	}
	return errBuilder.Build()
}

// checkFormatFunc checks the type of the given string.
func (t *TagFormatChecker) checkFormatFunc(tag *excel_import.ExcelImportTagAttr, str string) error {
	// if tcf is not set or str is empty, return true
	if len(tag.FCF) == 0 || len(str) == 0 {
		return nil
	}

	// use the tcf function to check the type
	if f, ok := t.tcfFuncMap[tag.FCF]; ok {
		return f(str)
	}

	return nil
}
