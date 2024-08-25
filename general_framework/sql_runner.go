package general_framework

import (
	util "excel_import/utils"
	"gorm.io/gorm"
)

const (
	defaultCacheSize = 1000
	defaultBatchSize = 1000
)

type SqlRunnerMiddleware struct {
	runner              *util.SqlSentencesRunner
	cache               []string
	cacheSize           int
	enableExecuteDirect bool
}

func NewSqlRunnerMiddleware(sqlPath string, db *gorm.DB, tableName string, enableExecuteDirect bool) *SqlRunnerMiddleware {
	if enableExecuteDirect && db == nil {
		panic("db is nil")
	}

	return &SqlRunnerMiddleware{
		runner:              util.NewSqlSentencesRunner(sqlPath, db, tableName),
		cacheSize:           defaultCacheSize,
		enableExecuteDirect: enableExecuteDirect,
		cache:               make([]string, 0, defaultCacheSize),
	}
}

func (s *SqlRunnerMiddleware) PreImportHandle(tx *gorm.DB, whole *RawWhole) error {
	// do nothing
	return nil
}

func (s *SqlRunnerMiddleware) PostImportSectionHandle(tx *gorm.DB, rc *RawContent) error {
	// generate insert sql if effect model is not nil
	im := rc.GetInsertModel()
	if im != nil {
		s.cache = append(s.cache, util.GenerateInsertSQLWithValues(s.runner.TableName(), im))
	}

	// generate update sql if effect update is exists
	upCond, whereCond := rc.GetUpdateCond()
	if len(upCond) > 0 && len(whereCond) > 0 {
		s.cache = append(s.cache, util.GenerateUpdateSQLWithValues(s.runner.TableName(), upCond, whereCond))
	}

	// write sql if cache is full
	if len(s.cache) >= s.cacheSize {
		if err := s.runner.WriteSqlSentences(s.cache); err != nil {
			return err
		}
		s.cache = s.cache[:0]
	}

	return nil
}

func (s *SqlRunnerMiddleware) PostImportHandle(tx *gorm.DB, whole *RawWhole) error {
	// write the rest of sql
	if len(s.cache) > 0 {
		if err := s.runner.WriteSqlSentences(s.cache); err != nil {
			return err
		}
	}

	// execute sql
	if s.enableExecuteDirect {
		return s.runner.RunSqlSentencesWithBatch(defaultBatchSize)
	}

	return nil
}
