package general_framework

import (
	util "excel_import/utils"
	"gorm.io/gorm"
)

const (
	defaultCacheSize = 1000
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
