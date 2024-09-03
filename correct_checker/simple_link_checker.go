package correct_checker

import (
	"errors"
	"gorm.io/gorm"
)

// LinkedTableWhere is used to get the where condition of the linked table.
type LinkedTableWhere func(m any) (any, map[string]any)

// SimpleLinkChecker is used to check the simple link.
// simple link is one table link to multiple tables in same number of columns and same column content.
// for example, table A link to table B and table C on A.resource_id = B.id and A.resource_id = C.id.
type SimpleLinkChecker struct {
	// linkedFunc is the where condition of the linked table.
	linkedFunc LinkedTableWhere
	// TableModel is the table model.
	TableModel any
	// RangeWhere is the range where condition.
	RangeWhere map[string]any
}

func NewSimpleLinkChecker(linkedFunc LinkedTableWhere, tableModel any) *SimpleLinkChecker {
	return &SimpleLinkChecker{
		linkedFunc: linkedFunc,
		TableModel: tableModel,
	}
}

func (s *SimpleLinkChecker) PreCollect(tx *gorm.DB) error {
	// get last id
	var lastID int64
	if err := tx.Model(s.TableModel).Select("id").Order("id desc").First(&lastID).Error; err != nil {
		return err
	}

	// construct range where
	s.RangeWhere = map[string]any{
		"id": gorm.Expr("> ?", lastID),
	}

	return nil
}

func (s *SimpleLinkChecker) CheckCorrect(tx *gorm.DB) error {
	// get the linked table where condition
	linkedWhere, linkedMap := s.linkedFunc(s.TableModel)

	// get the count
	var count int64
	if err := tx.Model(s.TableModel).Where(s.RangeWhere).Where(linkedWhere, linkedMap).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return errors.New("the linked table is not correct")
}
