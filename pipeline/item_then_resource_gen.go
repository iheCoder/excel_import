package pipeline

import "go/ast"

var (
	importsStmt = []string{
		"errors",
		"gorm.io/gorm",
		"",
		"excel_import/general_framework",
	}
)

type ItemResourceAstGenerator struct {
	f *ast.File
}

func (i *ItemResourceAstGenerator) CreateImportDecl() {
	i.f.Decls = append(i.f.Decls, CreateImportDecl(importsStmt))
}
