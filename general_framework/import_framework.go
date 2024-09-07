package general_framework

import (
	"context"
	"errors"
	"excel_import"
	"excel_import/features"
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
	middlewares      []GeneralMiddleware
	correctCheckers  []excel_import.CorrectnessChecker

	featureMgr *features.FeatureMgr
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

func WithSimpleModelFactory(m any) OptionFunc {
	return func(framework *ImportFramework) {
		framework.rowRawModel = util.NewSimpleModelFactory(m)
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

func WithMiddlewares(middlewares ...GeneralMiddleware) OptionFunc {
	return func(framework *ImportFramework) {
		framework.middlewares = append(framework.middlewares, middlewares...)
	}
}

func WithRowFilter(filter excel_import.RowFilter) OptionFunc {
	return func(framework *ImportFramework) {
		framework.control.RowFilter = filter
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
		featureMgr:       features.NewFeatureMgr(),
	}

	for _, option := range options {
		option(ki)
	}

	if ki.control.EnableBatch {
		ki.middlewares = append(ki.middlewares, newBatchSupportFeature(ki.control.BatchSize))
	}

	if ki.control.EnableTagFormatCheck {
		ki.featureMgr.EnableTagFormatChecker()
	}

	return ki
}

func (k *ImportFramework) WithOption(option OptionFunc) *ImportFramework {
	option(k)
	return k
}

// EnableCorrectnessCheck enable the correctness check.
// must be called before Import.
func (k *ImportFramework) EnableCorrectnessCheck(correctnessCheckers ...excel_import.CorrectnessChecker) error {
	k.correctCheckers = correctnessCheckers

	for _, checker := range k.correctCheckers {
		if err := checker.PreCollect(k.db); err != nil {
			return err
		}
	}

	return nil
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

	for _, middleware := range k.middlewares {
		if err = middleware.PreImportHandle(k.db, content); err != nil {
			fmt.Printf("middleware pre handle failed: %v\n", err)
			return err
		}
	}

	if err = k.importContent(content); err != nil {
		fmt.Printf("import content failed: %v\n", err)
		return err
	}

	for _, middleware := range k.middlewares {
		if err = middleware.PostHandle(k.db); err != nil {
			fmt.Printf("middleware post handle failed: %v\n", err)
			return err
		}
	}

	if err = k.postHandle(); err != nil {
		fmt.Printf("post handle failed: %v\n", err)
		return err
	}

	return nil
}

func (k *ImportFramework) parseContent(path string) (*RawWhole, error) {
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

	// filter the rows
	if k.control.RowFilter != nil {
		contents = util.FilterRows(contents, k.control.RowFilter)
	}

	// format the content
	for i, content := range contents {
		// if the content is less than the min column count, complete it with empty string
		if len(content) < k.rowRawModel.MinColumnCount() {
			content = append(content, make([]string, k.rowRawModel.MinColumnCount()-len(content))...)
		}

		// format the cell
		fc := util.FormatCell
		if k.control.CellFormatFunc != nil {
			fc = k.control.CellFormatFunc
		}
		for j, cell := range content {
			content[j] = fc(cell)
		}

		contents[i] = content
	}

	return contents
}

func (k *ImportFramework) parseRawWhole(contents [][]string) (*RawWhole, error) {
	whole := &RawWhole{}
	// parse model tags
	var tags []*excel_import.ExcelImportTagAttr
	if k.rowRawModel != nil {
		model := k.rowRawModel.GetModel()
		tags = util.ParseTag(model)
	}

	rawContents := make([]*RawContent, 0, len(contents))
	for i, content := range contents {
		// recognize the section type
		sectionType := k.recognizer(content)

		// parse the content into models
		var model any
		if k.rowRawModel != nil {
			model = k.rowRawModel.GetModel()
			if err := util.FillModelByTags(tags, model, content); err != nil {
				return nil, err
			}
		}

		rawContents = append(rawContents, &RawContent{
			SectionType: sectionType,
			Content:     content,
			Model:       model,
			Row:         i + k.control.StartRow,
			whole:       whole,
		})
	}

	whole.rawContents = rawContents
	whole.modelInfo = &ModelsInfo{
		excelModelTags: tags,
	}
	return whole, nil
}

func (k *ImportFramework) checkContent(whole *RawWhole) error {
	var err error
	var checkFailed bool
	for i, rc := range whole.rawContents {
		var terr error
		if k.control.EnableTagFormatCheck {
			if terr = k.checkFormatError(rc); terr != nil {
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

func (k *ImportFramework) checkFormatError(rc *RawContent) error {
	// no need to check if the model is nil
	if k.rowRawModel == nil {
		return nil
	}

	// check the content format by the model tags
	return k.featureMgr.CheckContents(rc.GetContent(), rc.GetModelTags())
}

func (k *ImportFramework) importContent(whole *RawWhole) error {
	k.progressReporter.StartProgress(len(whole.rawContents))

	if k.checkAllowImportParallel() {
		return k.importContentParallel(whole)
	}

	return k.importContentSerial(whole)
}

func (k *ImportFramework) importContentParallel(whole *RawWhole) error {
	maxParallel := k.control.MaxParallel
	eg, _ := errgroup.WithContext(context.Background())
	eg.SetLimit(maxParallel)

	for _, content := range whole.rawContents {
		sectionType := content.SectionType
		importer, ok := k.importers[sectionType]
		if !ok {
			fmt.Printf("importer not found for section type: %s, content: %s \n", sectionType, content.GetContent())
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

func (k *ImportFramework) importContentSerial(whole *RawWhole) error {
	for _, content := range whole.rawContents {
		sectionType := content.SectionType
		importer, ok := k.importers[sectionType]
		if !ok {
			fmt.Printf("importer not found for section type: %s, content: %s \n", sectionType, content.GetContent())
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

	// middleware post handle
	for _, middleware := range k.middlewares {
		if err := middleware.PostImportSectionHandle(k.db, content); err != nil {
			fmt.Printf("middleware post import section handle failed: %v\n", err)
			return err
		}
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

func (k *ImportFramework) CheckCorrect() error {
	for _, checker := range k.correctCheckers {
		if err := checker.CheckCorrect(k.db); err != nil {
			return err
		}
	}

	return nil
}
