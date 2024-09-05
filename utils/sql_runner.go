package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"os"
	"reflect"
	"strings"
)

const (
	defaultValueQuote   = "'"
	alternateValueQuote = "\""
)

type SqlSentencesRunner struct {
	sqlPath   string
	sqlFile   *os.File
	db        *gorm.DB
	tableName string
}

func NewSqlSentencesRunner(sqlPath string, db *gorm.DB, tableName string) *SqlSentencesRunner {
	return &SqlSentencesRunner{
		sqlPath:   sqlPath,
		db:        db,
		tableName: tableName,
	}
}

func (r *SqlSentencesRunner) GenerateSqlInsertSentences(model any) error {
	if r.sqlFile == nil {
		err := r.initSqlFile()
		if err != nil {
			return err
		}
	}

	// generate insert sql sentences
	sql := GenerateInsertSQLWithValues(r.tableName, model)

	// write to file
	_, err := r.sqlFile.WriteString(sql)
	if err != nil {
		return err
	}

	return nil
}

func (r *SqlSentencesRunner) TableName() string {
	return r.tableName
}

func (r *SqlSentencesRunner) WriteSqlSentences(sqls []string) error {
	if r.sqlFile == nil {
		err := r.initSqlFile()
		if err != nil {
			return err
		}
	}

	// write to file
	for _, sql := range sqls {
		_, err := r.sqlFile.WriteString(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SqlSentencesRunner) RunSqlSentencesWithBatch(batchSize int) error {
	// reset file pointer
	_, err := r.sqlFile.Seek(0, 0)
	if err != nil {
		return err
	}

	// read sql line from file
	ioReader := bufio.NewReader(r.sqlFile)
	var sqlSentences []string
	var readLastLine bool
	for {
		line, _, err := ioReader.ReadLine()
		if err == io.EOF {
			readLastLine = true
		} else if err != nil {
			return err
		}

		// collect sql sentences until batch size
		// if batch size is reached, run sql sentences
		sqlSentences = append(sqlSentences, string(line))

		if len(sqlSentences) == batchSize || readLastLine {
			if len(sqlSentences) == 0 {
				break
			}

			// run sql sentences
			sql := strings.Join(sqlSentences, "\n")
			if err = r.db.Exec(sql).Error; err != nil {
				return err
			}
			sqlSentences = nil
		}

		if readLastLine {
			break
		}
	}

	// run sql sentences with batch size
	return nil
}

func (r *SqlSentencesRunner) Close() error {
	if r.sqlFile != nil {
		return r.sqlFile.Close()
	}
	return nil
}

func GenerateInsertSQLWithValues(tableName string, v interface{}) string {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var columns []string
	var values []string

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)
		dbTag := fieldType.Tag.Get("db")

		// 如果字段是零值，跳过它
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}
			field = field.Elem()
		}

		if isZero(field) {
			continue
		}

		columns = append(columns, dbTag)
		values = append(values, formatValue(field))
	}

	columnsStr := strings.Join(columns, ", ")
	valuesStr := strings.Join(values, ", ")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);\n", tableName, columnsStr, valuesStr)
	return query
}

func GenerateInsertSqlWithMap(tableName string, m map[string]any) string {
	var columns []string
	var values []string

	for k, v := range m {
		columns = append(columns, k)
		values = append(values, formatValue(reflect.ValueOf(v)))
	}

	columnsStr := strings.Join(columns, ", ")
	valuesStr := strings.Join(values, ", ")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);\n", tableName, columnsStr, valuesStr)
	return query
}

func (r *SqlSentencesRunner) GenerateSqlUpdateSentences(updates, where map[string]any) error {
	if r.sqlFile == nil {
		err := r.initSqlFile()
		if err != nil {
			return err
		}
	}

	// generate update sql sentences
	sql := GenerateUpdateSQLWithValues(r.tableName, updates, where)

	// write to file
	_, err := r.sqlFile.WriteString(sql)
	if err != nil {
		return err
	}

	return nil
}

func GenerateUpdateSQLWithValues(tableName string, updates, where map[string]any) string {
	var updateStr []string
	var whereStr []string
	for k, v := range updates {
		updateStr = append(updateStr, fmt.Sprintf("%s = %s", k, formatValue(reflect.ValueOf(v))))
	}
	for k, v := range where {
		whereStr = append(whereStr, fmt.Sprintf("%s = %s", k, formatValue(reflect.ValueOf(v))))
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s;\n", tableName, strings.Join(updateStr, ", "), strings.Join(whereStr, " AND "))
	return query
}

// 判断是否为零值
func isZero(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func (r *SqlSentencesRunner) initSqlFile() error {
	file, err := os.Create(r.sqlPath)
	if err != nil {
		return err
	}

	r.sqlFile = file
	return nil
}

func formatValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return formatValueString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	case reflect.Bool:
		if v.Bool() {
			return "TRUE"
		}
		return "FALSE"
	default:
		// 可以根据需要处理更多类型
		return formatValueInterface(v.Interface())
	}
}

func formatValueInterface(v any) string {
	b, _ := json.Marshal(v)
	s := string(b)
	return formatValueString(s)
}

func formatValueString(s string) string {
	s = escapeSQLValue(s)
	return fmt.Sprintf("'%s'", s)
}

// escapeSQLValue escapes a string value for use in a SQL query.
func escapeSQLValue(value string) string {
	// Escape backslashes and single quotes
	value = strings.ReplaceAll(value, `\`, `\\`)
	// Escape single quotes
	value = strings.ReplaceAll(value, `'`, `''`)
	return value
}
