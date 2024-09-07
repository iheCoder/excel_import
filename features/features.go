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

// CheckContents checks the contents.
func (f *FeatureMgr) CheckContents(contents []string, tags []*excel_import.ExcelImportTagAttr) error {
	// if format check status is init, conclude the format check status
	if f.formatCheckStatus == formatCheckTagInit {
		f.formatCheckStatus = formatCheckTagNotExists
		for _, tag := range tags {
			if len(tag.FCF) > 0 {
				f.formatCheckStatus = formatCheckTagExists
				break
			}
		}
	}

	// if format check tag is not exists or tag format checker is nil, return nil
	if f.formatCheckStatus == formatCheckTagNotExists || f.tagFormatChecker == nil {
		return nil
	}

	// check contents
	return f.tagFormatChecker.CheckContents(contents, tags)
}

// RegisterFormatChecker registers the format checker.
func (f *FeatureMgr) RegisterFormatChecker(fcf excel_import.FormatCheckFunc, fc excel_import.FormatChecker) {
	if f.tagFormatChecker == nil {
		return
	}
	f.tagFormatChecker.RegisterFormatChecker(fcf, fc)
}
