package correct_checker

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

// LinkedTableWhere is used to get the where condition of the linked table.
type LinkedTableWhere func(m any) (any, string)

// SimpleLinkChecker is used to check the simple link.
// simple link is one table link to multiple tables in same number of columns and same column content.
// for example, table A link to table B and table C on A.resource_id = B.id and A.resource_id = C.id.
type SimpleLinkChecker struct {
	// linkedFunc is the where condition of the linked table.
	linkedFunc LinkedTableWhere
	// TableModel is the table model.
	TableModel any
	// RangeWhere is the range where condition.
	RangeWhere string
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
	s.RangeWhere = fmt.Sprintf("id > %d", lastID)

	return nil
}

func (s *SimpleLinkChecker) CheckCorrect(tx *gorm.DB) error {
	// get the table model in range
	models, err := s.GetTableModels(tx)
	if err != nil {
		return err
	}

	// iterate the models
	for _, model := range models {
		// get the linkedTable condition
		linkedTable, linkedWhere := s.linkedFunc(model)
		if linkedTable == nil {
			continue
		}

		// get the linked models
		if err = tx.Model(linkedTable).Where(linkedWhere).First(linkedTable).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New(fmt.Sprintf("linked table where %v not found", linkedWhere))
			}
			return err
		}
	}

	return nil
}

// GetTableModels get s.TableModel models in range.
func (s *SimpleLinkChecker) GetTableModels(tx *gorm.DB) ([]any, error) {
	// get model type
	modelType := reflect.TypeOf(s.TableModel).Elem()
	sliceType := reflect.SliceOf(modelType)

	// create model slice
	dbModels := reflect.New(sliceType).Elem()

	// query models
	if err := tx.Where(s.RangeWhere).Find(dbModels.Addr().Interface()).Error; err != nil {
		return nil, err
	}

	// convert to []any
	var models []any
	for i := 0; i < dbModels.Len(); i++ {
		models = append(models, dbModels.Index(i).Addr().Interface())
	}

	return models, nil
}
