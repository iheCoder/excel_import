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
	Name   string `excel:"index:0"`
	Career string `excel:"index:2"`
	Degree string `excel:"index:4"`
}
