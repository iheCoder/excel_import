package general_framework

import (
	"errors"
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

func (p *PartRecordContentChecker) CheckCorrect(tx *gorm.DB) error {
	for _, oceItem := range p.oce {
		// create new models
		n := len(oceItem.Items)
		models := createNewModels(oceItem.TableModel, n)

		// get max offset
		var maxOffset int
		for _, item := range oceItem.Items {
			if item.Offset > maxOffset {
				maxOffset = item.Offset
			}
		}

		// get maxOffset ids start from lastID
		var ids []int64
		err := tx.Model(oceItem.TableModel).Where("id > ?", oceItem.lastID).Order("id asc").Limit(maxOffset).Pluck("id", &ids).Error
		if err != nil {
			return err
		}

		if len(ids) != maxOffset {
			return errors.New("ids count is not equal to maxOffset")
		}

		// query models
		for i := 0; i < n; i++ {
			model := models[i]
			id := ids[oceItem.Items[i].Offset-1]
			if err = tx.First(model, id).Error; err != nil {
				return err
			}
		}

		// check models
		for i, item := range oceItem.Items {
			if err = util.CompareModel(models[i], item.ExpectedModel, oceItem.modelAttr); err != nil {
				return err
			}
		}
	}

	return nil
}

func createNewModels(model any, count int) []any {
	models := make([]any, count)
	for i := 0; i < count; i++ {
		models[i] = util.NewModel(model)
	}

	return models
}
