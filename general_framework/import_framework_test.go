package general_framework

import (
	"bufio"
	"errors"
	"excel_import/correct_checker"
	util "excel_import/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"strconv"
	"testing"
)

const (
	personSectionType RowType = "person"
	defaultDBPath             = "../testdata/test.db"
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

func TestImportFramework_ImportOneSectionWithCorrectCheck(t *testing.T) {
	path := "../testdata/excel_test_resource.xlsx"
	db := initDB()
	tx := db.Begin()

	cci := &correctCheckImporter{db: db}
	fac := &resourceFac{}
	framework := NewImporterOneSectionFramework(tx, cci, WithRowRawModel(fac))

	tableModel := &ResourceTestModel{}
	countChecker := correct_checker.NewRecordCountChecker(&correct_checker.ExpectedCountChange{TablesCount: []correct_checker.TableCountInfo{
		{CountDelta: 5, TableModel: tableModel},
	}})
	partRecordChecker := correct_checker.NewPartRecordContentChecker([]*correct_checker.OffsetContentExpected{
		{
			Items: []*correct_checker.OffsetContentExpectedItem{
				{
					Offset: 1,
					ExpectedModel: &ResourceTestModel{
						Name: "A",
					},
				},
				{
					Offset: 3,
					ExpectedModel: &ResourceTestModel{
						Name: "C",
					},
				},
				{
					Offset: 5,
					ExpectedModel: &ResourceTestModel{
						Name: "E",
					},
				},
			},
			TableModel: tableModel,
			ChkKey:     "on",
		},
	})
	if err := framework.EnableCorrectnessCheck(countChecker, partRecordChecker); err != nil {
		t.Fatal(err)
	}

	err := framework.Import(path)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// check correctness
	if err = framework.CheckCorrect(); err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	tx.Commit()
	t.Log("done")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(defaultDBPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

type correctCheckImporter struct {
	db *gorm.DB
}

func (di *correctCheckImporter) ImportSection(tx *gorm.DB, s *RawContent) error {
	// get model
	model, ok := s.Model.(*resourceExcelModel)
	if !ok {
		return errors.New("model is not ResourceTestModel")
	}

	// convert into ResourceTestModel
	resource := &ResourceTestModel{
		Name:         model.Name,
		ResourceType: model.ResourceType,
	}

	// insert model
	if err := tx.Create(resource).Error; err != nil {
		return err
	}

	return nil
}

type resourceExcelModel struct {
	Name         string `exi:"index:0"`
	ResourceType int32  `exi:"index:1"`
}

type resourceFac struct {
}

func (mf *resourceFac) GetModel() any {
	return &resourceExcelModel{}
}

func (mf *resourceFac) MinColumnCount() int {
	return 2
}

type ResourceTestModel struct {
	ID           int     `json:"id" gorm:"column:id"`
	Name         string  `json:"name" gorm:"column:name" exi:"chk:on"`
	ResourceType int32   `json:"resource_type" gorm:"column:resource_type"`
	ResourceID   int64   `json:"resource_id" gorm:"column:resource_id"`
	Sort         float64 `json:"sort" gorm:"column:sort"`
}

func (r *ResourceTestModel) TableName() string {
	return "resource"
}
