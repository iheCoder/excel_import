package general_framework

import (
	"errors"
	util "excel_import/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

// batchSupportFeature a middleware that supports batch sql execution.
// It will cache sqls until the batch size is reached, then execute them.
// unlike other middlewares, this is an in-framework middleware.
type batchSupportFeature struct {
	BatchSize int
	contents  []string
}

func newBatchSupportFeature(batchSize int) *batchSupportFeature {
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}

	return &batchSupportFeature{
		BatchSize: batchSize,
		contents:  make([]string, 0, batchSize),
	}
}

// AddModel add a model to the batch.
func (b *batchSupportFeature) AddModel(tx *gorm.DB, model any) error {
	tableName, err := getModelTableName(model)
	if err != nil {
		return err
	}

	sql := util.GenerateInsertSQLWithValues(tableName, model)

	return b.addSql(tx, sql)
}

func (b *batchSupportFeature) AddUpdate(tx *gorm.DB, model any, updateCond, whereCond map[string]interface{}) error {
	tableName, err := getModelTableName(model)
	if err != nil {
		return err
	}

	sql := util.GenerateUpdateSQLWithValues(tableName, updateCond, whereCond)
	return b.addSql(tx, sql)
}

func getModelTableName(model any) (string, error) {
	tabler, ok := model.(schema.Tabler)
	if !ok {
		return "", errors.New("model does not implement schema.Tabler")
	}

	return tabler.TableName(), nil
}

// addSql add a sql to the batch.
// if the batch size is reached, execute the batch.
func (b *batchSupportFeature) addSql(tx *gorm.DB, sql string) error {
	b.contents = append(b.contents, sql)
	if len(b.contents) >= b.BatchSize {
		return b.executeBatch(tx)
	}
	return nil
}

// executeBatch execute the batch.
func (b *batchSupportFeature) executeBatch(tx *gorm.DB) error {
	if len(b.contents) == 0 {
		return nil
	}

	sqls := b.contents
	b.contents = b.contents[:0]
	sql := strings.Join(sqls, "\n")

	if err := tx.Exec(sql).Error; err != nil {
		return err
	}

	return nil
}

func (b *batchSupportFeature) PreImportHandle(tx *gorm.DB, whole *RawWhole) error {
	return nil
}

func (b *batchSupportFeature) PostImportSectionHandle(tx *gorm.DB, rc *RawContent) error {
	// generate insert sql if effect model is not nil
	im := rc.GetInsertModel()
	if im != nil {
		if err := b.AddModel(tx, im); err != nil {
			return err
		}
	}

	// generate update sql if effect update is exists
	um, upCond, whereCond := rc.GetUpdateCond()
	if len(upCond) > 0 && len(whereCond) > 0 {
		if err := b.AddUpdate(tx, um, upCond, whereCond); err != nil {
			return err
		}
	}

	return nil
}

func (b *batchSupportFeature) PostHandle(tx *gorm.DB) error {
	return b.executeBatch(tx)
}
