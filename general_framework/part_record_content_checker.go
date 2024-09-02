package general_framework

import (
	"excel_import"
	util "excel_import/utils"
	"gorm.io/gorm"
)

// OffsetContentExpectedItem is used to check the offset content in insert.
type OffsetContentExpectedItem struct {
	// Offset is the offset of pre import last id.
	Offset int
	// ExpectedModel is the expected model.
	ExpectedModel any
}

type OffsetContentExpected struct {
	// Items is the offset content expected items.
	Items []*OffsetContentExpectedItem
	// TableModel is the table model.
	TableModel any

	// private fields
	modelAttr []*excel_import.ExcelImportTagAttr
	lastID    int64
}

type PartRecordContentChecker struct {
	oce []*OffsetContentExpected
}

func NewPartRecordContentChecker(oce []*OffsetContentExpected) *PartRecordContentChecker {
	for _, oceItem := range oce {
		if oceItem.TableModel == nil {
			panic("TableModel is nil")
		}

		oceItem.modelAttr = util.ParseTag(oceItem.TableModel)
	}

	return &PartRecordContentChecker{
		oce: oce,
	}
}

func (p *PartRecordContentChecker) PreCollect(tx *gorm.DB) error {
	for _, oceItem := range p.oce {
		// get the last id
		if err := tx.Model(oceItem.TableModel).Order("id desc").Limit(1).Pluck("id", &oceItem.lastID).Error; err != nil {
			return err
		}
	}

	return nil
}
