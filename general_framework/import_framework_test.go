package general_framework

import (
	"bufio"
	util "excel_import/utils"
	"gorm.io/gorm"
	"os"
	"strconv"
	"testing"
)

const (
	personSectionType RowType = "person"
)

func TestImportFramework_Import(t *testing.T) {
	stdi := &simpleTestDataImporter{}
	importers := map[RowType]SectionImporter{
		personSectionType: stdi,
	}

	framework := NewImporterFramework(nil, importers, psr, WithRowRawModel(stdi))
	path := "../testdata/excel_test_data.xlsx"

	err := framework.Import(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range stdi.persons {
		t.Log(p)
	}
	t.Log("done")
}

func TestImportFramework_ImportOneSection(t *testing.T) {
	stdi := &simpleTestDataImporter{}
	framework := NewImporterOneSectionFramework(nil, stdi, WithRowRawModel(stdi))
	path := "../testdata/excel_test_data.xlsx"

	err := framework.Import(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range stdi.persons {
		t.Log(p)
	}
	t.Log("done")
}

type simpleTestDataImporter struct {
	persons []*Person
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (di *simpleTestDataImporter) ImportSection(tx *gorm.DB, s *RawContent) error {
	di.persons = append(di.persons, s.Model.(*Person))

	return nil
}

func (di *simpleTestDataImporter) MinColumnCount() int {
	return 2
}

func (di *simpleTestDataImporter) GetModel() any {
	return &Person{}
}

func psr(s []string) RowType {
	return personSectionType
}

func TestImportFramework_ImportOneSectionWithTag(t *testing.T) {
	stdi := &simpleTestDataWithTagImporter{}
	fac := &personWithTagFac{}
	framework := NewImporterOneSectionFramework(nil, stdi, WithRowRawModel(fac))
	path := "../testdata/excel_test_tag.xlsx"

	err := framework.Import(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range stdi.persons {
		t.Log(p)
	}
	t.Log("done")
}

type simpleTestDataWithTagImporter struct {
	persons []*PersonWithTag
}

func (di *simpleTestDataWithTagImporter) ImportSection(tx *gorm.DB, s *RawContent) error {
	di.persons = append(di.persons, s.Model.(*PersonWithTag))
	return nil
}

type personWithTagFac struct {
}

func (mf *personWithTagFac) GetModel() any {
	return &PersonWithTag{}
}

func (mf *personWithTagFac) MinColumnCount() int {
	return 6
}

type PersonWithTag struct {
	Name   string `exi:"index:0"`
	Career string `exi:"index:2"`
	Degree string `exi:"index:4"`
}

func TestImportFramework_ImportOneSectionWithRewrite(t *testing.T) {
	path := "../testdata/excel_test_rewrite.xlsx"
	stdi := &simpleTestDataSupportMiddlewareImporter{}
	fac := &calculateExampleFac{}
	rewriteMiddleware := NewExcelRewriterMiddleware(path)

	framework := NewImporterOneSectionFramework(nil, stdi, WithRowRawModel(fac), WithMiddlewares(rewriteMiddleware))

	err := framework.Import(path)
	if err != nil {
		t.Fatal(err)
	}

	// read the rewritten content and check
	expectedSum := []int{2, 5, 12, 5}
	rewriteColIndex := 2
	content, err := util.ReadExcelContent(path)
	if err != nil {
		t.Fatal(err)
	}

	for i, s := range expectedSum {
		if content[i+1][rewriteColIndex] != strconv.Itoa(s) {
			t.Fatalf("sum is %s, expected %d", content[i][rewriteColIndex], s)
		}
	}
}

type calculateExample struct {
	X   int `exi:"index:0"`
	Y   int `exi:"index:1"`
	Sum int `exi:"index:2,rewrite:true"`
}

type simpleTestDataSupportMiddlewareImporter struct {
}

func (di *simpleTestDataSupportMiddlewareImporter) ImportSection(tx *gorm.DB, s *RawContent) error {
	model := s.Model.(*calculateExample)
	model.Sum = model.X + model.Y
	return nil
}

type calculateExampleFac struct {
}

func (mf *calculateExampleFac) GetModel() any {
	return &calculateExample{}
}

func (mf *calculateExampleFac) MinColumnCount() int {
	return 3
}

func TestImportFramework_ImportOneSectionWithSqlRunner(t *testing.T) {
	path := "../testdata/excel_test_sql.xlsx"
	stdi := &simpleTestDataSupportSqlRunnerMiddlewareImporter{}
	fac := &computerFac{}
	sqlPath := "../testdata/sql_middleware_test.sql"
	tableName := "computer"
	sqlMiddleware := NewSqlRunnerMiddleware(sqlPath, nil, tableName, false)

	framework := NewImporterOneSectionFramework(nil, stdi, WithRowRawModel(fac), WithMiddlewares(sqlMiddleware))

	err := framework.Import(path)
	if err != nil {
		t.Fatal(err)
	}

	// open the sql file
	sqlFile, err := os.Open(sqlPath)
	if err != nil {
		t.Fatal(err)
	}

	// check if match expected sql
	expectedSqls := []string{
		"INSERT INTO computer (cpu, gpu, main_board, ram, hard_disk) VALUES ('intel i5', '4060', 'rog', '海力士', 980);",
		"INSERT INTO computer (cpu, gpu, main_board, ram, hard_disk) VALUES ('amd 3501', '3060ti', '微星', '三星', 970);",
		"INSERT INTO computer (cpu, gpu, main_board, ram, hard_disk) VALUES ('intel i9', '4090', 'rog吹雪', '海力士', 990);",
	}

	sqlReader := bufio.NewReader(sqlFile)
	for _, expectedSql := range expectedSqls {
		sql, _, err := sqlReader.ReadLine()
		if err != nil {
			t.Fatal(err)
		}
		if string(sql) != expectedSql {
			t.Fatalf("sql is %s, expected %s", string(sql), expectedSql)
		}
	}

	t.Log("done")
}

type computer struct {
	CPU       string `exi:"index:0" db:"cpu"`
	GPU       string `exi:"index:1" db:"gpu"`
	MainBoard string `exi:"index:2" db:"main_board"`
	RAM       string `exi:"index:3" db:"ram"`
	HardDisk  int    `exi:"index:4" db:"hard_disk"`
}

type simpleTestDataSupportSqlRunnerMiddlewareImporter struct {
}

func (di *simpleTestDataSupportSqlRunnerMiddlewareImporter) ImportSection(tx *gorm.DB, s *RawContent) error {
	model := s.Model.(*computer)
	s.SetInsertModel(model)
	return nil
}

type computerFac struct {
}

func (mf *computerFac) GetModel() any {
	return &computer{}
}

func (mf *computerFac) MinColumnCount() int {
	return 5
}
