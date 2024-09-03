package correct_checker

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
	// ChkKey is the check key.
	ChkKey string

	// private fields
	modelAttr []*excel_import.ExcelImportTagAttr
	lastID    int64
}

// IDContentExpectedItem is used to check the id content in update.
type IDContentExpectedItem struct {
	// ID is the id of the model.
	ID int64
	// ExpectedModel is the expected model.
	ExpectedModel any
}

type IDContentExpected struct {
	// Items is the id content expected items.
	Items []*IDContentExpectedItem
	// TableModel is the table model.
	TableModel any
	// ChkKey is the check key.
	ChkKey string

	// private fields
	modelAttr []*excel_import.ExcelImportTagAttr
}

type PartRecordContentChecker struct {
	oce []*OffsetContentExpected
	ice []*IDContentExpected
}

func NewPartRecordContentChecker() *PartRecordContentChecker {
	return &PartRecordContentChecker{}
}

func (p *PartRecordContentChecker) SetIDContentExpected(ice []*IDContentExpected) {
	for _, iceItem := range ice {
		if iceItem.TableModel == nil {
			panic("TableModel is nil")
		}

		iceItem.modelAttr = util.ParseTag(iceItem.TableModel)
	}

	p.ice = ice
}

func (p *PartRecordContentChecker) SetOffsetContentExpected(oce []*OffsetContentExpected) {
	for _, oceItem := range oce {
		if oceItem.TableModel == nil {
			panic("TableModel is nil")
		}

		oceItem.modelAttr = util.ParseTag(oceItem.TableModel)
	}

	p.oce = oce
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
			if err = util.CompareModel(models[i], item.ExpectedModel, oceItem.modelAttr, oceItem.ChkKey); err != nil {
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
