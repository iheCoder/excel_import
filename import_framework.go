package excel_import

import (
	"encoding/csv"
	"errors"
	util "excel_import/utils"
	"fmt"
	"gorm.io/gorm"
	"os"
)

type RowType string
type ColumnType int

type SectionChecker interface {
	// checkValid checks the validity of the section.
	// if encounter an invalid section, return an error and record into the file.
	checkValid(s []string) error
}

type SectionImporter interface {
	// importSection imports the section.
	importSection(tx *gorm.DB, s []string) error
}

type SectionPostHandler interface {
	// postHandle post handle the section.
	postHandle(tx *gorm.DB) error
}

type sectionContent struct {
	rawContent []string
	studyOrder int
}

type ketFieldsOrder int
type sectionRecognizer func(s []string) RowType

var (
	errContentCheckFailed = errors.New("content check failed")
)

type importFramewrok struct {
	db                      *gorm.DB
	invalidSectionCsvWriter *csv.Writer
	checkers                map[RowType]SectionChecker
	importers               map[RowType]SectionImporter
	recognizer              sectionRecognizer
	postHandlers            map[RowType]SectionPostHandler
}

func WithPostHandlers(postHandlers map[RowType]SectionPostHandler) optionFunc {
	return func(ki *importFramewrok) {
		ki.postHandlers = postHandlers
	}
}

func WithCheckers(checkers map[RowType]SectionChecker) optionFunc {
	return func(ki *importFramewrok) {
		ki.checkers = checkers
	}
}

type optionFunc func(*importFramewrok)

func NewKetImporter(db *gorm.DB, importers map[RowType]SectionImporter, recognizer sectionRecognizer, options ...optionFunc) *importFramewrok {
	invalidSectionCsvWriter := initInvalidSectionCSVWriter()

	ki := &importFramewrok{
		db:                      db,
		invalidSectionCsvWriter: invalidSectionCsvWriter,
		importers:               importers,
		recognizer:              recognizer,
	}

	for _, option := range options {
		option(ki)
	}

	return ki
}

func initInvalidSectionCSVWriter() *csv.Writer {
	path := "invalid_section.csv"
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return csv.NewWriter(file)
}

type ketRawContent struct {
	sectionTypes []RowType
	content      [][]string
}

func (k *importFramewrok) Import(path string) error {
	defer k.invalidSectionCsvWriter.Flush()
	content, err := k.parseContent(path)
	if err != nil {
		fmt.Printf("read file content failed: %v\n", err)
		return err
	}

	if err = k.checkContent(content); err != nil {
		fmt.Printf("check content failed: %v\n", err)
		return err
	}

	if err = k.importContent(content); err != nil {
		fmt.Printf("import content failed: %v\n", err)
		return err
	}

	if err = k.postHandle(); err != nil {
		fmt.Printf("post handle failed: %v\n", err)
		return err
	}

	return nil
}

func (k *importFramewrok) parseContent(path string) (*ketRawContent, error) {
	content, err := util.ReadExcelContent(path)
	if err != nil {
		return nil, err
	}
	// skip the header
	content = content[1:]

	sectionTypes := make([]RowType, 0, len(content))
	for _, s := range content {
		sectionType := k.recognizer(s)
		sectionTypes = append(sectionTypes, sectionType)
	}

	return &ketRawContent{
		sectionTypes: sectionTypes,
		content:      content,
	}, nil
}

func (k *importFramewrok) checkContent(ketContents *ketRawContent) error {
	var err error
	var checkFailed bool
	for i, s := range ketContents.content {
		sectionType := ketContents.sectionTypes[i]
		checker, ok := k.checkers[sectionType]
		if !ok {
			continue
		}

		err = checker.checkValid(s)

		if err != nil {
			checkFailed = true
			if err = k.recordInvalidError(util.CombineErrors(i, err)); err != nil {
				return err
			}
		}
	}

	if checkFailed {
		return errContentCheckFailed
	}
	return nil
}

func (k *importFramewrok) importContent(ketContent *ketRawContent) error {
	for i, content := range ketContent.content {
		if i >= 10 {
			break
		}
		sectionType := ketContent.sectionTypes[i]
		importer, ok := k.importers[sectionType]
		if !ok {
			fmt.Printf("importer not found for section type: %s, content: %s \n", sectionType, content)
			continue
		}

		if err := importer.importSection(k.db, content); err != nil {
			fmt.Printf("import section failed: %v\n", err)
			return err
		}
	}

	return nil
}

func (k *importFramewrok) postHandle() error {
	for _, handler := range k.postHandlers {
		if err := handler.postHandle(k.db); err != nil {
			return err
		}
	}

	return nil
}

func (k *importFramewrok) recordInvalidError(err error) error {
	// Write the error into the csv file.
	return k.invalidSectionCsvWriter.Write([]string{err.Error()})
}
