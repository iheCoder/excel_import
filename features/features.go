package features

import "excel_import"

type FeatureMgr struct {
	// tagFormatChecker is the tag format checker.
	tagFormatChecker *TagFormatChecker
}

func NewFeatureMgr() *FeatureMgr {
	return &FeatureMgr{}
}

func (f *FeatureMgr) EnableTagFormatChecker() {
	f.tagFormatChecker = NewTagCommonFormatCheck()
}

func (f *FeatureMgr) CheckContents(contents []string, tags []*excel_import.ExcelImportTagAttr) error {
	return f.tagFormatChecker.CheckContents(contents, tags)
}
