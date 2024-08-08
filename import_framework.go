package excel_import

import (
	"encoding/csv"
	"errors"
	util "excel_import/utils"
	"fmt"
	"gorm.io/gorm"
	"os"
)

var (
	errContentCheckFailed = errors.New("content check failed")
)

type importFramework struct {
	db                      *gorm.DB
	invalidSectionCsvWriter *csv.Writer
	checkers                map[RowType]SectionChecker
	importers               map[RowType]SectionImporter
	recognizer              sectionRecognizer
	postHandlers            map[RowType]SectionPostHandler
}

func WithPostHandlers(postHandlers map[RowType]SectionPostHandler) optionFunc {
	return func(ki *importFramework) {
		ki.postHandlers = postHandlers
	}
}

func WithCheckers(checkers map[RowType]SectionChecker) optionFunc {
	return func(ki *importFramework) {
		ki.checkers = checkers
	}
}

func NewKetImporter(db *gorm.DB, importers map[RowType]SectionImporter, recognizer sectionRecognizer, options ...optionFunc) *importFramework {
	invalidSectionCsvWriter := initInvalidSectionCSVWriter()

	ki := &importFramework{
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

func (k *importFramework) Import(path string) error {
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

func (k *importFramework) parseContent(path string) (*rawContent, error) {
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

	return &rawContent{
		sectionTypes: sectionTypes,
		content:      content,
	}, nil
}

func (k *importFramework) checkContent(ketContents *rawContent) error {
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

func (k *importFramework) importContent(ketContent *rawContent) error {
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

func (k *importFramework) postHandle() error {
	for _, handler := range k.postHandlers {
		if err := handler.postHandle(k.db); err != nil {
			return err
		}
	}

	return nil
}

func (k *importFramework) recordInvalidError(err error) error {
	// Write the error into the csv file.
	return k.invalidSectionCsvWriter.Write([]string{err.Error()})
}
