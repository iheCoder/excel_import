package general_framework

import (
	"context"
	"errors"
	"excel_import"
	"excel_import/utils"
	"fmt"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

var (
	errContentCheckFailed                 = errors.New("content check failed")
	ImportFrameworkOneSectionType RowType = "import_framework_one_section"
)

type ImportFramework struct {
	db               *gorm.DB
	recorder         *util.UnexpectedRecorder
	checkers         map[RowType]SectionChecker
	importers        map[RowType]SectionImporter
	recognizer       SectionRecognizer
	postHandlers     map[RowType]excel_import.PostHandler
	rowRawModel      excel_import.RowModelFactory
	control          ImportControl
	progressReporter *util.ProgressReporter
}

func WithPostHandlers(postHandlers map[RowType]excel_import.PostHandler) OptionFunc {
	return func(framework *ImportFramework) {
		framework.postHandlers = postHandlers
	}
}

func WithCheckers(checkers map[RowType]SectionChecker) OptionFunc {
	return func(framework *ImportFramework) {
		framework.checkers = checkers
	}
}

func WithControl(control ImportControl) OptionFunc {
	return func(framework *ImportFramework) {
		framework.control = control
	}
}

func WithRowRawModel(rrm excel_import.RowModelFactory) OptionFunc {
	return func(framework *ImportFramework) {
		framework.rowRawModel = rrm
	}
}

func WithOneSectionPostHandlers(postHandler excel_import.PostHandler) OptionFunc {
	return func(framework *ImportFramework) {
		framework.postHandlers = map[RowType]excel_import.PostHandler{
			ImportFrameworkOneSectionType: postHandler,
		}
	}
}

func WithOneSectionCheckers(checker SectionChecker) OptionFunc {
	return func(framework *ImportFramework) {
		framework.checkers = map[RowType]SectionChecker{
			ImportFrameworkOneSectionType: checker,
		}
	}
}

func NewImporterFramework(db *gorm.DB, importers map[RowType]SectionImporter, recognizer SectionRecognizer, options ...OptionFunc) *ImportFramework {
	ki := &ImportFramework{
		db:               db,
		recorder:         util.NewDefaultUnexpectedRecorder(),
		importers:        importers,
		recognizer:       recognizer,
		control:          defaultImportControl,
		progressReporter: util.NewProgressReporter(true),
	}

	for _, option := range options {
		option(ki)
	}

	return ki
}

// NewImporterOneSectionFramework create a new ImportFramework with only one section type.
func NewImporterOneSectionFramework(db *gorm.DB, importer SectionImporter, options ...OptionFunc) *ImportFramework {
	importers := map[RowType]SectionImporter{
		ImportFrameworkOneSectionType: importer,
	}

	return NewImporterFramework(db, importers, func(s []string) RowType {
		return ImportFrameworkOneSectionType
	}, options...)
}

func (k *ImportFramework) Import(path string) error {
	defer k.recorder.Flush()
	defer k.progressReporter.Report()

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

func (k *ImportFramework) parseContent(path string) (*rawWhole, error) {
	content, err := util.ReadExcelContent(path)
	if err != nil {
		return nil, err
	}

	contents := k.preHandleRawContent(content)

	return k.parseRawWhole(contents)
}

func (k *ImportFramework) preHandleRawContent(contents [][]string) [][]string {
	// skip the header default
	contents = contents[k.control.StartRow:]

	// end row with func
	if k.control.Ef != nil {
		for i, content := range contents {
			if k.control.Ef(content) {
				contents = contents[:i]
				break
			}
		}
	}

	// format the content
	for i, content := range contents {
		// if the content is less than the min column count, complete it with empty string
		if len(content) < k.rowRawModel.MinColumnCount() {
			content = append(content, make([]string, k.rowRawModel.MinColumnCount()-len(content))...)
		}

		// format the cell
		for j, cell := range content {
			content[j] = util.FormatCell(cell)
		}

		contents[i] = content
	}

	return contents
}

func (k *ImportFramework) parseRawWhole(contents [][]string) (*rawWhole, error) {
	rawContents := make([]*RawContent, 0, len(contents))
	for i, content := range contents {
		// recognize the section type
		sectionType := k.recognizer(content)

		// parse the content into models
		var model any
		if k.rowRawModel != nil {
			model = k.rowRawModel.GetModel()
			if err := util.FillModelOrder(model, content); err != nil {
				return nil, err
			}
		}

		rawContents = append(rawContents, &RawContent{
			SectionType: sectionType,
			Content:     content,
			Model:       model,
			Row:         i + k.control.StartRow,
		})
	}

	return &rawWhole{
		rawContents: rawContents,
	}, nil
}

func (k *ImportFramework) checkContent(whole *rawWhole) error {
	var err error
	var checkFailed bool
	for i, rc := range whole.rawContents {
		var terr error
		if k.control.EnableTypeCheck {
			if terr = k.checkTypeError(rc); terr != nil {
				checkFailed = true
			}
		}

		sectionType := rc.SectionType
		checker, ok := k.checkers[sectionType]
		if !ok {
			continue
		}

		err = checker.CheckValid(rc)

		if err != nil || terr != nil {
			checkFailed = true
			if err = k.recorder.RecordCheckError(util.CombineErrors(i, terr, err)); err != nil {
				return err
			}
		}
	}

	if checkFailed {
		return errContentCheckFailed
	}
	return nil
}

func (k *ImportFramework) checkTypeError(rc *RawContent) error {
	// no need to check if the model is nil
	if k.rowRawModel == nil {
		return nil
	}

	// check the type of the model
	return util.CheckModelOrder(k.rowRawModel.GetModel(), rc.Content)
}

func (k *ImportFramework) importContent(whole *rawWhole) error {
	k.progressReporter.StartProgress(len(whole.rawContents))

	if k.checkAllowImportParallel() {
		return k.importContentParallel(whole)
	}

	return k.importContentSerial(whole)
}

func (k *ImportFramework) importContentParallel(whole *rawWhole) error {
	maxParallel := k.control.MaxParallel
	eg, _ := errgroup.WithContext(context.Background())
	eg.SetLimit(maxParallel)

	for _, content := range whole.rawContents {
		sectionType := content.SectionType
		importer, ok := k.importers[sectionType]
		if !ok {
			fmt.Printf("importer not found for section type: %s, content: %s \n", sectionType, content.Content)
			continue
		}

		gcontent := content
		eg.Go(func() error {
			if err := k.importSection(importer, gcontent); err != nil {
				return err
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (k *ImportFramework) importContentSerial(whole *rawWhole) error {
	for _, content := range whole.rawContents {
		sectionType := content.SectionType
		importer, ok := k.importers[sectionType]
		if !ok {
			fmt.Printf("importer not found for section type: %s, content: %s \n", sectionType, content.Content)
			continue
		}

		if err := k.importSection(importer, content); err != nil {
			return err
		}
	}

	return nil
}

func (k *ImportFramework) importSection(importer SectionImporter, content *RawContent) error {
	status := util.ProgressStatusSuccess
	defer k.progressReporter.CommitProgress(1, status)

	if err := importer.ImportSection(k.db, content); err != nil {
		status = util.ProgressStatusFailed
		fmt.Printf("import row %d section failed: %v\n", content.Row, err)
		k.recorder.RecordImportError(util.CombineErrors(content.Row, err))
		return err
	}

	return nil
}

func (k *ImportFramework) checkAllowImportParallel() bool {
	return k.control.EnableParallel && k.control.MaxParallel > 1
}

func (k *ImportFramework) postHandle() error {
	for _, handler := range k.postHandlers {
		if err := handler.PostHandle(k.db); err != nil {
			return err
		}
	}

	return nil
}
