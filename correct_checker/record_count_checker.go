package correct_checker

import (
	"fmt"
	"gorm.io/gorm"
)

type ExpectedCountChange struct {
	// TablesCount is the count info of tables
	TablesCount []TableCountInfo
}

type TableCountInfo struct {
	// CountDelta is the count delta
	CountDelta int64
	// TableModel is the model of the table
	TableModel any
	// preCount is the count before import
	preCount int64
	// RangeWhere is the range where condition
	RangeWhere string
}

type RecordCountChecker struct {
	// change is the expected count change
	change *ExpectedCountChange
}

// NewRecordCountChecker creates a new record count checker
func NewRecordCountChecker(change *ExpectedCountChange) *RecordCountChecker {
	return &RecordCountChecker{
		change: change,
	}
}

// PreCollect collects the pre import info
func (c *RecordCountChecker) PreCollect(tx *gorm.DB) error {
	if c.change == nil {
		return nil
	}

	// get the count
	var err error
	for i := range c.change.TablesCount {
		tableCount := &c.change.TablesCount[i]
		if tableCount.TableModel == nil {
			continue
		}

		// get the count
		tableCount.preCount, err = c.getCount(tx, tableCount)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RecordCountChecker) getCount(tx *gorm.DB, tableCount *TableCountInfo) (int64, error) {
	var count int64
	if len(tableCount.RangeWhere) > 0 {
		if err := tx.Model(tableCount.TableModel).Where(tableCount.RangeWhere).Count(&count).Error; err != nil {
			return 0, err
		}
	} else {
		if err := tx.Model(tableCount.TableModel).Count(&count).Error; err != nil {
			return 0, err
		}
	}

	return count, nil
}

// CheckCorrect checks the correctness of the import
func (c *RecordCountChecker) CheckCorrect(tx *gorm.DB) error {
	if c.change == nil {
		return nil
	}

	// get the count
	for i := range c.change.TablesCount {
		tableCount := &c.change.TablesCount[i]
		if tableCount.TableModel == nil {
			continue
		}

		// get the count
		count, err := c.getCount(tx, tableCount)
		if err != nil {
			return err
		}

		// check the count
		if count != tableCount.preCount+tableCount.CountDelta {
			return fmt.Errorf("table %v count is not correct, expect %v, but got %v", tableCount.TableModel, tableCount.preCount+tableCount.CountDelta, count)
		}
	}

	return nil
}
