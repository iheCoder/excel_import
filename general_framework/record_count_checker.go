package general_framework

import (
	"fmt"
	"gorm.io/gorm"
)

type ExpectedCountChange struct {
	// CountDelta is the count delta
	CountDelta int64
}

type RecordCountChecker struct {
	// preCount is the pre count
	preCount int64
	// tableModel is the model of the table
	tableModel any
	// change is the expected count change
	change *ExpectedCountChange
}

// NewRecordCountChecker creates a new record count checker
func NewRecordCountChecker(tableModel any, change *ExpectedCountChange) *RecordCountChecker {
	return &RecordCountChecker{
		tableModel: tableModel,
		change:     change,
	}
}

// PreCollect collects the pre import info
func (c *RecordCountChecker) PreCollect(tx *gorm.DB) error {
	if c.tableModel == nil {
		return nil
	}

	// get the count
	if err := tx.Model(c.tableModel).Count(&c.preCount).Error; err != nil {
		return err
	}

	return nil
}

// CheckCorrect checks the correctness of the import
func (c *RecordCountChecker) CheckCorrect(tx *gorm.DB) error {
	if c.tableModel == nil {
		return nil
	}

	// get the count
	var postCount int64
	if err := tx.Model(c.tableModel).Count(&postCount).Error; err != nil {
		return err
	}

	// check the count
	if c.change.CountDelta != postCount-c.preCount {
		return fmt.Errorf("expected count change %d, but got %d", c.change.CountDelta, postCount-c.preCount)
	}

	return nil
}
