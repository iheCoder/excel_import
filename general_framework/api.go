package general_framework

import (
	"excel_import"
	"gorm.io/gorm"
)

type SectionChecker interface {
	// CheckValid checks the validity of the section.
	// if encounter an invalid section, return an error and record into the file.
	CheckValid(s *RawContent) error
}

type SectionImporter interface {
	// ImportSection imports the section.
	ImportSection(tx *gorm.DB, s *RawContent) error
}

// used for recognize row section
type SectionRecognizer func(s []string) RowType

type GeneralPreHandler interface {
	// PreImportHandle pre import handle
	PreImportHandle(tx *gorm.DB, s *RawContent) error
}

type SectionImportPostHandler interface {
	// PostImportSectionHandle post import handle
	PostImportSectionHandle(tx *gorm.DB, s *RawContent) error
}

type GeneralMiddleware interface {
	GeneralPreHandler
	SectionImportPostHandler
	excel_import.PostHandler
}
