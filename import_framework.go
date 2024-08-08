package excel_import

import (
	"encoding/csv"
	"errors"
	"excel_import/utils"
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
	rowRawModel             RowModelFactory
	control                 importControl
}

func WithPostHandlers(postHandlers map[RowType]SectionPostHandler) optionFunc {
	return func(framework *importFramework) {
		framework.postHandlers = postHandlers
	}
}

func WithCheckers(checkers map[RowType]SectionChecker) optionFunc {
	return func(framework *importFramework) {
		framework.checkers = checkers
	}
}

func WithControl(control importControl) optionFunc {
	return func(framework *importFramework) {
		framework.control = control
	}
}

func WithRowRawModel(rrm RowModelFactory) optionFunc {
	return func(framework *importFramework) {
		framework.rowRawModel = rrm
	}
}

func NewImporterFramework(db *gorm.DB, importers map[RowType]SectionImporter, recognizer sectionRecognizer, options ...optionFunc) *importFramework {
	invalidSectionCsvWriter := initInvalidSectionCSVWriter()

	ki := &importFramework{
		db:                      db,
		invalidSectionCsvWriter: invalidSectionCsvWriter,
		importers:               importers,
		recognizer:              recognizer,
		control:                 defaultImportControl,
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

func (k *importFramework) parseContent(path string) (*rawWhole, error) {
	content, err := util.ReadExcelContent(path)
	if err != nil {
		return nil, err
	}
	// skip the header default
	content = content[k.control.startRow:]

	return k.parseRawWhole(content)
}

func (k *importFramework) parseRawWhole(contents [][]string) (*rawWhole, error) {
	rawContents := make([]*rawContent, 0, len(contents))
	for i, content := range contents {
		if k.control.ef != nil && k.control.ef(content) {
			break
		}

		// if the content is less than the min column count, complete it with empty string
		if len(content) < k.rowRawModel.minColumnCount() {
			content = append(content, make([]string, k.rowRawModel.minColumnCount()-len(content))...)
		}

		// recognize the section type
		sectionType := k.recognizer(content)

		// parse the content into models
		model := k.rowRawModel.getModel()
		if err := util.FillModelOrder(model, content); err != nil {
			return nil, err
		}

		rawContents = append(rawContents, &rawContent{
			sectionType: sectionType,
			content:     content,
			model:       model,
			row:         i + k.control.startRow,
		})
	}

	return &rawWhole{
		rawContents: rawContents,
	}, nil
}

func (k *importFramework) checkContent(whole *rawWhole) error {
	var err error
	var checkFailed bool
	for i, rc := range whole.rawContents {
		sectionType := rc.sectionType
		checker, ok := k.checkers[sectionType]
		if !ok {
			continue
		}

		err = checker.checkValid(rc)

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

func (k *importFramework) importContent(whole *rawWhole) error {
	for _, content := range whole.rawContents {
		sectionType := content.sectionType
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
