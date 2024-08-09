package general_framework

import "gorm.io/gorm"

type SectionChecker interface {
	// checkValid checks the validity of the section.
	// if encounter an invalid section, return an error and record into the file.
	checkValid(s *rawContent) error
}

type SectionImporter interface {
	// importSection imports the section.
	importSection(tx *gorm.DB, s *rawContent) error
}

type SectionPostHandler interface {
	// postHandle post handle the section.
	postHandle(tx *gorm.DB) error
}

// used for recognize row section
type sectionRecognizer func(s []string) RowType
