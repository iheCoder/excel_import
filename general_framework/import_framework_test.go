package general_framework

import (
	"gorm.io/gorm"
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

type simpleTestDataImporter struct {
	persons []*Person
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (di *simpleTestDataImporter) importSection(tx *gorm.DB, s *rawContent) error {
	di.persons = append(di.persons, s.model.(*Person))

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
