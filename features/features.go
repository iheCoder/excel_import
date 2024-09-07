package features

import "excel_import"

type formatCheckStatus int32

const (
	// formatCheckTagInit is the format check tag init.
	formatCheckTagInit formatCheckStatus = 0
	// formatCheckTagExists is the format check tag exists.
	formatCheckTagExists formatCheckStatus = 1
	// formatCheckTagNotExists is the format check tag not exists.
	formatCheckTagNotExists formatCheckStatus = 2
)

type FeatureMgr struct {
	// formatCheckStatus is the format check tag status.
	formatCheckStatus formatCheckStatus
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
	if f.formatCheckStatus == formatCheckTagInit {
		f.formatCheckStatus = formatCheckTagNotExists
		for _, tag := range tags {
			if len(tag.FCF) > 0 {
				f.formatCheckStatus = formatCheckTagExists
				break
			}
		}
	}

	if f.formatCheckStatus == formatCheckTagNotExists || f.tagFormatChecker == nil {
		return nil
	}

	return f.tagFormatChecker.CheckContents(contents, tags)
}
