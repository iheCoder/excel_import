package general_framework

import (
	"context"
	"errors"
	"excel_import/utils"
	"fmt"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

var (
	errContentCheckFailed = errors.New("content check failed")
)

type importFramework struct {
	db           *gorm.DB
	recorder     *util.UnexpectedRecorder
	checkers     map[RowType]SectionChecker
	importers    map[RowType]SectionImporter
	recognizer   sectionRecognizer
	postHandlers map[RowType]SectionPostHandler
	rowRawModel  RowModelFactory
	control      importControl
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
	ki := &importFramework{
		db:         db,
		recorder:   util.NewDefaultUnexpectedRecorder(),
		importers:  importers,
		recognizer: recognizer,
		control:    defaultImportControl,
	}

	for _, option := range options {
		option(ki)
	}

	return ki
}

func (k *importFramework) Import(path string) error {
	defer k.recorder.Flush()
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

	contents := k.preHandleRawContent(content)

	return k.parseRawWhole(contents)
}

func (k *importFramework) preHandleRawContent(contents [][]string) [][]string {
	// skip the header default
	contents = contents[k.control.startRow:]

	// end row with func
	if k.control.ef != nil {
		for i, content := range contents {
			if k.control.ef(content) {
				contents = contents[:i]
				break
			}
		}
	}

	// format the content
	for i, content := range contents {
		// if the content is less than the min column count, complete it with empty string
		if len(content) < k.rowRawModel.minColumnCount() {
			content = append(content, make([]string, k.rowRawModel.minColumnCount()-len(content))...)
		}

		// format the cell
		for j, cell := range content {
			content[j] = util.FormatCell(cell)
		}

		contents[i] = content
	}

	return contents
}

func (k *importFramework) parseRawWhole(contents [][]string) (*rawWhole, error) {
	rawContents := make([]*rawContent, 0, len(contents))
	for i, content := range contents {
		// recognize the section type
		sectionType := k.recognizer(content)

		// parse the content into models
		var model any
		if k.rowRawModel != nil {
			model = k.rowRawModel.getModel()
			if err := util.FillModelOrder(model, content); err != nil {
				return nil, err
			}
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
		var terr error
		if k.control.enableTypeCheck {
			if terr = k.checkTypeError(rc); terr != nil {
				checkFailed = true
			}
		}

		sectionType := rc.sectionType
		checker, ok := k.checkers[sectionType]
		if !ok {
			continue
		}

		err = checker.checkValid(rc)

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

func (k *importFramework) checkTypeError(rc *rawContent) error {
	// no need to check if the model is nil
	if k.rowRawModel == nil {
		return nil
	}

	// check the type of the model
	return util.CheckModelOrder(k.rowRawModel.getModel(), rc.content)
}

func (k *importFramework) importContent(whole *rawWhole) error {
	if k.checkAllowImportParallel() {
		return k.importContentParallel(whole)
	}

	return k.importContentSerial(whole)
}

func (k *importFramework) importContentParallel(whole *rawWhole) error {
	maxParallel := k.control.maxParallel
	eg, _ := errgroup.WithContext(context.Background())
	eg.SetLimit(maxParallel)

	for _, content := range whole.rawContents {
		sectionType := content.sectionType
		importer, ok := k.importers[sectionType]
		if !ok {
			fmt.Printf("importer not found for section type: %s, content: %s \n", sectionType, content.content)
			continue
		}

		gcontent := content
		eg.Go(func() error {
			if err := importer.importSection(k.db, gcontent); err != nil {
				k.recorder.RecordImportError(util.CombineErrors(gcontent.row, err))
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

func (k *importFramework) importContentSerial(whole *rawWhole) error {
	for _, content := range whole.rawContents {
		sectionType := content.sectionType
		importer, ok := k.importers[sectionType]
		if !ok {
			fmt.Printf("importer not found for section type: %s, content: %s \n", sectionType, content.content)
			continue
		}

		if err := importer.importSection(k.db, content); err != nil {
			fmt.Printf("import section failed: %v\n", err)
			return util.CombineErrors(content.row, err)
		}
	}

	return nil
}

func (k *importFramework) checkAllowImportParallel() bool {
	return k.control.enableParallel && k.control.maxParallel > 1
}

func (k *importFramework) postHandle() error {
	for _, handler := range k.postHandlers {
		if err := handler.postHandle(k.db); err != nil {
			return err
		}
	}

	return nil
}
