package pipeline

import "go/ast"

var (
	importsStmt = []string{
		"errors",
		"gorm.io/gorm",
		"",
		"excel_import/general_framework",
	}
	structName = "SectionImporter"
)

type ItemResourceAstGenerator struct {
	f *ast.File
}

func (i *ItemResourceAstGenerator) AddImportDecl() {
	i.f.Decls = append(i.f.Decls, CreateImportDecl(importsStmt))
}

func (i *ItemResourceAstGenerator) AddStructDecl() {
	i.f.Decls = append(i.f.Decls, CreateStructDecl(&StructInfo{
		Name: structName,
	}))
}
