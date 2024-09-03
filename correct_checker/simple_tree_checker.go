package correct_checker

import "excel_import"

type LeafContentExpectedItem struct {
	// ExpectedLeafModel is the expected leaf model.
	ExpectLeafModel any
	// ExpectedContentModels is the expected content models.
	ExpectContentModels []any
}

type LeafContentExpected struct {
	// Items is the leaf content expected items.
	Items []*LeafContentExpectedItem
	// TreeModel is the tree model.
	TreeModel any
	// ContentModel is the content model.
	ContentModel any
	// CheckKey is the check key for content model.
	CheckKey string

	// private fields
	treeModelAttr    []*excel_import.ExcelImportTagAttr
	contentModelAttr []*excel_import.ExcelImportTagAttr
}

// SimpleTreeChecker is the simple tree checker
// simple tree is a tree with only one parent and multiple children.
// besides the content only hang on the leaf node
type SimpleTreeChecker struct {
	// leafContentExpected is the leaf content expected.
	leafContentExpected *LeafContentExpected
}
